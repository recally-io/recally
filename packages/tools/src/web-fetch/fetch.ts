import { checkUrlTarget } from "@recally/capture";
import { AppError, DEFAULT_LIMITS, sha256Hex } from "@recally/domain";

// web_fetch (§6.2): ONE controlled HTTP GET + deterministic extraction +
// evidence persistence. It never starts a browser, never calls a secondary
// fetch service, never declares capture success — the agent decides.

export interface FetchedSource {
  finalUrl: string;
  status: number;
  contentType: string;
  body: Uint8Array;
  sha256: string;
}

export async function safeFetch(
  rawUrl: string,
  fetchFn: typeof fetch,
  opts: { maxBytes?: number; timeoutMs?: number; maxRedirects?: number } = {},
): Promise<FetchedSource> {
  const maxBytes = opts.maxBytes ?? DEFAULT_LIMITS.fetch.maxHtmlBytes;
  const timeoutMs = opts.timeoutMs ?? DEFAULT_LIMITS.fetch.timeoutMs;
  const maxRedirects = opts.maxRedirects ?? DEFAULT_LIMITS.fetch.maxRedirects;

  let url = checkUrlTarget(rawUrl).toString();

  for (let hop = 0; hop <= maxRedirects; hop++) {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort("timeout"), timeoutMs);
    let res: Response;

    try {
      res = await fetchFn(url, {
        signal: controller.signal,
        redirect: "manual",
        headers: {
          "user-agent": "recally-archive/0.1 (+personal reading archive)",
          accept: "text/html,application/xhtml+xml,application/pdf;q=0.9,*/*;q=0.5",
          // Only advertise encodings we can decode ourselves (workerd's
          // DecompressionStream has no br); otherwise brotli bytes reach the
          // parser as garbage.
          "accept-encoding": "gzip, deflate",
        },
      });
    } catch (err) {
      if (controller.signal.aborted) {
        throw new AppError("source_unavailable", `fetch timeout after ${timeoutMs}ms`, {
          retryable: true,
        });
      }

      throw new AppError(
        "source_unavailable",
        `fetch failed: ${err instanceof Error ? err.message : String(err)}`,
        { retryable: true },
      );
    } finally {
      clearTimeout(timer);
    }

    if (res.status >= 300 && res.status < 400) {
      const location = res.headers.get("location");

      if (!location) {
        throw new AppError("source_unavailable", `redirect ${res.status} without location`);
      }

      url = checkUrlTarget(new URL(location, url).toString()).toString();
      continue;
    }

    const contentType = res.headers.get("content-type") ?? "";
    let body = await readLimited(res, maxBytes);
    // A still-set content-encoding means the runtime did not decompress for
    // us — do it ourselves so downstream parsers see real bytes.
    const encoding = res.headers.get("content-encoding")?.toLowerCase();

    if (encoding === "gzip" || encoding === "deflate") {
      const ds = new DecompressionStream(encoding);
      body = new Uint8Array(
        await new Response(new Blob([new Uint8Array(body)]).stream().pipeThrough(ds)).arrayBuffer(),
      );
    }

    // brotli has no DecompressionStream in workerd — leave it as evidence and
    // let extraction report unavailable instead of producing garbage.
    return {
      finalUrl: url,
      status: res.status,
      contentType,
      body,
      sha256: await sha256Hex(body),
    };
  }

  throw new AppError("policy_denied", `more than ${maxRedirects} redirects`);
}

async function readLimited(res: Response, maxBytes: number): Promise<Uint8Array> {
  const reader = res.body?.getReader();

  if (!reader) return new Uint8Array(await res.arrayBuffer());
  const parts: Uint8Array[] = [];
  let total = 0;

  for (;;) {
    const { done, value } = await reader.read();

    if (done) break;
    total += value.byteLength;

    if (total > maxBytes) {
      await reader.cancel();
      throw new AppError("incomplete_content", `body exceeded ${maxBytes} bytes`, {
        nextAction: "archive_partial",
      });
    }

    parts.push(value);
  }

  const out = new Uint8Array(total);
  let offset = 0;

  for (const c of parts) {
    out.set(c, offset);
    offset += c.byteLength;
  }

  return out;
}

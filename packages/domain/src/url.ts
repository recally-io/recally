import { sha256Hex } from "./hash";

// Query params that never carry content meaning. Deliberately conservative:
// content-bearing params and fragments are preserved (plan §4.3).
const TRACKING_PARAMS = new Set([
  "fbclid",
  "gclid",
  "gclsrc",
  "dclid",
  "msclkid",
  "yclid",
  "twclid",
  "mc_cid",
  "mc_eid",
  "igshid",
  "spm",
  "scm",
  "_ga",
  "_hsenc",
  "_hsmi",
  "ref_src",
  "ref_url",
  "li_fat_id",
  "wickedid",
  "vero_id",
]);

export function normalizeUrl(raw: string): string {
  const url = new URL(raw);
  url.protocol = url.protocol.toLowerCase();

  if (url.protocol !== "http:" && url.protocol !== "https:") {
    throw new Error(`unsupported url scheme ${url.protocol}`);
  }

  url.hostname = url.hostname.toLowerCase();

  if (
    (url.protocol === "https:" && url.port === "443") ||
    (url.protocol === "http:" && url.port === "80")
  ) {
    url.port = "";
  }

  const kept: Array<[string, string]> = [];
  url.searchParams.forEach((value, key) => {
    if (!key.startsWith("utm_") && !TRACKING_PARAMS.has(key)) kept.push([key, value]);
  });
  kept.sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0));
  url.search = "";

  for (const [key, value] of kept) url.searchParams.append(key, value);

  return url.toString();
}

export async function normalizedUrlHash(normalizedUrl: string): Promise<string> {
  return sha256Hex(normalizedUrl);
}

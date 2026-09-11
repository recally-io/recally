// Thin fetch client against the app worker. Dev bypass token comes from
// .dev.vars (DEV_BYPASS_TOKEN) via a dev-only header — real auth is Access JWT.
const headers = () => {
  const h: Record<string, string> = { "content-type": "application/json" };
  const t = import.meta.env.VITE_API_TOKEN;
  if (t) h.authorization = `Bearer ${t}`;
  return h;
};

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { ...headers(), ...(init?.headers ?? {}) },
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error((body as { message?: string }).message ?? `${res.status}`);
  }
  return res.json() as Promise<T>;
}

export interface Item {
  id: string;
  original_url: string;
  title: string | null;
  saved_at: string;
  read_status: string;
  current_snapshot_id: string | null;
}

export const api = {
  listItems: (cursor?: string) =>
    req<{ items: Item[]; next_cursor: string | null }>(
      `/api/v1/items${cursor ? `?cursor=${encodeURIComponent(cursor)}` : ""}`,
    ),
  getItem: (id: string) => req<Item>(`/api/v1/items/${id}`),
  saveUrl: (url: string, note?: string) =>
    req<{ item_id: string; job_id: string; status_url: string }>("/api/v1/items", {
      method: "POST",
      headers: { "Idempotency-Key": crypto.randomUUID() },
      body: JSON.stringify({
        source: { kind: "url", url },
        note: note ?? undefined,
      }),
    }),
  saveText: (text: string, sourceUrl?: string) =>
    req<{ item_id: string; job_id: string }>("/api/v1/items", {
      method: "POST",
      headers: { "Idempotency-Key": crypto.randomUUID() },
      body: JSON.stringify({
        source: { kind: "manual", text, source_url: sourceUrl },
      }),
    }),
  getJob: (id: string) =>
    req<{
      id: string;
      status: string;
      result: string | null;
      error: string | null;
    }>(`/api/v1/jobs/${id}`),
  recapture: (itemId: string) =>
    req<{ job_id: string }>(`/api/v1/items/${itemId}/captures`, {
      method: "POST",
    }),
  search: (q: string) =>
    req<{
      results: Array<{
        item_id: string;
        title: string | null;
        snippet: string;
      }>;
    }>(`/api/v1/search?q=${encodeURIComponent(q)}`),
  article: (revisionId: string) =>
    req<{ article_md: string; content_revision_id: string }>(`/api/v1/content/${revisionId}`),
};

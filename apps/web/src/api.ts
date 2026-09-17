// Thin fetch client against the app worker. Requests without a token resolve
// to the default library; VITE_API_TOKEN sets an optional rcl_ Bearer token.
const headers = () => {
  const h: Record<string, string> = { "content-type": "application/json" };
  const t = import.meta.env.VITE_API_TOKEN;

  if (t) h.authorization = `Bearer ${t}`;

  return h;
};

export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
  ) {
    super(message);
  }
}

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { ...headers(), ...init?.headers },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    const b = body as { message?: string; code?: string; error?: string };
    throw new ApiError(res.status, b.code ?? b.error ?? "", b.message ?? `${res.status}`);
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
  capture_generation: number;
  content_quality: string | null;
  resource_quality: string | null;
  latest_job: { id: string; kind: string; status: string } | null;
}

export interface Snapshot {
  id: string;
  captured_at: string;
  capture_method: string;
  content_quality: string;
  resource_quality: string;
  final_url: string;
  content_revision_id: string | null;
}

export interface Note {
  id: string;
  body: string;
  version: number;
  updated_at: string;
}

export interface Artifact {
  id: string;
  type: string;
  status: string;
  output: string;
  created_at: string;
}

export interface SummaryOutput {
  short_summary?: string;
  key_points?: string[];
  topics?: string[];
  entities?: string[];
  caveats?: string[];
}

export interface ItemDetail {
  item: Item;
  snapshots: Snapshot[];
  notes: Note[];
  artifacts: Artifact[];
}

export interface Job {
  id: string;
  kind: string;
  status: string;
  item_id: string;
  attempt_count: number;
  error: string | null;
  created_at: string;
  updated_at: string;
  item_title?: string | null;
}

export interface JobEvent {
  sequence: number;
  kind: string;
  reason_code: string | null;
  summary: string | null;
  created_at: string;
}

export interface JobDetail extends Job {
  outcome_code: string | null;
  events: JobEvent[];
}

export interface SearchHit {
  item_id: string;
  title: string | null;
  original_url: string;
  snippet: string;
  chunk_id?: string;
}

export interface Me {
  library_id: string;
  library_name: string;
  actor: string;
  email: string | null;
  via: string;
}

export interface ApiToken {
  id: string;
  name: string;
  scopes: string;
  last_used_at: string | null;
  expires_at: string | null;
  created_at: string;
}

export const api = {
  me: () => req<Me>("/api/v1/me"),
  listTokens: () => req<{ tokens: ApiToken[] }>("/api/v1/tokens"),
  createToken: (name: string, scopes: string[], expires_in_days: number | null) =>
    req<{ token_id: string; token: string }>("/api/v1/tokens", {
      method: "POST",
      body: JSON.stringify({ name, scopes, expires_in_days }),
    }),
  revokeToken: (id: string) => req<{ ok: true }>(`/api/v1/tokens/${id}`, { method: "DELETE" }),
  listItems: (cursor?: string) =>
    req<{ items: Item[]; next_cursor: string | null }>(
      `/api/v1/items${cursor ? `?cursor=${encodeURIComponent(cursor)}` : ""}`,
    ),
  getItem: (id: string) => req<ItemDetail>(`/api/v1/items/${id}`),
  saveUrl: (url: string, note?: string) =>
    req<{ item_id: string; job_id: string; status_url: string }>("/api/v1/items", {
      method: "POST",
      headers: { "Idempotency-Key": crypto.randomUUID() },
      body: JSON.stringify({
        source: { kind: "url", url },
        note: note ?? undefined,
      }),
    }),
  patchItem: (id: string, body: { title?: string; read_status?: string }) =>
    req<{ ok: true }>(`/api/v1/items/${id}`, {
      method: "PATCH",
      body: JSON.stringify(body),
    }),
  deleteItem: (id: string) => req<{ ok: true }>(`/api/v1/items/${id}`, { method: "DELETE" }),
  recapture: (itemId: string) =>
    req<{ job_id: string }>(`/api/v1/items/${itemId}/captures`, { method: "POST" }),
  listJobs: () => req<{ jobs: Job[] }>("/api/v1/jobs"),
  getJob: (id: string) => req<JobDetail>(`/api/v1/jobs/${id}`),
  cancelJob: (id: string) => req<{ ok: true }>(`/api/v1/jobs/${id}/cancel`, { method: "POST" }),
  search: (q: string) =>
    req<{ results: SearchHit[]; mode: string }>(`/api/v1/search?q=${encodeURIComponent(q)}`),
  article: (revisionId: string) =>
    req<{ article_md: string | null; blocks: unknown[] }>(`/api/v1/content/${revisionId}`),
  addNote: (itemId: string, body: string) =>
    req<{ note_id: string; version: number }>(`/api/v1/items/${itemId}/notes`, {
      method: "POST",
      body: JSON.stringify({ body }),
    }),
  patchNote: (id: string, body: string, version: number) =>
    req<{ ok: true; version: number }>(`/api/v1/notes/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ body, version }),
    }),
  readingEvent: (itemId: string, kind: "open" | "progress" | "finish") =>
    req<{ ok: true }>(`/api/v1/items/${itemId}/reading-events`, {
      method: "POST",
      body: JSON.stringify({ kind }),
    }),
  createShare: (itemId: string, snapshotId: string, contentRevisionId?: string | null) =>
    req<{ share_id: string; url: string; expires_at: string | null }>("/api/v1/shares", {
      method: "POST",
      body: JSON.stringify({
        item_id: itemId,
        snapshot_id: snapshotId,
        content_revision_id: contentRevisionId ?? undefined,
        include_full_text: Boolean(contentRevisionId),
        allowed_artifact_ids: [],
      }),
    }),
};

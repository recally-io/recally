import { itemListSchema, itemViewSchema, jobViewSchema, type JobView } from "@recally/contracts";
import { z } from "zod";

// Thin fetch client against the app worker. Requests without a token resolve
// to the default library; VITE_API_TOKEN sets an optional rcl_ Bearer token.
const headers = () => {
  const h: Record<string, string> = { "content-type": "application/json" };
  const t = import.meta.env.VITE_API_TOKEN;

  if (t) h.authorization = `Bearer ${t}`;

  return h;
};

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { ...headers(), ...init?.headers },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));

    const parsed = z.object({ message: z.string().optional() }).safeParse(body);

    const b = parsed.success ? parsed.data : {};
    throw new Error(b.message ?? `${res.status}`);
  }

  return res.json();
}

export function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

const snapshotSchema = z.object({
  id: z.string(),
  captured_at: z.string(),
  capture_method: z.string(),
  content_quality: z.string(),
  resource_quality: z.string(),
  final_url: z.string(),
  content_revision_id: z.string().nullable(),
});

const noteSchema = z.object({ id: z.string(), body: z.string(), version: z.number() });

const artifactSchema = z.object({ type: z.string(), status: z.string(), output: z.string() });

// Item detail returns the stored row, without the list endpoint's quality/job joins.
const itemDetailSchema = z.object({
  item: itemViewSchema
    .omit({ content_quality: true, resource_quality: true })
    .partial({ latest_job: true }),
  snapshots: z.array(snapshotSchema),
  notes: z.array(noteSchema),
  artifacts: z.array(artifactSchema),
});

const jobSchema = jobViewSchema.omit({ updated_at: true, outcome_code: true });

const jobListSchema = z.object({
  jobs: z.array(jobSchema.extend({ item_title: z.string().nullable().optional() })),
});

const jobDetailSchema = jobSchema.extend({
  events: z.array(
    z.object({
      sequence: z.number(),
      kind: z.string(),
      reason_code: z.string().nullable(),
      summary: z.string().nullable(),
    }),
  ),
});

export type Snapshot = z.infer<typeof snapshotSchema>;

export type Note = z.infer<typeof noteSchema>;

export type ItemDetail = z.infer<typeof itemDetailSchema>;

export type JobListEntry = Omit<JobView, "updated_at" | "outcome_code"> & {
  item_title?: string | null | undefined;
};

export type JobDetail = z.infer<typeof jobDetailSchema>;

export interface SearchHit {
  item_id: string;
  title: string | null;
  original_url: string;
  snippet: string;
  chunk_id?: string;
}

export interface Me {
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
    req(`/api/v1/items${cursor ? `?cursor=${encodeURIComponent(cursor)}` : ""}`).then((body) =>
      itemListSchema.parse(body),
    ),
  getItem: (id: string) => req(`/api/v1/items/${id}`).then((body) => itemDetailSchema.parse(body)),
  saveUrl: (url: string, note?: string) =>
    req<{ item_id: string }>("/api/v1/items", {
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
  recapture: (itemId: string) =>
    req<{ job_id: string }>(`/api/v1/items/${itemId}/captures`, { method: "POST" }),
  listJobs: () => req("/api/v1/jobs").then((body) => jobListSchema.parse(body)),
  getJob: (id: string) => req(`/api/v1/jobs/${id}`).then((body) => jobDetailSchema.parse(body)),
  cancelJob: (id: string) => req<{ ok: true }>(`/api/v1/jobs/${id}/cancel`, { method: "POST" }),
  search: (q: string) =>
    req<{ results: SearchHit[]; mode: string }>(`/api/v1/search?q=${encodeURIComponent(q)}`),
  article: (revisionId: string) =>
    req<{ article_md: string | null }>(`/api/v1/content/${revisionId}`),
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

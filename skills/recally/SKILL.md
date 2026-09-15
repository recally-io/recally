---
name: recally
description: Save URLs to and search a Recally personal reading archive over HTTP. Use when the user wants to archive a link, fetch archived article text, search their saved reading, or add notes to saved items. Requires a Recally API key (rcl_…).
---

# Recally API

Recally is a personal reading archive: it stores the actual archived text of saved
URLs (verified against the source), not just bookmarks. All endpoints are JSON
under the base URL `https://recally.io/api/v1` (or the deployment's own origin).

## Auth

Every request needs a Bearer token created in the app's Settings page:

```
Authorization: Bearer rcl_xxxxxxxx…
Content-Type: application/json
```

Tokens are scoped. An agent typically wants `items:read`, `items:write`,
`search:read`; add `notes:write` only if the user asks you to write notes.

## Save a URL

```
POST /api/v1/items
Idempotency-Key: <random uuid>          # REQUIRED — retries must reuse it

{"source": {"kind": "url", "url": "https://example.com/post"}, "note": "optional"}
```

→ `201 {"item_id", "job_id", "status_url"}`

Capture is asynchronous — an agent fetches the page, verifies it, and commits a
snapshot. Poll the job:

```
GET /api/v1/jobs/<job_id>
```

→ `{"status": "queued|running|succeeded|failed", "events": […]}`

`failed` is normal for paywalled/dead pages; report the `error` field honestly.
A succeeded capture may still be `partial` (some content unreachable) — check
`content_quality` on the item.

## List items

```
GET /api/v1/items?cursor=<cursor>&limit=50
```

→ `{"items": [{"id","title","original_url","saved_at","read_status",
"content_quality","latest_job"}], "next_cursor"}`

## Item detail

```
GET /api/v1/items/<item_id>
```

→ `{"item", "snapshots", "notes", "artifacts"}`

`snapshots[]` carries `content_revision_id` — fetch the archived article with:

```
GET /api/v1/content/<content_revision_id>
```

→ `{"article_md", "blocks"}` — markdown + JSONL blocks. This is the committed
snapshot text. Do NOT fetch the live URL instead and present it as the archive.

## Search

```
GET /api/v1/search?q=<query>
```

→ `{"results": [{"item_id","title","original_url","snippet"}], "mode":"fts"}`

Full-text over committed snapshots only. Snippets contain `<b>` markers.
Summaries/ask-with-citations land in a later milestone.

## Notes

```
POST /api/v1/items/<item_id>/notes        {"body": "…"}
PATCH /api/v1/notes/<note_id>             {"body": "…", "version": <int>}
```

PATCH is version-checked; on `conflict` re-read the item and retry once.

## Other actions

- `POST /api/v1/items/<id>/captures` — re-capture (new snapshot, new generation)
- `PATCH /api/v1/items/<id>` — `{"read_status": "read"|"unread", "title": …}`
- `POST /api/v1/jobs/<id>/cancel` — cancel a running job
- `POST /api/v1/shares` — public share link for a snapshot

## Errors

`{"error", "message", "code?"}` — `forbidden` = token missing scope,
`conflict` = stale version/job state,
`not_found` = deleted or never existed.

# Source type signals

## Document-like (usually `web_fetch` first)

Blogs, docs, news articles, RFCs, specs, wiki pages. The HTML usually contains the full text.

## App-like records (check `site_list` first)

A URL that denotes a structured object more than a page:

- `github.com/{owner}/{repo}` → repo record (metadata, README, files list)
- `github.com/.../issues|pull/N` → issue/PR record
- `news.ycombinator.com/item?id=...` → HN item + comments
- Social posts/statuses/threads

The adapter returns a structured record with provenance. If its output is thin or errors, fall back to `web_fetch` or `browser_*` — do not treat adapter output as gospel.

## Dynamic pages (expect `browser_*`)

- Obvious app shells (SPA mounts, empty body + bundle scripts).
- Content behind expand/scroll interaction.
- Multi-page articles where pagination requires JS.

## Not worth an automatic fetch

- User-submitted pasted content — archive it directly.
- Schemes/types outside http(s) or clearly unsupported MIME — `finish(outcome="failed")` early.

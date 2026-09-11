# Completeness heuristics

Judge completeness from evidence, not from length or status codes.

## Signals a fetch result is probably the real article

- Title candidates match the URL's apparent topic.
- Body starts with a lede/first paragraph, not navigation.
- Body ends with a conclusion, signature, or footer — not mid-sentence.
- Code blocks and tables have structure, not truncated fragments.
- Link list is plausible for the site (few, in-article).

## Signals something is missing

- Body is a paragraph-long summary while the site clearly hosts full text.
- "Read more", "Show more", "展开全文" controls observed.
- Content ends at an apparent paywall boundary or "Sign in to continue".
- Page is a JS shell: near-empty body, large scripts, app mount points.
- Article visibly continues on a second page ("Page 2", print view).

## Signals to stop

- HTTP 401/403, captcha interstitial, Cloudflare challenge page.
- Login redirect or "sign in to read".
- Repeated identical observations across actions.

When in doubt between `partial` and `complete`, choose `partial` and list the gap.

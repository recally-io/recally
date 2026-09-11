import { Link } from "@tanstack/react-router";

export function LandingPage() {
  return (
    <div className="min-h-screen bg-paper">
      <header className="mx-auto flex max-w-5xl items-center justify-between px-6 py-5">
        <span className="font-serif text-[17px] font-bold">
          Re<span className="text-accent">call</span>ly
        </span>
        <Link
          to="/app"
          className="rounded-lg bg-accent px-4 py-2 text-[13px] font-semibold text-white hover:opacity-90"
        >
          Open app →
        </Link>
      </header>

      <main className="mx-auto max-w-5xl px-6">
        <section className="py-20 text-center">
          <p className="mb-4 font-mono text-[11px] font-semibold tracking-[0.15em] text-accent uppercase">
            personal reading archive
          </p>
          <h1 className="mx-auto max-w-2xl font-serif text-[42px] leading-[1.15] tracking-tight">
            Save the link.
            <br />
            Keep the text.
          </h1>
          <p className="mx-auto mt-5 max-w-xl text-[15px] leading-relaxed text-ink-2">
            Recally archives what you read — not a bookmark, the actual content, verified against
            its source and searchable forever. Built for humans, readable by agents.
          </p>
          <div className="mt-8 flex items-center justify-center gap-3">
            <Link
              to="/app"
              className="rounded-lg bg-ink px-5 py-2.5 text-[14px] font-semibold text-white hover:opacity-90"
            >
              Sign in
            </Link>
            <a
              href="/skill.md"
              className="rounded-lg border border-line bg-surface px-5 py-2.5 text-[14px] text-ink-2 hover:border-ink-3"
            >
              Agent API ↗
            </a>
          </div>
        </section>

        <section className="grid gap-px overflow-hidden rounded-xl border border-line bg-line sm:grid-cols-3">
          {[
            {
              n: "01",
              t: "Capture",
              d: "Paste a URL. A capture agent picks the right path — direct fetch, a site adapter, or a real browser — and commits a verified snapshot.",
            },
            {
              n: "02",
              t: "Keep",
              d: "The archived text is immutable and rendered from source evidence, never model output. Partial captures say so; nothing is silently synthetic.",
            },
            {
              n: "03",
              t: "Recall",
              d: "Full-text search over every committed snapshot, with honest quality labels and provenance down to the run that produced it.",
            },
          ].map((s) => (
            <div key={s.n} className="bg-surface p-7">
              <p className="font-mono text-[11px] text-ink-3">{s.n}</p>
              <h2 className="mt-2 mb-2 font-serif text-lg">{s.t}</h2>
              <p className="text-[13px] leading-relaxed text-ink-2">{s.d}</p>
            </div>
          ))}
        </section>

        <section className="mt-12 overflow-hidden rounded-xl border border-line bg-ink text-paper">
          <div className="p-7">
            <p className="font-mono text-[11px] font-semibold tracking-[0.15em] text-accent-soft uppercase">
              for agents
            </p>
            <h2 className="mt-2 mb-3 font-serif text-xl">Your archive, in your agent's hands</h2>
            <p className="max-w-2xl text-[13.5px] leading-relaxed text-paper/80">
              Create an API key in Settings, then point your agent at the skill file — it learns how
              to save URLs, read archived content, and search your library.
            </p>
          </div>
          <pre className="overflow-x-auto border-t border-white/10 bg-black/30 p-5 font-mono text-[12px] leading-relaxed text-paper/90">
            {`curl https://recally.io/skill.md

curl -H "Authorization: Bearer rcl_…" \\
     -d '{"source":{"kind":"url","url":"https://…"}}' \\
     https://recally.io/api/v1/items`}
          </pre>
        </section>

        <footer className="py-10 text-center text-[12px] text-ink-3">
          Runs entirely on Cloudflare — Workers, D1, R2, Vectorize.
        </footer>
      </main>
    </div>
  );
}

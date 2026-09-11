import { Link } from "@tanstack/react-router";
import { useCallback, useEffect, useRef, useState } from "react";
import { api, type SearchHit } from "../api";
import { DomainChip, Snippet } from "../components/bits";

export function SearchPage() {
  const [q, setQ] = useState("");
  const [hits, setHits] = useState<SearchHit[]>([]);
  const [searched, setSearched] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const run = useCallback(async (query: string) => {
    if (!query.trim()) return;
    setError(null);
    try {
      const r = await api.search(query.trim());
      setHits(r.results);
      setSearched(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }, []);

  // ?q= deep link.
  useEffect(() => {
    inputRef.current?.focus();
    const initial = new URLSearchParams(location.search).get("q");
    if (initial) {
      setQ(initial);
      void run(initial);
    }
  }, [run]);

  return (
    <div className="mx-auto max-w-3xl px-8 py-10">
      <form
        className="mb-2 flex items-center gap-3 border-b-2 border-ink pb-2.5"
        onSubmit={(e) => {
          e.preventDefault();
          void run(q);
        }}
      >
        <input
          ref={inputRef}
          className="w-full bg-transparent font-serif text-[26px] leading-tight text-ink outline-none placeholder:text-ink-3"
          placeholder="Search your archive…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <kbd className="font-mono text-[11px] whitespace-nowrap text-ink-3">⏎ search</kbd>
      </form>
      <p className="mb-8 text-[11.5px] text-ink-3">
        Full-text over committed snapshots only. Ask-with-citations lands with M4.
      </p>
      {error && <p className="mb-4 text-sm text-danger">{error}</p>}
      <ul>
        {hits.map((h) => (
          <li key={h.chunk_id ?? h.item_id} className="border-b border-line py-3.5">
            <Link
              to="/app/items/$itemId"
              params={{ itemId: h.item_id }}
              className="text-[14px] font-semibold hover:text-accent"
            >
              {h.title ?? h.original_url}
            </Link>
            <div className="mt-0.5 text-[11.5px] text-ink-3">
              <DomainChip url={h.original_url} />
            </div>
            <p className="mt-1.5 text-[12.5px] leading-relaxed">
              <Snippet text={h.snippet} />
            </p>
          </li>
        ))}
      </ul>
      {searched && hits.length === 0 && (
        <p className="py-12 text-center text-sm text-ink-3">No matches in the archive.</p>
      )}
    </div>
  );
}

import { Link } from "@tanstack/react-router";
import { useState } from "react";
import { api } from "../api";

export function SearchPage() {
  const [q, setQ] = useState("");
  const [results, setResults] = useState<
    Array<{ item_id: string; title: string | null; snippet: string }>
  >([]);
  const [error, setError] = useState<string | null>(null);

  const run = async () => {
    setError(null);
    try {
      const r = await api.search(q);
      setResults(r.results);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  };

  return (
    <div>
      <form
        className="mb-6 flex gap-2"
        onSubmit={(e) => {
          e.preventDefault();
          void run();
        }}
      >
        <input
          className="flex-1 rounded border px-3 py-2 text-sm"
          placeholder="Search your archive…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <button type="submit" className="rounded bg-black px-4 py-2 text-sm text-white">
          search
        </button>
      </form>
      {error && <p className="text-sm text-red-600">{error}</p>}
      <ul className="divide-y">
        {results.map((r) => (
          <li key={r.item_id} className="py-3">
            <Link
              to="/items/$itemId"
              params={{ itemId: r.item_id }}
              className="font-medium hover:underline"
            >
              {r.title ?? r.item_id}
            </Link>
            <p className="mt-0.5 text-sm text-neutral-600">{r.snippet}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}

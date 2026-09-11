import { Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { api, type Item } from "../api";

export function LibraryPage() {
  const [items, setItems] = useState<Item[]>([]);
  const [url, setUrl] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    api
      .listItems()
      .then((r) => setItems(r.items))
      .catch((e) => setError(e.message));
  }, []);

  const refresh = () =>
    api
      .listItems()
      .then((r) => setItems(r.items))
      .catch((e) => setError(e.message));

  const save = async () => {
    if (!url.trim()) return;
    setBusy(true);
    setError(null);
    try {
      await api.saveUrl(url.trim());
      setUrl("");
      setTimeout(refresh, 500);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div>
      <form
        className="mb-6 flex gap-2"
        onSubmit={(e) => {
          e.preventDefault();
          void save();
        }}
      >
        <input
          className="flex-1 rounded border px-3 py-2 text-sm"
          placeholder="Paste a URL to archive…"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button
          type="submit"
          className="rounded bg-black px-4 py-2 text-sm text-white disabled:opacity-50"
          disabled={busy}
        >
          save
        </button>
      </form>
      {error && <p className="mb-4 text-sm text-red-600">{error}</p>}
      <ul className="divide-y">
        {items.map((it) => (
          <li key={it.id} className="py-3">
            <Link
              to="/items/$itemId"
              params={{ itemId: it.id }}
              className="font-medium hover:underline"
            >
              {it.title ?? it.original_url}
            </Link>
            <div className="mt-0.5 text-xs text-neutral-500">
              {new Date(it.saved_at).toLocaleString()} · {it.read_status}
              {it.current_snapshot_id ? " · archived" : " · capturing…"}
            </div>
          </li>
        ))}
        {items.length === 0 && (
          <li className="py-8 text-center text-sm text-neutral-400">nothing saved yet</li>
        )}
      </ul>
    </div>
  );
}

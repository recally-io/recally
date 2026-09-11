import { useParams } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { api, type Item } from "../api";

export function ItemPage() {
  const { itemId } = useParams({ strict: false }) as { itemId: string };
  const [item, setItem] = useState<Item | null>(null);
  const [article, setArticle] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api
      .getItem(itemId)
      .then(setItem)
      .catch((e) => setError(e.message));
  }, [itemId]);

  useEffect(() => {
    // Latest committed revision id rides along on the item detail response in a
    // later milestone; for now the reading pane loads lazily once present.
    const rev = (item as unknown as { current_content_revision_id?: string })
      ?.current_content_revision_id;
    if (rev)
      api
        .article(rev)
        .then((a) => setArticle(a.article_md))
        .catch(() => {});
  }, [item]);

  if (error) return <p className="text-sm text-red-600">{error}</p>;
  if (!item) return <p className="text-sm text-neutral-400">loading…</p>;

  return (
    <div>
      <div className="mb-4">
        <h1 className="text-lg font-semibold">{item.title ?? item.original_url}</h1>
        <a
          className="text-xs text-neutral-500 underline"
          href={item.original_url}
          target="_blank"
          rel="noreferrer"
        >
          {item.original_url}
        </a>
        <div className="mt-1 text-xs text-neutral-500">
          saved {new Date(item.saved_at).toLocaleString()} · {item.read_status}
        </div>
        <button
          type="button"
          className="mt-2 rounded border px-3 py-1 text-xs"
          onClick={() => void api.recapture(itemId)}
        >
          re-capture
        </button>
      </div>
      {article ? (
        <article className="prose whitespace-pre-wrap text-sm leading-6">{article}</article>
      ) : (
        <p className="text-sm text-neutral-400">
          {item.current_snapshot_id
            ? "loading article…"
            : "archive in progress — check back shortly"}
        </p>
      )}
    </div>
  );
}

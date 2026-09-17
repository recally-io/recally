import { useEffect, useState } from "react";
import { api, type ItemDetail, type Snapshot } from "../../api";
import { Badge, DomainChip, timeAgo } from "../bits";
import { Article } from "../markdown";

export function ReaderPane({
  item,
  snap,
  active,
}: {
  item: ItemDetail["item"];
  snap: Snapshot | undefined;
  active: boolean;
}) {
  const [article, setArticle] = useState<string | null>(null);
  const [articleMissing, setArticleMissing] = useState(false);
  const revId = snap?.content_revision_id ?? null;

  useEffect(() => {
    if (!revId) return;
    api
      .article(revId)
      .then((a) => setArticle(a.article_md))
      .catch(() => setArticleMissing(true));
  }, [revId]);

  return (
    <div className="border-r border-line bg-surface px-8 py-10 lg:px-14">
      <div className="mb-3 font-mono text-[11px] font-semibold tracking-[0.12em] text-accent uppercase">
        <DomainChip url={item.original_url} />
      </div>
      <h1 className="mb-2 font-serif text-3xl leading-tight tracking-tight">
        {item.title ?? item.original_url}
      </h1>
      <div className="mb-8 flex flex-wrap items-center gap-2.5 text-[12.5px] text-ink-3">
        <span>captured {snap ? timeAgo(snap.captured_at) : "—"}</span>
        {snap && <Badge label={snap.content_quality} />}
        {snap?.capture_method && (
          <span className="font-mono text-[11px]">via {snap.capture_method}</span>
        )}
        <a
          href={item.original_url}
          target="_blank"
          rel="noreferrer"
          className="text-ink-3 underline decoration-line underline-offset-2 hover:text-ink"
        >
          original ↗
        </a>
      </div>

      {active && (
        <p className="mb-6 rounded-lg bg-live-soft px-4 py-3 text-[13px] text-live">
          {item.latest_job?.kind} {item.latest_job?.status} — the archived text appears here when
          the run commits a snapshot.
        </p>
      )}

      {article ? (
        <Article md={article} />
      ) : snap && !revId ? (
        <EmptyNote text="Snapshot committed without a content revision — nothing was archived as readable text." />
      ) : articleMissing ? (
        <EmptyNote text="Archived text unavailable for this revision." />
      ) : revId ? (
        <EmptyNote text="loading article…" />
      ) : (
        !active && <EmptyNote text="No snapshot yet — the archive has not been captured." />
      )}
    </div>
  );
}

function EmptyNote({ text }: { text: string }) {
  return <p className="rounded-lg bg-paper px-4 py-6 text-center text-sm text-ink-3">{text}</p>;
}

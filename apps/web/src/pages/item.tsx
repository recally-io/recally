import { Link, useParams } from "@tanstack/react-router";
import { useCallback, useEffect, useRef, useState } from "react";
import { api, type ItemDetail, type Note, type Snapshot, type SummaryOutput } from "../api";
import { Badge, DomainChip, timeAgo } from "../components/bits";
import { Article } from "../components/markdown";

const JOB_ACTIVE = new Set(["queued", "running", "pending"]);

export function ItemPage() {
  const { itemId } = useParams({ strict: false }) as { itemId: string };
  const [detail, setDetail] = useState<ItemDetail | null>(null);
  const [article, setArticle] = useState<string | null>(null);
  const [articleMissing, setArticleMissing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [toast, setToast] = useState<string | null>(null);
  const opened = useRef(false);

  const refresh = useCallback(
    () =>
      api
        .getItem(itemId)
        .then(setDetail)
        .catch((e) => setError(e.message)),
    [itemId],
  );

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const snap: Snapshot | undefined = detail?.snapshots.find(
    (s) => s.id === detail.item.current_snapshot_id,
  );

  const revId = snap?.content_revision_id ?? null;

  useEffect(() => {
    if (!revId) return;
    api
      .article(revId)
      .then((a) => setArticle(a.article_md))
      .catch(() => setArticleMissing(true));
  }, [revId]);

  // Record an "open" reading event once per view (§9.4 signals feed ranking later).
  useEffect(() => {
    if (detail && !opened.current) {
      opened.current = true;
      void api.readingEvent(itemId, "open").catch(() => {});
    }
  }, [detail, itemId]);

  // Poll while a job for this item is still in flight.
  const active = detail?.item.latest_job && JOB_ACTIVE.has(detail.item.latest_job.status);
  useEffect(() => {
    if (!active) return;
    const t = setInterval(() => void refresh(), 4000);

    return () => clearInterval(t);
  }, [active, refresh]);

  const flash = (m: string) => {
    setToast(m);
    setTimeout(() => setToast(null), 3000);
  };

  const summary = parseSummary(detail?.artifacts);

  if (error) return <p className="p-8 text-sm text-danger">{error}</p>;

  if (!detail) return <p className="p-8 text-sm text-ink-3">loading…</p>;
  const { item } = detail;

  return (
    <div className="grid lg:grid-cols-[1fr_300px]">
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
            {detail.item.latest_job!.kind} {detail.item.latest_job!.status} — the archived text
            appears here when the run commits a snapshot.
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

      <aside className="bg-paper px-6 py-7">
        {toast && (
          <p className="mb-4 rounded-lg bg-accent-soft px-3 py-2 text-xs text-accent">{toast}</p>
        )}
        {summary && (
          <>
            <H>Summary</H>
            <p className="text-[12.5px] leading-relaxed text-ink-2">{summary.short_summary}</p>
            {summary.key_points?.length ? (
              <>
                <H>Key points</H>
                <ul className="list-disc space-y-1.5 pl-4 text-[12.5px] leading-relaxed text-ink-2 marker:text-accent">
                  {summary.key_points.map((k) => (
                    <li key={k}>{k}</li>
                  ))}
                </ul>
              </>
            ) : null}
          </>
        )}
        <H>Provenance</H>
        <dl className="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1.5 text-xs">
          <dt className="text-ink-3">snapshot</dt>
          <dd className="font-mono text-ink-2">{snap?.id.slice(0, 8) ?? "—"}</dd>
          <dt className="text-ink-3">revision</dt>
          <dd className="font-mono text-ink-2">{revId?.slice(0, 8) ?? "—"}</dd>
          <dt className="text-ink-3">generation</dt>
          <dd className="text-ink-2">{item.capture_generation}</dd>
          <dt className="text-ink-3">final url</dt>
          <dd className="truncate text-ink-2">{snap?.final_url ?? item.original_url}</dd>
        </dl>

        <H>My note</H>
        <NoteEditor itemId={itemId} notes={detail.notes} onSaved={refresh} flash={flash} />

        <H>Actions</H>
        <div className="flex flex-wrap gap-2">
          <Action
            label={item.read_status === "read" ? "mark unread" : "mark read"}
            onClick={async () => {
              await api.patchItem(itemId, {
                read_status: item.read_status === "read" ? "unread" : "read",
              });
              await refresh();
            }}
          />
          <Action
            label="re-capture"
            onClick={async () => {
              await api.recapture(itemId);
              flash("capture queued");
              await refresh();
            }}
          />
          {snap && (
            <Action
              label="share link"
              onClick={async () => {
                const s = await api.createShare(itemId, snap.id, revId);
                await navigator.clipboard.writeText(`${location.origin}${s.url}`);
                flash("share link copied");
              }}
            />
          )}
        </div>
        <div className="mt-6">
          <Link to="/app" className="text-xs text-ink-3 hover:text-ink">
            ← back to library
          </Link>
        </div>
      </aside>
    </div>
  );
}

function H({ children }: { children: React.ReactNode }) {
  return (
    <h3 className="mt-6 mb-2 font-mono text-[11px] font-semibold tracking-[0.12em] text-ink-3 uppercase first:mt-0">
      {children}
    </h3>
  );
}

function EmptyNote({ text }: { text: string }) {
  return <p className="rounded-lg bg-paper px-4 py-6 text-center text-sm text-ink-3">{text}</p>;
}

function Action({ label, onClick }: { label: string; onClick: () => Promise<void> }) {
  const [busy, setBusy] = useState(false);

  return (
    <button
      type="button"
      disabled={busy}
      className="rounded-md border border-line bg-surface px-3 py-1.5 text-[12.5px] text-ink-2 hover:border-ink-3 disabled:opacity-50"
      onClick={() => {
        setBusy(true);
        onClick().finally(() => setBusy(false));
      }}
    >
      {label}
    </button>
  );
}

function NoteEditor({
  itemId,
  notes,
  onSaved,
  flash,
}: {
  itemId: string;
  notes: Note[];
  onSaved: () => Promise<void>;
  flash: (m: string) => void;
}) {
  const note = notes[0];
  const [body, setBody] = useState(note?.body ?? "");
  const [dirty, setDirty] = useState(false);
  useEffect(() => {
    if (!dirty) setBody(note?.body ?? "");
  }, [note?.body, dirty]);

  const save = async () => {
    if (note) await api.patchNote(note.id, body, note.version);
    else if (body.trim()) await api.addNote(itemId, body);
    else return;
    setDirty(false);
    flash("note saved");
    await onSaved();
  };

  return (
    <div>
      <textarea
        className="w-full rounded-lg border border-line bg-surface p-3 text-[12.5px] leading-relaxed text-ink-2 outline-none focus:border-accent"
        rows={3}
        placeholder="Why did you save this?"
        value={body}
        onChange={(e) => {
          setBody(e.target.value);
          setDirty(true);
        }}
      />
      {dirty && (
        <button
          type="button"
          className="mt-1.5 rounded-md bg-accent px-3 py-1 text-xs font-semibold text-white"
          onClick={() => void save()}
        >
          save note
        </button>
      )}
    </div>
  );
}

function parseSummary(artifacts: ItemDetail["artifacts"] | undefined): SummaryOutput | null {
  const a = artifacts?.find((x) => x.type === "summary" && x.status === "verified");

  if (!a) return null;

  try {
    return JSON.parse(a.output) as SummaryOutput;
  } catch {
    return null;
  }
}

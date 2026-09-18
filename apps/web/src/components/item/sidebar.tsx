import { Link } from "@tanstack/react-router";
import { useState } from "react";
import { api, type ItemDetail, type Snapshot } from "../../api";
import { NoteEditor } from "./note-editor";
import { parseSummary } from "./summary";

export function Sidebar({
  detail,
  snap,
  refresh,
}: {
  detail: ItemDetail;
  snap: Snapshot | undefined;
  refresh: () => Promise<void>;
}) {
  const [toast, setToast] = useState<string | null>(null);
  const { item } = detail;
  const itemId = item.id;
  const revId = snap?.content_revision_id ?? null;
  const summary = parseSummary(detail.artifacts);

  const flash = (message: string) => {
    setToast(message);
    setTimeout(() => setToast(null), 3000);
  };

  return (
    <aside className="bg-paper px-6 py-7">
      {toast && (
        <p className="mb-4 rounded-lg bg-accent-soft px-3 py-2 text-xs text-accent">{toast}</p>
      )}
      {summary && (
        <>
          <H>Summary</H>
          <p className="text-[12.5px] leading-relaxed text-ink-2">{summary.short_summary}</p>
          {summary.key_points.length ? (
            <>
              <H>Key points</H>
              <ul className="list-disc space-y-1.5 pl-4 text-[12.5px] leading-relaxed text-ink-2 marker:text-accent">
                {summary.key_points.map((k) => (
                  <li key={k.claim}>{k.claim}</li>
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
  );
}

function H({ children }: { children: React.ReactNode }) {
  return (
    <h3 className="mt-6 mb-2 font-mono text-[11px] font-semibold tracking-[0.12em] text-ink-3 uppercase first:mt-0">
      {children}
    </h3>
  );
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

export function timeAgo(iso: string): string {
  const s = Math.max(0, (Date.now() - new Date(iso).getTime()) / 1000);

  if (s < 60) return "just now";

  if (s < 3600) return `${Math.floor(s / 60)}m ago`;

  if (s < 86400) return `${Math.floor(s / 3600)}h ago`;

  if (s < 86400 * 30) return `${Math.floor(s / 86400)}d ago`;

  return new Date(iso).toLocaleDateString();
}

export function domainOf(url: string): string {
  try {
    return new URL(url).hostname.replace(/^www\./, "");
  } catch {
    return url;
  }
}

export function DomainChip({ url }: { url: string }) {
  const d = domainOf(url);

  return (
    <span className="inline-flex items-center gap-1">
      <span className="inline-flex h-3.5 w-3.5 items-center justify-center rounded-[3px] bg-accent-soft font-sans text-[8px] font-bold leading-none text-accent">
        {d[0]?.toUpperCase() ?? "?"}
      </span>
      {d}
    </span>
  );
}

const BADGE_STYLES: Record<string, string> = {
  complete: "bg-accent-soft text-accent",
  partial: "bg-warn-soft text-warn",
  unavailable: "bg-warn-soft text-warn",
  unknown: "bg-paper text-ink-3",
  running: "bg-live-soft text-live",
  queued: "bg-live-soft text-live",
  pending: "bg-paper text-ink-3",
  succeeded: "bg-accent-soft text-accent",
  failed: "bg-warn-soft text-warn",
  canceled: "bg-paper text-ink-3",
  cancel_requested: "bg-warn-soft text-warn",
  workflow_error: "bg-warn-soft text-warn",
};

export function Badge({ label }: { label: string }) {
  const cls = BADGE_STYLES[label] ?? "bg-paper text-ink-3";
  const live = label === "running" || label === "queued";

  return (
    <span
      className={`inline-flex items-center rounded px-1.5 py-0.5 font-sans text-[10px] font-semibold uppercase tracking-wide ${cls}`}
    >
      {live && <span className="mr-1 animate-pulse text-[8px]">●</span>}
      {label.replaceAll("_", " ")}
    </span>
  );
}

/** FTS snippets carry literal <b>/</b> marks from SQLite snippet(); split on
// them and render React nodes — everything else is plain text. */
export function Snippet({ text }: { text: string }) {
  const parts = text.split(/<\/?b>/g);

  return (
    <span className="text-ink-2">
      {parts.map((p, i) =>
        i % 2 === 1 ? (
          // biome-ignore lint/suspicious/noArrayIndexKey: split parts of a fixed string; order never changes
          <mark key={`m${i}`} className="rounded-sm bg-mark px-0.5 text-ink">
            {p}
          </mark>
        ) : (
          p
        ),
      )}
    </span>
  );
}

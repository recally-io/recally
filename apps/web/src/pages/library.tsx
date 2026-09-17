import { Link } from "@tanstack/react-router";
import { useCallback, useEffect, useMemo, useState } from "react";
import { api, type Item } from "../api";
import { Badge, DomainChip, domainOf, timeAgo } from "../components/bits";

const JOB_ACTIVE = new Set(["queued", "running", "pending", "cancel_requested"]);

type Filter = "all" | "unread" | "partial";

export function LibraryPage() {
  const [items, setItems] = useState<Item[]>([]);
  const [filter, setFilter] = useState<Filter>("all");
  const [error, setError] = useState<string | null>(null);
  const [loaded, setLoaded] = useState(false);

  const refresh = useCallback(
    () =>
      api
        .listItems()
        .then((r) => {
          setItems(r.items);
          setLoaded(true);
        })
        .catch((e) => setError(e.message)),
    [],
  );

  useEffect(() => {
    void refresh();
  }, [refresh]);

  // Poll while any item has an in-flight job so capture progress shows live.
  const hasActive = items.some((i) => i.latest_job && JOB_ACTIVE.has(i.latest_job.status));
  useEffect(() => {
    if (!hasActive) return;
    const t = setInterval(() => void refresh(), 4000);

    return () => clearInterval(t);
  }, [hasActive, refresh]);

  const filtered = useMemo(
    () =>
      items.filter((i) => {
        if (filter === "unread") return i.read_status !== "read";

        if (filter === "partial") return i.content_quality === "partial";

        return true;
      }),
    [items, filter],
  );

  const unread = items.filter((i) => i.read_status !== "read").length;

  const chip = (f: Filter, label: string) => (
    <button
      key={f}
      type="button"
      onClick={() => setFilter(f)}
      className={`rounded-full border px-3 py-1 text-xs ${
        filter === f
          ? "border-accent bg-accent-soft font-semibold text-accent"
          : "border-line bg-surface text-ink-2 hover:border-ink-3"
      }`}
    >
      {label}
    </button>
  );

  return (
    <div className="mx-auto max-w-4xl">
      <div className="flex items-baseline justify-between px-6 pt-6 pb-3">
        <h1 className="font-serif text-2xl">Library</h1>
        <span className="text-xs text-ink-3">
          {items.length} saved · {unread} unread
        </span>
      </div>
      <div className="flex gap-2 border-b border-line px-6 pb-3">
        {chip("all", "All")}
        {chip("unread", "Unread")}
        {chip("partial", "Partial")}
      </div>
      {error && <p className="px-6 py-3 text-sm text-danger">{error}</p>}
      <ul>
        {filtered.map((it) => {
          const capturing = it.latest_job && JOB_ACTIVE.has(it.latest_job.status);
          const failed = it.latest_job?.status === "failed" && !it.current_snapshot_id;

          return (
            <li key={it.id}>
              <Link
                to="/app/items/$itemId"
                params={{ itemId: it.id }}
                className="grid grid-cols-[4px_1fr_auto] gap-x-4 border-b border-line bg-surface px-6 py-3.5 hover:bg-paper/60"
              >
                <i
                  className={`w-1 self-stretch rounded-sm ${
                    capturing
                      ? "bg-live"
                      : failed
                        ? "bg-warn"
                        : it.read_status !== "read"
                          ? "bg-accent"
                          : ""
                  }`}
                />
                <div className="min-w-0">
                  <div
                    className={`truncate text-[14.5px] leading-snug ${
                      it.read_status === "read" ? "font-normal text-ink-2" : "font-semibold"
                    }`}
                  >
                    {it.title ?? domainOf(it.original_url)}
                  </div>
                  <div className="mt-1 flex items-center gap-2 text-xs text-ink-3">
                    <DomainChip url={it.original_url} />
                    <span>·</span>
                    <span>saved {timeAgo(it.saved_at)}</span>
                    {capturing && (
                      <span className="text-live">
                        · {it.latest_job!.kind} {it.latest_job!.status}
                      </span>
                    )}
                    {failed && <span className="text-warn">· capture failed</span>}
                  </div>
                </div>
                <div className="pt-0.5 text-right">
                  {capturing ? (
                    <Badge label="running" />
                  ) : it.content_quality ? (
                    <Badge label={it.content_quality} />
                  ) : null}
                </div>
              </Link>
            </li>
          );
        })}
      </ul>
      {loaded && filtered.length === 0 && (
        <p className="py-16 text-center text-sm text-ink-3">
          {items.length === 0
            ? "Nothing saved yet — paste a URL above to archive it."
            : "No items match this filter."}
        </p>
      )}
    </div>
  );
}

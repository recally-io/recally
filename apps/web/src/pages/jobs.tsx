import { Link } from "@tanstack/react-router";
import { useCallback, useEffect, useState } from "react";
import { api, errorMessage, type JobListEntry, type JobDetail } from "../api";
import { Badge, timeAgo } from "../components/bits";
import { isJobLive, isJobCancellable } from "../components/job-status";
import { useRefreshWhileActive } from "../components/use-refresh-while-active";

export function JobsPage() {
  const [jobs, setJobs] = useState<JobListEntry[]>([]);
  const [open, setOpen] = useState<JobDetail | null>(null);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(
    () =>
      api
        .listJobs()
        .then((r) => setJobs(r.jobs))
        .catch((e) => setError(errorMessage(e))),
    [],
  );

  useEffect(() => {
    void refresh();
  }, [refresh]);

  useRefreshWhileActive(
    jobs.some((j) => isJobLive(j.status)),
    refresh,
  );

  const inspect = (id: string) =>
    api
      .getJob(id)
      .then(setOpen)
      .catch((e) => setError(errorMessage(e)));

  return (
    <div className="mx-auto max-w-4xl px-6 py-6">
      <h1 className="mb-4 font-serif text-2xl">Jobs</h1>
      {error && <p className="mb-3 text-sm text-danger">{error}</p>}
      <ul className="divide-y divide-line border-y border-line">
        {jobs.map((j) => (
          <li key={j.id}>
            <button
              type="button"
              className="grid w-full grid-cols-[90px_80px_1fr_auto_auto] items-center gap-4 px-2 py-3 text-left hover:bg-surface"
              onClick={() => void inspect(j.id)}
            >
              <span className="font-mono text-[11px] text-ink-3">{j.id.slice(0, 8)}</span>
              <span className="text-xs font-medium text-ink-2">{j.kind}</span>
              <span className="truncate text-[12.5px] text-ink-2">
                {j.item_title ?? j.item_id?.slice(0, 8) ?? "—"}
                {j.error && <span className="text-warn"> — {j.error}</span>}
              </span>
              <Badge label={j.status} />
              <span className="text-[11px] text-ink-3">{timeAgo(j.created_at)}</span>
            </button>
          </li>
        ))}
        {jobs.length === 0 && (
          <li className="py-12 text-center text-sm text-ink-3">No jobs yet.</li>
        )}
      </ul>

      {open && (
        <div className="mt-6 rounded-xl border border-line bg-surface p-5">
          <div className="mb-3 flex items-center gap-3">
            <span className="font-mono text-xs text-ink-3">{open.id}</span>
            <Badge label={open.status} />
            <span className="text-xs text-ink-3">
              {open.kind} · attempts {open.attempt_count}
            </span>
            {isJobCancellable(open.status) && (
              <button
                type="button"
                className="ml-auto rounded-md border border-line px-2.5 py-1 text-xs text-ink-2 hover:border-warn hover:text-warn"
                onClick={() => void api.cancelJob(open.id).then(() => inspect(open.id))}
              >
                cancel
              </button>
            )}
            <button
              type="button"
              className="ml-auto text-xs text-ink-3 hover:text-ink"
              style={isJobCancellable(open.status) ? { marginLeft: 8 } : undefined}
              onClick={() => setOpen(null)}
            >
              close
            </button>
          </div>
          {open.item_id && (
            <Link
              to="/app/items/$itemId"
              params={{ itemId: open.item_id }}
              className="mb-3 inline-block text-xs text-accent underline decoration-line underline-offset-2"
            >
              view item →
            </Link>
          )}
          <ul className="max-h-80 space-y-1 overflow-y-auto font-mono text-[11px] leading-relaxed text-ink-2">
            {open.events.map((e) => (
              <li key={e.sequence} className="flex gap-3">
                <span className="w-6 shrink-0 text-right text-ink-3">{e.sequence}</span>
                <span className="w-28 shrink-0 text-ink-3">{e.kind}</span>
                <span className="min-w-0">{e.summary ?? e.reason_code ?? ""}</span>
              </li>
            ))}
            {open.events.length === 0 && <li className="text-ink-3">no events recorded</li>}
          </ul>
        </div>
      )}
    </div>
  );
}

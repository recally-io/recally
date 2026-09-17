import { useParams } from "@tanstack/react-router";
import { useCallback, useEffect, useRef, useState } from "react";
import { api, errorMessage, type ItemDetail, type Snapshot } from "../api";
import { ReaderPane } from "../components/item/reader-pane";
import { Sidebar } from "../components/item/sidebar";
import { isJobLive } from "../components/job-status";
import { useRefreshWhileActive } from "../components/use-refresh-while-active";

export function ItemPage() {
  const { itemId } = useParams({ from: "/app/items/$itemId" });
  const [detail, setDetail] = useState<ItemDetail | null>(null);
  const [error, setError] = useState<string | null>(null);
  const opened = useRef(false);

  const refresh = useCallback(
    () =>
      api
        .getItem(itemId)
        .then(setDetail)
        .catch((e) => setError(errorMessage(e))),
    [itemId],
  );

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const snap: Snapshot | undefined = detail?.snapshots.find(
    (s) => s.id === detail.item.current_snapshot_id,
  );

  // Record an "open" reading event once per view (§9.4 signals feed ranking later).
  useEffect(() => {
    if (detail && !opened.current) {
      opened.current = true;
      void api.readingEvent(itemId, "open").catch(() => {});
    }
  }, [detail, itemId]);

  const active = Boolean(detail?.item.latest_job && isJobLive(detail.item.latest_job.status));
  useRefreshWhileActive(active, refresh);

  if (error) return <p className="p-8 text-sm text-danger">{error}</p>;

  if (!detail) return <p className="p-8 text-sm text-ink-3">loading…</p>;

  return (
    <div className="grid lg:grid-cols-[1fr_300px]">
      <ReaderPane item={detail.item} snap={snap} active={active} />
      <Sidebar detail={detail} snap={snap} refresh={refresh} />
    </div>
  );
}

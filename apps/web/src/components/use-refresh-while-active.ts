import { useEffect } from "react";

export function useRefreshWhileActive(
  enabled: boolean,
  refresh: () => void | Promise<void>,
  intervalMs = 4000,
): void {
  useEffect(() => {
    if (!enabled) return;

    const timer = setInterval(() => void refresh(), intervalMs);

    return () => clearInterval(timer);
  }, [enabled, refresh, intervalMs]);
}

import {
  createRootRoute,
  createRoute,
  Link,
  Outlet,
  useNavigate,
  useRouterState,
} from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { api } from "./api";
import { ItemPage } from "./pages/item";
import { JobsPage } from "./pages/jobs";
import { LandingPage } from "./pages/landing";
import { LibraryPage } from "./pages/library";
import { SearchPage } from "./pages/search";
import { SettingsPage } from "./pages/settings";

function SaveBar() {
  const navigate = useNavigate();
  const [value, setValue] = useState("");
  const [busy, setBusy] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        inputRef.current?.focus();
      }
    };

    window.addEventListener("keydown", onKey);

    return () => window.removeEventListener("keydown", onKey);
  }, []);

  const save = async () => {
    const url = value.trim();

    if (!url || busy) return;
    setBusy(true);

    try {
      const r = await api.saveUrl(url);
      setValue("");
      void navigate({ to: "/app/items/$itemId", params: { itemId: r.item_id } });
    } catch {
      setBusy(false);
    }
  };

  return (
    <form
      className="flex flex-1 items-center gap-2 rounded-lg border border-line bg-paper px-3 py-1.5"
      onSubmit={(e) => {
        e.preventDefault();
        void save();
      }}
    >
      <input
        ref={inputRef}
        className="w-full bg-transparent text-[13px] text-ink outline-none placeholder:text-ink-3"
        placeholder="Paste a URL to archive it…"
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      {busy ? (
        <span className="font-mono text-[11px] text-live">saving…</span>
      ) : (
        <kbd className="rounded border border-line border-b-2 bg-surface px-1.5 py-0.5 font-mono text-[11px] text-ink-3">
          ⌘K
        </kbd>
      )}
    </form>
  );
}

function AppShell() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });

  const link = (to: string, label: string, exact = false) => {
    const on = exact ? pathname === to : pathname.startsWith(to);

    return (
      <Link
        to={to}
        className={`rounded-md px-2.5 py-1 text-[12.5px] ${
          on ? "bg-accent-soft font-semibold text-accent" : "text-ink-2 hover:text-ink"
        }`}
      >
        {label}
      </Link>
    );
  };

  return (
    <div className="min-h-screen">
      <header className="sticky top-0 z-10 flex items-center gap-4 border-b border-line bg-surface px-5 py-2.5">
        <Link to="/app" className="font-serif text-[15px] font-bold">
          Re<span className="text-accent">call</span>ly
        </Link>
        <SaveBar />
        <nav className="flex gap-0.5">
          {link("/app", "Library", true)}
          {link("/app/search", "Search")}
          {link("/app/jobs", "Jobs")}
          {link("/app/settings", "Settings")}
        </nav>
      </header>
      <Outlet />
    </div>
  );
}

const rootRoute = createRootRoute({ component: Outlet });

const landingRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: LandingPage,
});

const appRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/app",
  component: AppShell,
});

const libraryRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/",
  component: LibraryPage,
});

const itemRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/items/$itemId",
  component: ItemPage,
});

const searchRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/search",
  component: SearchPage,
});

const jobsRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/jobs",
  component: JobsPage,
});

const settingsRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/settings",
  component: SettingsPage,
});

export const routeTree = rootRoute.addChildren([
  landingRoute,
  appRoute.addChildren([libraryRoute, itemRoute, searchRoute, jobsRoute, settingsRoute]),
]);

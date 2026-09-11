import { createRootRoute, createRoute, Link, Outlet } from "@tanstack/react-router";
import { ItemPage } from "./pages/item";
import { JobsPage } from "./pages/jobs";
import { LibraryPage } from "./pages/library";
import { SearchPage } from "./pages/search";

const rootRoute = createRootRoute({
  component: () => (
    <div className="mx-auto max-w-3xl px-4 py-6 font-sans">
      <header className="mb-6 flex items-center justify-between">
        <Link to="/" className="text-xl font-semibold">
          recally
        </Link>
        <nav className="flex gap-4 text-sm text-neutral-500">
          <Link to="/search">search</Link>
          <Link to="/jobs">jobs</Link>
        </nav>
      </header>
      <Outlet />
    </div>
  ),
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: LibraryPage,
});

const itemRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/items/$itemId",
  component: ItemPage,
});

const searchRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/search",
  component: SearchPage,
});

const jobsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/jobs",
  component: JobsPage,
});

export const routeTree = rootRoute.addChildren([indexRoute, itemRoute, searchRoute, jobsRoute]);

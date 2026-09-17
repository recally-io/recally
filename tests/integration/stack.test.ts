// Opt-in live-stack integration test.
//
//   ALCHEMY_INTEG=1 pnpm test tests/integration/stack.test.ts
//
// Spawns `alchemy deploy` for an isolated stage, drives the deployed API
// worker over HTTP, then destroys the stage. Requires Cloudflare
// credentials (`alchemy profile edit --add Cloudflare`), built web assets
// (pnpm build), and DEV_LIBRARY_ID in .env for the anonymous /health
// assertion.

import { execFileSync } from "node:child_process";
import { afterAll, beforeAll, expect, it } from "vitest";

const enabled = !!process.env.ALCHEMY_INTEG;

const stage = "alchemy-integ";

const alchemyBin = "node_modules/.bin/alchemy";

let apiUrl: string | undefined;

const alchemy = (args: string[]) =>
  execFileSync(alchemyBin, args, {
    stdio: ["ignore", "pipe", "inherit"],
    encoding: "utf8",
  });

beforeAll(() => {
  if (!enabled) return;
  const out = alchemy(["deploy", "--adopt", "--stage", stage, "--yes"]);
  // Stack outputs are printed at the end of the deploy; pick the api URL.
  const match = out.match(/https:\/\/recally-api[-a-z0-9.]*\.workers\.dev/);

  if (!match) throw new Error(`could not find api url in deploy output:\n${out}`);
  apiUrl = match[0];
}, 900_000);

afterAll(() => {
  if (!enabled || process.env.NO_DESTROY) return;
  alchemy(["destroy", "--stage", stage, "--yes"]);
}, 600_000);

it.runIf(enabled)(
  "deployed stack serves health, shares, and SPA from one origin",
  async () => {
    expect(apiUrl).toBeTruthy();

    // Public health endpoint inside the auth'd /api/v1 mount; with
    // DEV_LIBRARY_ID set (single-user dev default) it resolves anonymously.
    const health = await fetch(`${apiUrl}/api/v1/health`);
    expect(health.status).toBe(200);

    // Public share surface: unknown token → the domain 404, not the SPA shell.
    const share = await fetch(`${apiUrl}/s/does-not-exist`);
    expect(share.status).toBe(404);

    // SPA fallback: a deep link must serve the app shell.
    const deepLink = await fetch(`${apiUrl}/items/some-id`);
    expect(deepLink.status).toBe(200);
    const body = await deepLink.text();
    expect(body).toContain("<!doctype html>");
  },
  120_000,
);

// Opt-in live-stack integration test.
//
//   ALCHEMY_INTEG=1 pnpm test tests/integration/stack.test.ts
//
// Spawns `alchemy deploy` for an isolated stage, drives the deployed API
// worker over HTTP, then destroys the stage. Requires Cloudflare
// credentials (`alchemy profile edit --add Cloudflare`) and built web assets
// (`pnpm build`).

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
  const out = alchemy(["deploy", "--stage", stage, "--yes"]);
  // Stack outputs are printed as an object literal at the end of the deploy.
  const match = out.match(/apiUrl:\s*'(https:\/\/[^']+)'/);

  if (!match) throw new Error(`could not find api url in deploy output:\n${out}`);
  apiUrl = match[1];
}, 900_000);

afterAll(() => {
  if (!enabled || process.env.NO_DESTROY) return;
  alchemy(["destroy", "--stage", stage, "--yes"]);
}, 600_000);

it.runIf(enabled)(
  "deployed stack serves the SPA, the public share surface, and the auth gate",
  async () => {
    expect(apiUrl).toBeTruthy();

    // SPA fallback: a deep link must serve the app shell (assets binding).
    const deepLink = await fetch(`${apiUrl}/items/some-id`);
    expect(deepLink.status).toBe(200);
    expect(await deepLink.text()).toContain("<!doctype html>");

    // Public share surface: an unknown token → the domain 404, not the SPA
    // shell. Answering at all proves the D1 binding works and migrations ran —
    // the route joins shares → items.
    const share = await fetch(`${apiUrl}/s/does-not-exist`);
    expect(share.status).toBe(404);
    expect(await share.json()).toMatchObject({ code: "not_found" });

    // Auth gate: a stage starts with no library provisioned, so the API mount
    // answers needs_login. With DEV_LIBRARY_ID bound and its library row
    // present (the single-user setup) the same request is 200.
    const health = await fetch(`${apiUrl}/api/v1/health`);
    expect(health.status).toBe(401);
    expect(await health.json()).toMatchObject({ code: "needs_login" });
  },
  120_000,
);

// Browser Run adapter (§5.8, ADR-04). Owns session lifecycle: launch, bounded
// actions, always close in finally — disconnect is not teardown. Pi decides
// each action; this class only executes and reports observations.
//
// M0-T06 verifies on real cloud: launch/close semantics, idle expiry, billing
// during model waits.

interface BrowserSessionEnv {
  BROWSER: Fetcher;
}

interface PageObservation {
  url: string;
  title: string;
  // Text-tree summary for the agent — real a11y tree shape is M0 work.
  text: string;
  nodeRefs: Array<{ id: string; kind: string; label: string }>;
}

export class BrowserSession {
  private constructor(
    private readonly browser: Awaited<ReturnType<typeof launch>>,
    private readonly page: any,
  ) {}

  static async open(env: BrowserSessionEnv, url: string): Promise<BrowserSession> {
    const { launch } = await import("@cloudflare/playwright");
    const browser = await launch(env.BROWSER);
    const page = await browser.newPage();
    await page.goto(url, { waitUntil: "domcontentloaded", timeout: 30_000 });

    return new BrowserSession(browser, page);
  }

  async observe(): Promise<PageObservation> {
    const title: string = await this.page.title();
    const url: string = this.page.url();

    const text: string = await this.page.evaluate(
      "document.body ? document.body.innerText.slice(0, 20000) : ''",
    );

    // Minimal actionable targets: links and buttons with labels. Real node-ref
    // discipline (refs die on DOM change) is enforced by the tool layer.
    const nodeRefs: PageObservation["nodeRefs"] = await this.page.evaluate(`(() => {
      const els = document.querySelectorAll('a[href], button, [role="button"]');
      return [...els].slice(0, 100).map((el, i) => ({
        id: 'n' + i,
        kind: el.tagName.toLowerCase(),
        label: (el.innerText || el.getAttribute('aria-label') || '').slice(0, 120),
      }));
    })()`);

    return { url, title, text, nodeRefs };
  }

  async scroll(amount = 0): Promise<void> {
    await this.page.evaluate("window.scrollBy(0, arguments[0] || window.innerHeight)", amount);
  }

  async click(selector: string): Promise<void> {
    await this.page.click(selector, { timeout: 10_000 });
  }

  async waitFor(selector: string): Promise<void> {
    await this.page.waitForSelector(selector, { timeout: 10_000 });
  }

  async navigate(url: string): Promise<void> {
    await this.page.goto(url, {
      waitUntil: "domcontentloaded",
      timeout: 30_000,
    });
  }

  async renderedDom(): Promise<string> {
    return this.page.content();
  }

  async screenshot(): Promise<Uint8Array> {
    return new Uint8Array(await this.page.screenshot({ type: "webp" }));
  }

  async close(): Promise<void> {
    await this.browser.close();
  }
}

import type { launch } from "@cloudflare/playwright";

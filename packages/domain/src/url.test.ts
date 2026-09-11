import { describe, expect, it } from "vitest";
import { normalizeUrl } from "./url";

describe("normalizeUrl", () => {
  it("lowercases host and strips default port", () => {
    expect(normalizeUrl("HTTPS://EXAMPLE.com:443/a")).toBe("https://example.com/a");
  });

  it("strips utm_* and tracking params, sorts the rest", () => {
    expect(normalizeUrl("https://x.com/p?b=2&utm_source=tw&a=1&fbclid=z")).toBe(
      "https://x.com/p?a=1&b=2",
    );
  });

  it("rejects non-http(s)", () => {
    expect(() => normalizeUrl("ftp://x.com/a")).toThrow();
  });
});

import { describe, expect, it } from "vitest";
import { parseSummary } from "./summary";

describe("archived summary boundary", () => {
  it("accepts production key points with evidence references", () => {
    const summary = parseSummary([
      {
        type: "summary",
        status: "verified",
        output: JSON.stringify({
          short_summary: "Summary",
          key_points: [{ claim: "Claim", evidence_block_ids: ["b1"] }],
        }),
      },
    ]);

    expect(summary?.key_points).toEqual([{ claim: "Claim", evidence_block_ids: ["b1"] }]);
  });

  it.each(["{", JSON.stringify({ short_summary: "Summary", key_points: ["Old string format"] })])(
    "rejects malformed or drifted payload: %s",
    (output) => {
      expect(parseSummary([{ type: "summary", status: "verified", output }])).toBeNull();
    },
  );

  it("does not render unverified summaries", () => {
    expect(parseSummary([{ type: "summary", status: "failed", output: "{}" }])).toBeNull();
  });
});

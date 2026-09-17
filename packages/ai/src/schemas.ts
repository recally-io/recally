import { z } from "zod";

// Independent verifier output (plan §7.3): judged from content evidence, not
// the capture agent's self-report.
export const verificationSchema = z.object({
  verdict: z.enum(["complete", "partial", "wrong"]),
  confidence: z.enum(["high", "medium", "low"]),
  problems: z.array(z.string()).default([]),
  missingParts: z.array(z.string()).default([]),
});

export type Verification = z.infer<typeof verificationSchema>;

// Summary artifact payload (plan §10.2).
export const summaryOutputSchema = z.object({
  short_summary: z.string(),
  key_points: z.array(
    z.object({
      claim: z.string(),
      evidence_block_ids: z.array(z.string()).default([]),
    }),
  ),
  topics: z.array(z.string()).default([]),
  entities: z.array(z.string()).default([]),
  caveats: z.array(z.string()).default([]),
});

export type SummaryOutput = z.infer<typeof summaryOutputSchema>;

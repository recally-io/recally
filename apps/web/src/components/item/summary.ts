import { summaryOutputSchema, type SummaryOutput } from "@recally/ai/schemas";
import type { ItemDetail } from "../../api";

export function parseSummary(artifacts: ItemDetail["artifacts"]): SummaryOutput | null {
  const artifact = artifacts.find(
    (entry) => entry.type === "summary" && entry.status === "verified",
  );

  if (!artifact) return null;

  try {
    const result = summaryOutputSchema.safeParse(JSON.parse(artifact.output));

    return result.success ? result.data : null;
  } catch {
    return null;
  }
}

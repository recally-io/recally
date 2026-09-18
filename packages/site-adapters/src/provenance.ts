import { nowIso } from "@recally/domain";
import type { AdapterResult } from "./types";

export function adapterProvenance(endpointUrls: string[]): AdapterResult["provenance"] {
  return { endpointUrls, fetchedAt: nowIso() };
}

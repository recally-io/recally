// Initial application limits (plan §16.1). These are project defaults, not platform limits.
export const DEFAULT_LIMITS = {
  fetch: {
    maxRedirects: 5,
    timeoutMs: 15_000,
    maxHtmlBytes: 5 * 1024 * 1024,
    maxBodyUrlsPerItem: 3,
  },
  assets: {
    maxAssetBytes: 5 * 1024 * 1024,
    maxAssetsPerItem: 30,
    maxTotalAssetBytes: 50 * 1024 * 1024,
    concurrency: 3,
  },
  agent: {
    // Live-tuned on dev: a real capture takes ~15-20 model calls (site
    // probing, retries after tool errors, block-range exploration) — 12
    // turned every retryable failure into budget_exceeded.
    maxModelCalls: 24,
    maxInputTokens: 64_000,
    maxOutputTokens: 8_000,
    maxFormatRepairs: 1,
    maxNoProgressTurns: 2,
  },
  browser: {
    maxEpisodesPerRun: 2,
    maxSessionMs: 120_000,
    maxActions: 12,
  },
  run: {
    maxWallClockMs: 10 * 60 * 1000,
    maxNetworkRetries: 2,
  },
  analysis: {
    maxAnalyzeTokens: 60_000,
    askMaxEvidenceBlocks: 12,
    maxStepResultBytes: 64 * 1024,
  },
  search: {
    ftsCandidates: 40,
    vectorCandidates: 40,
    vectorTopK: 40,
    rrfK: 60,
    chunkTokens: 600,
    chunkOverlap: 100,
  },
  concurrency: {
    captureGlobal: 4,
    browserGlobal: 2,
    perOrigin: 1,
  },
  retention: {
    failedAttemptDays: 7,
    shareDefaultDays: 7,
    backupKeepDays: 30,
  },
  outbox: {
    batchSize: 50,
  },
} as const;

export type Limits = typeof DEFAULT_LIMITS;

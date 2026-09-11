-- Recally core schema (D1 / SQLite).
-- IDs are TEXT (uuid). Timestamps are ISO-8601 UTC strings.
-- Blueprint per docs/plan.md §4.2: ownership, dedupe, immutability, job, and search constraints.

CREATE TABLE libraries (
  id            TEXT PRIMARY KEY,
  name          TEXT NOT NULL,
  timezone      TEXT NOT NULL DEFAULT 'UTC',
  access_sub    TEXT,
  settings      TEXT NOT NULL DEFAULT '{}',
  created_at    TEXT NOT NULL,
  updated_at    TEXT NOT NULL
);

CREATE TABLE items (
  id                  TEXT PRIMARY KEY,
  library_id          TEXT NOT NULL REFERENCES libraries(id),
  original_url        TEXT NOT NULL,
  normalized_url      TEXT NOT NULL,
  title               TEXT,
  saved_at            TEXT NOT NULL,
  first_read_at       TEXT,
  last_read_at        TEXT,
  read_status         TEXT NOT NULL DEFAULT 'unread'
                      CHECK (read_status IN ('unread', 'reading', 'read', 'archived')),
  current_snapshot_id TEXT,
  capture_generation  INTEGER NOT NULL DEFAULT 1,
  deleted_at          TEXT,
  created_at          TEXT NOT NULL,
  updated_at          TEXT NOT NULL
);
CREATE INDEX idx_items_library_saved ON items(library_id, saved_at, id);

-- Effective URL mappings; a fresh save after delete gets a new item id.
CREATE TABLE item_urls (
  id                  TEXT PRIMARY KEY,
  library_id          TEXT NOT NULL REFERENCES libraries(id),
  item_id             TEXT NOT NULL REFERENCES items(id),
  normalized_url_hash TEXT NOT NULL,
  normalized_url      TEXT NOT NULL,
  is_active           INTEGER NOT NULL DEFAULT 1,
  created_at          TEXT NOT NULL
);
CREATE UNIQUE INDEX uq_item_urls_active
  ON item_urls(library_id, normalized_url_hash) WHERE is_active = 1;

CREATE TABLE jobs (
  id                   TEXT PRIMARY KEY,
  library_id           TEXT NOT NULL REFERENCES libraries(id),
  kind                 TEXT NOT NULL
                       CHECK (kind IN ('ingest', 'enrich', 'index', 'digest', 'export', 'purge')),
  item_id              TEXT REFERENCES items(id),
  status               TEXT NOT NULL DEFAULT 'queued'
                       CHECK (status IN ('queued', 'running', 'retry_wait', 'succeeded',
                                         'needs_input', 'failed', 'cancelled')),
  operation_key        TEXT,
  workflow_instance_id TEXT,
  generation           INTEGER NOT NULL DEFAULT 1,
  parent_job_id        TEXT REFERENCES jobs(id),
  attempt_count        INTEGER NOT NULL DEFAULT 0,
  next_attempt_at      TEXT,
  payload              TEXT NOT NULL DEFAULT '{}',
  result               TEXT,
  error                TEXT,
  cancel_requested     INTEGER NOT NULL DEFAULT 0,
  created_at           TEXT NOT NULL,
  updated_at           TEXT NOT NULL
);
CREATE UNIQUE INDEX uq_jobs_operation_key ON jobs(operation_key) WHERE operation_key IS NOT NULL;
CREATE UNIQUE INDEX uq_jobs_workflow_instance
  ON jobs(workflow_instance_id) WHERE workflow_instance_id IS NOT NULL;
CREATE INDEX idx_jobs_dispatch ON jobs(status, next_attempt_at, id);

CREATE TABLE outbox (
  id            TEXT PRIMARY KEY,
  library_id    TEXT NOT NULL REFERENCES libraries(id),
  job_id        TEXT NOT NULL REFERENCES jobs(id),
  kind          TEXT NOT NULL,
  payload       TEXT NOT NULL DEFAULT '{}',
  available_at  TEXT NOT NULL,
  dispatched_at TEXT,
  created_at    TEXT NOT NULL
);
CREATE INDEX idx_outbox_due ON outbox(available_at, id) WHERE dispatched_at IS NULL;

CREATE TABLE idempotency_records (
  id           TEXT PRIMARY KEY,
  actor        TEXT NOT NULL,
  library_id   TEXT NOT NULL,
  route        TEXT NOT NULL,
  idem_key     TEXT NOT NULL,
  payload_hash TEXT NOT NULL,
  status_code  INTEGER,
  response     TEXT,
  created_at   TEXT NOT NULL,
  UNIQUE (actor, route, library_id, idem_key)
);

CREATE TABLE capture_runs (
  id                    TEXT PRIMARY KEY,
  library_id            TEXT NOT NULL REFERENCES libraries(id),
  item_id               TEXT NOT NULL REFERENCES items(id),
  job_id                TEXT NOT NULL REFERENCES jobs(id),
  generation            INTEGER NOT NULL,
  -- Revision fencing (§9.6): a run replays under the exact revisions it was
  -- created with; upgrades only affect new runs.
  agent_runtime_version TEXT NOT NULL,
  skill_revision        TEXT NOT NULL,
  toolset_version       TEXT NOT NULL,
  policy_version        TEXT NOT NULL,
  pipeline_version      TEXT NOT NULL,
  model_config_version  TEXT NOT NULL,
  pi_state_key          TEXT, -- latest serialized agent state object in R2
  requested_at         TEXT NOT NULL,
  target_url           TEXT NOT NULL,
  status               TEXT NOT NULL DEFAULT 'queued'
                       CHECK (status IN ('queued', 'running', 'succeeded', 'needs_input',
                                         'failed', 'cancelled')),
  outcome_code         TEXT,
  budget               TEXT NOT NULL DEFAULT '{}',
  created_at           TEXT NOT NULL,
  updated_at           TEXT NOT NULL
);
CREATE INDEX idx_capture_runs_item ON capture_runs(library_id, item_id, generation);

CREATE TABLE capture_attempts (
  id         TEXT PRIMARY KEY,
  library_id TEXT NOT NULL REFERENCES libraries(id),
  run_id     TEXT NOT NULL REFERENCES capture_runs(id),
  attempt_no INTEGER NOT NULL,
  kind       TEXT NOT NULL CHECK (kind IN ('fetch', 'browser', 'manual')),
  status     TEXT NOT NULL DEFAULT 'running'
             CHECK (status IN ('running', 'succeeded', 'failed')),
  error_code TEXT,
  started_at TEXT NOT NULL,
  ended_at   TEXT,
  UNIQUE (run_id, attempt_no)
);

CREATE TABLE agent_events (
  id           TEXT PRIMARY KEY,
  library_id   TEXT NOT NULL REFERENCES libraries(id),
  run_id       TEXT NOT NULL REFERENCES capture_runs(id),
  attempt_id   TEXT REFERENCES capture_attempts(id),
  sequence     INTEGER NOT NULL,
  kind         TEXT NOT NULL
               CHECK (kind IN ('observation', 'decision', 'tool_call', 'tool_result', 'validation')),
  reason_code  TEXT,
  summary      TEXT,
  evidence_ref TEXT,
  usage        TEXT,
  created_at   TEXT NOT NULL,
  UNIQUE (run_id, attempt_id, sequence)
);

-- Agent state checkpoints (§6.6): small index in D1, serialized Pi state
-- itself lives in R2 at state_key. Revision columns pin what the state means.
CREATE TABLE agent_states (
  id          TEXT PRIMARY KEY,
  library_id  TEXT NOT NULL REFERENCES libraries(id),
  run_id      TEXT NOT NULL REFERENCES capture_runs(id),
  sequence    INTEGER NOT NULL,
  state_key   TEXT NOT NULL,
  created_at  TEXT NOT NULL,
  UNIQUE (run_id, sequence)
);

CREATE TABLE snapshots (
  id                  TEXT PRIMARY KEY,
  library_id          TEXT NOT NULL REFERENCES libraries(id),
  item_id             TEXT NOT NULL REFERENCES items(id),
  capture_run_id      TEXT NOT NULL REFERENCES capture_runs(id),
  captured_at         TEXT NOT NULL,
  final_url           TEXT NOT NULL,
  capture_method      TEXT NOT NULL CHECK (capture_method IN ('fetch', 'browser', 'manual')),
  source_manifest_key TEXT NOT NULL,
  manifest_hash       TEXT NOT NULL,
  content_quality     TEXT NOT NULL DEFAULT 'unknown'
                      CHECK (content_quality IN ('complete', 'partial', 'unavailable', 'unknown')),
  resource_quality    TEXT NOT NULL DEFAULT 'unknown'
                      CHECK (resource_quality IN ('complete', 'partial', 'not_applicable', 'unknown')),
  committed_at        TEXT NOT NULL,
  created_at          TEXT NOT NULL
);
CREATE INDEX idx_snapshots_item ON snapshots(library_id, item_id, captured_at);

CREATE TABLE source_documents (
  id           TEXT PRIMARY KEY,
  library_id   TEXT NOT NULL REFERENCES libraries(id),
  run_id       TEXT NOT NULL REFERENCES capture_runs(id),
  attempt_id   TEXT REFERENCES capture_attempts(id),
  snapshot_id  TEXT REFERENCES snapshots(id),
  url          TEXT NOT NULL,
  kind         TEXT NOT NULL
               CHECK (kind IN ('response_body', 'rendered_dom', 'screenshot', 'observation',
                               'response_metadata', 'decision', 'adapter_record', 'manual')),
  content_type TEXT,
  r2_key       TEXT NOT NULL UNIQUE,
  sha256       TEXT NOT NULL,
  byte_size    INTEGER NOT NULL,
  created_at   TEXT NOT NULL
);

CREATE TABLE content_revisions (
  id                     TEXT PRIMARY KEY,
  library_id             TEXT NOT NULL REFERENCES libraries(id),
  snapshot_id            TEXT NOT NULL REFERENCES snapshots(id),
  item_id                TEXT NOT NULL REFERENCES items(id),
  extractor_version      TEXT NOT NULL,
  selected_source_blocks TEXT NOT NULL DEFAULT '[]',
  body_hash              TEXT NOT NULL,
  article_key            TEXT NOT NULL,
  blocks_key             TEXT NOT NULL,
  language               TEXT,
  token_count            INTEGER,
  created_at             TEXT NOT NULL
);
CREATE INDEX idx_content_revisions_snapshot ON content_revisions(library_id, snapshot_id);

CREATE TABLE assets (
  id           TEXT PRIMARY KEY,
  library_id   TEXT NOT NULL REFERENCES libraries(id),
  snapshot_id  TEXT REFERENCES snapshots(id),
  item_id      TEXT REFERENCES items(id),
  role         TEXT NOT NULL DEFAULT 'optional' CHECK (role IN ('required', 'optional')),
  original_url TEXT,
  content_type TEXT,
  r2_key       TEXT,
  sha256       TEXT,
  byte_size    INTEGER,
  status       TEXT NOT NULL DEFAULT 'pending'
               CHECK (status IN ('pending', 'stored', 'missing', 'failed')),
  created_at   TEXT NOT NULL
);

CREATE TABLE chunks (
  id                  TEXT PRIMARY KEY,
  library_id          TEXT NOT NULL REFERENCES libraries(id),
  content_revision_id TEXT NOT NULL REFERENCES content_revisions(id),
  item_id             TEXT NOT NULL REFERENCES items(id),
  ordinal             INTEGER NOT NULL,
  block_range         TEXT NOT NULL,
  text                TEXT NOT NULL,
  text_hash           TEXT NOT NULL,
  tokenizer_version   TEXT NOT NULL,
  embedding_version   TEXT,
  vector_id           TEXT,
  index_state         TEXT NOT NULL DEFAULT 'pending'
                      CHECK (index_state IN ('pending', 'submitted', 'queryable', 'stale', 'removed')),
  created_at          TEXT NOT NULL,
  UNIQUE (content_revision_id, ordinal)
);
CREATE INDEX idx_chunks_revision ON chunks(library_id, content_revision_id, ordinal);

CREATE TABLE ai_artifacts (
  id                        TEXT PRIMARY KEY,
  library_id                TEXT NOT NULL REFERENCES libraries(id),
  item_id                   TEXT NOT NULL REFERENCES items(id),
  type                      TEXT NOT NULL
                            CHECK (type IN ('summary', 'digest', 'imported_answer', 'qa')),
  status                    TEXT NOT NULL DEFAULT 'pending'
                            CHECK (status IN ('pending', 'verified', 'unverified_import', 'failed')),
  input_snapshot_id         TEXT REFERENCES snapshots(id),
  input_content_revision_id TEXT REFERENCES content_revisions(id),
  model_id                  TEXT,
  prompt_version            TEXT,
  pipeline_version          TEXT,
  operation_key             TEXT,
  output                    TEXT NOT NULL DEFAULT '{}',
  citations                 TEXT NOT NULL DEFAULT '[]',
  coverage                  TEXT,
  usage                     TEXT,
  created_at                TEXT NOT NULL
);
CREATE UNIQUE INDEX uq_ai_artifacts_op ON ai_artifacts(operation_key) WHERE operation_key IS NOT NULL;
CREATE INDEX idx_ai_artifacts_item ON ai_artifacts(library_id, item_id, type);

CREATE TABLE notes (
  id         TEXT PRIMARY KEY,
  library_id TEXT NOT NULL REFERENCES libraries(id),
  item_id    TEXT NOT NULL REFERENCES items(id),
  body       TEXT NOT NULL,
  version    INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE note_revisions (
  id         TEXT PRIMARY KEY,
  note_id    TEXT NOT NULL REFERENCES notes(id),
  version    INTEGER NOT NULL,
  body       TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE (note_id, version)
);

CREATE TABLE collections (
  id         TEXT PRIMARY KEY,
  library_id TEXT NOT NULL REFERENCES libraries(id),
  name       TEXT NOT NULL,
  created_at TEXT NOT NULL,
  UNIQUE (library_id, name)
);

CREATE TABLE item_collections (
  library_id    TEXT NOT NULL,
  item_id       TEXT NOT NULL REFERENCES items(id),
  collection_id TEXT NOT NULL REFERENCES collections(id),
  PRIMARY KEY (collection_id, item_id)
);

CREATE TABLE tags (
  id         TEXT PRIMARY KEY,
  library_id TEXT NOT NULL REFERENCES libraries(id),
  name       TEXT NOT NULL,
  source     TEXT NOT NULL DEFAULT 'user' CHECK (source IN ('user', 'ai')),
  created_at TEXT NOT NULL,
  UNIQUE (library_id, name, source)
);

CREATE TABLE item_tags (
  library_id TEXT NOT NULL,
  item_id    TEXT NOT NULL REFERENCES items(id),
  tag_id     TEXT NOT NULL REFERENCES tags(id),
  PRIMARY KEY (tag_id, item_id)
);

CREATE TABLE reading_events (
  id          TEXT PRIMARY KEY,
  library_id  TEXT NOT NULL REFERENCES libraries(id),
  item_id     TEXT NOT NULL REFERENCES items(id),
  kind        TEXT NOT NULL CHECK (kind IN ('open', 'finish', 'progress')),
  occurred_at TEXT NOT NULL,
  meta        TEXT NOT NULL DEFAULT '{}'
);
CREATE INDEX idx_reading_events_time ON reading_events(library_id, occurred_at);

CREATE TABLE shares (
  id                  TEXT PRIMARY KEY,
  library_id          TEXT NOT NULL REFERENCES libraries(id),
  item_id             TEXT NOT NULL REFERENCES items(id),
  snapshot_id         TEXT NOT NULL REFERENCES snapshots(id),
  content_revision_id TEXT REFERENCES content_revisions(id),
  token_hash          TEXT NOT NULL UNIQUE,
  include_full_text   INTEGER NOT NULL DEFAULT 0,
  allowed_artifacts   TEXT NOT NULL DEFAULT '[]',
  expires_at          TEXT,
  revoked_at          TEXT,
  created_at          TEXT NOT NULL
);

CREATE TABLE api_tokens (
  id           TEXT PRIMARY KEY,
  library_id   TEXT NOT NULL REFERENCES libraries(id),
  name         TEXT NOT NULL,
  token_hash   TEXT NOT NULL UNIQUE,
  scopes       TEXT NOT NULL,
  last_used_at TEXT,
  expires_at   TEXT,
  revoked_at   TEXT,
  created_at   TEXT NOT NULL
);

CREATE TABLE usage_reservations (
  id            TEXT PRIMARY KEY,
  library_id    TEXT NOT NULL REFERENCES libraries(id),
  operation_key TEXT NOT NULL,
  invocation_id TEXT NOT NULL UNIQUE,
  kind          TEXT NOT NULL,
  amount        INTEGER NOT NULL,
  state         TEXT NOT NULL DEFAULT 'reserved'
                CHECK (state IN ('reserved', 'settled', 'released', 'unknown')),
  created_at    TEXT NOT NULL,
  settled_at    TEXT
);

CREATE TABLE usage_events (
  id            TEXT PRIMARY KEY,
  library_id    TEXT NOT NULL REFERENCES libraries(id),
  invocation_id TEXT NOT NULL,
  kind          TEXT NOT NULL,
  amount        INTEGER NOT NULL,
  job_id        TEXT,
  run_id        TEXT,
  meta          TEXT NOT NULL DEFAULT '{}',
  created_at    TEXT NOT NULL
);

-- Logical concurrency slots: acquire = DELETE expired + conditional INSERT; count rows per name.
CREATE TABLE resource_slots (
  name        TEXT NOT NULL,
  holder      TEXT NOT NULL,
  epoch       INTEGER NOT NULL,
  lease_until TEXT NOT NULL,
  created_at  TEXT NOT NULL,
  PRIMARY KEY (name, holder)
);
CREATE INDEX idx_resource_slots_lease ON resource_slots(name, lease_until);

-- Write barrier for consistent backup; enforced by the single write layer, not the UI.
CREATE TABLE library_maintenance (
  library_id TEXT PRIMARY KEY REFERENCES libraries(id),
  mode       TEXT NOT NULL DEFAULT 'normal' CHECK (mode IN ('normal', 'write_barrier')),
  holder     TEXT,
  since      TEXT,
  note       TEXT
);

CREATE TABLE deletion_tombstones (
  id          TEXT PRIMARY KEY,
  library_id  TEXT NOT NULL,
  object_type TEXT NOT NULL,
  object_id   TEXT NOT NULL,
  deleted_at  TEXT NOT NULL,
  purge_after TEXT
);
CREATE INDEX idx_tombstones_obj ON deletion_tombstones(library_id, object_type, object_id);

CREATE TABLE backup_runs (
  id           TEXT PRIMARY KEY,
  library_id   TEXT NOT NULL REFERENCES libraries(id),
  status       TEXT NOT NULL DEFAULT 'running'
               CHECK (status IN ('running', 'complete', 'failed')),
  manifest_key TEXT,
  started_at   TEXT NOT NULL,
  completed_at TEXT,
  error        TEXT
);

-- FTS5 index over pre-tokenized text (app-layer CJK segmentation; see docs/plan.md §11.1).
-- entity_type: chunk | note | artifact | item_title. entity_id points back at the source row.
CREATE VIRTUAL TABLE search_index USING fts5(
  library_id UNINDEXED,
  entity_type UNINDEXED,
  entity_id UNINDEXED,
  text
);

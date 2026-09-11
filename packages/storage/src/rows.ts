import type {
  ArtifactStatus,
  ArtifactType,
  ContentQuality,
  IndexState,
  JobKind,
  JobStatus,
  ReadStatus,
  ResourceQuality,
} from "@recally/domain";

export interface LibraryRow {
  id: string;
  name: string;
  timezone: string;
  access_sub: string | null;
  settings: string;
  created_at: string;
  updated_at: string;
}

export interface ItemRow {
  id: string;
  library_id: string;
  original_url: string;
  normalized_url: string;
  title: string | null;
  saved_at: string;
  first_read_at: string | null;
  last_read_at: string | null;
  read_status: ReadStatus;
  current_snapshot_id: string | null;
  capture_generation: number;
  deleted_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface CaptureRunRow {
  id: string;
  library_id: string;
  item_id: string;
  job_id: string;
  generation: number;
  agent_runtime_version: string;
  skill_revision: string;
  toolset_version: string;
  policy_version: string;
  pipeline_version: string;
  model_config_version: string;
  pi_state_key: string | null;
  requested_at: string;
  target_url: string;
  status: string;
  outcome_code: string | null;
  budget: string;
  created_at: string;
  updated_at: string;
}

export interface JobRow {
  id: string;
  library_id: string;
  kind: JobKind;
  item_id: string | null;
  status: JobStatus;
  operation_key: string | null;
  workflow_instance_id: string | null;
  generation: number;
  parent_job_id: string | null;
  attempt_count: number;
  next_attempt_at: string | null;
  payload: string;
  result: string | null;
  error: string | null;
  cancel_requested: number;
  created_at: string;
  updated_at: string;
}

export interface OutboxRow {
  id: string;
  library_id: string;
  job_id: string;
  kind: string;
  payload: string;
  available_at: string;
  dispatched_at: string | null;
  created_at: string;
}

export interface SnapshotRow {
  id: string;
  library_id: string;
  item_id: string;
  capture_run_id: string;
  captured_at: string;
  final_url: string;
  capture_method: string;
  source_manifest_key: string;
  manifest_hash: string;
  content_quality: ContentQuality;
  resource_quality: ResourceQuality;
  committed_at: string;
  created_at: string;
}

export interface ContentRevisionRow {
  id: string;
  library_id: string;
  snapshot_id: string;
  item_id: string;
  extractor_version: string;
  selected_source_blocks: string;
  body_hash: string;
  article_key: string;
  blocks_key: string;
  language: string | null;
  token_count: number | null;
  created_at: string;
}

export interface ArtifactRow {
  id: string;
  library_id: string;
  item_id: string;
  type: ArtifactType;
  status: ArtifactStatus;
  input_snapshot_id: string | null;
  input_content_revision_id: string | null;
  model_id: string | null;
  prompt_version: string | null;
  pipeline_version: string | null;
  operation_key: string | null;
  output: string;
  citations: string;
  coverage: string | null;
  usage: string | null;
  created_at: string;
}

export interface ChunkRow {
  id: string;
  library_id: string;
  content_revision_id: string;
  item_id: string;
  ordinal: number;
  block_range: string;
  text: string;
  text_hash: string;
  tokenizer_version: string;
  embedding_version: string | null;
  vector_id: string | null;
  index_state: IndexState;
  created_at: string;
}

export interface ApiTokenRow {
  id: string;
  library_id: string;
  name: string;
  token_hash: string;
  scopes: string;
  last_used_at: string | null;
  expires_at: string | null;
  revoked_at: string | null;
  created_at: string;
}

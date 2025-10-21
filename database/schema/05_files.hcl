// Files Schema
// This file defines the files table for S3 resource mapping

schema "public" {
}

table "files" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = false
    type = uuid
  }

  column "original_url" {
    null = false
    type = text
  }

  column "s3_key" {
    null = false
    type = text
  }

  column "s3_url" {
    null = true
    type = text
  }

  column "file_name" {
    null = true
    type = text
  }

  column "file_type" {
    null = false
    type = varchar(255)
  }

  column "file_size" {
    null = true
    type = bigint
  }

  column "file_hash" {
    null = true
    type = text
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'")
  }

  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "files_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  index "unique_original_url" {
    unique  = true
    columns = [column.original_url]
  }

  index "unique_s3_key" {
    unique  = true
    columns = [column.s3_key]
  }

  index "idx_original_url" {
    columns = [column.original_url]
  }

  index "idx_s3_url" {
    columns = [column.s3_url]
  }

  index "idx_file_hash" {
    columns = [column.file_hash]
  }

  index "idx_file_type" {
    columns = [column.file_type]
  }

  index "idx_metadata" {
    columns = [column.metadata]
    type    = GIN
    ops     = sql("jsonb_path_ops")
  }

  index "idx_user_id" {
    columns = [column.user_id]
  }
}

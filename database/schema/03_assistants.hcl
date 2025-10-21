// Assistant System Schema
// This file defines AI assistant functionality with conversation management


table "assistants" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = true
    type = uuid
  }

  column "name" {
    null = false
    type = varchar(255)
  }

  column "description" {
    null = true
    type = text
  }

  column "system_prompt" {
    null = true
    type = text
  }

  column "model" {
    null = false
    type = varchar(32)
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'::JSONB")
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

  foreign_key "assistants_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  index "assistants_uuid_key" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_assistants_user_created_at" {
    columns = [column.user_id, column.created_at]
  }
}

table "assistant_threads" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = true
    type = uuid
  }

  column "assistant_id" {
    null = true
    type = uuid
  }

  column "name" {
    null = false
    type = varchar(255)
  }

  column "description" {
    null = true
    type = text
  }

  column "system_prompt" {
    null = true
    type = text
  }

  column "model" {
    null = false
    type = varchar(32)
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'::JSONB")
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

  foreign_key "assistant_threads_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  foreign_key "assistant_threads_assistant_id_fkey" {
    columns     = [column.assistant_id]
    ref_columns = [table.assistants.column.uuid]
  }

  index "assistant_threads_uuid_key" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_user_assistant_created_at" {
    columns = [column.user_id, column.assistant_id, column.created_at]
  }

  index "idx_assistant_threads_assistant_created_at" {
    columns = [column.assistant_id, column.created_at]
  }
}

table "assistant_messages" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = true
    type = uuid
  }

  column "assistant_id" {
    null = true
    type = uuid
  }

  column "thread_id" {
    null = true
    type = uuid
  }

  column "model" {
    null = true
    type = varchar(32)
  }

  column "role" {
    null = false
    type = varchar(255)
  }

  column "text" {
    null = true
    type = text
  }

  column "prompt_token" {
    null = true
    type = integer
  }

  column "completion_token" {
    null = true
    type = integer
  }

  column "embeddings" {
    null = true
    type = sql("vector(1536)")
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'::JSONB")
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

  foreign_key "assistant_messages_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  foreign_key "assistant_messages_assistant_id_fkey" {
    columns     = [column.assistant_id]
    ref_columns = [table.assistants.column.uuid]
  }

  foreign_key "assistant_messages_thread_id_fkey" {
    columns     = [column.thread_id]
    ref_columns = [table.assistant_threads.column.uuid]
  }

  index "assistant_messages_uuid_key" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_assistant_messages_user_created_at" {
    columns = [column.user_id, column.created_at]
  }

  index "idx_assistant_messages_assistant_created_at" {
    columns = [column.assistant_id, column.created_at]
  }

  index "idx_assistant_messages_thread_created_at" {
    columns = [column.thread_id, column.created_at]
  }
}

table "assistant_attachments" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = true
    type = uuid
  }

  column "assistant_id" {
    null = true
    type = uuid
  }

  column "thread_id" {
    null = true
    type = uuid
  }

  column "name" {
    null = true
    type = varchar(255)
  }

  column "type" {
    null = true
    type = varchar(255)
  }

  column "url" {
    null = true
    type = varchar(512)
  }

  column "size" {
    null = true
    type = integer
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'::JSONB")
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

  foreign_key "assistant_attachments_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  foreign_key "assistant_attachments_assistant_id_fkey" {
    columns     = [column.assistant_id]
    ref_columns = [table.assistants.column.uuid]
  }

  foreign_key "assistant_attachments_thread_id_fkey" {
    columns     = [column.thread_id]
    ref_columns = [table.assistant_threads.column.uuid]
  }

  index "assistant_attachments_uuid_key" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_assistant_attachments_user_created_at" {
    columns = [column.user_id, column.created_at]
  }

  index "idx_assistant_attachments_assistant_created_at" {
    columns = [column.assistant_id, column.created_at]
  }

  index "idx_assistant_attachments_thread_created_at" {
    columns = [column.thread_id, column.created_at]
  }
}

// IMPORTANT: This table name has a typo (3 d's) which must be preserved for compatibility
table "assistant_embedddings" {
  schema = schema.public

  column "id" {
    null = false
    type = serial
  }

  column "uuid" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = true
    type = uuid
  }

  column "attachment_id" {
    null = true
    type = uuid
  }

  column "text" {
    null = false
    type = text
  }

  column "embeddings" {
    null = true
    type = sql("vector(1536)")
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'::JSONB")
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

  foreign_key "assistant_embedddings_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  index "assistant_embedddings_uuid_key" {
    unique  = true
    columns = [column.uuid]
  }

  index "idx_user" {
    columns = [column.user_id]
  }

  index "idx_attachment" {
    columns = [column.attachment_id]
  }
}

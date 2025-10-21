// Content and Bookmark Schema
// This file defines both legacy content tables and modern bookmark system tables


// ============================================================================
// Legacy Content System Tables
// ============================================================================

table "content" {
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

  column "type" {
    null = false
    type = varchar(50)
  }

  column "title" {
    null = false
    type = text
  }

  column "description" {
    null = true
    type = text
  }

  column "url" {
    null = true
    type = text
  }

  column "domain" {
    null = true
    type = text
  }

  column "s3_key" {
    null = true
    type = text
  }

  column "summary" {
    null = true
    type = text
  }

  column "content" {
    null = true
    type = text
  }

  column "html" {
    null = true
    type = text
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'")
  }

  column "is_favorite" {
    null    = true
    type    = boolean
    default = sql("false")
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

  foreign_key "content_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  index "idx_content_user_id" {
    columns = [column.user_id]
  }

  index "idx_content_type" {
    columns = [column.type]
  }

  index "idx_content_url" {
    columns = [column.url]
  }

  index "idx_content_domain" {
    columns = [column.domain]
  }

  index "idx_content_created_at" {
    columns = [column.created_at]
  }

  index "idx_content_metadata" {
    columns = [column.metadata]
    type    = GIN
    // Note: jsonb_path_ops operator class will be added to migration manually if needed
  }
}

table "content_tags" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "name" {
    null = false
    type = varchar(50)
  }

  column "user_id" {
    null = false
    type = uuid
  }

  column "usage_count" {
    null    = true
    type    = integer
    default = sql("0")
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

  foreign_key "content_tags_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  index "content_tags_name_user_id_key" {
    unique  = true
    columns = [column.name, column.user_id]
  }

  index "idx_content_tags_name" {
    columns = [column.name]
  }

  index "idx_content_tags_user_id" {
    columns = [column.user_id]
  }
}

table "content_tags_mapping" {
  schema = schema.public

  column "content_id" {
    null = false
    type = uuid
  }

  column "tag_id" {
    null = false
    type = uuid
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
    columns = [column.content_id, column.tag_id]
  }

  foreign_key "content_tags_mapping_content_id_fkey" {
    columns     = [column.content_id]
    ref_columns = [table.content.column.id]
    on_delete   = CASCADE
  }

  foreign_key "content_tags_mapping_tag_id_fkey" {
    columns     = [column.tag_id]
    ref_columns = [table.content_tags.column.id]
    on_delete   = CASCADE
  }
}

table "content_folders" {
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

  column "name" {
    null = false
    type = varchar(100)
  }

  column "parent_id" {
    null = true
    type = uuid
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

  foreign_key "content_folders_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  foreign_key "content_folders_parent_id_fkey" {
    columns     = [column.parent_id]
    ref_columns = [table.content_folders.column.id]
  }

  index "idx_content_folders_user_id" {
    columns = [column.user_id]
  }
}

table "content_folders_mapping" {
  schema = schema.public

  column "content_id" {
    null = false
    type = uuid
  }

  column "folder_id" {
    null = false
    type = uuid
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
    columns = [column.content_id, column.folder_id]
  }

  foreign_key "content_folders_mapping_content_id_fkey" {
    columns     = [column.content_id]
    ref_columns = [table.content.column.id]
    on_delete   = CASCADE
  }

  foreign_key "content_folders_mapping_folder_id_fkey" {
    columns     = [column.folder_id]
    ref_columns = [table.content_folders.column.id]
    on_delete   = CASCADE
  }
}

table "content_share" {
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

  column "content_id" {
    null = true
    type = uuid
  }

  column "expires_at" {
    null    = true
    type    = timestamptz
    default = sql("NULL")
  }

  column "created_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "content_share_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  foreign_key "content_share_content_id_fkey" {
    columns     = [column.content_id]
    ref_columns = [table.content.column.id]
    on_delete   = CASCADE
  }

  index "content_share_user_id_idx" {
    columns = [column.user_id]
  }

  index "content_share_content_id_idx" {
    columns = [column.content_id]
  }
}

// ============================================================================
// Modern Bookmark System Tables
// ============================================================================

table "bookmark_content" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "type" {
    null = false
    type = varchar(50)
  }

  column "url" {
    null = false
    type = text
  }

  column "user_id" {
    null    = true
    type    = uuid
    default = sql("NULL")
  }

  column "title" {
    null = true
    type = text
  }

  column "description" {
    null = true
    type = text
  }

  column "domain" {
    null = true
    type = text
  }

  column "s3_key" {
    null = true
    type = text
  }

  column "summary" {
    null = true
    type = text
  }

  column "content" {
    null = true
    type = text
  }

  column "html" {
    null = true
    type = text
  }

  column "tags" {
    null    = true
    type    = sql("VARCHAR(50)[]")
    default = sql("'{}'")
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

  index "bookmark_content_url_user_id_key" {
    unique  = true
    columns = [column.url, column.user_id]
  }

  index "idx_bookmark_content_type" {
    columns = [column.type]
  }

  index "idx_bookmark_content_url" {
    columns = [column.url]
  }

  index "idx_bookmark_content_domain" {
    columns = [column.domain]
  }

  index "idx_bookmark_content_created_at" {
    columns = [column.created_at]
  }

  index "idx_bookmark_content_metadata" {
    columns = [column.metadata]
    type    = GIN
    // Note: jsonb_path_ops operator class will be added to migration manually if needed
  }
}

table "bookmarks" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "user_id" {
    null = true
    type = uuid
  }

  column "content_id" {
    null = true
    type = uuid
  }

  column "is_favorite" {
    null    = false
    type    = boolean
    default = sql("FALSE")
  }

  column "is_archive" {
    null    = false
    type    = boolean
    default = sql("FALSE")
  }

  column "metadata" {
    null    = true
    type    = jsonb
    default = sql("'{}'")
  }

  column "created_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "bookmarks_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
    on_delete   = CASCADE
  }

  foreign_key "bookmarks_content_id_fkey" {
    columns     = [column.content_id]
    ref_columns = [table.bookmark_content.column.id]
  }

  index "idx_bookmarks_user_created_at" {
    columns = [column.user_id, column.created_at]
  }

  index "idx_bookmarks_favorite" {
    columns = [column.user_id, column.is_favorite]
  }

  index "idx_bookmarks_archive" {
    columns = [column.user_id, column.is_archive]
  }

  index "idx_bookmarks_metadata" {
    columns = [column.metadata]
    type    = GIN
    // Note: jsonb_path_ops operator class will be added to migration manually if needed
  }
}

table "bookmark_tags" {
  schema = schema.public

  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "name" {
    null = false
    type = varchar(50)
  }

  column "user_id" {
    null = false
    type = uuid
  }

  column "created_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "bookmark_tags_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
    on_delete   = CASCADE
  }

  index "bookmark_tags_user_id_name_key" {
    unique  = true
    columns = [column.user_id, column.name]
  }

  index "idx_bookmark_tags_name" {
    columns = [column.name]
  }
}

table "bookmark_tags_mapping" {
  schema = schema.public

  column "bookmark_id" {
    null = false
    type = uuid
  }

  column "tag_id" {
    null = false
    type = uuid
  }

  column "created_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.bookmark_id, column.tag_id]
  }

  foreign_key "bookmark_tags_mapping_bookmark_id_fkey" {
    columns     = [column.bookmark_id]
    ref_columns = [table.bookmarks.column.id]
    on_delete   = CASCADE
  }

  foreign_key "bookmark_tags_mapping_tag_id_fkey" {
    columns     = [column.tag_id]
    ref_columns = [table.bookmark_tags.column.id]
    on_delete   = CASCADE
  }

  index "idx_bookmark_tags_mapping_tag_id" {
    columns = [column.tag_id]
  }
}

table "bookmark_share" {
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

  column "bookmark_id" {
    null = true
    type = uuid
  }

  column "expires_at" {
    null    = true
    type    = timestamptz
    default = sql("NULL")
  }

  column "created_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = true
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "bookmark_share_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.uuid]
  }

  foreign_key "bookmark_share_bookmark_id_fkey" {
    columns     = [column.bookmark_id]
    ref_columns = [table.bookmarks.column.id]
    on_delete   = CASCADE
  }

  index "bookmark_share_user_id_bookmark_id_key" {
    unique  = true
    columns = [column.user_id, column.bookmark_id]
  }

  index "idx_bookmark_share_user_id" {
    columns = [column.user_id]
  }

  index "idx_bookmark_share_content_id" {
    columns = [column.bookmark_id]
  }
}

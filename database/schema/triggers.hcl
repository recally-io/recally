// Database Triggers
// This file defines reusable trigger functions

schema "public" {
}

// update_updated_at_column - Automatically updates the updated_at timestamp
// This function is used by triggers on 18 tables across the schema
function "update_updated_at_column" {
  schema = schema.public
  lang   = PLpgSQL
  return = trigger
  as     = <<-SQL
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
  SQL
}

// Trigger definitions for cache table
trigger "update_cache_updated_at" {
  on     = table.cache
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

// Trigger definitions for users table
trigger "update_users_updated_at" {
  on     = table.users
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

// Trigger definitions for auth tables
trigger "update_oauth_connections_updated_at" {
  on     = table.auth_user_oauth_connections
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_auth_api_keys_updated_at" {
  on     = table.auth_api_keys
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

// Trigger definitions for assistant tables
trigger "update_assistants_updated_at" {
  on     = table.assistants
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_assistant_threads_updated_at" {
  on     = table.assistant_threads
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_assistant_messages_updated_at" {
  on     = table.assistant_messages
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_assistant_attachments_updated_at" {
  on     = table.assistant_attachments
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_assistant_embedddings_updated_at" {
  on     = table.assistant_embedddings
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

// Trigger definitions for legacy content tables
trigger "update_content_updated_at" {
  on     = table.content
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_content_tags_updated_at" {
  on     = table.content_tags
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_content_tags_mapping_updated_at" {
  on     = table.content_tags_mapping
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_content_folders_updated_at" {
  on     = table.content_folders
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_content_folders_mapping_updated_at" {
  on     = table.content_folders_mapping
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_content_share_updated_at" {
  on     = table.content_share
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

// Trigger definitions for modern bookmark tables
trigger "update_bookmark_content_updated_at" {
  on     = table.bookmark_content
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_bookmarks_updated_at" {
  on     = table.bookmarks
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_bookmark_tags_updated_at" {
  on     = table.bookmark_tags
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_bookmark_tags_mapping_updated_at" {
  on     = table.bookmark_tags_mapping
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

trigger "update_bookmark_share_updated_at" {
  on     = table.bookmark_share
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

// Trigger definitions for files table
trigger "update_files_updated_at" {
  on     = table.files
  before = true
  events = ["UPDATE"]
  foreach = ROW
  execute {
    function = function.update_updated_at_column
  }
}

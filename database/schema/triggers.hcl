// Trigger Function and Trigger Definitions
//
// NOTE: Functions and triggers require Atlas login for schema inspection
// These will be added manually to the generated migration SQL

// Trigger function for automatic updated_at timestamp updates
// Will be added to migration SQL as:
//
// CREATE OR REPLACE FUNCTION update_updated_at_column()
// RETURNS TRIGGER AS $$
// BEGIN
//     NEW.updated_at = CURRENT_TIMESTAMP;
//     RETURN NEW;
// END;
// $$ LANGUAGE plpgsql;

// TRIGGER ATTACHMENTS
// These triggers will be added manually to the generated migration SQL:
//
// For each table with an updated_at column, create a trigger:
// CREATE TRIGGER update_{table}_updated_at
//   BEFORE UPDATE ON {table}
//   FOR EACH ROW
//   EXECUTE FUNCTION update_updated_at_column();
//
// Tables requiring triggers (22 total):
// - cache
// - users
// - auth_user_oauth_connections
// - auth_api_keys
// - assistants
// - assistant_threads
// - assistant_messages
// - assistant_attachments
// - assistant_embedddings (note: typo with 3 d's is intentional)
// - content
// - content_tags
// - content_tags_mapping
// - content_folders
// - content_folders_mapping
// - content_share
// - files
// - bookmark_content
// - bookmarks
// - bookmark_tags
// - bookmark_tags_mapping
// - bookmark_share
// - auth_revoked_tokens

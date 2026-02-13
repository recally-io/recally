-- ============================================================================
-- Initial Schema Migration
-- Generated: 2025-10-21 22:52:51
--
-- This migration creates the complete database schema for Recally including:
-- - 17 tables across 4 domain areas
-- - Vector extension for embeddings
-- - Trigger function for automatic updated_at timestamps
-- - BM25 indexes for full-text search (ParadeDB)
-- ============================================================================

-- ============================================================================
-- Extensions
-- ============================================================================

-- Vector extension for embeddings (pgvector)
-- Note: ParadeDB includes pg_search pre-installed
CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================================================
-- Trigger Function
-- ============================================================================

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- Tables
-- ============================================================================

-- Infrastructure
-- ----------------------------------------------------------------------------

CREATE TABLE "cache" (
    "id" SERIAL NOT NULL,
    "domain" VARCHAR(255) NOT NULL,
    "key" TEXT NOT NULL,
    "value" JSONB NOT NULL,
    "expires_at" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Authentication & Authorization
-- ----------------------------------------------------------------------------

CREATE TABLE "users" (
    "id" SERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "username" VARCHAR(255),
    "password_hash" TEXT,
    "email" VARCHAR(255),
    "phone" VARCHAR(50),
    "activate_assistant_id" UUID,
    "activate_thread_id" UUID,
    "status" VARCHAR(255) NOT NULL DEFAULT 'pending',
    "settings" JSONB DEFAULT '{}'::JSONB,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "auth_user_oauth_connections" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "provider" VARCHAR(50) NOT NULL,
    "provider_user_id" VARCHAR(255) NOT NULL,
    "provider_email" VARCHAR(255),
    "access_token" TEXT,
    "refresh_token" TEXT,
    "token_expires_at" TIMESTAMPTZ,
    "provider_data" JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "auth_api_keys" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "key_prefix" VARCHAR(8) NOT NULL,
    "key_hash" VARCHAR(255) NOT NULL,
    "scopes" TEXT[] NOT NULL,
    "expires_at" TIMESTAMPTZ,
    "last_used_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "auth_revoked_tokens" (
    "jti" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "expires_at" TIMESTAMPTZ NOT NULL,
    "revoked_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "reason" VARCHAR(100),
    PRIMARY KEY ("jti")
);

-- Assistants
-- ----------------------------------------------------------------------------

CREATE TABLE "assistants" (
    "id" SERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "system_prompt" TEXT,
    "model" VARCHAR(32) NOT NULL,
    "metadata" JSONB DEFAULT '{}'::JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "assistant_threads" (
    "id" SERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID,
    "assistant_id" UUID,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "system_prompt" TEXT,
    "model" VARCHAR(32) NOT NULL,
    "metadata" JSONB DEFAULT '{}'::JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "assistant_messages" (
    "id" SERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID,
    "assistant_id" UUID,
    "thread_id" UUID,
    "model" VARCHAR(32),
    "role" VARCHAR(255) NOT NULL,
    "text" TEXT,
    "prompt_token" INTEGER,
    "completion_token" INTEGER,
    "embeddings" vector(1536),
    "metadata" JSONB DEFAULT '{}'::JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "assistant_attachments" (
    "id" SERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID,
    "assistant_id" UUID,
    "thread_id" UUID,
    "name" VARCHAR(255),
    "type" VARCHAR(255),
    "url" VARCHAR(512),
    "size" INTEGER,
    "metadata" JSONB DEFAULT '{}'::JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "assistant_embedddings" (
    "id" SERIAL NOT NULL,
    "uuid" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID,
    "attachment_id" UUID,
    "text" TEXT NOT NULL,
    "embeddings" vector(1536),
    "metadata" JSONB DEFAULT '{}'::JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Legacy Content System
-- ----------------------------------------------------------------------------

CREATE TABLE "content" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "type" VARCHAR(50) NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT,
    "url" TEXT,
    "domain" TEXT,
    "s3_key" TEXT,
    "summary" TEXT,
    "content" TEXT,
    "html" TEXT,
    "metadata" JSONB DEFAULT '{}',
    "is_favorite" BOOLEAN DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "content_tags" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR(50) NOT NULL,
    "user_id" UUID NOT NULL,
    "usage_count" INTEGER DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "content_tags_mapping" (
    "content_id" UUID NOT NULL,
    "tag_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("content_id", "tag_id")
);

CREATE TABLE "content_folders" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "name" VARCHAR(100) NOT NULL,
    "parent_id" UUID,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "content_folders_mapping" (
    "content_id" UUID NOT NULL,
    "folder_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("content_id", "folder_id")
);

CREATE TABLE "content_share" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "content_id" UUID,
    "expires_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Modern Bookmark System
-- ----------------------------------------------------------------------------

CREATE TABLE "bookmark_content" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "type" VARCHAR(50) NOT NULL,
    "url" TEXT NOT NULL,
    "user_id" UUID DEFAULT NULL,
    "title" TEXT,
    "description" TEXT,
    "domain" TEXT,
    "s3_key" TEXT,
    "summary" TEXT,
    "content" TEXT,
    "html" TEXT,
    "tags" VARCHAR(50)[] DEFAULT '{}',
    "metadata" JSONB DEFAULT '{}',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "bookmarks" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID,
    "content_id" UUID,
    "is_favorite" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_archive" BOOLEAN NOT NULL DEFAULT FALSE,
    "metadata" JSONB DEFAULT '{}',
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "bookmark_tags" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR(50) NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "bookmark_tags_mapping" (
    "bookmark_id" UUID NOT NULL,
    "tag_id" UUID NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("bookmark_id", "tag_id")
);

CREATE TABLE "bookmark_share" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "bookmark_id" UUID,
    "expires_at" TIMESTAMPTZ DEFAULT NULL,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Files
-- ----------------------------------------------------------------------------

CREATE TABLE "files" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "original_url" TEXT NOT NULL,
    "s3_key" TEXT NOT NULL,
    "s3_url" TEXT,
    "file_name" TEXT,
    "file_type" VARCHAR(255) NOT NULL,
    "file_size" BIGINT,
    "file_hash" TEXT,
    "metadata" JSONB DEFAULT '{}',
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- ============================================================================
-- Unique Indexes
-- ============================================================================
-- NOTE: Unique indexes MUST be created before foreign keys that reference them

-- cache
CREATE UNIQUE INDEX "uni_cache_domain_key" ON "cache" ("domain", "key");

-- users
CREATE UNIQUE INDEX "users_uuid_key" ON "users" ("uuid");
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
CREATE UNIQUE INDEX "idx_users_email" ON "users" (LOWER(email)) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX "idx_users_phone" ON "users" ("phone") WHERE phone IS NOT NULL;
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username") WHERE username IS NOT NULL;

-- auth_user_oauth_connections
CREATE UNIQUE INDEX "uq_oauth_connection" ON "auth_user_oauth_connections" ("provider", "provider_user_id");

-- auth_api_keys
CREATE UNIQUE INDEX "uq_user_key_name" ON "auth_api_keys" ("user_id", "name");

-- auth_revoked_tokens
CREATE UNIQUE INDEX "uq_revoked_token" ON "auth_revoked_tokens" ("user_id", "jti");

-- content_tags
CREATE UNIQUE INDEX "content_tags_name_user_id_key" ON "content_tags" ("name", "user_id");

-- bookmark_content
CREATE UNIQUE INDEX "bookmark_content_url_user_id_key" ON "bookmark_content" ("url", "user_id");

-- bookmark_tags
CREATE UNIQUE INDEX "bookmark_tags_user_id_name_key" ON "bookmark_tags" ("user_id", "name");

-- bookmark_share
CREATE UNIQUE INDEX "bookmark_share_user_id_bookmark_id_key" ON "bookmark_share" ("user_id", "bookmark_id");

-- files
CREATE UNIQUE INDEX "unique_original_url" ON "files" ("original_url");
CREATE UNIQUE INDEX "unique_s3_key" ON "files" ("s3_key");

-- ============================================================================
-- Foreign Keys
-- ============================================================================

-- auth_user_oauth_connections
ALTER TABLE "auth_user_oauth_connections" ADD CONSTRAINT "auth_user_oauth_connections_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid") ON DELETE CASCADE;

-- auth_api_keys
ALTER TABLE "auth_api_keys" ADD CONSTRAINT "auth_api_keys_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid") ON DELETE CASCADE;

-- auth_revoked_tokens
ALTER TABLE "auth_revoked_tokens" ADD CONSTRAINT "auth_revoked_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");

-- content
ALTER TABLE "content" ADD CONSTRAINT "content_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");

-- content_tags
ALTER TABLE "content_tags" ADD CONSTRAINT "content_tags_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");

-- content_tags_mapping
ALTER TABLE "content_tags_mapping" ADD CONSTRAINT "content_tags_mapping_content_id_fkey" FOREIGN KEY ("content_id") REFERENCES "content" ("id") ON DELETE CASCADE;
ALTER TABLE "content_tags_mapping" ADD CONSTRAINT "content_tags_mapping_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "content_tags" ("id") ON DELETE CASCADE;

-- content_folders
ALTER TABLE "content_folders" ADD CONSTRAINT "content_folders_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");
ALTER TABLE "content_folders" ADD CONSTRAINT "content_folders_parent_id_fkey" FOREIGN KEY ("parent_id") REFERENCES "content_folders" ("id");

-- content_folders_mapping
ALTER TABLE "content_folders_mapping" ADD CONSTRAINT "content_folders_mapping_content_id_fkey" FOREIGN KEY ("content_id") REFERENCES "content" ("id") ON DELETE CASCADE;
ALTER TABLE "content_folders_mapping" ADD CONSTRAINT "content_folders_mapping_folder_id_fkey" FOREIGN KEY ("folder_id") REFERENCES "content_folders" ("id") ON DELETE CASCADE;

-- content_share
ALTER TABLE "content_share" ADD CONSTRAINT "content_share_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");
ALTER TABLE "content_share" ADD CONSTRAINT "content_share_content_id_fkey" FOREIGN KEY ("content_id") REFERENCES "content" ("id") ON DELETE CASCADE;

-- bookmarks
ALTER TABLE "bookmarks" ADD CONSTRAINT "bookmarks_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid") ON DELETE CASCADE;
ALTER TABLE "bookmarks" ADD CONSTRAINT "bookmarks_content_id_fkey" FOREIGN KEY ("content_id") REFERENCES "bookmark_content" ("id");

-- bookmark_tags
ALTER TABLE "bookmark_tags" ADD CONSTRAINT "bookmark_tags_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid") ON DELETE CASCADE;

-- bookmark_tags_mapping
ALTER TABLE "bookmark_tags_mapping" ADD CONSTRAINT "bookmark_tags_mapping_bookmark_id_fkey" FOREIGN KEY ("bookmark_id") REFERENCES "bookmarks" ("id") ON DELETE CASCADE;
ALTER TABLE "bookmark_tags_mapping" ADD CONSTRAINT "bookmark_tags_mapping_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "bookmark_tags" ("id") ON DELETE CASCADE;

-- bookmark_share
ALTER TABLE "bookmark_share" ADD CONSTRAINT "bookmark_share_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");
ALTER TABLE "bookmark_share" ADD CONSTRAINT "bookmark_share_bookmark_id_fkey" FOREIGN KEY ("bookmark_id") REFERENCES "bookmarks" ("id") ON DELETE CASCADE;

-- files
ALTER TABLE "files" ADD CONSTRAINT "files_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("uuid");

-- ============================================================================
-- Check Constraints
-- ============================================================================

-- users
ALTER TABLE "users" ADD CONSTRAINT "users_email_check" CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$');
ALTER TABLE "users" ADD CONSTRAINT "users_contact_check" CHECK (email IS NOT NULL OR phone IS NOT NULL OR username IS NOT NULL);

-- auth_api_keys
ALTER TABLE "auth_api_keys" ADD CONSTRAINT "ck_key_expiry" CHECK (expires_at IS NULL OR expires_at > created_at);

-- auth_revoked_tokens
ALTER TABLE "auth_revoked_tokens" ADD CONSTRAINT "ck_token_revocation" CHECK (expires_at > revoked_at);

-- ============================================================================
-- Regular Indexes
-- ============================================================================

-- auth_user_oauth_connections
CREATE INDEX "idx_oauth_user_id" ON "auth_user_oauth_connections" ("user_id");
CREATE INDEX "idx_oauth_provider_lookup" ON "auth_user_oauth_connections" ("provider", "provider_user_id");
CREATE INDEX "idx_oauth_token_expiry" ON "auth_user_oauth_connections" ("token_expires_at") WHERE token_expires_at IS NOT NULL;

-- auth_api_keys
CREATE INDEX "idx_auth_api_keys_prefix" ON "auth_api_keys" ("key_prefix");
CREATE INDEX "idx_auth_api_keys_user" ON "auth_api_keys" ("user_id");
CREATE INDEX "idx_auth_api_keys_expiry" ON "auth_api_keys" ("expires_at") WHERE expires_at IS NOT NULL;

-- auth_revoked_tokens
CREATE INDEX "idx_auth_revoked_tokens_expiry" ON "auth_revoked_tokens" ("expires_at");
CREATE INDEX "idx_auth_revoked_tokens_user" ON "auth_revoked_tokens" ("user_id");

-- content
CREATE INDEX "idx_content_user_id" ON "content" ("user_id");
CREATE INDEX "idx_content_type" ON "content" ("type");
CREATE INDEX "idx_content_url" ON "content" ("url");
CREATE INDEX "idx_content_domain" ON "content" ("domain");
CREATE INDEX "idx_content_created_at" ON "content" ("created_at");
CREATE INDEX "idx_content_metadata" ON "content" USING GIN ("metadata");

-- content_tags
CREATE INDEX "idx_content_tags_name" ON "content_tags" ("name");
CREATE INDEX "idx_content_tags_user_id" ON "content_tags" ("user_id");

-- content_folders
CREATE INDEX "idx_content_folders_user_id" ON "content_folders" ("user_id");

-- content_share
CREATE INDEX "content_share_user_id_idx" ON "content_share" ("user_id");
CREATE INDEX "content_share_content_id_idx" ON "content_share" ("content_id");

-- bookmark_content
CREATE INDEX "idx_bookmark_content_type" ON "bookmark_content" ("type");
CREATE INDEX "idx_bookmark_content_url" ON "bookmark_content" ("url");
CREATE INDEX "idx_bookmark_content_domain" ON "bookmark_content" ("domain");
CREATE INDEX "idx_bookmark_content_created_at" ON "bookmark_content" ("created_at");
CREATE INDEX "idx_bookmark_content_metadata" ON "bookmark_content" USING GIN ("metadata");

-- bookmarks
CREATE INDEX "idx_bookmarks_user_created_at" ON "bookmarks" ("user_id", "created_at");
CREATE INDEX "idx_bookmarks_favorite" ON "bookmarks" ("user_id", "is_favorite");
CREATE INDEX "idx_bookmarks_archive" ON "bookmarks" ("user_id", "is_archive");
CREATE INDEX "idx_bookmarks_metadata" ON "bookmarks" USING GIN ("metadata");

-- bookmark_tags
CREATE INDEX "idx_bookmark_tags_name" ON "bookmark_tags" ("name");

-- bookmark_tags_mapping
CREATE INDEX "idx_bookmark_tags_mapping_tag_id" ON "bookmark_tags_mapping" ("tag_id");

-- bookmark_share
CREATE INDEX "idx_bookmark_share_user_id" ON "bookmark_share" ("user_id");
CREATE INDEX "idx_bookmark_share_content_id" ON "bookmark_share" ("bookmark_id");

-- files
CREATE INDEX "idx_original_url" ON "files" ("original_url");
CREATE INDEX "idx_s3_url" ON "files" ("s3_url");
CREATE INDEX "idx_file_hash" ON "files" ("file_hash");
CREATE INDEX "idx_file_type" ON "files" ("file_type");
CREATE INDEX "idx_metadata" ON "files" USING GIN ("metadata");
CREATE INDEX "idx_user_id" ON "files" ("user_id");

-- ============================================================================
-- Triggers
-- ============================================================================

-- Automatic updated_at timestamp triggers (21 tables)
CREATE TRIGGER update_cache_updated_at
    BEFORE UPDATE ON "cache"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON "users"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_auth_user_oauth_connections_updated_at
    BEFORE UPDATE ON "auth_user_oauth_connections"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_auth_api_keys_updated_at
    BEFORE UPDATE ON "auth_api_keys"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_updated_at
    BEFORE UPDATE ON "content"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_tags_updated_at
    BEFORE UPDATE ON "content_tags"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_tags_mapping_updated_at
    BEFORE UPDATE ON "content_tags_mapping"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_folders_updated_at
    BEFORE UPDATE ON "content_folders"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_folders_mapping_updated_at
    BEFORE UPDATE ON "content_folders_mapping"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_share_updated_at
    BEFORE UPDATE ON "content_share"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookmark_content_updated_at
    BEFORE UPDATE ON "bookmark_content"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookmarks_updated_at
    BEFORE UPDATE ON "bookmarks"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookmark_tags_updated_at
    BEFORE UPDATE ON "bookmark_tags"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookmark_tags_mapping_updated_at
    BEFORE UPDATE ON "bookmark_tags_mapping"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookmark_share_updated_at
    BEFORE UPDATE ON "bookmark_share"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_files_updated_at
    BEFORE UPDATE ON "files"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- ============================================================================
-- BM25 Full-Text Search Indexes (ParadeDB)
-- ============================================================================

-- BM25 index on content table
-- https://docs.paradedb.com/documentation/indexing/create_index
CREATE INDEX idx_content_bm25_search
ON content
USING bm25 (id, title, description, summary, content, metadata)
WITH (key_field = 'id');

-- BM25 index on bookmark_content table
CREATE INDEX idx_bookmark_content_bm25_search
ON bookmark_content
USING bm25(id, title, description, summary, content, metadata)
WITH (key_field = 'id');
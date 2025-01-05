-- Drop content folder mapping table and its trigger
DROP TABLE IF EXISTS content_folders_mapping;

-- Drop content folders table and its trigger
DROP TABLE IF EXISTS content_folders;

-- Drop content tags mapping table and its trigger
DROP TABLE IF EXISTS content_tags_mapping;

-- Drop content tags table and its trigger
DROP TABLE IF EXISTS content_tags;

-- Drop content table and its indexes/trigger
DROP INDEX IF EXISTS idx_content_bm25_search;
DROP INDEX IF EXISTS idx_content_metadata;
DROP INDEX IF EXISTS idx_content_created_at;
DROP INDEX IF EXISTS idx_content_type;
DROP INDEX IF EXISTS idx_content_user_id;
DROP TABLE IF EXISTS content;

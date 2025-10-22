DROP TRIGGER IF EXISTS update_files_updated_at ON files;
DROP INDEX IF EXISTS idx_original_url;
DROP INDEX IF EXISTS idx_s3_url;
DROP INDEX IF EXISTS idx_file_hash;
DROP INDEX IF EXISTS idx_file_type;
DROP INDEX IF EXISTS idx_metadata;
DROP TABLE IF EXISTS files;

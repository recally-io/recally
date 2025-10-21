DROP TRIGGER IF EXISTS update_content_share_updated_at ON content_share;
DROP INDEX IF EXISTS content_share_content_id_idx;
DROP INDEX IF EXISTS content_share_user_id_idx;
DROP TABLE IF EXISTS content_share;

CREATE TABLE IF NOT EXISTS content_share (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(uuid),
    content_id uuid REFERENCES content(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS content_share_user_id_idx ON content_share(user_id);
CREATE INDEX IF NOT EXISTS content_share_content_id_idx ON content_share(content_id);

DROP TRIGGER IF EXISTS update_content_share_updated_at ON content_share;
CREATE TRIGGER update_content_share_updated_at
    BEFORE UPDATE ON content_share
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

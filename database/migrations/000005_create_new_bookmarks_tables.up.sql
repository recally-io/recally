CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid),
    url TEXT NOT NULL,
    title VARCHAR(255),
    summary TEXT,
    summary_embeddings vector (1536) NOT NULL,
    content TEXT,
    content_embeddings vector (1536) NOT NULL,
    html TEXT,
    metadata JSONB,
    screenshot TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_created_at ON bookmarks (user_id, updated_at);
CREATE INDEX IF NOT EXISTS idx_uuid ON bookmarks (uuid);
CREATE INDEX IF NOT EXISTS idx_url ON bookmarks (url);

-- Create PGVector Extension
CREATE EXTENSION IF NOT EXISTS vector;

--- Create vector table
CREATE TABLE
    IF NOT EXISTS text_embeddings (
        id BIGSERIAL PRIMARY KEY,
        metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
        user_id string NOT NULL,
        text TEXT NOT NULL,
        embeddings vector (1536) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS idx_user ON text_embeddings (user_id);

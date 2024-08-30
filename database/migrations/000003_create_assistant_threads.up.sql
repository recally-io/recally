-- Create PGVector Extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create Tables Assistant
CREATE TABLE IF NOT EXISTS assistants (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    system_prompt TEXT,
    model VARCHAR(32) NOT NULL,
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_created_at ON assistants (user_id, created_at);


CREATE TABLE IF NOT EXISTS assistant_threads (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid),
    assistant_id UUID REFERENCES assistants(uuid),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    system_prompt TEXT,
    model VARCHAR(32) NOT NULL,
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_assistant_created_at ON assistant_threads (user_id, assistant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_assistant_created_at ON assistant_threads (assistant_id, created_at);


CREATE TABLE IF NOT EXISTS assistant_messages (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid),
    assistant_id UUID REFERENCES assistants(uuid),
    thread_id UUID REFERENCES assistant_threads(uuid),
    model VARCHAR(32),
    role VARCHAR(255) NOT NULL,
    text TEXT,
    prompt_token INTEGER,
    completion_token INTEGER,
    embeddings vector (1536),
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_created_at ON assistant_messages (user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_assistant_created_at ON assistant_messages (assistant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_thread_created_at ON assistant_messages (thread_id, created_at);


CREATE TABLE IF NOT EXISTS assistant_attachments (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid),
    assistant_id UUID REFERENCES assistants(uuid),
    thread_id UUID REFERENCES assistant_threads(uuid),
    name VARCHAR(255),
    type VARCHAR(255),
    url VARCHAR(512),
    size INTEGER,
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_created_at ON assistant_attachments (user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_assistant_created_at ON assistant_attachments (assistant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_thread_created_at ON assistant_attachments (thread_id, created_at);


-- Create Text Embeddings Table
CREATE TABLE IF NOT EXISTS assistant_embedddings (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(uuid),
    attachment_id UUID,
    text TEXT NOT NULL,
    embeddings vector (1536) NOT NULL,
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user ON assistant_embedddings (user_id);
CREATE INDEX IF NOT EXISTS idx_attachment ON assistant_embedddings (attachment_id);

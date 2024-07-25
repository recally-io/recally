CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    github VARCHAR(255) UNIQUE,
    google VARCHAR(255) UNIQUE,
    telegram VARCHAR(255) UNIQUE,
    activate_assistant_id UUID,
    activate_thread_id UUID,
    status VARCHAR(255) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
    model VARCHAR(32) NOT NULL,
    is_long_term_memory BOOLEAN DEFAULT FALSE,
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
    thread_id UUID REFERENCES assistant_threads(uuid),
    model VARCHAR(32),
    token int,
    role VARCHAR(255) NOT NULL,
    text TEXT,
    attachments UUID[],
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_thread_created_at ON assistant_messages (user_id, thread_id, created_at);
CREATE INDEX IF NOT EXISTS idx_thread_created_at ON assistant_messages (thread_id, created_at);


CREATE TABLE IF NOT EXISTS assistant_attachments (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid),
    entity VARCHAR(255) NOT NULL,
    entity_id UUID,
    file_type VARCHAR(255),
    file_url VARCHAR(255),
    size INTEGER,
    metadata JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_entity_created_at ON assistant_attachments (user_id, entity, entity_id, created_at);
CREATE INDEX IF NOT EXISTS idx_entity_id ON assistant_attachments (entity_id);


CREATE TABLE IF NOT EXISTS assistant_embedddings (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(uuid),
    message_id UUID,
    attachment_id UUID,
    text TEXT NOT NULL,
    embeddings vector (1536) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_message ON assistant_embedddings (message_id);
CREATE INDEX IF NOT EXISTS idx_attachment ON assistant_embedddings (attachment_id);
CREATE INDEX IF NOT EXISTS idx_user ON assistant_embedddings (user_id);

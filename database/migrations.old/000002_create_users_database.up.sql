CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(255),
    password_hash TEXT,
    email VARCHAR(255) UNIQUE,
    github VARCHAR(255) UNIQUE,
    google VARCHAR(255) UNIQUE,
    telegram VARCHAR(255) UNIQUE,
    activate_assistant_id UUID,
    activate_thread_id UUID,
    status VARCHAR(255) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

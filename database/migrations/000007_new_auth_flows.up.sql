-- This migration sets up the authentication and authorization schema
-- It includes user management, OAuth integration, verification systems, and RBAC

ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE users DROP COLUMN IF EXISTS github;
ALTER TABLE users DROP COLUMN IF EXISTS google;
ALTER TABLE users DROP COLUMN IF EXISTS telegram;
ALTER TABLE users ADD COLUMN IF NOT EXISTS settings JSONB DEFAULT '{}'::JSONB;
DO $$ BEGIN
    ALTER TABLE users ADD CONSTRAINT users_email_check 
        CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
    EXCEPTION WHEN duplicate_object THEN NULL;
END $$;
DO $$ BEGIN
    ALTER TABLE users ADD CONSTRAINT users_contact_check 
        CHECK (email IS NOT NULL OR phone IS NOT NULL OR username IS NOT NULL);
    EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- CREATE TABLE IF NOT EXISTS users (
--     id SERIAL PRIMARY KEY,
--     uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
--     username VARCHAR(255),
--     email VARCHAR(255) UNIQUE,
--     phone VARCHAR(50) UNIQUE,
--     password_hash TEXT,
--     activate_assistant_id UUID,
--     activate_thread_id UUID,
--     status VARCHAR(255) NOT NULL DEFAULT 'pending',
--     settings JSONB DEFAULT '{}'::JSONB,
--     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

-- Add partial indexes for performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (LOWER(email)) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users (phone) WHERE phone IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username) WHERE username IS NOT NULL;

-- Add trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- OAuth Connections: Manages third-party authentication providers
-- Enables users to link multiple OAuth providers to their account
-- Stores necessary tokens and provider-specific data
CREATE TABLE IF NOT EXISTS auth_user_oauth_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,           -- OAuth provider name (e.g., 'google', 'github')
    provider_user_id VARCHAR(255) NOT NULL,  -- User's ID in the provider's system
    provider_email VARCHAR(255),
    access_token TEXT,                       -- OAuth access token for API calls
    refresh_token TEXT,                      -- Token for refreshing access_token
    token_expires_at TIMESTAMPTZ,
    provider_data JSONB,                     -- Flexible storage for provider-specific information
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_oauth_connection UNIQUE(provider, provider_user_id)
);

CREATE INDEX IF NOT EXISTS idx_oauth_user_id ON auth_user_oauth_connections(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_provider_lookup ON auth_user_oauth_connections(provider, provider_user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_token_expiry ON auth_user_oauth_connections(token_expires_at) WHERE token_expires_at IS NOT NULL;

DROP TRIGGER IF EXISTS update_oauth_connections_updated_at ON auth_user_oauth_connections;
CREATE TRIGGER update_oauth_connections_updated_at
    BEFORE UPDATE ON auth_user_oauth_connections
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- API Key Management: Enables programmatic access to the API
-- Supports multiple active keys per user with different scopes
-- Implements key rotation and activity tracking
CREATE TABLE IF NOT EXISTS auth_api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,           -- User-defined name for the key
    key_prefix VARCHAR(8) NOT NULL,       -- For key identification without exposing full key
    key_hash VARCHAR(255) NOT NULL,       -- Securely stored complete key hash
    scopes TEXT[] NOT NULL,               -- Array of permissions granted to this key
    expires_at TIMESTAMPTZ,               -- Optional key expiration
    last_used_at TIMESTAMPTZ,             -- Tracks key usage
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_user_key_name UNIQUE(user_id, name),
    CONSTRAINT ck_key_expiry CHECK (expires_at IS NULL OR expires_at > created_at)
);

CREATE INDEX IF NOT EXISTS idx_auth_api_keys_prefix ON auth_api_keys(key_prefix);
CREATE INDEX IF NOT EXISTS idx_auth_api_keys_user ON auth_api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_api_keys_expiry ON auth_api_keys(expires_at) WHERE expires_at IS NOT NULL;

DROP TRIGGER IF EXISTS update_auth_api_keys_updated_at ON auth_api_keys;
CREATE TRIGGER update_auth_api_keys_updated_at
    BEFORE UPDATE ON auth_api_keys
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Token Revocation: Implements secure logout and token invalidation
-- Tracks revoked JWT tokens until their natural expiration
-- Supports various revocation reasons for audit purposes
CREATE TABLE IF NOT EXISTS auth_revoked_tokens (
    jti UUID PRIMARY KEY,                 -- Unique JWT identifier
    user_id UUID NOT NULL REFERENCES users(uuid),
    expires_at TIMESTAMPTZ NOT NULL,      -- When token would naturally expire
    revoked_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason VARCHAR(100),                  -- Reason for revocation (logout, security concern, etc.)
    CONSTRAINT uq_revoked_token UNIQUE (user_id, jti),
    CONSTRAINT ck_token_revocation CHECK (expires_at > revoked_at)
);

CREATE INDEX IF NOT EXISTS idx_auth_revoked_tokens_expiry ON auth_revoked_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_auth_revoked_tokens_user ON auth_revoked_tokens(user_id);

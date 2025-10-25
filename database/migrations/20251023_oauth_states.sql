-- OAuth state management for CSRF protection
CREATE TABLE IF NOT EXISTS auth_oauth_states (
    state VARCHAR(64) PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    redirect_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- Index for cleanup queries
CREATE INDEX idx_oauth_states_expiry ON auth_oauth_states(expires_at);

-- Automatic cleanup of expired states (run via cron or scheduled job)
-- DELETE FROM auth_oauth_states WHERE expires_at < NOW();

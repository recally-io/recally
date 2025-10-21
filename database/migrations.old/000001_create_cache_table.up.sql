-- Create a table to store cache data
CREATE TABLE IF NOT EXISTS cache (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uni_cache_domain_key ON cache (domain, key);

-- Create a table to store cache data
CREATE TABLE
    IF NOT EXISTS cache (
        id SERIAL PRIMARY KEY, -- Unique identifier for each cache entry
        key TEXT UNIQUE NOT NULL, -- The cache key
        VALUE TEXT NOT NULL, -- The cached value
        expires_at TIMESTAMP, -- The expiration timestamp for the cache entry
        created_at TIMESTAMP, -- The timestamp when the cache entry was created
        updated_at TIMESTAMP -- The timestamp when the cache entry was last updated
    );
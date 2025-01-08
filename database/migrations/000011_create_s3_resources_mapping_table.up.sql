-- insert the dummy user
INSERT INTO users (username, password_hash, status) 
VALUES ('dummy_user', 'dummy_hash', 'active');

-- Create a table to store data of files in S3
CREATE TABLE IF NOT EXISTS files (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(uuid),
    original_url TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    s3_url TEXT, -- s3 public URL
    file_name TEXT,
    file_type VARCHAR(255) NOT NULL,
    file_size BIGINT,  -- size in bytes
    file_hash TEXT, -- For duplicate detection
    metadata JSONB DEFAULT '{}',    -- Flexible metadata storage
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT unique_original_url UNIQUE(original_url),
    CONSTRAINT unique_s3_key UNIQUE(s3_key)
);


-- add trigger to update the updated_at column
CREATE TRIGGER update_files_updated_at
BEFORE UPDATE ON files
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- add indeses for frequently queried columns
CREATE INDEX idx_original_url ON files (original_url);
CREATE INDEX idx_s3_url ON files (s3_url);
CREATE INDEX idx_file_hash ON files (file_hash);
CREATE INDEX idx_file_type ON files (file_type);
CREATE INDEX idx_metadata ON files USING GIN (metadata jsonb_path_ops);
CREATE INDEX idx_user_id ON files (user_id);

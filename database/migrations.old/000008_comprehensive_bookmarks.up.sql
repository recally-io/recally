-- Content table (main bookmarks table)
CREATE TABLE content (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(uuid),
    type VARCHAR(50) NOT NULL, -- type of content (bookmark, rss, newsletter, pdf, epub, image, podcast, video, etc.)
    title TEXT NOT NULL,
    description TEXT,
    url TEXT, -- URL of the bookmark,  news article, etc.
    domain TEXT, -- domain of the URL
    s3_key TEXT, -- S3 key for storing content
    summary TEXT, -- AI generated summary
    content TEXT, -- markdown content
    html TEXT, -- html content
    metadata JSONB DEFAULT '{}',
    is_favorite BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_content_user_id ON content(user_id);
CREATE INDEX idx_content_type ON content(type);
CREATE INDEX idx_content_url ON content(url);
CREATE INDEX idx_content_domain ON content(domain);
CREATE INDEX idx_content_created_at ON content(created_at);
CREATE INDEX idx_content_metadata ON content USING gin(metadata jsonb_path_ops);

-- BM25 index on content
-- https://docs.paradedb.com/documentation/indexing/create_index
CREATE INDEX idx_content_bm25_search ON content
USING bm25 (id, title, description, summary, content, metadata)
WITH (key_field='id');

DROP TRIGGER IF EXISTS update_content_updated_at ON content;
CREATE TRIGGER update_content_updated_at
    BEFORE UPDATE ON content
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Tags table
CREATE TABLE content_tags (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    user_id uuid NOT NULL REFERENCES users(uuid),
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, user_id)
);

CREATE INDEX idx_content_tags_name ON content_tags(name);
CREATE INDEX idx_content_tags_user_id ON content_tags(user_id);

DROP TRIGGER IF EXISTS update_content_tags_updated_at ON content_tags;
CREATE TRIGGER update_content_tags_updated_at
    BEFORE UPDATE ON content_tags
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Content tags relationship
CREATE TABLE content_tags_mapping (
    content_id uuid REFERENCES content(id) ON DELETE CASCADE,
    tag_id uuid REFERENCES content_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (content_id, tag_id)
);

DROP TRIGGER IF EXISTS update_content_tags_mapping_updated_at ON content_tags_mapping;
CREATE TRIGGER update_content_tags_mapping_updated_at
    BEFORE UPDATE ON content_tags_mapping
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Folders for organization
CREATE TABLE content_folders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(uuid),
    name VARCHAR(100) NOT NULL,
    parent_id uuid REFERENCES content_folders(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_content_folders_user_id ON content_folders(user_id);

DROP TRIGGER IF EXISTS update_content_folders_updated_at ON content_folders;
CREATE TRIGGER update_content_folders_updated_at
    BEFORE UPDATE ON content_folders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Content folder relationship
CREATE TABLE content_folders_mapping (
    content_id uuid REFERENCES content(id) ON DELETE CASCADE,
    folder_id uuid REFERENCES content_folders(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (content_id, folder_id)
);

DROP TRIGGER IF EXISTS update_content_folders_mapping_updated_at ON content_folders_mapping;
CREATE TRIGGER update_content_folders_mapping_updated_at
    BEFORE UPDATE ON content_folders_mapping
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

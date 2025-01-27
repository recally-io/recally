CREATE TABLE bookmark_content(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL, -- type of bookmark_content(bookmark, pdf, epub, image, podcast, video, etc.)
    url TEXT NOT NULL, -- URL of the bookmark
    user_id uuid DEFAULT NULL, -- when user is null, this content is shared, when it's not null, it will belong to the user
    title TEXT,
    description TEXT,
    domain TEXT, -- domain of the URL
    s3_key TEXT, -- S3 key for storing raw content like pdf, epub, video, etc.
    summary TEXT, -- AI generated summary
    content TEXT, -- content in markdown format
    html TEXT, -- html content for web page
    tags VARCHAR(50)[] DEFAULT '{}', -- tags for the content by default empty array
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(url, user_id)
  );

-- Indexes
CREATE INDEX idx_bookmark_content_type ON bookmark_content(type);
CREATE INDEX idx_bookmark_content_url ON bookmark_content(url);
CREATE INDEX idx_bookmark_content_domain ON bookmark_content(domain);
CREATE INDEX idx_bookmark_content_created_at ON bookmark_content(created_at);
CREATE INDEX idx_bookmark_content_metadata ON bookmark_content USING gin(metadata jsonb_path_ops);

-- BM25 index on bookmark_content
-- https://docs.paradedb.com/documentation/indexing/create_index
CREATE INDEX idx_bookmark_content_bm25_search ON bookmark_content USING bm25(id, title, description, summary, content, metadata) WITH (key_field = 'id');
CREATE TRIGGER update_bookmark_content_updated_at BEFORE UPDATE ON bookmark_content FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- new bookmarks table
DROP TABLE IF EXISTS bookmarks;
CREATE TABLE bookmarks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(uuid) ON DELETE CASCADE,
    content_id UUID REFERENCES bookmark_content(id),
    is_favorite BOOLEAN DEFAULT FALSE,
    is_archive BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT FALSE,
    reading_progress INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
  );

CREATE INDEX idx_bookmarks_user_created_at ON bookmarks(user_id, created_at);
CREATE INDEX idx_bookmarks_favorite ON bookmarks(user_id, is_favorite);
CREATE INDEX idx_bookmarks_archive ON bookmarks(user_id, is_archive);
CREATE INDEX idx_bookmarks_public ON bookmarks(user_id, is_public);
CREATE INDEX idx_bookmarks_metadata ON bookmarks USING gin(metadata jsonb_path_ops);
CREATE TRIGGER update_bookmarks_updated_at BEFORE UPDATE ON bookmarks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- bookmark_tags table
CREATE TABLE bookmark_tags (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    user_id uuid NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
  );

CREATE INDEX idx_bookmark_tags_name ON bookmark_tags(name);
CREATE TRIGGER update_bookmark_tags_updated_at BEFORE UPDATE ON bookmark_tags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- -- bookmark_content tags relationship
CREATE TABLE bookmark_tags_mapping(
    bookmark_id uuid REFERENCES bookmarks(id) ON DELETE CASCADE,
    tag_id uuid REFERENCES bookmark_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    PRIMARY KEY(bookmark_id, tag_id)
  );

CREATE INDEX idx_bookmark_tags_mapping_tag_id ON bookmark_tags_mapping(tag_id);
CREATE TRIGGER update_bookmark_tags_mapping_updated_at BEFORE UPDATE ON bookmark_tags_mapping FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

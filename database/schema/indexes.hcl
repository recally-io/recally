// Database Extensions and Special Indexes
// This file defines PostgreSQL extensions and ParadeDB-specific BM25 indexes

schema "public" {
}

// Vector extension for pgvector support
extension "vector" {
  schema  = schema.public
  version = null
  comment = "vector data type and ivfflat and hnsw access methods"
}

// BM25 Indexes for Full-Text Search using ParadeDB
// These indexes use ParadeDB's custom BM25 index type which Atlas may not recognize natively
// Therefore they are defined in SQL blocks

// BM25 index for legacy content table
sql {
  up = <<-SQL
    CREATE INDEX IF NOT EXISTS idx_content_bm25_search ON content
    USING bm25 (id, title, description, summary, content, metadata)
    WITH (key_field='id');
  SQL

  down = <<-SQL
    DROP INDEX IF EXISTS idx_content_bm25_search;
  SQL
}

// BM25 index for bookmark_content table
sql {
  up = <<-SQL
    CREATE INDEX IF NOT EXISTS idx_bookmark_content_bm25_search ON bookmark_content
    USING bm25(id, title, description, summary, content, metadata)
    WITH (key_field = 'id');
  SQL

  down = <<-SQL
    DROP INDEX IF EXISTS idx_bookmark_content_bm25_search;
  SQL
}

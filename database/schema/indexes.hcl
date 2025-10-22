// Database Extensions and Special Indexes
// This file defines PostgreSQL extensions and ParadeDB-specific BM25 indexes


// Vector extension for pgvector support
// NOTE: ParadeDB includes vector extension pre-installed
// Extension management is handled by ParadeDB infrastructure
// The vector type is available for use in column definitions

// BM25 Indexes for Full-Text Search using ParadeDB
// NOTE: Atlas does not support sql blocks in schema files (only in migration files)
// These BM25 indexes will be added manually to the generated migration SQL:
//
// CREATE INDEX IF NOT EXISTS idx_content_bm25_search ON content
// USING bm25 (id, title, description, summary, content, metadata)
// WITH (key_field='id');
//
// CREATE INDEX IF NOT EXISTS idx_bookmark_content_bm25_search ON bookmark_content
// USING bm25(id, title, description, summary, content, metadata)
// WITH (key_field = 'id');

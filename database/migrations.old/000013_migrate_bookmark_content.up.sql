INSERT INTO bookmark_content(id, type, url, user_id, title, description, domain, s3_key, summary, content, html, created_at, updated_at) 
SELECT id, 'bookmark', url, user_id, title, description, domain, s3_key, summary, content, html, created_at, updated_at FROM content;


INSERT INTO bookmarks(id, user_id, content_id, created_at, updated_at)
SELECT id, user_id, id, created_at, updated_at FROM content;

INSERT INTO bookmark_tags(id, name, user_id, created_at, updated_at)
SELECT id, name, user_id, created_at, updated_at FROM content_tags;

INSERT INTO bookmark_tags_mapping(bookmark_id, tag_id, created_at, updated_at)
SELECT content_id, tag_id, created_at, updated_at FROM content_tags_mapping;

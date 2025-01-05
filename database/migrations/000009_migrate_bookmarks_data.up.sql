INSERT INTO content (
    user_id,
    type,
    title,
    url,
    domain,
    summary,
    content,
    html,
    metadata,
    created_at,
    updated_at
)
SELECT 
    user_id,
    'bookmark' as type,
    COALESCE(title, '') as title,
    url,
    -- Extract domain from URL using regex
    REGEXP_REPLACE(
        REGEXP_REPLACE(url, '^https?://(?:www\.)?([^/]+).*', '\1'),
        ':[0-9]+$', ''
    ) as domain,
    summary,
    content,
    html,
    metadata,
    created_at,
    updated_at
FROM bookmarks;

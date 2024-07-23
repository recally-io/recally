-- name: InsertTextEmbedding :exec
INSERT INTO text_embeddings (user_id, text, embeddings, metadata)
VALUES ($1, $2, $3, $4);

-- name: GetTextEmbeddingById :one
SELECT id, metadata, user_id, text, created_at, updated_at, embeddings
FROM text_embeddings
WHERE id = $1;

-- name: DeleteTextEmbeddingById :exec
DELETE FROM text_embeddings
WHERE id = $1;

-- name: SimilaritySearchByL2Distance :many
SELECT id, metadata, user_id, text, created_at, updated_at, embeddings <-> $2 AS score 
FROM text_embeddings
WHERE user_id = $1 AND embeddings <-> $2
ORDER BY embeddings <-> $2  LIMIT $3;

-- name: SimilaritySearchByL2DistanceWithFilter :many
SELECT id, metadata, user_id, text, created_at, updated_at, embeddings <-> $2 AS score 
FROM text_embeddings
WHERE user_id = $1 AND embeddings <-> $2
    AND metadata @> $3::jsonb
ORDER BY embeddings <-> $2  LIMIT $4;

-- name: SimilaritySearchByCosineDistance :many
SELECT id, metadata, user_id, text, created_at, updated_at, 1 - (embeddings <=> $2) AS score 
FROM text_embeddings 
WHERE user_id = $1 AND embeddings <=> $2
ORDER BY 1 - (embedding <=> $2) LIMIT $3;

-- name: SimilaritySearchByCosineDistanceWithFilter :many
SELECT id, metadata, user_id, text, created_at, updated_at, 1 - (embeddings <=> $2) AS score 
FROM text_embeddings 
WHERE user_id = $1 AND embeddings <=> $2
    AND metadata @> $3::jsonb
ORDER BY 1 - (embedding <=> $2) LIMIT $4;

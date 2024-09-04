-- CRUD for assistant_message_embedddings
-- name: CreateAssistantEmbedding :exec
INSERT INTO assistant_embedddings (uuid, user_id, attachment_id, text, embeddings, metadata)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: IsAssistantEmbeddingExists :one
SELECT EXISTS (
    SELECT 1
    FROM assistant_embedddings
    WHERE uuid = $1
) AS exists;

-- name: DeleteAssistantEmbeddings :exec
DELETE FROM assistant_embedddings WHERE id = $1;

-- name: DeleteAssistantEmbeddingsByAttachmentId :exec
DELETE FROM assistant_embedddings WHERE attachment_id = $1;

-- name: DeleteAssistantEmbeddingsByThreadId :exec
DELETE FROM assistant_embedddings em
USING assistant_attachments aa
WHERE aa.uuid = em.attachment_id AND aa.thread_id = $1;

-- name: DeleteAssistantEmbeddingsByAssistantId :exec
-- This is a bit tricky, because we need to delete all embeddings for all threads of the assistant
DELETE FROM assistant_embedddings
USING assistant_attachments aa
WHERE aa.uuid = em.attachment_id AND aa.assistant_id = $1;

-- name: SimilaritySearchByThreadId :many
SELECT em.*
FROM assistant_embedddings em
JOIN assistant_attachments att ON em.attachment_id = att.uuid
JOIN assistant_threads th ON (th.uuid = att.thread_id OR th.assistant_id = att.assistant_id)
WHERE th.uuid = $1 AND (em.embeddings <=> $2 < 0.5)
ORDER BY 1 - (em.embeddings <=> $2) DESC
LIMIT $3;

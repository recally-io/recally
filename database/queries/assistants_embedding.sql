-- CRUD for assistant_message_embedddings
-- name: CreateAssistantEmbedding :exec
INSERT INTO assistant_embedddings (user_id, attachment_id, text, embeddings, metadata)
VALUES ($1, $2, $3, $4, $5);

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
SELECT em.id, em.text, em.metadata, 1 - (em.embeddings <=> $2) AS score
FROM assistant_embedddings em
JOIN assistant_attachments att ON em.attachment_id = att.uuid
JOIN assistant_threads th ON (th.uuid = att.thread_id OR th.assistant_id = att.assistant_id)
WHERE th.uuid = $1
    AND em.embeddings <=> $2
ORDER BY score DESC
LIMIT $3;

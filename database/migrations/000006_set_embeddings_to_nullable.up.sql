ALTER TABLE assistant_embedddings ALTER COLUMN embeddings TYPE vector(1536), ALTER COLUMN embeddings DROP NOT NULL;

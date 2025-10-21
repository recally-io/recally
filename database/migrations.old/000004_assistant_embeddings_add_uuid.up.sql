ALTER TABLE assistant_embedddings ADD COLUMN IF NOT EXISTS uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid();

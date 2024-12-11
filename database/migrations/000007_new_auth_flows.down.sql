-- Drop tables and their associated triggers
DROP TABLE IF EXISTS auth_revoked_tokens;
DROP TABLE IF EXISTS auth_api_keys;
DROP TABLE IF EXISTS auth_user_oauth_connections;

-- Drop triggers and functions
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_oauth_connections_updated_at ON auth_user_oauth_connections;
DROP TRIGGER IF EXISTS update_api_keys_updated_at ON auth_api_keys;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes on users table
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_status;

-- Remove constraints from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_check;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_contact_check;

-- Remove columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users DROP COLUMN IF EXISTS settings;

-- Restore original columns (these were dropped in the up migration)
ALTER TABLE users ADD COLUMN IF NOT EXISTS github VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS google VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS telegram VARCHAR(255);
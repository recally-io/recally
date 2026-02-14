-- Fix users_email_check constraint that used double-backslash escaping,
-- causing the regex to require a literal backslash in email addresses.
-- Use [.] character class instead of \. to avoid escaping issues.
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_email_check";
ALTER TABLE "users" ADD CONSTRAINT "users_email_check" CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+[.][A-Za-z]{2,}$');

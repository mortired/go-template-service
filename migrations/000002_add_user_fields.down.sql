BEGIN;

-- Rollback additional fields addition
ALTER TABLE users.users 
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS is_active,
DROP COLUMN IF EXISTS last_login_at;

-- Remove index
DROP INDEX IF EXISTS idx_users_is_active;

COMMIT;

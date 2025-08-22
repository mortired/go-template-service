BEGIN;

-- Rollback users table and schema creation
DROP TABLE IF EXISTS users.users CASCADE;
DROP SCHEMA IF EXISTS users CASCADE;

COMMIT;

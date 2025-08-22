BEGIN;

-- Rollback logs schema creation
DROP TABLE IF EXISTS logs.user_actions CASCADE;
DROP SCHEMA IF EXISTS logs CASCADE;

COMMIT;

BEGIN;

-- Rollback authentication schema creation
DROP TABLE IF EXISTS auth.sessions CASCADE;
DROP SCHEMA IF EXISTS auth CASCADE;

COMMIT;

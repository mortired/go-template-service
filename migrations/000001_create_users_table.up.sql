BEGIN;

-- Create schema for users
CREATE SCHEMA IF NOT EXISTS users;

-- Create users table in users schema
CREATE TABLE IF NOT EXISTS users.users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for fast search
CREATE INDEX IF NOT EXISTS idx_users_email ON users.users(email);

-- Create index on name for search
CREATE INDEX IF NOT EXISTS idx_users_name ON users.users(name);

-- Table comments
COMMENT ON SCHEMA users IS 'Schema for user tables';
COMMENT ON TABLE users.users IS 'System users table';
COMMENT ON COLUMN users.users.id IS 'Unique user identifier';
COMMENT ON COLUMN users.users.name IS 'User name';
COMMENT ON COLUMN users.users.email IS 'User email (unique)';
COMMENT ON COLUMN users.users.created_at IS 'Record creation date';
COMMENT ON COLUMN users.users.updated_at IS 'Last record update date';

COMMIT;

BEGIN;

-- Create schema for authentication
CREATE SCHEMA IF NOT EXISTS auth;

-- Create user sessions table
CREATE TABLE IF NOT EXISTS auth.sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users.users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for sessions table
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON auth.sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON auth.sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON auth.sessions(expires_at);

-- Schema and table comments
COMMENT ON SCHEMA auth IS 'Schema for authentication and authorization tables';
COMMENT ON TABLE auth.sessions IS 'User sessions table';
COMMENT ON COLUMN auth.sessions.id IS 'Unique session identifier';
COMMENT ON COLUMN auth.sessions.user_id IS 'User reference';
COMMENT ON COLUMN auth.sessions.token IS 'Session token (unique)';
COMMENT ON COLUMN auth.sessions.expires_at IS 'Session expiration date';
COMMENT ON COLUMN auth.sessions.created_at IS 'Session creation date';
COMMENT ON COLUMN auth.sessions.updated_at IS 'Last session update date';

COMMIT;

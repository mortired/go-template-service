BEGIN;

-- Create schema for logs
CREATE SCHEMA IF NOT EXISTS logs;

-- Create user actions log table
CREATE TABLE IF NOT EXISTS logs.user_actions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users.users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for logs table
CREATE INDEX IF NOT EXISTS idx_user_actions_user_id ON logs.user_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_action ON logs.user_actions(action);
CREATE INDEX IF NOT EXISTS idx_user_actions_created_at ON logs.user_actions(created_at);
CREATE INDEX IF NOT EXISTS idx_user_actions_details_gin ON logs.user_actions USING GIN (details);

-- Schema and table comments
COMMENT ON SCHEMA logs IS 'Schema for logging tables';
COMMENT ON TABLE logs.user_actions IS 'User actions log table';
COMMENT ON COLUMN logs.user_actions.id IS 'Unique log record identifier';
COMMENT ON COLUMN logs.user_actions.user_id IS 'User reference (can be NULL for anonymous actions)';
COMMENT ON COLUMN logs.user_actions.action IS 'Action type (e.g.: login, logout, create_user)';
COMMENT ON COLUMN logs.user_actions.details IS 'Additional action details in JSON format';
COMMENT ON COLUMN logs.user_actions.ip_address IS 'User IP address';
COMMENT ON COLUMN logs.user_actions.user_agent IS 'Browser User-Agent';
COMMENT ON COLUMN logs.user_actions.created_at IS 'Action date and time';

COMMIT;

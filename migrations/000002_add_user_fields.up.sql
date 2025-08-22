BEGIN;

-- Add additional fields to users table
ALTER TABLE users.users 
ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP WITH TIME ZONE;

-- Create index on activity status
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users.users(is_active);

-- Comments for new fields
COMMENT ON COLUMN users.users.phone IS 'User phone number';
COMMENT ON COLUMN users.users.is_active IS 'User activity status';
COMMENT ON COLUMN users.users.last_login_at IS 'Last system login date';

COMMIT;

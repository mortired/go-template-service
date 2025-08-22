# 📊 Migrations

## 📋 Overview

The migration system allows managing database schema in a versioned way, ensuring safe database structure updates across different environments.

## 🚀 Usage

### CLI Tool

```bash
# Apply all migrations
go run ./cmd/migrate -command=up

# Show migration status
go run ./cmd/migrate -command=status

# Rollback all migrations
go run ./cmd/migrate -command=down

# Apply specific migration
go run ./cmd/migrate -command=up -version=1

# Rollback specific migration
go run ./cmd/migrate -command=down -version=1
```

### Configuration

```env
# Database URL for migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/users?sslmode=disable

# Or separate parameters
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=users
```

## 📁 Migration Structure

```
migrations/
├── 000001_create_users_table.up.sql    # Create users table
├── 000001_create_users_table.down.sql  # Rollback users table creation
├── 000002_add_user_fields.up.sql       # Add user fields
├── 000002_add_user_fields.down.sql     # Rollback field addition
├── 000003_create_auth_schema.up.sql    # Create authentication schema
├── 000003_create_auth_schema.down.sql  # Rollback authentication schema
├── 000004_create_logs_schema.up.sql    # Create logging schema
└── 000004_create_logs_schema.down.sql  # Rollback logging schema
```

## 📝 Creating New Migrations

### 1. Creating Migration Files

```bash
# Create new migration
touch migrations/000005_add_user_roles.up.sql
touch migrations/000005_add_user_roles.down.sql
```

### 2. Writing SQL

```sql
-- 000005_add_user_roles.up.sql
ALTER TABLE users.users 
ADD COLUMN role VARCHAR(20) DEFAULT 'user' NOT NULL;

CREATE INDEX idx_users_role ON users.users(role);
```

```sql
-- 000005_add_user_roles.down.sql
DROP INDEX IF EXISTS idx_users_role;
ALTER TABLE users.users DROP COLUMN IF EXISTS role;
```

## 🔧 Best Practices

### 1. File Naming
- Use version prefix: `000001_`, `000002_`
- Describe action in name: `create_users_table`
- Separate into `.up.sql` and `.down.sql`

### 2. Migration Content
- Each migration should be atomic
- Include indexes in migrations
- Add comments for complex operations

### 3. Migration Rollback
- Always create `.down.sql` files
- Test migration rollbacks
- Consider table dependencies

## 🔗 Related Documentation

- [Database Setup](database.md) - PostgreSQL and migrations
- [Development](development.md) - development guide

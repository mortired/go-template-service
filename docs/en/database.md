# 🗄️ Database

## 📋 Overview

The application uses **PostgreSQL** as the primary database with a migration system for schema management. The database is divided into schemas for better data organization.

## 🏗️ Database Architecture

### PostgreSQL Schemas

The project uses PostgreSQL schema separation:

- **`users`** - user management
- **`auth`** - authentication and authorization
- **`logs`** - action logging

### Table Structure

#### `users` Schema

```sql
-- Users table
CREATE TABLE users.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimization
CREATE INDEX idx_users_username ON users.users(username);
CREATE INDEX idx_users_email ON users.users(email);
CREATE INDEX idx_users_active ON users.users(is_active);
```

#### `auth` Schema

```sql
-- Sessions table
CREATE TABLE auth.sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users.users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Authentication indexes
CREATE INDEX idx_sessions_token ON auth.sessions(token);
CREATE INDEX idx_sessions_user_id ON auth.sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON auth.sessions(expires_at);
```

#### `logs` Schema

```sql
-- User actions log table
CREATE TABLE logs.user_actions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users.users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Logging indexes
CREATE INDEX idx_user_actions_user_id ON logs.user_actions(user_id);
CREATE INDEX idx_user_actions_action ON logs.user_actions(action);
CREATE INDEX idx_user_actions_created_at ON logs.user_actions(created_at);
```

## ⚙️ Configuration

### Environment Variables

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_NAME=users
```

### Configuration Structure

```go
type DatabaseConfig struct {
    Host     string // Database host
    Port     int    // Database port
    Username string // Database username
    Password string // Database password
    Database string // Database name
}
```

### Database Connection

```go
// internal/infrastructure/postgres/postgres.go
type DB struct {
    *sql.DB
}

func New(cfg Config) (*DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Connection check
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{DB: db}, nil
}
```

## 🚀 Migrations

### Overview

The application uses a migration system for database schema management. Migrations are stored in the `migrations/` folder and executed using a CLI tool.

### Migration Structure

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

### Migration CLI Tool

#### Installation and Usage

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

#### Migration Configuration

```go
// internal/infrastructure/migration/config.go
type Config struct {
    DatabaseURL string // Database connection URL
    MigrationsPath string // Path to migration files
    TableName string // Migration tracking table name
}
```

### Migration Examples

#### Creating Users Table

```sql
-- 000001_create_users_table.up.sql
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users.users(username);
CREATE INDEX idx_users_email ON users.users(email);
```

```sql
-- 000001_create_users_table.down.sql
DROP TABLE IF EXISTS users.users;
DROP SCHEMA IF EXISTS users;
```

#### Adding User Fields

```sql
-- 000002_add_user_fields.up.sql
ALTER TABLE users.users 
ADD COLUMN first_name VARCHAR(50),
ADD COLUMN last_name VARCHAR(50),
ADD COLUMN is_active BOOLEAN DEFAULT true;

CREATE INDEX idx_users_active ON users.users(is_active);
```

```sql
-- 000002_add_user_fields.down.sql
ALTER TABLE users.users 
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS is_active;

DROP INDEX IF EXISTS idx_users_active;
```

### Migration Best Practices

#### 1. File Naming
- Use version prefix: `000001_`, `000002_`
- Describe action in name: `create_users_table`
- Separate into `.up.sql` and `.down.sql`

#### 2. Migration Content
- Each migration should be atomic
- Include indexes in migrations
- Add comments for complex operations

#### 3. Migration Rollback
- Always create `.down.sql` files
- Test migration rollbacks
- Consider table dependencies

## 🔧 Repositories

### Repository Interface

```go
// internal/repository/repository.go
type User interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id int) (*model.User, error)
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, limit, offset int) ([]*model.User, error)
}
```

### Repository Implementation

```go
// internal/repository/user/repository.go
type Repository struct {
    db *postgres.DB
}

func New(db *postgres.DB) repository.User {
    return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users.users (username, email, password_hash, first_name, last_name)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRowContext(ctx, query,
        user.Username, user.Email, user.PasswordHash,
        user.FirstName, user.LastName,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id int) (*model.User, error) {
    query := `
        SELECT id, username, email, password_hash, first_name, last_name, 
               is_active, created_at, updated_at
        FROM users.users
        WHERE id = $1
    `
    
    user := &model.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash,
        &user.FirstName, &user.LastName, &user.IsActive,
        &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return user, nil
}
```

## 🧪 Testing

### Repository Unit Tests

```go
func TestUserRepository_Create(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := user.New(db)
    user := &model.User{
        Username:     "testuser",
        Email:        "test@example.com",
        PasswordHash: "hashed_password",
    }
    
    // Act
    err := repo.Create(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.NotZero(t, user.CreatedAt)
}
```

### Integration Tests

```go
func TestDatabaseIntegration(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    
    // Act & Assert
    err := db.Ping()
    assert.NoError(t, err)
    
    // Test migrations
    err = runMigrations(db)
    assert.NoError(t, err)
}
```

## 🔐 Security

### Database Connection

- Use SSL in production
- Restrict database access by IP
- Use separate user for application
- Regularly update passwords

### SQL Injection

- Always use parameterized queries
- Avoid dynamic SQL construction
- Validate input data

### Secure Query Example

```go
// ✅ Secure - parameterized query
query := "SELECT * FROM users.users WHERE username = $1"
err := db.QueryRowContext(ctx, query, username).Scan(&user)

// ❌ Insecure - string concatenation
query := fmt.Sprintf("SELECT * FROM users.users WHERE username = '%s'", username)
err := db.QueryRowContext(ctx, query).Scan(&user)
```

## 📊 Monitoring

### Database Health Check

```go
// Database health check
func (db *DB) HealthCheck() error {
    return db.Ping()
}
```

### Query Logging

```go
// Log slow queries
func (db *DB) logSlowQuery(query string, duration time.Duration) {
    if duration > 100*time.Millisecond {
        log.Printf("Slow query detected: %s (duration: %v)", query, duration)
    }
}
```

## 🔗 Related Documentation

- [Architecture](architecture.md) - architecture principles
- [Configuration](configuration.md) - environment variables setup
- [API Reference](api.md) - API endpoints documentation
- [Development](development.md) - development guide

# 🗄️ База данных

## 📋 Обзор

Приложение использует **PostgreSQL** как основную базу данных с системой миграций для управления схемой. База данных разделена на схемы для лучшей организации данных.

## 🏗️ Архитектура базы данных

### Схемы PostgreSQL

Проект использует разделение таблиц по схемам PostgreSQL:

- **`users`** - управление пользователями
- **`auth`** - аутентификация и авторизация
- **`logs`** - логирование действий

### Структура таблиц

#### Схема `users`

```sql
-- Таблица пользователей
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

-- Индексы для оптимизации
CREATE INDEX idx_users_username ON users.users(username);
CREATE INDEX idx_users_email ON users.users(email);
CREATE INDEX idx_users_active ON users.users(is_active);
```

#### Схема `auth`

```sql
-- Таблица сессий
CREATE TABLE auth.sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users.users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для аутентификации
CREATE INDEX idx_sessions_token ON auth.sessions(token);
CREATE INDEX idx_sessions_user_id ON auth.sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON auth.sessions(expires_at);
```

#### Схема `logs`

```sql
-- Таблица логов действий
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

-- Индексы для логирования
CREATE INDEX idx_user_actions_user_id ON logs.user_actions(user_id);
CREATE INDEX idx_user_actions_action ON logs.user_actions(action);
CREATE INDEX idx_user_actions_created_at ON logs.user_actions(created_at);
```

## ⚙️ Конфигурация

### Переменные окружения

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_NAME=users
```

### Структура конфигурации

```go
type DatabaseConfig struct {
    Host     string // Хост базы данных
    Port     int    // Порт базы данных
    Username string // Имя пользователя БД
    Password string // Пароль БД
    Database string // Имя базы данных
}
```

### Подключение к базе данных

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
    
    // Проверка подключения
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{DB: db}, nil
}
```

## 🚀 Миграции

### Обзор

Приложение использует систему миграций для управления схемой базы данных. Миграции хранятся в папке `migrations/` и выполняются с помощью CLI инструмента.

### Структура миграций

```
migrations/
├── 000001_create_users_table.up.sql    # Создание таблицы пользователей
├── 000001_create_users_table.down.sql  # Откат создания таблицы пользователей
├── 000002_add_user_fields.up.sql       # Добавление полей пользователей
├── 000002_add_user_fields.down.sql     # Откат добавления полей
├── 000003_create_auth_schema.up.sql    # Создание схемы аутентификации
├── 000003_create_auth_schema.down.sql  # Откат схемы аутентификации
├── 000004_create_logs_schema.up.sql    # Создание схемы логирования
└── 000004_create_logs_schema.down.sql  # Откат схемы логирования
```

### CLI инструмент миграций

#### Установка и использование

```bash
# Применить все миграции
go run ./cmd/migrate -command=up

# Показать статус миграций
go run ./cmd/migrate -command=status

# Откатить все миграции
go run ./cmd/migrate -command=down

# Применить конкретную миграцию
go run ./cmd/migrate -command=up -version=1

# Откатить конкретную миграцию
go run ./cmd/migrate -command=down -version=1
```

#### Конфигурация миграций

```go
// internal/infrastructure/migration/config.go
type Config struct {
    DatabaseURL string // URL подключения к БД
    MigrationsPath string // Путь к файлам миграций
    TableName string // Имя таблицы для отслеживания миграций
}
```

### Примеры миграций

#### Создание таблицы пользователей

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

#### Добавление полей пользователей

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

### Лучшие практики миграций

#### 1. Именование файлов
- Используйте префикс с номером версии: `000001_`, `000002_`
- Описывайте действие в названии: `create_users_table`
- Разделяйте на `.up.sql` и `.down.sql`

#### 2. Содержимое миграций
- Каждая миграция должна быть атомарной
- Включайте индексы в миграции
- Добавляйте комментарии для сложных операций

#### 3. Откат миграций
- Всегда создавайте `.down.sql` файлы
- Тестируйте откат миграций
- Учитывайте зависимости между таблицами

## 🔧 Репозитории

### Интерфейс репозитория

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

### Реализация репозитория

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

## 🧪 Тестирование

### Unit тесты репозиториев

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

### Integration тесты

```go
func TestDatabaseIntegration(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    
    // Act & Assert
    err := db.Ping()
    assert.NoError(t, err)
    
    // Тестирование миграций
    err = runMigrations(db)
    assert.NoError(t, err)
}
```

## 🔐 Безопасность

### Подключение к базе данных

- Используйте SSL в продакшне
- Ограничьте доступ к базе данных по IP
- Используйте отдельного пользователя для приложения
- Регулярно обновляйте пароли

### SQL инъекции

- Всегда используйте параметризованные запросы
- Избегайте динамического построения SQL
- Валидируйте входные данные

### Пример безопасного запроса

```go
// ✅ Безопасно - параметризованный запрос
query := "SELECT * FROM users.users WHERE username = $1"
err := db.QueryRowContext(ctx, query, username).Scan(&user)

// ❌ Небезопасно - конкатенация строк
query := fmt.Sprintf("SELECT * FROM users.users WHERE username = '%s'", username)
err := db.QueryRowContext(ctx, query).Scan(&user)
```

## 📊 Мониторинг

### Проверка состояния базы данных

```go
// Health check для базы данных
func (db *DB) HealthCheck() error {
    return db.Ping()
}
```

### Логирование запросов

```go
// Логирование медленных запросов
func (db *DB) logSlowQuery(query string, duration time.Duration) {
    if duration > 100*time.Millisecond {
        log.Printf("Slow query detected: %s (duration: %v)", query, duration)
    }
}
```

## 🔗 Связанная документация

- [Архитектура](architecture.md) - принципы архитектуры
- [Конфигурация](configuration.md) - настройка переменных окружения
- [API Reference](api.md) - документация API endpoints
- [Разработка](development.md) - руководство по разработке

# 📊 Миграции

## 📋 Обзор

Система миграций позволяет управлять схемой базы данных версионированным способом, обеспечивая безопасное обновление структуры БД в разных окружениях.

## 🚀 Использование

### CLI инструмент

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

### Конфигурация

```env
# Database URL для миграций
DATABASE_URL=postgres://postgres:postgres@localhost:5432/users?sslmode=disable

# Или отдельные параметры
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=users
```

## 📁 Структура миграций

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

## 📝 Создание новых миграций

### 1. Создание файлов миграции

```bash
# Создать новую миграцию
touch migrations/000005_add_user_roles.up.sql
touch migrations/000005_add_user_roles.down.sql
```

### 2. Написание SQL

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

## 🔧 Лучшие практики

### 1. Именование файлов
- Используйте префикс с номером версии: `000001_`, `000002_`
- Описывайте действие в названии: `create_users_table`
- Разделяйте на `.up.sql` и `.down.sql`

### 2. Содержимое миграций
- Каждая миграция должна быть атомарной
- Включайте индексы в миграции
- Добавляйте комментарии для сложных операций

### 3. Откат миграций
- Всегда создавайте `.down.sql` файлы
- Тестируйте откат миграций
- Учитывайте зависимости между таблицами

## 🔗 Связанная документация

- [Настройка базы данных](database.md) - PostgreSQL и миграции
- [Разработка](development.md) - руководство по разработке

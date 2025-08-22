# ⚙️ Конфигурация приложения

## 📋 Обзор

Приложение использует многоуровневую систему конфигурации с поддержкой переменных окружения, файлов `.env` и автоматической загрузкой через dependency injection.

## 🏗️ Архитектура конфигурации

### Уровни конфигурации

1. **Переменные окружения** - основной источник настроек
2. **Файл .env** - локальная разработка (опционально)
3. **Структуры конфигурации** - типизированный доступ к настройкам
4. **Dependency Injection** - автоматическое внедрение через fx

### Приоритет загрузки

1. Системные переменные окружения (высший приоритет)
2. Файл `.env` в корне проекта
3. Значения по умолчанию (низший приоритет)

## 📋 Структура конфигурации

### ServerConfig
```go
type ServerConfig struct {
    Port int    // Порт сервера
    Host string // Хост сервера
}
```

**Переменные окружения:**
- `SERVER_PORT` - порт сервера (по умолчанию: 8080)
- `SERVER_HOST` - хост сервера (по умолчанию: localhost)

### DatabaseConfig
```go
type DatabaseConfig struct {
    Host     string // Хост базы данных
    Port     int    // Порт базы данных
    Username string // Имя пользователя БД
    Password string // Пароль БД
    Database string // Имя базы данных
}
```

**Переменные окружения:**
- `DB_HOST` - хост базы данных (по умолчанию: localhost)
- `DB_PORT` - порт базы данных (по умолчанию: 5432)
- `DB_USERNAME` - имя пользователя БД (по умолчанию: postgres)
- `DB_PASSWORD` - пароль БД (по умолчанию: пустая строка)
- `DB_NAME` - имя базы данных (по умолчанию: users)

### HMACConfig
```go
type HMACConfig struct {
    ClientSecrets []HMACClientSecret // Секреты клиентов
    RouteRights   HMACRouteRights    // Права доступа к маршрутам
    Algorithm     string             // Алгоритм хеширования
    MaxAge        int                // Максимальное время жизни подписи
    Required      bool               // Обязательная аутентификация
}
```

**Переменные окружения:**
- `HMAC_CLIENT_SECRETS` - JSON массив секретов клиентов
- `HMAC_ROUTE_RIGHTS` - JSON объект прав доступа к маршрутам
- `HMAC_ALGORITHM` - алгоритм хеширования (по умолчанию: sha256)
- `HMAC_MAX_AGE` - максимальное время жизни подписи в секундах (по умолчанию: 300)
- `HMAC_REQUIRED` - обязательная аутентификация (по умолчанию: true)

### LoggingConfig
```go
type LoggingConfig struct {
    Level     string // Уровень логирования
    FxDebug   bool   // Включить дебаг режим fx
    FxVerbose bool   // Включить подробное логирование fx
    ShowGraph bool   // Показать граф зависимостей fx
    FxEnabled bool   // Включить/отключить логирование fx полностью
}
```

**Переменные окружения:**
- `LOG_LEVEL` - уровень логирования (по умолчанию: info)
- `FX_DEBUG` - включить дебаг режим fx (по умолчанию: false)
- `FX_VERBOSE` - включить подробное логирование fx (по умолчанию: false)
- `FX_SHOW_GRAPH` - показать граф зависимостей fx (по умолчанию: false)
- `FX_ENABLED` - включить/отключить логирование fx (по умолчанию: true)

## 🚀 Использование

### 1. Настройка переменных окружения

Создайте файл `.env` в корне проекта:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_NAME=users

# HMAC Authentication
HMAC_CLIENT_SECRETS='[{"clientid":"web","secret":"your-secret-key","department":"web","descr":"web client"}]'
HMAC_ROUTE_RIGHTS='{"web":["/api/v1/users"]}'
HMAC_ALGORITHM=sha256
HMAC_MAX_AGE=300
HMAC_REQUIRED=true

# Logging Configuration
LOG_LEVEL=info
FX_DEBUG=false
FX_VERBOSE=false
FX_SHOW_GRAPH=false
FX_ENABLED=true

# App Debug Mode
SERVICE_DEBUG=true
```

### 2. Загрузка конфигурации

Конфигурация автоматически загружается через провайдер в `bootstrap.Setup()`:

```go
// ProvideConfig создает конфигурацию приложения
func ProvideConfig() *config.Config {
    return config.LoadConfig()
}
```

### 3. Использование в коде

```go
// В сервисе
type UserService struct {
    config *config.Config
    repo   repository.User
}

func NewUserService(config *config.Config, repo repository.User) *UserService {
    return &UserService{
        config: config,
        repo:   repo,
    }
}

// В контроллере
type UserController struct {
    config  *config.Config
    service service.User
}

func NewUserController(config *config.Config, service service.User) *UserController {
    return &UserController{
        config:  config,
        service: service,
    }
}
```

## 🔧 Утилиты конфигурации

### Загрузка .env файла

```go
import "users/internal/infrastructure/config"

// Загрузить .env файл
err := config.LoadEnv()
if err != nil {
    log.Printf("Warning: failed to load .env file: %v", err)
}
```

### Получение переменных окружения

```go
// Получить переменную с значением по умолчанию
host := config.GetEnvOrDefault("DB_HOST", "localhost")
port := config.GetEnvAsIntOrDefault("DB_PORT", 5432)
debug := config.GetEnvAsBoolOrDefault("SERVICE_DEBUG", false)
```

### Структурированная конфигурация

```go
// Загрузить полную конфигурацию
cfg := config.LoadConfig()

// Использовать конкретные секции
serverPort := cfg.Server.Port
dbHost := cfg.Database.Host
logLevel := cfg.Logging.Level
```

## 🌍 Окружения

### Development (Разработка)

```env
# Development Configuration
SERVICE_DEBUG=true
FX_DEBUG=true
FX_VERBOSE=true
LOG_LEVEL=debug

# Local Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=users_dev

# Local Server
SERVER_HOST=localhost
SERVER_PORT=8080
```

### Production (Продакшн)

```env
# Production Configuration
SERVICE_DEBUG=false
FX_DEBUG=false
FX_VERBOSE=false
LOG_LEVEL=info

# Production Database
DB_HOST=production-db.example.com
DB_PORT=5432
DB_USERNAME=app_user
DB_PASSWORD=secure_password
DB_NAME=users_prod

# Production Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### Testing (Тестирование)

```env
# Testing Configuration
SERVICE_DEBUG=true
FX_DEBUG=false
FX_VERBOSE=false
LOG_LEVEL=warn

# Test Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=users_test

# Test Server
SERVER_HOST=localhost
SERVER_PORT=8081
```

## 🔐 Безопасность

### Переменные окружения

- **Никогда не коммитьте** файлы `.env` в репозиторий
- Используйте `.env.example` для примеров конфигурации
- Храните секреты в системных переменных окружения
- Используйте менеджеры секретов в продакшне

### HMAC конфигурация

```env
# Пример безопасной HMAC конфигурации
HMAC_CLIENT_SECRETS='[{"clientid":"web","secret":"your-very-secure-secret-key","department":"web","descr":"web client"}]'
HMAC_ROUTE_RIGHTS='{"web":["/api/v1/users"]}'
HMAC_ALGORITHM=sha256
HMAC_MAX_AGE=300
HMAC_REQUIRED=true
```

### База данных

```env
# Безопасная конфигурация БД
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=app_user
DB_PASSWORD=secure_password_with_special_chars
DB_NAME=users
```

## 🐛 Отладка

### Включение режима отладки

```bash
# Включить режим отладки
export SERVICE_DEBUG=true

# Запустить приложение
go run ./cmd/app/main.go
```

### Проверка конфигурации

```go
// В коде для отладки
func debugConfig(cfg *config.Config) {
    fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
    fmt.Printf("Log Level: %s\n", cfg.Logging.Level)
    fmt.Printf("HMAC Required: %t\n", cfg.HMAC.Required)
}
```

### Логирование конфигурации

```go
// Логировать конфигурацию при запуске
log.Printf("Configuration loaded: Server=%s:%d, DB=%s:%d/%s",
    cfg.Server.Host, cfg.Server.Port,
    cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
```

## 🔗 Связанная документация

- [Архитектура](architecture.md) - принципы архитектуры
- [Настройка базы данных](database.md) - PostgreSQL конфигурация
- [API Reference](api.md) - документация API endpoints
- [Разработка](development.md) - руководство по разработке

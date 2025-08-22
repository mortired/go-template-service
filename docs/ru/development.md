# 🚀 Руководство по разработке

## 📋 Обзор

Это руководство поможет вам начать разработку с Users Service, включая настройку окружения, добавление новых функций и лучшие практики.

## 🛠️ Настройка окружения разработки

### Требования

- **Go 1.19+**
- **PostgreSQL 12+**
- **Git**
- **VS Code** (рекомендуется)

### Установка зависимостей

```bash
# Клонирование репозитория
git clone <repository-url>
cd users

# Установка Go зависимостей
go mod download

# Проверка установки
go version
```

### Настройка базы данных

```bash
# Создание базы данных
createdb users

# Применение миграций
go run ./cmd/migrate -command=up
```

### Конфигурация для разработки

Создайте файл `.env` в корне проекта:

```env
# Development Configuration
SERVICE_DEBUG=true
FX_DEBUG=true
FX_VERBOSE=true
LOG_LEVEL=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=users

# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# HMAC Authentication
HMAC_CLIENT_SECRETS='[{"client_id":"web","secret":"dev-secret-key"}]'
HMAC_ROUTE_RIGHTS='{"GET:/api/v1/users":["web"],"POST:/api/v1/users":["web"]}'
HMAC_ALGORITHM=sha256
HMAC_MAX_AGE=300
HMAC_REQUIRED=true
```

## 🔧 Запуск приложения

### Базовый запуск

```bash
# Запуск приложения
go run ./cmd/app/main.go
```

### Отладка в VS Code

1. Откройте проект в VS Code
2. Нажмите **F5** или используйте конфигурацию "Launch Users Service"
3. Приложение запустится и будет ждать сигнала завершения (**Ctrl+C**)

### Отладка с дополнительными опциями

```bash
# Включить режим отладки
SERVICE_DEBUG=true go run ./cmd/app/main.go

# Включить подробное логирование fx
FX_DEBUG=true FX_VERBOSE=true go run ./cmd/app/main.go

# Запуск с профилированием
go run -cpuprofile=cpu.prof ./cmd/app/main.go
```

## 🏗️ Структура проекта

### Добавление новых компонентов

#### 1. Модель данных

```go
// internal/model/user/user.go
package user

import "time"

type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"` // Не экспортируем пароль
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

#### 2. Репозиторий

```go
// internal/repository/user/repository.go
package user

import (
    "context"
    "users/internal/model/user"
    "users/internal/infrastructure/postgres"
)

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
```

#### 3. Сервис

```go
// internal/service/user/user.go
package user

import (
    "context"
    "users/internal/model/user"
    "users/internal/repository"
)

type Service struct {
    repo repository.User
}

func New(repo repository.User) service.User {
    return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {
    // Валидация
    if err := s.validateUser(user); err != nil {
        return err
    }
    
    // Хеширование пароля
    if err := s.hashPassword(user); err != nil {
        return err
    }
    
    // Сохранение в БД
    return s.repo.Create(ctx, user)
}
```

#### 4. Контроллер

```go
// internal/controller/user/user.go
package user

import (
    "net/http"
    "strconv"
    "users/internal/model/user"
    "users/internal/service"
    "github.com/labstack/echo/v4"
)

type Controller struct {
    service service.User
}

func New(service service.User) *Controller {
    return &Controller{service: service}
}

func (c *Controller) GetUsers(ctx echo.Context) error {
    // Получение параметров запроса
    limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
    offset, _ := strconv.Atoi(ctx.QueryParam("offset"))
    
    // Получение данных из сервиса
    users, err := c.service.GetUsers(ctx.Request().Context(), limit, offset)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusOK, users)
}
```

#### 5. Провайдер в Bootstrap

```go
// internal/bootstrap/repository.go
func ProvideUserRepository(pg *postgres.DB) repository.User {
    return userRepo.New(pg)
}

// internal/bootstrap/service.go
func ProvideUserService(repo repository.User) service.User {
    return user.New(repo)
}

// internal/bootstrap/controller.go
func ProvideUserController(service service.User) *user.Controller {
    return user.New(service)
}
```

#### 6. Регистрация в Setup

```go
// internal/bootstrap/setup.go
func Setup() *appcore.Application {
    options := []appcore.Option{
        // Infrastructure
        appcore.PostgresModule,
        appcore.HMACModule,

        // Repositories
        appcore.Provide(ProvideUserRepository),

        // Services
        appcore.Provide(ProvideUserService),

        // Controllers
        appcore.Provide(ProvideUserController),

        // Echo module
        appcore.EchoModule,

        // Router
        appcore.Invoke(router.SetupRoutes),

        // HTTP Server
        appcore.EchoServer,
    }

    return appcore.New(options...)
}
```

#### 7. Маршрутизация

```go
// internal/router/router.go
func SetupRoutes(e *echo.Echo, userController *user.Controller) {
    api := e.Group("/api/v1")
    
    users := api.Group("/users")
    users.GET("", userController.GetUsers)
    users.GET("/:id", userController.GetUser)
    users.POST("", userController.CreateUser)
    users.PUT("/:id", userController.UpdateUser)
    users.DELETE("/:id", userController.DeleteUser)
}
```

## 🧪 Тестирование

### Unit тесты

```go
// internal/service/user/user_test.go
package user

import (
    "context"
    "testing"
    "users/internal/model/user"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := New(mockRepo)
    
    user := &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Integration тесты

```go
// internal/service/user/user_integration_test.go
package user

import (
    "context"
    "testing"
    "users/internal/model/user"
    "users/internal/repository/user"
    "github.com/stretchr/testify/assert"
)

func TestUserServiceIntegration(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := user.New(db)
    service := New(repo)
    
    user := &model.User{
        Username: "integration_test",
        Email:    "integration@example.com",
        Password: "password123",
    }
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### Запуск тестов

```bash
# Запустить все тесты
go test ./...

# Запустить тесты с покрытием
go test -cover ./...

# Запустить тесты конкретного пакета
go test ./internal/service/user

# Запустить тесты с подробным выводом
go test -v ./internal/service/user
```

## 🔧 Инструменты разработки

### Линтинг и форматирование

```bash
# Форматирование кода
go fmt ./...

# Линтинг кода
golangci-lint run

# Проверка импортов
goimports -w .
```

### Генерация документации

```bash
# Генерация документации API
swag init -g cmd/app/main.go

# Просмотр документации
godoc -http=:6060
```

### Профилирование

```bash
# CPU профилирование
go run -cpuprofile=cpu.prof ./cmd/app/main.go
go tool pprof cpu.prof

# Memory профилирование
go run -memprofile=mem.prof ./cmd/app/main.go
go tool pprof mem.prof
```

## 📝 Лучшие практики

### 1. Код стиль

- Следуйте [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Используйте `gofmt` для форматирования
- Добавляйте комментарии к экспортируемым функциям
- Используйте осмысленные имена переменных и функций

### 2. Обработка ошибок

```go
// ✅ Хорошо - возврат ошибки с контекстом
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// ❌ Плохо - игнорирование ошибок
if err != nil {
    return err
}
```

### 3. Контекст

```go
// ✅ Хорошо - использование контекста
func (s *Service) GetUser(ctx context.Context, id int) (*model.User, error) {
    return s.repo.GetByID(ctx, id)
}

// ❌ Плохо - игнорирование контекста
func (s *Service) GetUser(id int) (*model.User, error) {
    return s.repo.GetByID(context.Background(), id)
}
```

### 4. Валидация

```go
// ✅ Хорошо - валидация входных данных
func (s *Service) CreateUser(ctx context.Context, user *model.User) error {
    if user.Username == "" {
        return errors.New("username is required")
    }
    
    if !isValidEmail(user.Email) {
        return errors.New("invalid email format")
    }
    
    return s.repo.Create(ctx, user)
}
```

### 5. Логирование

```go
// ✅ Хорошо - структурированное логирование
log.Printf("User created: id=%d, username=%s", user.ID, user.Username)

// ❌ Плохо - неструктурированное логирование
log.Printf("User %s created with ID %d", user.Username, user.ID)
```

## 🔄 Git Workflow

### Создание новой функции

```bash
# Создание новой ветки
git checkout -b feature/add-user-validation

# Внесение изменений
# ... редактирование файлов ...

# Добавление изменений
git add .

# Коммит изменений
git commit -m "feat: add user validation

- Add email format validation
- Add username uniqueness check
- Add password strength validation"

# Пуш изменений
git push origin feature/add-user-validation
```

### Создание Pull Request

1. Создайте Pull Request на GitHub
2. Добавьте описание изменений
3. Укажите связанные issues
4. Попросите code review
5. После одобрения выполните merge

## 🐛 Отладка

### Логирование

```go
// Добавление отладочных логов
log.Printf("DEBUG: Processing user creation for username: %s", user.Username)

// Логирование ошибок
if err != nil {
    log.Printf("ERROR: Failed to create user: %v", err)
    return err
}
```

### Отладка в VS Code

1. Установите Go extension для VS Code
2. Настройте точки останова (breakpoints)
3. Запустите отладку (F5)
4. Используйте Debug Console для проверки переменных

### Отладка HTTP запросов

```bash
# Тестирование API endpoints
curl -v http://localhost:8080/health

# Тестирование с HMAC аутентификацией
curl -H "Authorization: HMAC <signature>" \
     -v http://localhost:8080/api/v1/users
```

## 🔗 Связанная документация

- [Архитектура](architecture.md) - принципы архитектуры
- [Конфигурация](configuration.md) - настройка переменных окружения
- [Настройка базы данных](database.md) - PostgreSQL и миграции
- [API Reference](api.md) - документация API endpoints

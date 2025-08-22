# 🧪 Тестирование

## 📋 Обзор

Users Service использует комплексный подход к тестированию, включающий unit тесты, integration тесты и end-to-end тесты.

## 🎯 Стратегии тестирования

### 1. Unit тесты
- Тестирование отдельных функций и методов
- Использование моков для изоляции зависимостей
- Быстрое выполнение
- Высокое покрытие кода

### 2. Integration тесты
- Тестирование взаимодействия между компонентами
- Использование реальной базы данных
- Проверка HTTP endpoints
- Тестирование миграций

### 3. End-to-End тесты
- Тестирование полного пользовательского сценария
- Проверка работы всей системы
- Тестирование в реальных условиях

## 🚀 Запуск тестов

### Базовые команды

```bash
# Запустить все тесты
go test ./...

# Запустить тесты с покрытием
go test -cover ./...

# Запустить тесты конкретного пакета
go test ./internal/service/user

# Запустить тесты с подробным выводом
go test -v ./internal/service/user

# Запустить тесты с покрытием и генерацией отчета
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Тестирование с флагами

```bash
# Запустить только unit тесты
go test -tags=unit ./...

# Запустить только integration тесты
go test -tags=integration ./...

# Запустить тесты с таймаутом
go test -timeout=30s ./...

# Запустить тесты параллельно
go test -parallel=4 ./...
```

## 📝 Unit тесты

### Пример теста сервиса

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

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*model.User), args.Error(1)
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
    assert.NotZero(t, user.ID)
    mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_ValidationError(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := New(mockRepo)
    
    user := &model.User{
        Username: "", // Невалидный username
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "username is required")
    mockRepo.AssertNotCalled(t, "Create")
}
```

### Пример теста контроллера

```go
// internal/controller/user/user_test.go
package user

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "users/internal/model/user"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *model.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestUserController_CreateUser(t *testing.T) {
    // Arrange
    mockService := &MockUserService{}
    controller := New(mockService)
    
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/users", 
        strings.NewReader(`{"username":"test","email":"test@example.com","password":"password123"}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    
    mockService.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
    
    // Act
    err := controller.CreateUser(c)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, rec.Code)
    
    var response map[string]interface{}
    json.Unmarshal(rec.Body.Bytes(), &response)
    assert.Equal(t, "test", response["username"])
    mockService.AssertExpectations(t)
}
```

## 🔗 Integration тесты

### Настройка тестовой базы данных

```go
// internal/service/user/user_integration_test.go
package user

import (
    "context"
    "database/sql"
    "testing"
    "users/internal/model/user"
    "users/internal/repository/user"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
    // Подключение к тестовой базе данных
    db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/users_test?sslmode=disable")
    require.NoError(t, err)
    
    // Применение миграций
    err = runMigrations(db)
    require.NoError(t, err)
    
    return db
}

func TestUserServiceIntegration(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    defer db.Close()
    
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
    
    // Проверяем, что пользователь действительно создан
    createdUser, err := service.GetUser(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Username, createdUser.Username)
    assert.Equal(t, user.Email, createdUser.Email)
}
```

### Тестирование HTTP endpoints

```go
// internal/controller/user/user_integration_test.go
package user

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "users/internal/bootstrap"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestUserControllerIntegration(t *testing.T) {
    // Arrange
    app := bootstrap.Setup()
    defer app.Stop(context.Background())
    
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/api/v1/users", 
        strings.NewReader(`{"username":"test","email":"test@example.com","password":"password123"}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    
    // Act
    e.ServeHTTP(rec, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, rec.Code)
    
    var response map[string]interface{}
    json.Unmarshal(rec.Body.Bytes(), &response)
    assert.Equal(t, "test", response["username"])
}
```

## 🧹 Test Helpers

### Утилиты для тестирования

```go
// internal/testutil/helpers.go
package testutil

import (
    "database/sql"
    "testing"
    "users/internal/infrastructure/postgres"
)

// SetupTestDB создает тестовую базу данных
func SetupTestDB(t *testing.T) *postgres.DB {
    db, err := postgres.New(postgres.Config{
        Host:     "localhost",
        Port:     5432,
        Username: "postgres",
        Password: "postgres",
        Database: "users_test",
    })
    
    if err != nil {
        t.Fatalf("Failed to setup test database: %v", err)
    }
    
    return db
}

// CleanupTestDB очищает тестовую базу данных
func CleanupTestDB(t *testing.T, db *postgres.DB) {
    if err := db.Close(); err != nil {
        t.Errorf("Failed to close test database: %v", err)
    }
}

// CreateTestUser создает тестового пользователя
func CreateTestUser(t *testing.T, db *postgres.DB) *model.User {
    user := &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // Создание пользователя в БД
    // ...
    
    return user
}
```

## 📊 Покрытие кода

### Генерация отчетов о покрытии

```bash
# Генерация отчета о покрытии
go test -coverprofile=coverage.out ./...

# Просмотр отчета в браузере
go tool cover -html=coverage.out -o coverage.html

# Просмотр отчета в консоли
go tool cover -func=coverage.out
```

### Целевые показатели покрытия

- **Unit тесты**: 80%+
- **Integration тесты**: 70%+
- **Общее покрытие**: 75%+

## 🔧 Test Configuration

### Переменные окружения для тестов

```env
# Test Configuration
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USERNAME=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=users_test

# Test Logging
TEST_LOG_LEVEL=warn
```

### Конфигурация тестовой базы данных

```go
// internal/testutil/config.go
package testutil

type TestConfig struct {
    Database DatabaseConfig
    Logging  LoggingConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

func LoadTestConfig() *TestConfig {
    return &TestConfig{
        Database: DatabaseConfig{
            Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
            Port:     getEnvAsIntOrDefault("TEST_DB_PORT", 5432),
            Username: getEnvOrDefault("TEST_DB_USERNAME", "postgres"),
            Password: getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
            Database: getEnvOrDefault("TEST_DB_NAME", "users_test"),
        },
        Logging: LoggingConfig{
            Level: getEnvOrDefault("TEST_LOG_LEVEL", "warn"),
        },
    }
}
```

## 🚨 Best Practices

### 1. Именование тестов

```go
// ✅ Хорошо - описательные имена
func TestUserService_CreateUser_Success(t *testing.T) { ... }
func TestUserService_CreateUser_ValidationError(t *testing.T) { ... }
func TestUserService_CreateUser_DatabaseError(t *testing.T) { ... }

// ❌ Плохо - неописательные имена
func TestCreateUser(t *testing.T) { ... }
func TestCreateUser2(t *testing.T) { ... }
```

### 2. Структура тестов

```go
// ✅ Хорошо - AAA pattern (Arrange, Act, Assert)
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := New(mockRepo)
    user := &model.User{...}
    
    // Act
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### 3. Использование моков

```go
// ✅ Хорошо - настройка ожиданий
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
mockRepo.On("GetByID", mock.Anything, 1).Return(&model.User{...}, nil)

// ❌ Плохо - отсутствие ожиданий
mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
```

### 4. Очистка ресурсов

```go
// ✅ Хорошо - очистка после тестов
func TestUserService(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // тесты...
}
```

## 🔗 Связанная документация

- [Архитектура](architecture.md) - принципы архитектуры
- [Разработка](development.md) - руководство по разработке
- [Настройка базы данных](database.md) - PostgreSQL и миграции
- [API Reference](api.md) - документация API endpoints

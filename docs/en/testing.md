# 🧪 Testing

## 📋 Overview

Users Service uses a comprehensive testing approach including unit tests, integration tests, and end-to-end tests.

## 🎯 Testing Strategies

### 1. Unit Tests
- Testing individual functions and methods
- Using mocks to isolate dependencies
- Fast execution
- High code coverage

### 2. Integration Tests
- Testing interaction between components
- Using real database
- HTTP endpoint verification
- Migration testing

### 3. End-to-End Tests
- Testing complete user scenarios
- System-wide verification
- Real-world condition testing

## 🚀 Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/service/user

# Run tests with verbose output
go test -v ./internal/service/user

# Run tests with coverage report generation
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Testing with Flags

```bash
# Run only unit tests
go test -tags=unit ./...

# Run only integration tests
go test -tags=integration ./...

# Run tests with timeout
go test -timeout=30s ./...

# Run tests in parallel
go test -parallel=4 ./...
```

## 📝 Unit Tests

### Service Test Example

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
        Username: "", // Invalid username
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

### Controller Test Example

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

## 🔗 Integration Tests

### Test Database Setup

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
    // Connect to test database
    db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/users_test?sslmode=disable")
    require.NoError(t, err)
    
    // Apply migrations
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
    
    // Verify user was actually created
    createdUser, err := service.GetUser(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Username, createdUser.Username)
    assert.Equal(t, user.Email, createdUser.Email)
}
```

### HTTP Endpoint Testing

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

### Testing Utilities

```go
// internal/testutil/helpers.go
package testutil

import (
    "database/sql"
    "testing"
    "users/internal/infrastructure/postgres"
)

// SetupTestDB creates test database
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

// CleanupTestDB cleans up test database
func CleanupTestDB(t *testing.T, db *postgres.DB) {
    if err := db.Close(); err != nil {
        t.Errorf("Failed to close test database: %v", err)
    }
}

// CreateTestUser creates test user
func CreateTestUser(t *testing.T, db *postgres.DB) *model.User {
    user := &model.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // Create user in database
    // ...
    
    return user
}
```

## 📊 Code Coverage

### Coverage Report Generation

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View report in browser
go tool cover -html=coverage.out -o coverage.html

# View report in console
go tool cover -func=coverage.out
```

### Coverage Targets

- **Unit tests**: 80%+
- **Integration tests**: 70%+
- **Overall coverage**: 75%+

## 🔧 Test Configuration

### Environment Variables for Tests

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

### Test Database Configuration

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

### 1. Test Naming

```go
// ✅ Good - descriptive names
func TestUserService_CreateUser_Success(t *testing.T) { ... }
func TestUserService_CreateUser_ValidationError(t *testing.T) { ... }
func TestUserService_CreateUser_DatabaseError(t *testing.T) { ... }

// ❌ Bad - non-descriptive names
func TestCreateUser(t *testing.T) { ... }
func TestCreateUser2(t *testing.T) { ... }
```

### 2. Test Structure

```go
// ✅ Good - AAA pattern (Arrange, Act, Assert)
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

### 3. Using Mocks

```go
// ✅ Good - setting expectations
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
mockRepo.On("GetByID", mock.Anything, 1).Return(&model.User{...}, nil)

// ❌ Bad - no expectations
mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
```

### 4. Resource Cleanup

```go
// ✅ Good - cleanup after tests
func TestUserService(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // tests...
}
```

## 🔗 Related Documentation

- [Architecture](architecture.md) - architecture principles
- [Development](development.md) - development guide
- [Database Setup](database.md) - PostgreSQL and migrations
- [API Reference](api.md) - API endpoints documentation

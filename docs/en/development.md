# 🚀 Development Guide

## 📋 Overview

This guide will help you start developing with Users Service, including environment setup, adding new features, and best practices.

## 🛠️ Development Environment Setup

### Requirements

- **Go 1.19+**
- **PostgreSQL 12+**
- **Git**
- **VS Code** (recommended)

### Installing Dependencies

```bash
# Clone repository
git clone <repository-url>
cd users

# Install Go dependencies
go mod download

# Check installation
go version
```

### Database Setup

```bash
# Create database
createdb users

# Apply migrations
go run ./cmd/migrate -command=up
```

### Development Configuration

Create a `.env` file in the project root:

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

## 🔧 Running the Application

### Basic Run

```bash
# Run application
go run ./cmd/app/main.go
```

### Debugging in VS Code

1. Open project in VS Code
2. Press **F5** or use "Launch Users Service" configuration
3. Application will start and wait for termination signal (**Ctrl+C**)

### Debugging with Additional Options

```bash
# Enable debug mode
SERVICE_DEBUG=true go run ./cmd/app/main.go

# Enable verbose fx logging
FX_DEBUG=true FX_VERBOSE=true go run ./cmd/app/main.go

# Run with profiling
go run -cpuprofile=cpu.prof ./cmd/app/main.go
```

## 🏗️ Project Structure

### Adding New Components

#### 1. Data Model

```go
// internal/model/user/user.go
package user

import "time"

type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"` // Don't export password
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

#### 2. Repository

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

#### 3. Service

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
    // Validation
    if err := s.validateUser(user); err != nil {
        return err
    }
    
    // Password hashing
    if err := s.hashPassword(user); err != nil {
        return err
    }
    
    // Save to database
    return s.repo.Create(ctx, user)
}
```

#### 4. Controller

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
    // Get request parameters
    limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
    offset, _ := strconv.Atoi(ctx.QueryParam("offset"))
    
    // Get data from service
    users, err := c.service.GetUsers(ctx.Request().Context(), limit, offset)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusOK, users)
}
```

#### 5. Provider in Bootstrap

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

#### 6. Registration in Setup

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

#### 7. Routing

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

## 🧪 Testing

### Unit Tests

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

### Integration Tests

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

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/service/user

# Run tests with verbose output
go test -v ./internal/service/user
```

## 🔧 Development Tools

### Linting and Formatting

```bash
# Code formatting
go fmt ./...

# Code linting
golangci-lint run

# Import checking
goimports -w .
```

### Documentation Generation

```bash
# Generate API documentation
swag init -g cmd/app/main.go

# View documentation
godoc -http=:6060
```

### Profiling

```bash
# CPU profiling
go run -cpuprofile=cpu.prof ./cmd/app/main.go
go tool pprof cpu.prof

# Memory profiling
go run -memprofile=mem.prof ./cmd/app/main.go
go tool pprof mem.prof
```

## 📝 Best Practices

### 1. Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Add comments to exported functions
- Use meaningful variable and function names

### 2. Error Handling

```go
// ✅ Good - return error with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// ❌ Bad - ignoring errors
if err != nil {
    return err
}
```

### 3. Context

```go
// ✅ Good - using context
func (s *Service) GetUser(ctx context.Context, id int) (*model.User, error) {
    return s.repo.GetByID(ctx, id)
}

// ❌ Bad - ignoring context
func (s *Service) GetUser(id int) (*model.User, error) {
    return s.repo.GetByID(context.Background(), id)
}
```

### 4. Validation

```go
// ✅ Good - input validation
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

### 5. Logging

```go
// ✅ Good - structured logging
log.Printf("User created: id=%d, username=%s", user.ID, user.Username)

// ❌ Bad - unstructured logging
log.Printf("User %s created with ID %d", user.Username, user.ID)
```

## 🔄 Git Workflow

### Creating New Feature

```bash
# Create new branch
git checkout -b feature/add-user-validation

# Make changes
# ... edit files ...

# Add changes
git add .

# Commit changes
git commit -m "feat: add user validation

- Add email format validation
- Add username uniqueness check
- Add password strength validation"

# Push changes
git push origin feature/add-user-validation
```

### Creating Pull Request

1. Create Pull Request on GitHub
2. Add description of changes
3. Reference related issues
4. Request code review
5. After approval, perform merge

## 🐛 Debugging

### Logging

```go
// Adding debug logs
log.Printf("DEBUG: Processing user creation for username: %s", user.Username)

// Error logging
if err != nil {
    log.Printf("ERROR: Failed to create user: %v", err)
    return err
}
```

### Debugging in VS Code

1. Install Go extension for VS Code
2. Set breakpoints
3. Start debugging (F5)
4. Use Debug Console to check variables

### Debugging HTTP Requests

```bash
# Testing API endpoints
curl -v http://localhost:8080/health

# Testing with HMAC authentication
curl -H "Authorization: HMAC <signature>" \
     -v http://localhost:8080/api/v1/users
```

## 🔗 Related Documentation

- [Architecture](architecture.md) - architecture principles
- [Configuration](configuration.md) - environment variables setup
- [Database Setup](database.md) - PostgreSQL and migrations
- [API Reference](api.md) - API endpoints documentation

# 🏗️ Application Architecture

## 📋 Overview

Users Service is built on **Clean Architecture** principles using modern technology stack:

- **Go 1.19+** - primary development language
- **fx (Uber FX)** - dependency injection and lifecycle management
- **Echo** - HTTP framework
- **PostgreSQL** - database
- **HMAC** - API authentication

## 🏛️ Architecture Principles

### 1. Clean Architecture
The application is divided into layers with clear responsibility boundaries:

- **Presentation Layer** (Controllers) - HTTP request handling
- **Business Logic Layer** (Services) - business logic
- **Data Access Layer** (Repositories) - data access
- **Infrastructure Layer** - external dependencies (DB, logging, migrations)
- **Bootstrap Layer** - setup and composition of all layers

### 2. Dependency Inversion Principle
- High-level modules don't depend on low-level modules
- Both depend on abstractions
- Abstractions don't depend on details
- Dependencies point inward (towards business logic)

### 3. Allowed and Forbidden Imports

#### ✅ Allowed:
```
internal/controller/* → internal/service/*
internal/service/* → internal/repository/*
internal/repository/* → internal/model/*
internal/bootstrap/* → internal/* (all layers)
internal/infrastructure/appcore/* → internal/infrastructure/*
internal/router/* → internal/controller/*
```

#### ❌ Forbidden:
```
internal/infrastructure/* → internal/controller/*
internal/infrastructure/* → internal/service/*
internal/infrastructure/* → internal/repository/*
internal/repository/* → internal/service/*
internal/service/* → internal/controller/*
internal/model/* → internal/* (any other layers)
```

## 📁 Project Structure

```
users/
├── cmd/                   # Application entry points
│   ├── app/               # Main application
│   └── migrate/           # Migration tool
├── internal/              # Internal application code
│   ├── bootstrap/         # fx container setup
│   ├── controller/        # HTTP controllers
│   ├── infrastructure/    # Infrastructure components
│   │   ├── appcore/       # Application core (fx + echo)
│   │   ├── authhmac/      # HMAC authentication
│   │   ├── config/        # Configuration utilities
│   │   ├── hmac/          # HMAC utilities
│   │   ├── logging/       # Logging system
│   │   ├── middleware/    # HTTP middleware
│   │   ├── migration/     # Migration management
│   │   ├── postgres/      # PostgreSQL connection
│   │   └── response/      # HTTP response handling
│   ├── model/             # Data models
│   ├── provider/          # Dependency providers
│   ├── repository/        # Data access layer
│   ├── router/            # Routing
│   └── service/           # Business logic
├── migrations/            # Database SQL migrations
├── docs/                  # Documentation
├── deploy/                # Deployment files
├── env.example            # Environment variables example
├── go.mod                 # Go modules
└── go.sum                 # Go dependencies
```

## 🚀 AppCore - Application Core

### Overview
AppCore provides ready-to-use modules for common application components, encapsulating the complexity of dependency injection and configuration.

### Available Modules

#### PostgresModule
Provides PostgreSQL database connection:

```go
app := appcore.New(
    appcore.PostgresModule,
    // other options...
)
```

**Environment variables:**
- `DB_HOST` (default: "localhost")
- `DB_PORT` (default: 5432)
- `DB_USERNAME` (default: "postgres")
- `DB_PASSWORD` (default: "")
- `DB_NAME` (default: "users")

#### EchoModule
Provides HTTP server with Echo framework:

```go
app := appcore.New(
    appcore.EchoModule,
    appcore.EchoServer,
    // other options...
)
```

**Environment variables:**
- `SERVER_PORT` (default: 8080)
- `SERVER_HOST` (default: "localhost")

#### HMACModule
Provides HMAC authentication for API endpoints:

```go
app := appcore.New(
    appcore.HMACModule,
    // other options...
)
```

**Environment variables:**
- `HMAC_CLIENT_SECRETS` - JSON array of client secrets
- `HMAC_ROUTE_RIGHTS` - JSON object of route access rights
- `HMAC_ALGORITHM` (default: "sha256")
- `HMAC_MAX_AGE` (default: 300 seconds)
- `HMAC_REQUIRED` (default: true)

### Debug Mode
AppCore supports debug mode via `SERVICE_DEBUG` environment variable:

```bash
# Enable debug mode
export SERVICE_DEBUG=true
```

When debug mode is enabled:
- fx logging is enabled with debug and verbose options
- Additional debug information is output at startup
- "🔧 AppCore: Debug mode enabled" is printed to console

### Application Lifecycle

1. **Initialization** - creating fx.App with given options
2. **Configuration** - loading environment variables
3. **Dependency Injection** - setting up dependencies
4. **Startup** - starting HTTP server and other components
5. **Shutdown** - graceful shutdown on signal

## 🔧 Bootstrap - Application Setup

### Overview
Bootstrap contains all application initialization and providers for dependency injection.

### Structure
```
internal/bootstrap/
├── setup.go       # Main application setup
├── controller.go  # Controller providers
├── service.go     # Service providers
└── repository.go  # Repository providers
```

### Main Setup() Function
```go
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

### Providers

#### Repository Provider
```go
func ProvideUserRepository(pg *postgres.DB) repository.User {
    return userRepo.New(pg)
}
```

#### Service Provider
```go
func ProvideUserService(repo repository.User) service.User {
    return user.New(repo)
}
```

#### Controller Provider
```go
func ProvideUserController(service service.User) *user.Controller {
    return user.New(service)
}
```

### Dependency Lifecycle

1. **Infrastructure** - PostgreSQL, HMAC, logging
2. **Repositories** - data access
3. **Services** - business logic
4. **Controllers** - HTTP handlers
5. **Router** - routing
6. **HTTP Server** - server startup

## 🎯 Architecture Benefits

### 1. Separation of Concerns
- Each layer has clearly defined responsibilities
- Easy to test and maintain
- Simple to add new features

### 2. Framework Independence
- Business logic doesn't depend on external libraries
- Easy to replace technologies (e.g., Echo with Gin)
- Testing without external dependencies

### 3. Dependency Injection
- Automatic dependency management
- Easy testing with mocks
- Flexible configuration

### 4. Scalability
- Easy to add new modules
- Simple to split into microservices
- Horizontal scaling

## 🧪 Testing

### Unit Tests
- Testing business logic without external dependencies
- Using mocks for repositories
- Isolated service testing

### Integration Tests
- Testing with real database
- HTTP endpoint verification
- Migration testing

### Service Test Example
```go
func TestUserService_GetUsers(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // Act
    users, err := service.GetUsers()
    
    // Assert
    assert.NoError(t, err)
    assert.Len(t, users, 0)
}
```

## 🔗 Related Documentation

- [Configuration](configuration.md) - environment variables setup
- [Database Setup](database.md) - PostgreSQL and migrations
- [API Reference](api.md) - API endpoints documentation
- [Development](development.md) - development guide

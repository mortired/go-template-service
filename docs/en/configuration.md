# ⚙️ Application Configuration

## 📋 Overview

The application uses a multi-level configuration system with support for environment variables, `.env` files, and automatic loading through dependency injection.

## 🏗️ Configuration Architecture

### Configuration Levels

1. **Environment variables** - primary source of settings
2. **`.env` file** - local development (optional)
3. **Configuration structures** - typed access to settings
4. **Dependency Injection** - automatic injection through fx

### Loading Priority

1. System environment variables (highest priority)
2. `.env` file in project root
3. Default values (lowest priority)

## 📋 Configuration Structure

### ServerConfig
```go
type ServerConfig struct {
    Port int    // Server port
    Host string // Server host
}
```

**Environment variables:**
- `SERVER_PORT` - server port (default: 8080)
- `SERVER_HOST` - server host (default: localhost)

### DatabaseConfig
```go
type DatabaseConfig struct {
    Host     string // Database host
    Port     int    // Database port
    Username string // Database username
    Password string // Database password
    Database string // Database name
}
```

**Environment variables:**
- `DB_HOST` - database host (default: localhost)
- `DB_PORT` - database port (default: 5432)
- `DB_USERNAME` - database username (default: postgres)
- `DB_PASSWORD` - database password (default: empty string)
- `DB_NAME` - database name (default: users)

### HMACConfig
```go
type HMACConfig struct {
    ClientSecrets []HMACClientSecret // Client secrets
    RouteRights   HMACRouteRights    // Route access rights
    Algorithm     string             // Hashing algorithm
    MaxAge        int                // Maximum signature lifetime
    Required      bool               // Required authentication
}
```

**Environment variables:**
- `HMAC_CLIENT_SECRETS` - JSON array of client secrets
- `HMAC_ROUTE_RIGHTS` - JSON object of route access rights
- `HMAC_ALGORITHM` - hashing algorithm (default: sha256)
- `HMAC_MAX_AGE` - maximum signature lifetime in seconds (default: 300)
- `HMAC_REQUIRED` - required authentication (default: true)

### LoggingConfig
```go
type LoggingConfig struct {
    Level     string // Logging level
    FxDebug   bool   // Enable fx debug mode
    FxVerbose bool   // Enable fx verbose logging
    ShowGraph bool   // Show fx dependency graph
    FxEnabled bool   // Enable/disable fx logging completely
}
```

**Environment variables:**
- `LOG_LEVEL` - logging level (default: info)
- `FX_DEBUG` - enable fx debug mode (default: false)
- `FX_VERBOSE` - enable fx verbose logging (default: false)
- `FX_SHOW_GRAPH` - show fx dependency graph (default: false)
- `FX_ENABLED` - enable/disable fx logging (default: true)

## 🚀 Usage

### 1. Setting Environment Variables

Create a `.env` file in the project root:

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

### 2. Loading Configuration

Configuration is automatically loaded through the provider in `bootstrap.Setup()`:

```go
// ProvideConfig creates application configuration
func ProvideConfig() *config.Config {
    return config.LoadConfig()
}
```

### 3. Using in Code

```go
// In service
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

// In controller
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

## 🔧 Configuration Utilities

### Loading .env File

```go
import "users/internal/infrastructure/config"

// Load .env file
err := config.LoadEnv()
if err != nil {
    log.Printf("Warning: failed to load .env file: %v", err)
}
```

### Getting Environment Variables

```go
// Get variable with default value
host := config.GetEnvOrDefault("DB_HOST", "localhost")
port := config.GetEnvAsIntOrDefault("DB_PORT", 5432)
debug := config.GetEnvAsBoolOrDefault("SERVICE_DEBUG", false)
```

### Structured Configuration

```go
// Load full configuration
cfg := config.LoadConfig()

// Use specific sections
serverPort := cfg.Server.Port
dbHost := cfg.Database.Host
logLevel := cfg.Logging.Level
```

## 🌍 Environments

### Development

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

### Production

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

### Testing

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

## 🔐 Security

### Environment Variables

- **Never commit** `.env` files to repository
- Use `.env.example` for configuration examples
- Store secrets in system environment variables
- Use secret managers in production

### HMAC Configuration

```env
# Example of secure HMAC configuration
HMAC_CLIENT_SECRETS='[{"clientid":"web","secret":"your-very-secure-secret-key","department":"web","descr":"web client"}]'
HMAC_ROUTE_RIGHTS='{"web":["/api/v1/users"]}'
HMAC_ALGORITHM=sha256
HMAC_MAX_AGE=300
HMAC_REQUIRED=true
```

### Database

```env
# Secure database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=app_user
DB_PASSWORD=secure_password_with_special_chars
DB_NAME=users
```

## 🐛 Debugging

### Enabling Debug Mode

```bash
# Enable debug mode
export SERVICE_DEBUG=true

# Run application
go run ./cmd/app/main.go
```

### Checking Configuration

```go
// In code for debugging
func debugConfig(cfg *config.Config) {
    fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
    fmt.Printf("Log Level: %s\n", cfg.Logging.Level)
    fmt.Printf("HMAC Required: %t\n", cfg.HMAC.Required)
}
```

### Configuration Logging

```go
// Log configuration at startup
log.Printf("Configuration loaded: Server=%s:%d, DB=%s:%d/%s",
    cfg.Server.Host, cfg.Server.Port,
    cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
```

## 🔗 Related Documentation

- [Architecture](architecture.md) - architecture principles
- [Database Setup](database.md) - PostgreSQL configuration
- [API Reference](api.md) - API endpoints documentation
- [Development](development.md) - development guide

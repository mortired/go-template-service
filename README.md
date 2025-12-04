# 🚀 Users Service

A modern user management service built with Go, featuring Clean Architecture, Dependency Injection, and best practices.

## 📋 Description

This project demonstrates a Go application architecture using:
- **fx** - for dependency injection and lifecycle management
- **Echo** - for HTTP server
- **PostgreSQL** - for data storage
- **Clean Architecture** - for application layer separation
- **Database Migrations** - for database schema management
- **HMAC Authentication** - for secure API authentication

## 📚 Documentation

Complete documentation is available in the `docs/` folder:

- 🇷🇺 **[Russian Documentation](docs/ru/)** - Complete documentation in Russian
- 🇺🇸 **[English Documentation](docs/en/)** - Complete documentation in English

### Main documentation sections:

- [🏗️ Architecture](docs/en/architecture.md) - Architecture principles and project structure
- [⚙️ Configuration](docs/en/configuration.md) - Environment variables setup
- [🗄️ Database](docs/en/database.md) - PostgreSQL and migrations
- [🌐 API Reference](docs/en/api.md) - API endpoints documentation
- [🚀 Development](docs/en/development.md) - Development guide

## 🚀 Quick Start

### Requirements
- **Go 1.25+**
- **PostgreSQL 12+**

### Installation and Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mortired/go-template-service.git
   cd go-template-service
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup database**
   ```bash
   # Option 1: Using PostgreSQL CLI
   createdb users
   
   # Option 2: Using psql
   psql -c "CREATE DATABASE users;"
   
   # Option 3: Using pgAdmin or other GUI tools
   # Create database named "users" through your preferred PostgreSQL client
   ```

4. **Create configuration file**
   ```bash
   cp env.example .env
   # Edit .env file with your settings
   ```

5. **Apply migrations**
   ```bash
   go run ./cmd/migrate -command=up
   ```

6. **Run the application**
   ```bash
   go run ./cmd/app/main.go
   ```

## 🏗️ Project Structure

```
users/
├── cmd/                    # Application entry points
│   ├── app/               # Main application
│   └── migrate/           # Migration tool
├── internal/              # Internal application code
│   ├── bootstrap/         # fx container setup
│   ├── controller/        # HTTP controllers
│   ├── infrastructure/    # Infrastructure components
│   ├── model/             # Data models
│   ├── repository/        # Data access layer
│   ├── router/            # Routing
│   └── service/           # Business logic
├── migrations/            # Database SQL migrations
├── docs/                  # Documentation
└── deploy/                # Deployment files
```

## 🌐 API

### Endpoints

| Method | Path | Description | Authentication |
|--------|------|-------------|----------------|
| `GET` | `/api/v1/users` | Get list of users | HMAC |
| `POST` | `/api/v1/users` | Create new user | HMAC |

### Request Examples

```bash
# Get list of users (requires HMAC authentication)
curl -H "Authorization: HMAC <signature>" http://localhost:8080/api/v1/users

# Create new user (requires HMAC authentication)
curl -X POST \
  -H "Authorization: HMAC <signature>" \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}' \
  http://localhost:8080/api/v1/users
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 🔧 Development

### Debugging

```bash
# Enable debug mode
SERVICE_DEBUG=true go run ./cmd/app/main.go

# Enable verbose fx logging
FX_DEBUG=true FX_VERBOSE=true go run ./cmd/app/main.go
```

### Migrations

```bash
# Apply all migrations
go run ./cmd/migrate -command=up

# Show migration status
go run ./cmd/migrate -command=status

# Rollback all migrations
go run ./cmd/migrate -command=down
```

## 🚀 Deployment

### Docker

```bash
# Build image
docker build -t users-service .

# Run container
docker run -p 8080:8080 --env-file .env users-service
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests
5. Update documentation
6. Create a Pull Request

## 📞 Support

If you have questions or issues:

1. Check the [documentation](docs/)
2. Create an [Issue](https://github.com/mortired/go-template-service/issues)
3. Refer to [CHANGELOG.md](CHANGELOG.md) for information about changes

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.
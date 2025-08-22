# 📚 Users Service Documentation

Welcome to the Users Service documentation! This project demonstrates modern Go application architecture using Clean Architecture, Dependency Injection, and best development practices.

## 🌍 Documentation Languages

- 🇷🇺 [Русский](../ru/) - Complete documentation in Russian
- 🇺🇸 [English](.) - Complete documentation in English

## 🚀 Quick Start

### Requirements
- **Go 1.19+**
- **PostgreSQL 12+**
- **Environment variables** (see [Configuration](configuration.md))

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
   # Copy example configuration
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

## 📖 Documentation Structure

### 🇷🇺 Russian Documentation

#### 🏗️ Architecture
- [Принципы архитектуры](../ru/architecture.md) - Core principles and project structure
- [AppCore](../ru/appcore.md) - Application core usage
- [Bootstrap](../ru/bootstrap.md) - Setup and initialization

#### ⚙️ Configuration & Setup
- [Конфигурация](../ru/configuration.md) - Environment variables and settings
- [Настройка базы данных](../ru/database.md) - PostgreSQL and migrations
- [Миграции](../ru/migrations.md) - Database schema management

#### 🚀 Development
- [API Reference](../ru/api.md) - API endpoints documentation
- [Разработка](../ru/development.md) - Development guide
- [Тестирование](../ru/testing.md) - Testing strategies

### 🇺🇸 English Documentation

#### 🏗️ Architecture
- [Architecture Principles](architecture.md) - Core principles and project structure
- [AppCore](appcore.md) - Application core usage
- [Bootstrap](bootstrap.md) - Setup and initialization

#### ⚙️ Configuration & Setup
- [Configuration](configuration.md) - Environment variables and settings
- [Database Setup](database.md) - PostgreSQL and migrations
- [Migrations](migrations.md) - Database schema management

#### 🚀 Development
- [API Reference](api.md) - API endpoints documentation
- [Development](development.md) - Development guide
- [Testing](testing.md) - Testing strategies

## 🎯 Key Features

- **🏗️ Clean Architecture** - Clear separation of application layers
- **🔧 Dependency Injection** - Automatic dependency management
- **🔐 HMAC Authentication** - Secure API authentication
- **🗄️ PostgreSQL** - Reliable data storage
- **📊 Migrations** - Database schema management
- **📝 Logging** - Structured logging
- **🧪 Testing** - Unit and integration tests

## 🔗 Useful Links

- [GitHub Repository](https://github.com/mortired/go-template-service)
- [Issues](https://github.com/mortired/go-template-service/issues)
- [CHANGELOG](../../CHANGELOG.md)

## 📞 Support

If you have questions or issues:

1. Check the [documentation](.) or [документация](../ru/)
2. Create an [Issue](https://github.com/mortired/go-template-service/issues)
3. Refer to [CHANGELOG](../../CHANGELOG.md) for information about changes

---

**Version:** 1.0.0  
**Last Updated:** 2025  
**Go Version:** 1.19+

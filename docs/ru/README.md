# 📚 Документация Users Service

Добро пожаловать в документацию сервиса управления пользователями! Этот проект демонстрирует современную архитектуру Go приложения с использованием Clean Architecture, Dependency Injection и лучших практик разработки.

## 🌍 Языки документации

- 🇷🇺 [Русский](../ru/) - Полная документация на русском языке
- 🇺🇸 [English](../en/) - Complete documentation in English

## 🚀 Быстрый старт

### Требования
- **Go 1.19+**
- **PostgreSQL 12+**
- **Переменные окружения** (см. [Конфигурация](configuration.md))

### Установка и запуск

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/mortired/go-template-service.git
   cd go-template-service
   ```

2. **Установите зависимости**
   ```bash
   go mod tidy
   ```

3. **Настройте базу данных**
   ```bash
   # Вариант 1: Используя PostgreSQL CLI
   createdb users
   
   # Вариант 2: Используя psql
   psql -c "CREATE DATABASE users;"
   
   # Вариант 3: Используя pgAdmin или другие GUI инструменты
   # Создайте базу данных с именем "users" через ваш предпочитаемый PostgreSQL клиент
   ```

4. **Создайте файл конфигурации**  
   Скопируйте пример конфигурации и  
   Отредактируйте .env файл под ваши настройки
   ```bash
   cp env.example .env
   ```

5. **Примените миграции**
   ```bash
   go run ./cmd/migrate -command=up
   ```

6. **Запустите приложение**
   ```bash
   go run ./cmd/app/main.go
   ```

## 📖 Структура документации

### 🇷🇺 Русская документация

#### 🏗️ Архитектура
- [Принципы архитектуры](architecture.md) - Основные принципы и структура проекта
- [AppCore](architecture.md#🚀-appcore---ядро-приложения) - Использование ядра приложения
- [Bootstrap](architecture.md#🔧-bootstrap---настройка-приложения) - Настройка и инициализация

#### ⚙️ Конфигурация и настройка
- [Конфигурация](configuration.md) - Переменные окружения и настройки
- [Настройка базы данных](database.md) - PostgreSQL и миграции
- [Миграции](migrations.md) - Управление схемой базы данных

#### 🚀 Разработка
- [API Reference](api.md) - Документация API endpoints
- [Разработка](development.md) - Руководство по разработке
- [Тестирование](testing.md) - Стратегии тестирования

### 🇺🇸 English Documentation

#### 🏗️ Architecture
- [Architecture Principles](../en/architecture.md) - Core principles and project structure
- [AppCore](../en/architecture.md#🚀-appcore---application-core) - Application core usage
- [Bootstrap](../en/architecture.md#🔧-bootstrap---application-setup) - Setup and initialization

#### ⚙️ Configuration & Setup
- [Configuration](../en/configuration.md) - Environment variables and settings
- [Database Setup](../en/database.md) - PostgreSQL and migrations
- [Migrations](../en/migrations.md) - Database schema management

#### 🚀 Development
- [API Reference](../en/api.md) - API endpoints documentation
- [Development](../en/development.md) - Development guide
- [Testing](../en/testing.md) - Testing strategies

## 🎯 Основные возможности

- **🏗️ Clean Architecture** - Четкое разделение слоев приложения
- **🔧 Dependency Injection** - Автоматическое управление зависимостями
- **🔐 HMAC Authentication** - Безопасная аутентификация API
- **🗄️ PostgreSQL** - Надежное хранение данных
- **📊 Миграции** - Управление схемой базы данных
- **📝 Логирование** - Структурированное логирование
- **🧪 Тестирование** - Unit и integration тесты

## 🔗 Полезные ссылки

- [GitHub Repository](https://github.com/mortired/go-template-service)
- [Issues](https://github.com/mortired/go-template-service/issues)
- [CHANGELOG](../../CHANGELOG.md)

## 📞 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [документацию](../ru/) или [documentation](../en/)
2. Создайте [Issue](https://github.com/mortired/go-template-service/issues)
3. Обратитесь к [CHANGELOG](../../CHANGELOG.md) для информации об изменениях

---

**Версия:** 1.0.0  
**Последнее обновление:** 2025  
**Go версия:** 1.19+

# 🏗️ Архитектура приложения

## 📋 Обзор

Users Service построен на принципах **Clean Architecture** с использованием современного стека технологий:

- **Go 1.19+** - основной язык разработки
- **fx (Uber FX)** - dependency injection и управление жизненным циклом
- **Echo** - HTTP фреймворк
- **PostgreSQL** - база данных
- **HMAC** - аутентификация API

## 🏛️ Принципы архитектуры

### 1. Clean Architecture
Приложение разделено на слои с четкими границами ответственности:

- **Presentation Layer** (Контроллеры) - обработка HTTP запросов
- **Business Logic Layer** (Сервисы) - бизнес-логика
- **Data Access Layer** (Репозитории) - доступ к данным
- **Infrastructure Layer** - внешние зависимости (БД, логирование, миграции)
- **Bootstrap Layer** - настройка и композиция всех слоев

### 2. Dependency Inversion Principle
- Высокоуровневые модули не зависят от низкоуровневых
- Оба зависят от абстракций
- Абстракции не зависят от деталей
- Зависимости направлены внутрь (к бизнес-логике)

### 3. Разрешенные и запрещенные импорты

#### ✅ Разрешено:
```
internal/controller/* → internal/service/*
internal/service/* → internal/repository/*
internal/repository/* → internal/model/*
internal/bootstrap/* → internal/* (все слои)
internal/infrastructure/appcore/* → internal/infrastructure/*
internal/router/* → internal/controller/*
```

#### ❌ Запрещено:
```
internal/infrastructure/* → internal/controller/*
internal/infrastructure/* → internal/service/*
internal/infrastructure/* → internal/repository/*
internal/repository/* → internal/service/*
internal/service/* → internal/controller/*
internal/model/* → internal/* (любые другие слои)
```

## 📁 Структура проекта

```
users/
├── cmd/                    # Точки входа в приложение
│   ├── app/               # Основное приложение
│   └── migrate/           # Инструмент миграций
├── internal/              # Внутренний код приложения
│   ├── bootstrap/         # Настройка fx контейнера
│   ├── controller/        # HTTP контроллеры
│   ├── infrastructure/    # Инфраструктурные компоненты
│   │   ├── appcore/       # Ядро приложения (fx + echo)
│   │   ├── authhmac/      # HMAC аутентификация
│   │   ├── config/        # Утилиты конфигурации
│   │   ├── hmac/          # HMAC утилиты
│   │   ├── logging/       # Система логирования
│   │   ├── middleware/    # HTTP middleware
│   │   ├── migration/     # Управление миграциями
│   │   ├── postgres/      # Подключение к PostgreSQL
│   │   └── response/      # Обработка HTTP ответов
│   ├── model/             # Модели данных
│   ├── provider/          # Провайдеры зависимостей
│   ├── repository/        # Слой доступа к данным
│   ├── router/            # Маршрутизация
│   └── service/           # Бизнес-логика
├── migrations/            # SQL миграции базы данных
├── docs/                  # Документация
├── deploy/                # Файлы для развертывания
├── env.example            # Пример переменных окружения
├── go.mod                 # Go модули
└── go.sum                 # Go зависимости
```

## 🚀 AppCore - Ядро приложения

### Обзор
AppCore предоставляет готовые к использованию модули для общих компонентов приложения, инкапсулируя сложность dependency injection и конфигурации.

### Доступные модули

#### PostgresModule
Предоставляет подключение к базе данных PostgreSQL:

```go
app := appcore.New(
    appcore.PostgresModule,
    // другие опции...
)
```

**Переменные окружения:**
- `DB_HOST` (по умолчанию: "localhost")
- `DB_PORT` (по умолчанию: 5432)
- `DB_USERNAME` (по умолчанию: "postgres")
- `DB_PASSWORD` (по умолчанию: "")
- `DB_NAME` (по умолчанию: "users")

#### EchoModule
Предоставляет HTTP сервер с фреймворком Echo:

```go
app := appcore.New(
    appcore.EchoModule,
    appcore.EchoServer,
    // другие опции...
)
```

**Переменные окружения:**
- `SERVER_PORT` (по умолчанию: 8080)
- `SERVER_HOST` (по умолчанию: "localhost")

#### HMACModule
Предоставляет HMAC аутентификацию для API endpoints:

```go
app := appcore.New(
    appcore.HMACModule,
    // другие опции...
)
```

**Переменные окружения:**
- `HMAC_CLIENT_SECRETS` - JSON массив секретов клиентов
- `HMAC_ROUTE_RIGHTS` - JSON объект прав доступа к маршрутам
- `HMAC_ALGORITHM` (по умолчанию: "sha256")
- `HMAC_MAX_AGE` (по умолчанию: 300 секунд)
- `HMAC_REQUIRED` (по умолчанию: true)

### Режим отладки
AppCore поддерживает режим отладки через переменную окружения `SERVICE_DEBUG`:

```bash
# Включить режим отладки
export SERVICE_DEBUG=true
```

Когда режим отладки включен:
- Логирование fx включено с опциями debug и verbose
- Дополнительная отладочная информация выводится при запуске
- В консоль выводится "🔧 AppCore: Debug mode enabled"

### Жизненный цикл приложения

1. **Инициализация** - создание fx.App с заданными опциями
2. **Конфигурация** - загрузка переменных окружения
3. **Dependency Injection** - настройка зависимостей
4. **Запуск** - старт HTTP сервера и других компонентов
5. **Остановка** - graceful shutdown при получении сигнала

## 🔧 Bootstrap - Настройка приложения

### Обзор
Bootstrap содержит всю инициализацию приложения и провайдеры для dependency injection.

### Структура
```
internal/bootstrap/
├── setup.go       # Основная настройка приложения
├── controller.go  # Провайдеры для контроллеров
├── service.go     # Провайдеры для сервисов
└── repository.go  # Провайдеры для репозиториев
```

### Основная функция Setup()
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

### Провайдеры

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

### Жизненный цикл зависимостей

1. **Infrastructure** - PostgreSQL, HMAC, логирование
2. **Repositories** - доступ к данным
3. **Services** - бизнес-логика
4. **Controllers** - HTTP обработчики
5. **Router** - маршрутизация
6. **HTTP Server** - запуск сервера

## 🎯 Преимущества архитектуры

### 1. Разделение ответственности
- Каждый слой имеет четко определенную ответственность
- Легко тестировать и поддерживать
- Простое добавление новых функций

### 2. Независимость от фреймворков
- Бизнес-логика не зависит от внешних библиотек
- Легко заменить технологии (например, Echo на Gin)
- Тестирование без внешних зависимостей

### 3. Dependency Injection
- Автоматическое управление зависимостями
- Легкое тестирование с моками
- Гибкая конфигурация

### 4. Масштабируемость
- Легко добавлять новые модули
- Простое разделение на микросервисы
- Горизонтальное масштабирование

## 🧪 Тестирование

### Unit тесты
- Тестирование бизнес-логики без внешних зависимостей
- Использование моков для репозиториев
- Изолированное тестирование сервисов

### Integration тесты
- Тестирование с реальной базой данных
- Проверка HTTP endpoints
- Тестирование миграций

### Пример теста сервиса
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

## 🔗 Связанная документация

- [Конфигурация](configuration.md) - настройка переменных окружения
- [Настройка базы данных](database.md) - PostgreSQL и миграции
- [API Reference](api.md) - документация API endpoints
- [Разработка](development.md) - руководство по разработке

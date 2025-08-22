# 🌐 API Reference

## 📋 Overview

Users Service предоставляет RESTful API для управления пользователями с поддержкой HMAC аутентификации.

## 🔐 Authentication

### HMAC Authentication

API использует HMAC (Hash-based Message Authentication Code) для аутентификации запросов.

#### Как это работает

1. **Клиент** создает подпись запроса используя секретный ключ
2. **Сервер** проверяет подпись и аутентифицирует запрос
3. **Доступ** предоставляется только с правильной подписью

#### Создание подписи

```bash
# Пример создания HMAC подписи
echo -n "GET/api/v1/users" | openssl dgst -sha256 -hmac "your-secret-key"
```

#### Authorization header

```http
Authorization: HMAC <signature>
```

### HMAC Configuration

```env
# Секреты клиентов
HMAC_CLIENT_SECRETS='[{"clientid":"web","secret":"your-secret-key","department":"web","descr":"web client"}]'

# Права доступа к маршрутам
HMAC_ROUTE_RIGHTS='{"web":["/api/v1/users"]}'

# Настройки аутентификации
HMAC_ALGORITHM=sha256
HMAC_MAX_AGE=300
HMAC_REQUIRED=true
```

## 📡 Endpoints

### Base URL

```
http://localhost:8080
```

### Users API

#### GET /api/v1/users

Получение списка пользователей.

**Authentication:** HMAC

**Query parameters:**
- `name` (string, optional) - фильтр по имени пользователя (поиск по подстроке)
- `id` (int, optional) - фильтр по ID пользователя

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "email": "jane@example.com"
  }
]
```

**Example requests:**
```bash
# Получить всех пользователей
curl -H "Authorization: HMAC <signature>" http://localhost:8080/api/v1/users

# Фильтрация по имени
curl -H "Authorization: HMAC <signature>" "http://localhost:8080/api/v1/users?name=john"

# Фильтрация по ID
curl -H "Authorization: HMAC <signature>" "http://localhost:8080/api/v1/users?id=1"
```

#### POST /api/v1/users

Создание нового пользователя.

**Authentication:** HMAC

**Request body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Example request:**
```bash
curl -X POST \
  -H "Authorization: HMAC <signature>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }' \
  http://localhost:8080/api/v1/users
```

## 📊 Data Models

### User

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Fields:**
- `id` (int) - уникальный идентификатор
- `name` (string) - имя пользователя
- `email` (string) - email адрес (уникальный)

### CreateUserRequest

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Fields:**
- `name` (string, required) - имя пользователя (1-255 символов)
- `email` (string, required) - email адрес (уникальный)

### Filter

```json
{
  "name": "john",
  "id": 1
}
```

**Fields:**
- `name` (string, optional) - фильтр по имени (поиск по подстроке)
- `id` (int, optional) - фильтр по ID

## 🚨 Error Codes

### HTTP Status Codes

- `200 OK` - успешный запрос
- `201 Created` - ресурс создан
- `400 Bad Request` - неверный запрос
- `401 Unauthorized` - не авторизован
- `403 Forbidden` - доступ запрещен
- `404 Not Found` - ресурс не найден
- `422 Unprocessable Entity` - ошибка валидации
- `500 Internal Server Error` - внутренняя ошибка сервера

### Error Response

```json
{
  "error": {
    "type": "validation_error",
    "title": "Validation Failed",
    "detail": "Неверный формат email",
    "instance": "/api/v1/users",
    "status": 422
  }
}
```

**Error fields:**
- `type` (string) - тип ошибки
- `title` (string) - заголовок ошибки
- `detail` (string) - детальное описание
- `instance` (string) - путь к ресурсу
- `status` (int) - HTTP статус код

### Error Types

- `validation_error` - ошибка валидации данных
- `not_found` - ресурс не найден
- `unauthorized` - не авторизован
- `forbidden` - доступ запрещен
- `internal_error` - внутренняя ошибка сервера

## 🔧 Usage Examples

### Complete HMAC Example

```bash
#!/bin/bash

# Configuration
API_URL="http://localhost:8080"
SECRET_KEY="your-secret-key"
METHOD="GET"
PATH="/api/v1/users"

# Create signature
SIGNATURE=$(echo -n "${METHOD}${PATH}" | openssl dgst -sha256 -hmac "${SECRET_KEY}" | cut -d' ' -f2)

# Execute request
curl -H "Authorization: HMAC ${SIGNATURE}" \
     "${API_URL}${PATH}"
```

### JavaScript Example

```javascript
const crypto = require('crypto');

function createHMACSignature(method, path, secretKey) {
    const message = method + path;
    return crypto.createHmac('sha256', secretKey)
                 .update(message)
                 .digest('hex');
}

async function getUsers() {
    const method = 'GET';
    const path = '/api/v1/users';
    const secretKey = 'your-secret-key';
    
    const signature = createHMACSignature(method, path, secretKey);
    
    const response = await fetch('http://localhost:8080/api/v1/users', {
        headers: {
            'Authorization': `HMAC ${signature}`
        }
    });
    
    return response.json();
}
```

### Python Example

```python
import hashlib
import requests

def create_hmac_signature(method, path, secret_key):
    message = method + path
    signature = hashlib.hmac.new(
        secret_key.encode('utf-8'),
        message.encode('utf-8'),
        hashlib.sha256
    ).hexdigest()
    return signature

def get_users():
    method = 'GET'
    path = '/api/v1/users'
    secret_key = 'your-secret-key'
    
    signature = create_hmac_signature(method, path, secret_key)
    
    headers = {
        'Authorization': f'HMAC {signature}'
    }
    
    response = requests.get(
        'http://localhost:8080/api/v1/users',
        headers=headers
    )
    
    return response.json()
```

## 📝 Logging

### Request Logging

Все запросы логируются с информацией:
- HTTP метод
- Путь
- Статус код
- Время выполнения
- IP адрес
- User Agent

### Error Logging

Ошибки логируются с детальной информацией:
- Тип ошибки
- Стек вызовов
- Контекст запроса
- Пользователь (если авторизован)

## 🔗 Related Documentation

- [Architecture](architecture.md) - архитектурные принципы
- [Configuration](configuration.md) - настройка переменных окружения
- [Database Setup](database.md) - PostgreSQL и миграции
- [Development](development.md) - руководство по разработке

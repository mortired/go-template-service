# 🌐 API Reference

## 📋 Overview

Users Service provides a RESTful API for user management with HMAC authentication support.

## 🔐 Authentication

### HMAC Authentication

The API uses HMAC (Hash-based Message Authentication Code) for request authentication.

#### How it works

1. **Client** creates request signature using secret key
2. **Server** verifies signature and authenticates request
3. **Access** is granted only with correct signature

#### Creating signature

```bash
# Example of creating HMAC signature
echo -n "GET/api/v1/users" | openssl dgst -sha256 -hmac "your-secret-key"
```

#### Authorization header

```http
Authorization: HMAC <signature>
```

### HMAC Configuration

```env
# Client secrets
HMAC_CLIENT_SECRETS='[{"clientid":"web","secret":"your-secret-key","department":"web","descr":"web client"}]'

# Route access rights
HMAC_ROUTE_RIGHTS='{"web":["/api/v1/users"]}'

# Authentication settings
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

Get list of users.

**Authentication:** HMAC

**Query parameters:**
- `name` (string, optional) - filter by user name (substring search)
- `id` (int, optional) - filter by user ID

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
# Get all users
curl -H "Authorization: HMAC <signature>" http://localhost:8080/api/v1/users

# Filter by name
curl -H "Authorization: HMAC <signature>" "http://localhost:8080/api/v1/users?name=john"

# Filter by ID
curl -H "Authorization: HMAC <signature>" "http://localhost:8080/api/v1/users?id=1"
```

#### POST /api/v1/users

Create new user.

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
- `id` (int) - unique identifier
- `name` (string) - user name
- `email` (string) - email address (unique)

### CreateUserRequest

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Fields:**
- `name` (string, required) - user name (1-255 characters)
- `email` (string, required) - email address (unique)

### Filter

```json
{
  "name": "john",
  "id": 1
}
```

**Fields:**
- `name` (string, optional) - filter by name (substring search)
- `id` (int, optional) - filter by ID

## 🚨 Error Codes

### HTTP Status Codes

- `200 OK` - successful request
- `201 Created` - resource created
- `400 Bad Request` - invalid request
- `401 Unauthorized` - not authorized
- `403 Forbidden` - access denied
- `404 Not Found` - resource not found
- `422 Unprocessable Entity` - validation error
- `500 Internal Server Error` - internal server error

### Error Response

```json
{
  "error": {
    "type": "validation_error",
    "title": "Validation Failed",
    "detail": "Invalid email format",
    "instance": "/api/v1/users",
    "status": 422
  }
}
```

**Error fields:**
- `type` (string) - error type
- `title` (string) - error title
- `detail` (string) - detailed description
- `instance` (string) - resource path
- `status` (int) - HTTP status code

### Error Types

- `validation_error` - data validation error
- `not_found` - resource not found
- `unauthorized` - not authorized
- `forbidden` - access denied
- `internal_error` - internal server error

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

All requests are logged with information:
- HTTP method
- Path
- Status code
- Execution time
- IP address
- User Agent

### Error Logging

Errors are logged with detailed information:
- Error type
- Call stack
- Request context
- User (if authorized)

## 🔗 Related Documentation

- [Architecture](architecture.md) - architecture principles
- [Configuration](configuration.md) - environment variables setup
- [Database Setup](database.md) - PostgreSQL and migrations
- [Development](development.md) - development guide

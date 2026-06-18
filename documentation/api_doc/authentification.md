# Authentication Service

The Authentication Service handles user registration, login, logout, and session management.

---

## Overview

Authentication is based on server-side sessions.

After a successful login:

* A new session is created and stored in the database.
* A `session_id` cookie is sent to the client.
* The session remains valid for 24 hours.
* Protected endpoints require a valid session cookie.

---

## Session Management

### Cookie Information

| Name       | Type   | HttpOnly | Expiration |
| ---------- | ------ | -------- | ---------- |
| session_id | Cookie | Yes      | 24 Hours   |

Example:

```http
Cookie: session_id=550e8400-e29b-41d4-a716-446655440000
```

---

## Rate Limiting

The following endpoints are rate limited:

| Endpoint       | Limit                     |
| -------------- | ------------------------- |
| POST /login    | 1 request every 2 seconds |
| POST /register | 1 request every 2 seconds |

---

# Endpoints

## Register User

Creates a new user account.

### Endpoint

```http
POST /register
```

### Request Body

| Field            | Type    | Required |
| ---------------- | ------- | -------- |
| nickname         | string  | Yes      |
| first_name       | string  | Yes      |
| last_name        | string  | Yes      |
| age              | integer | Yes      |
| gender           | string  | Yes      |
| email            | string  | Yes      |
| password         | string  | Yes      |
| confirm_password | string  | Yes      |

### Example Request

```json
{
  "nickname": "john_doe",
  "first_name": "John",
  "last_name": "Doe",
  "gender": "male",
  "age": "25",
  "email": "john.doe@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

### Validation Rules

| Field            | Rule                                    |
| ---------------- | --------------------------------------- |
| nickname         | 2 - 50 characters                       |
| first_name       | 2 - 50 characters                       |
| last_name        | 2 - 50 characters                       |
| age              | 1 - 99                                  |
| email            | Valid email format (5 - 100 characters) |
| password         | 6 - 20 characters                       |
| confirm_password | Must match password                     |

### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "registration sucess",
  "data": {
    "nickname": "johndoe",
    "email": "john@example.com"
  }
}
```

### Error Responses

#### Missing Required Field

```json
{
  "status_code": 400,
  "message": "email is required"
}
```

#### Invalid Data

```json
{
  "status_code": 400,
  "message": "invalid age",
  "data": ". name (valid) : 2 ~ 50 chars ..."
}
```

#### Email Already Exists

```json
{
  "status_code": 400,
  "message": "Email already exist"
}
```

#### Username Already Taken

```json
{
  "status_code": 400,
  "message": "Username already taken"
}
```

#### Password Confirmation Failed

```json
{
  "status_code": 400,
  "message": "Password not confirmed"
}
```

---

## Login

Authenticates a user using either email or username.

### Endpoint

```http
POST /login
```

### Request Body

| Field      | Type   | Required |
| ---------- | ------ | -------- |
| identifier | string | Yes      |
| password   | string | Yes      |

### Example Request

```json
{
  "identifier": "john.doe@example.com",
  "password": "SecurePass123!"
}
```

or

```json
{
  "identifier": "john_doe",
  "password": "SecurePass123!"
}
```

### Success Response

**Status:** `201 Created`

```json
{
  "status_code": 201,
  "message": "Ilogin Sucess"
}
```

### Set-Cookie Header

```http
Set-Cookie: session_id=<uuid>; HttpOnly; Path=/;
```

### Error Responses

#### Missing Credentials

```json
{
  "status_code": 400,
  "message": "bad credentials"
}
```

#### Invalid Credentials

```json
{
  "status_code": 401,
  "message": "Invalid email/username or password."
}
```

#### Internal Error

```json
{
  "status_code": 500,
  "message": "Internal Server Error"
}
```

### Notes

* Login accepts either email or username.
* Existing sessions are removed before creating a new one.
* Only one active session per user is allowed.
* Session expiration time is 24 hours.

---

## Logout

Destroys the current session and removes the session cookie.

### Endpoint

```http
POST /logout
```

### Authentication

Requires a valid `session_id` cookie.

### Success Response

**Status:** `201 Created`

```json
{
  "status_code": 201,
  "message": "log out succes"
}
```

### Cookie Removal

```http
Set-Cookie: session_id=; Max-Age=-1; HttpOnly; Path=/
```

### Error Responses

#### Invalid Method

```json
{
  "status_code": 405,
  "message": "Method not allowed"
}
```

#### Invalid Path

```json
{
  "status_code": 404,
  "message": "path not found"
}
```

---

# Standard Response Format

All API responses follow the same structure:

```json
{
  "status_code": 200,
  "message": "success message",
  "data": {}
}
```

| Field       | Description                         |
| ----------- | ----------------------------------- |
| status_code | HTTP status code                    |
| message     | Response message                    |
| data        | Additional response data (optional) |

---

# Authentication Flow

```text
Register# login api

    ↓
Login
    ↓
Session Created
    ↓
session_id Cookie Returned
    ↓
Access Protected Endpoints
    ↓
Logout
    ↓
Session Deleted
```

---

# Security Notes

* Passwords are hashed using bcrypt before storage.
* Session identifiers are generated using UUIDs.
* Session cookies are marked as HttpOnly.
* User passwords are never returned in API responses.
* Sessions expire automatically after 24 hours.


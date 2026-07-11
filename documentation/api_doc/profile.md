# User Profile & Information API

## Table of Contents
1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Endpoints](#endpoints)
4. [Error Handling](#error-handling)

---

## Overview

The User Profile API provides endpoints to:
- Retrieve current user's authentication status
- Get user profile information
- Access user public profiles
- Track user presence (online/offline)

---

## Authentication

Most profile endpoints require a valid session.

### Required Cookie

| Name | Type | Required |
|------|------|----------|
| session_id | Cookie | Yes |

**Example:**
```http
Cookie: session_id=550e8400-e29b-41d4-a716-446655440000
```

---

## Endpoints

### Get Current User (Me)

Retrieves current authenticated user's information.

#### Endpoint
```http
GET /api/me
```

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "authenticated": true,
    "id": 1,
    "nickname": "john_doe",
    "last_seen": "2026-06-20T10:30:00Z"
  }
}
```

#### Response Data

| Field | Type | Description |
|-------|------|-------------|
| authenticated | boolean | User is authenticated |
| id | integer | User's ID |
| nickname | string | User's nickname |
| last_seen | string | Last activity timestamp |

#### Error Response

**Not Authenticated:**
```json
{
  "status_code": 401,
  "message": "login required"
}
```

---

### Get User by ID

Retrieves public profile information for a specific user.

#### Endpoint
```http
GET /api/users/{id}
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | User ID |

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "user data get success",
  "data": {
    "id": 2,
    "nickname": "jane_doe",
    "firstname": "Jane",
    "lastname": "Doe",
    "age": 24,
    "gender": "female"
  }
}
```

#### Response Data

| Field | Type | Description |
|-------|------|-------------|
| id | integer | User's ID |
| nickname | string | Display name |
| firstname | string | First name |
| lastname | string | Last name |
| age | integer | User's age |
| gender | string | Gender (male/female/other) |

**Note:** Email and password are never returned.

#### Error Responses

**User Not Found:**
```json
{
  "status_code": 404,
  "message": "user not found"
}
```

**Invalid ID:**
```json
{
  "status_code": 400,
  "message": "Invalid user ID"
}
```

---

## User Profile Data Structure

### Full User Object

```json
{
  "id": 1,
  "nickname": "john_doe",
  "firstname": "John",
  "lastname": "Doe",
  "age": 25,
  "gender": "male",
  "email": "john@example.com",
  "password": "$2a$10$...",
  "last_seen": "2026-06-20T10:30:00Z"
}
```

### Public Profile (Returned to Others)

```json
{
  "id": 1,
  "nickname": "john_doe",
  "firstname": "John",
  "lastname": "Doe",
  "age": 25,
  "gender": "male"
}
```

**Sensitive Fields Hidden:**
- Email
- Password hash
- Last seen (only available to self)

---

## Presence Information

### User Status Tracking

The system tracks user presence through:

1. **WebSocket Connection**: User connected to WebSocket
2. **Last Seen**: Timestamp of last activity
3. **Online/Offline**: Determined by active WebSocket connection

### Getting User Status

**In Conversation List:**
```json
{
  "profile": {
    "id": 2,
    "nickname": "jane_doe",
    ...
  },
  "conversation": {
    "status": "online",
    "lastSeen": "Just now"
  }
}
```

**In Conversations:**
```json
{
  "users": [
    {
      "id": 2,
      "nickname": "jane_doe",
      "status": "online"
    },
    {
      "id": 3,
      "nickname": "bob_smith",
      "status": "offline"
    }
  ]
}
```

---

## Frontend Integration

### Display User Profile

```javascript
// Get current user
async function getCurrentUser() {
  const response = await fetch(`${serverUri}/me`, {
    credentials: 'include'
  });
  
  const result = await response.json();
  
  if (response.ok) {
    return result.data;  // {authenticated, id, nickname, last_seen}
  } else {
    throw new Error(result.message);
  }
}

// Get other user's profile
async function getUserProfile(userId) {
  const response = await fetch(`${serverUri}/users/${userId}`, {
    credentials: 'include'
  });
  
  const result = await response.json();
  
  if (response.ok) {
    return result.data;  // Public profile
  } else {
    throw new Error(result.message);
  }
}
```

### Display User Information in Comments

```javascript
// In comment render
const comment = {
  id: 1,
  userId: 2,
  nickname: "jane_doe",
  text: "Great post!",
  ...
};

const html = `
  <div class="comment">
    <span class="comment-username">${sanitize(comment.nickname)}</span>
    <span class="comment-text">${sanitize(comment.text)}</span>
  </div>
`;
```

### Display User Information in Posts

```javascript
// In post header
const post = {
  id: 4,
  userId: 1,
  nickname: "john_doe",
  title: "My First Post",
  ...
};

const html = `
  <div class="post-header">
    <h3>${sanitize(post.title)}</h3>
    <span class="post-meta">
      Posted by ${post.nickname}
    </span>
  </div>
`;
```

### Track User Presence

```javascript
// WebSocket events provide presence updates
const onlineUsers = new Set();

ws.onmessage = (event) => {
  const {event: eventType, data} = JSON.parse(event.data);
  
  switch(eventType) {
    case 'init':
      data.forEach(userId => onlineUsers.add(userId));
      updatePresenceUI();
      break;
      
    case 'client_connect':
      onlineUsers.add(data);
      updateUserPresence(data, 'online');
      break;
      
    case 'client_disconnect':
      onlineUsers.delete(data);
      updateUserPresence(data, 'offline');
      break;
  }
};

function updateUserPresence(userId, status) {
  // Update UI indicators
  const indicator = document.querySelector(
    `[data-user-id="${userId}"] .presence-indicator`
  );
  if (indicator) {
    indicator.classList.toggle('online', status === 'online');
  }
}
```

---

## Security Considerations

### Data Exposure

**Never expose to client:**
- Password hashes
- Email addresses (of other users)
- Session tokens
- Private personal information

**Safe to expose:**
- Nickname
- First/Last name
- Age
- Gender
- Public profile info

### Authentication Checks

All endpoints that require authentication should:
1. Validate session cookie
2. Verify session not expired
3. Check user exists
4. Return 401 if invalid

---

## Error Handling

### Standard Error Response

```json
{
  "status_code": <HTTP_STATUS>,
  "message": "Error description",
  "data": null
}
```

### Common HTTP Status Codes

| Status | Meaning | Action |
|--------|---------|--------|
| 200 | Success | Use returned data |
| 400 | Bad request | Check request format |
| 401 | Unauthorized | Redirect to login |
| 404 | Not found | User doesn't exist |
| 500 | Server error | Retry or contact support |

### Error Examples

**Unauthenticated Request:**
```bash
$ curl http://localhost:8080/api/me
{"status_code": 401, "message": "login required"}
```

**Invalid User ID:**
```bash
$ curl http://localhost:8080/api/users/invalid
{"status_code": 400, "message": "Invalid user ID"}
```

**User Not Found:**
```bash
$ curl http://localhost:8080/api/users/9999
{"status_code": 404, "message": "user not found"}
```

---

## Example Flows

### User Registration & Login Flow

```
1. Register
   POST /api/register
   {nickname, email, password, ...}
   ← User created

2. Login
   POST /api/login
   {identifier, password}
   ← Session cookie set

3. Get Current User
   GET /api/me
   Cookie: session_id=...
   ← {id, nickname, authenticated}

4. Load Feed
   GET /api/posts
   Cookie: session_id=...
   ← Posts with author nicknames
   
5. View Other User Profile
   GET /api/users/42
   ← Public profile info
```

### User Profile in Comments Flow

```
1. View Post
   GET /api/posts/4
   ← Post with comments

2. Each Comment Contains
   {
     userId: 2,
     nickname: "jane_doe",
     text: "Comment text"
   }

3. Hover/Click on Username (Future)
   GET /api/users/2
   ← Full public profile
   ← Can initiate chat or view profile
```

### User Presence in Chat Flow

```
1. Connect to Chat
   WebSocket /ws
   GET /api/conversations
   ← List of users with status

2. Receive Presence Updates
   WebSocket: client_connect (userId)
   WebSocket: client_disconnect (userId)
   → Update UI with online/offline

3. View User Profile in Chat
   GET /api/users/{userId}
   ← Public profile info
```


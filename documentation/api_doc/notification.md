# Real-Time Messaging & Notifications API

## Table of Contents
1. [Overview](#overview)
2. [WebSocket Connection](#websocket-connection)
3. [Real-Time Events](#real-time-events)
4. [Message API](#message-api)
5. [Conversations API](#conversations-api)
6. [Error Handling](#error-handling)

---

## Overview

The Real-Time Messaging system provides:
- **WebSocket-based instant messaging** between users
- **Real-time notifications** for events (new messages, user presence)
- **Conversation management** with message history
- **User presence tracking** (online/offline status)

### Technology
- **Protocol**: WebSocket (ws://)
- **Format**: JSON
- **Authentication**: Session cookie validation

---

## WebSocket Connection

### Establish Connection

#### URL
```
ws://localhost:8080/ws
wss://api.example.com/ws  (production)
```

#### Connection Process

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleMessage(message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket closed');
  // Attempt reconnect with exponential backoff
};
```

### Authentication

The server validates the session cookie from the initial WebSocket upgrade request:

```http
GET /ws HTTP/1.1
Connection: Upgrade
Upgrade: websocket
Cookie: session_id=550e8400-e29b-41d4-a716-446655440000
```

**Invalid Session Response:**
```http
HTTP/1.1 401 Unauthorized
{"error": "not authenticated"}
```

---

## Real-Time Events

### Server-to-Client Events

#### 1. `init`

Sent when user first connects. Contains list of currently online users.

```json
{
  "event": "init",
  "data": [1, 3, 5, 7]
}
```

**Usage:** Initialize online users list

---

#### 2. `client_connect`

Notifies when another user comes online.

```json
{
  "event": "client_connect",
  "data": 2
}
```

**Data:** User ID of user who came online

**Usage:** Update presence indicator (green dot)

---

#### 3. `client_disconnect`

Notifies when another user goes offline.

```json
{
  "event": "client_disconnect",
  "data": 2
}
```

**Data:** User ID of user who went offline

**Usage:** Update presence indicator (gray dot)

---

#### 4. `new_message`

Incoming private message from another user.

```json
{
  "event": "new_message",
  "data": {
    "senderId": 1,
    "conversationId": 42,
    "text": "Hey, how are you?",
    "timestamp": "2026-06-20T10:30:00Z",
    "isMine": false,
    "isNewConversation": false
  }
}
```

**Data Fields:**
| Field | Type | Description |
|-------|------|-------------|
| senderId | integer | User ID of message sender |
| conversationId | integer | Conversation ID |
| text | string | Message content |
| timestamp | string | ISO timestamp |
| isMine | boolean | False (always from others) |
| isNewConversation | boolean | True if first message in conversation |

**Usage:** Display new message, play notification sound

---

#### 5. `new_post`

Notification of new post created by another user.

```json
{
  "event": "new_post",
  "data": {
    "postId": 42,
    "title": "Amazing Discovery",
    "author": "john_doe"
  }
}
```

**Usage:** Optional - show "new post" indicator in feed

---

#### 6. `force_logout`

Server forces logout of all tabs/devices (multi-tab sync).

```json
{
  "event": "force_logout",
  "data": {}
}
```

**Causes:**
- User logged out on another device
- Admin action
- Session invalidated

**Usage:** Clear session, redirect to login page

---

### Client-to-Server Events

#### Send Message

```json
{
  "type": "message",
  "data": {
    "recipientId": 2,
    "text": "Hello!",
    "conversationId": 42
  }
}
```

**Implementation:**
```javascript
function sendMessage(recipientId, text, conversationId) {
  ws.send(JSON.stringify({
    type: 'message',
    data: {
      recipientId,
      text,
      conversationId
    }
  }));
}
```

---

## Message API

### Send Message (HTTP + WebSocket)

#### Endpoint
```http
POST /api/messages
```

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| receiver_id | integer | Yes | User ID of recipient |
| text | string | Yes | Message content |
| conversation_id | integer | No | Existing conversation ID |

#### Example Request

**New Conversation:**
```json
{
  "receiver_id": 2,
  "text": "Hi, this is my first message to you"
}
```

**Existing Conversation:**
```json
{
  "receiver_id": 2,
  "text": "Hi again!",
  "conversation_id": 42
}
```

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "Message sent",
  "data": {
    "id": 123,
    "conversationId": 42,
    "senderId": 1,
    "receiverId": 2,
    "text": "Hi, this is my first message to you",
    "createdAt": "2026-06-20T10:30:00Z",
    "isRead": 0
  }
}
```

#### Error Responses

**Receiver Not Found:**
```json
{
  "status_code": 400,
  "message": "Receiver not found"
}
```

**Cannot Message Self:**
```json
{
  "status_code": 400,
  "message": "Cannot message yourself"
}
```

**Invalid Conversation:**
```json
{
  "status_code": 403,
  "message": "Not allowed"
}
```

#### Rate Limit
1 request per 100 milliseconds

---

## Conversations API

### Get All Conversations

Retrieves list of conversations for the current user, sorted by last message.

#### Endpoint
```http
GET /api/conversations?offset=0&limit=30
```

#### Query Parameters

| Parameter | Type | Default | Max |
|-----------|------|---------|-----|
| offset | integer | 0 | - |
| limit | integer | 30 | 30 |

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "ok",
  "data": [
    {
      "profile": {
        "id": 2,
        "nickname": "jane_doe",
        "firstname": "Jane",
        "lastname": "Doe",
        "age": 24,
        "gender": "female"
      },
      "conversation": {
        "conversationId": 42,
        "date": "2 hours ago",
        "lastMessage": "That's great!",
        "status": "online",
        "lastSeen": "2 hours ago",
        "unreadCount": 0,
        "lastSender": "jane_doe"
      }
    },
    {
      "profile": {
        "id": 3,
        "nickname": "bob_smith",
        ...
      },
      "conversation": {
        "conversationId": 43,
        "date": "1 day ago",
        "lastMessage": "See you tomorrow!",
        "status": "offline",
        "lastSeen": "1 day ago",
        "unreadCount": 2,
        "lastSender": "me"
      }
    }
  ]
}
```

#### Response Data

**Per Conversation:**
| Field | Type | Description |
|-------|------|-------------|
| profile | object | Other user's profile info |
| conversation.conversationId | integer | Conversation ID |
| conversation.date | string | Last message date (human readable) |
| conversation.lastMessage | string | Preview of last message |
| conversation.status | string | "online" or "offline" |
| conversation.lastSeen | string | When user was last seen |
| conversation.unreadCount | integer | Unread messages |
| conversation.lastSender | string | "me" or other user's nickname |

#### Sorting

1. **By Last Message**: Conversations with recent messages appear first
2. **Alphabetical**: New conversations (no messages) sorted alphabetically by nickname

#### Rate Limit
1 request per 3 seconds

---

### Get Conversation by ID

Retrieves messages from a specific conversation.

#### Endpoint
```http
GET /api/conversation/{conversationId}?offset=0&limit=10
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| conversationId | integer | Conversation ID |

#### Query Parameters
| Parameter | Type | Default | Max |
|-----------|------|---------|-----|
| offset | integer | 0 | - |
| limit | integer | 10 | 50 |

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "ok",
  "data": {
    "conversationId": 42,
    "messages": [
      {
        "id": 123,
        "senderId": 1,
        "senderNickname": "john_doe",
        "text": "Hi there!",
        "createdAt": "2026-06-20T09:00:00Z",
        "timeAgo": "1 hour ago",
        "isRead": 1
      },
      {
        "id": 124,
        "senderId": 2,
        "senderNickname": "jane_doe",
        "text": "Hey! How are you?",
        "createdAt": "2026-06-20T09:05:00Z",
        "timeAgo": "1 hour ago",
        "isRead": 1
      }
    ]
  }
}
```

#### Message Ordering
Messages returned in reverse chronological order (newest first). Combine with pagination to load older messages.

#### Pagination Example

```javascript
// Load initial messages (10 most recent)
GET /api/conversation/42?offset=0&limit=10

// Load older messages (next 10)
GET /api/conversation/42?offset=10&limit=10

// Load even older
GET /api/conversation/42?offset=20&limit=10
```

#### Error Responses

**Conversation Not Found:**
```json
{
  "status_code": 404,
  "message": "Conversation not found"
}
```

**Not Authorized (Not Part of Conversation):**
```json
{
  "status_code": 403,
  "message": "not allowed"
}
```

#### Rate Limit
1 request per 3 seconds

---

## Frontend Integration

### Connection Lifecycle

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');

// Event handlers
ws.addEventListener('open', () => {
  console.log('WebSocket connected');
  updatePresenceIndicators();
});

ws.addEventListener('message', (event) => {
  const {event: eventType, data} = JSON.parse(event.data);
  
  switch(eventType) {
    case 'init':
      initializeOnlineUsers(data);
      break;
    case 'client_connect':
      addOnlineUser(data);
      break;
    case 'client_disconnect':
      removeOnlineUser(data);
      break;
    case 'new_message':
      displayNewMessage(data);
      break;
    case 'force_logout':
      handleForceLogout();
      break;
  }
});

ws.addEventListener('close', () => {
  console.log('WebSocket disconnected');
  // Attempt reconnect
  setTimeout(reconnect, 3000);
});
```

### Sending Messages

```javascript
async function sendMessage(recipientId, text) {
  try {
    const response = await fetch(`${serverUri}/messages`, {
      method: 'POST',
      credentials: 'include',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        receiver_id: recipientId,
        text,
        conversation_id: currentConversationId
      })
    });
    
    const result = await response.json();
    if (response.ok) {
      displaySentMessage(result.data);
    } else {
      showError(result.message);
    }
  } catch (err) {
    console.error('Send failed:', err);
  }
}
```

---

## Error Handling

### WebSocket Connection Errors

**Unauthorized Connection:**
```
Status 401: Session cookie invalid or missing
→ Redirect to login
```

**Connection Lost:**
```
Implement exponential backoff reconnection:
Retry 1: 3 seconds
Retry 2: 6 seconds
Retry 3: 12 seconds
...
Max backoff: 60 seconds
```

### Message Errors

**Invalid Recipients:**
```json
{
  "status_code": 400,
  "message": "Receiver not found"
}
```

**Rate Limited:**
```json
{
  "status_code": 429,
  "message": "Too many requests"
}
```

---

## Performance Considerations

### Message Pagination

Load messages in batches to avoid loading entire conversation history:

```javascript
// Good: Load 10-20 messages at a time
GET /api/conversation/42?limit=10

// Avoid: Loading thousands of messages
GET /api/conversation/42?limit=10000
```

### Scroll Detection

Implement "load more" on scroll to top:

```javascript
const chatContainer = document.querySelector('.messages');
chatContainer.addEventListener('scroll', () => {
  if (chatContainer.scrollTop === 0) {
    loadOlderMessages();
  }
});
```

### Connection Stability

- Implement heartbeat/ping-pong
- Handle disconnections gracefully
- Persist unread count locally
- Retry failed sends


# Real-Time Messaging & Notifications API

## Table of Contents
1. [Overview](#overview)
2. [WebSocket Events](#websocket-events)
3. [Conversations API](#conversations-api)
4. [Message API](#message-api)
5. [Error Handling](#error-handling)

---

## Overview

The Real-Time Messaging system provides:
- **WebSocket-based instant messaging** between users
- **Real-time notifications** for events (new messages, user presence)
- **Conversation management** with message history
- **User presence tracking** (online/offline status)

### Technology
- **Protocol**: WebSocket (ws://) for real-time, REST for history
- **Format**: JSON
- **Authentication**: Session cookie validation

---

## WebSocket Events

> **📖 For complete WebSocket event documentation, see [websocket-events.md](./websocket-events.md)**

Key events for messaging:

- **`init`** - Initial online users list upon connection
- **`client_connect`** - User came online
- **`client_disconnect`** - User went offline
- **`new_message`** - Incoming private message
- **`new_post`** - New forum post created
- **`force_logout`** - Force user logout (multi-tab sync)

See [websocket-events.md](./websocket-events.md) for detailed event schemas, examples, and implementation patterns.

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


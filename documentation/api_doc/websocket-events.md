# WebSocket Events Reference Guide

## Table of Contents
1. [Overview](#overview)
2. [Connection](#connection)
3. [Server-Sent Events](#server-sent-events)
4. [Client-Sent Events](#client-sent-events)
5. [Event Examples](#event-examples)
6. [Error Handling](#error-handling)
7. [Best Practices](#best-practices)

---

## Overview

The WebSocket service provides real-time, bidirectional communication between server and clients. All events are transmitted as JSON objects.

### WebSocket URL
```
ws://localhost:8080/ws          (development)
wss://api.example.com/ws        (production)
```

### Protocol Format
```json
{
  "event": "event_name",
  "data": {}
}
```

---

## Connection

### Establishing Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.addEventListener('open', () => {
  console.log('Connected');
});

ws.addEventListener('close', () => {
  console.log('Disconnected');
});

ws.addEventListener('error', (error) => {
  console.error('WebSocket error:', error);
});
```

### Authentication

The server validates the session cookie from the initial WebSocket upgrade:

```http
GET /ws HTTP/1.1
Connection: Upgrade
Upgrade: websocket
Cookie: session_id=550e8400-e29b-41d4-a716-446655440000
```

**Valid:** Connection established ✅
**Invalid:** HTTP 401 Unauthorized ❌

### Reconnection Strategy

Implement exponential backoff for automatic reconnection:

```javascript
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 10;

function reconnect() {
  if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
    console.error('Max reconnection attempts reached');
    return;
  }
  
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000);
  console.log(`Reconnecting in ${delay}ms...`);
  
  setTimeout(() => {
    reconnectAttempts++;
    ws = new WebSocket('ws://localhost:8080/ws');
    setupWebSocket();
  }, delay);
}
```

---

## Server-Sent Events

### 1. `init`

**Sent When:** User first connects to WebSocket

**Purpose:** Provide list of currently online users

**Event Data:**
```json
{
  "event": "init",
  "data": [1, 3, 5, 7, 12]
}
```

**Data Type:** Array of user IDs (integers)

**Usage:**
```javascript
ws.addEventListener('message', (event) => {
  const {event: eventType, data} = JSON.parse(event.data);
  
  if (eventType === 'init') {
    const onlineUsers = new Set(data);
    updateOnlineIndicators(onlineUsers);
  }
});
```

**Typical Timing:** Immediately upon connection, before any other events

---

### 2. `client_connect`

**Sent When:** Another user connects to WebSocket

**Purpose:** Notify that a user is now online

**Event Data:**
```json
{
  "event": "client_connect",
  "data": 2
}
```

**Data Type:** User ID (integer)

**Usage:**
```javascript
case 'client_connect':
  onlineUsers.add(data);
  updateUserPresence(data, 'online');
  showNotification(`User ${data} is now online`);
  break;
```

**Where Used:**
- Chat sidebar - user presence indicator
- Conversation list - show online status
- Real-time notifications

---

### 3. `client_disconnect`

**Sent When:** A connected user disconnects from WebSocket

**Purpose:** Notify that a user is now offline

**Event Data:**
```json
{
  "event": "client_disconnect",
  "data": 2
}
```

**Data Type:** User ID (integer)

**Usage:**
```javascript
case 'client_disconnect':
  onlineUsers.delete(data);
  updateUserPresence(data, 'offline');
  showNotification(`User ${data} is now offline`);
  break;
```

**UI Impact:**
- Change presence indicator from green to gray
- Update last seen time
- Disable real-time typing indicators

---

### 4. `new_message`

**Sent When:** Authenticated user receives a private message

**Purpose:** Deliver incoming message in real-time

**Event Data:**
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
| senderId | integer | User ID of sender |
| conversationId | integer | Conversation ID |
| text | string | Message content |
| timestamp | string | ISO 8601 timestamp |
| isMine | boolean | Always false (from server) |
| isNewConversation | boolean | First message in conversation |

**Usage:**
```javascript
case 'new_message':
  if (data.isNewConversation) {
    // Add new conversation to list
    addNewConversation(data.senderId);
  }
  
  // Display message in active conversation
  if (currentConversationId === data.conversationId) {
    displayMessage({
      text: data.text,
      sender: data.senderId,
      timestamp: data.timestamp
    });
  }
  
  // Show notification
  showNotification(`New message from user ${data.senderId}`);
  
  // Play sound
  playMessageSound();
  
  // Update unread count
  incrementUnreadCount(data.conversationId);
  break;
```

**Notifications:**
- Toast notification with sender name
- Sound alert (if enabled)
- Update conversation last message preview
- Increment unread badge

---

### 5. `new_post`

**Sent When:** Another user creates a new post (broadcasts to all)

**Purpose:** Notify users of new forum content

**Event Data:**
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

**Data Fields:**
| Field | Type | Description |
|-------|------|-------------|
| postId | integer | ID of new post |
| title | string | Post title |
| author | string | Nickname of author |

**Usage:**
```javascript
case 'new_post':
  console.log(`New post from ${data.author}: ${data.title}`);
  
  // Option 1: Show notification
  showNotification(`${data.author} posted: ${data.title}`);
  
  // Option 2: Auto-refresh feed
  refreshPostsFeed();
  
  // Option 3: Show "New posts available" banner
  showNewPostsBanner();
  break;
```

**Typical Implementation:** Currently optional - shows notification but feed refresh is manual

---

### 6. `force_logout`

**Sent When:** User is force logged out (server action)

**Purpose:** Terminate session on all tabs/devices

**Event Data:**
```json
{
  "event": "force_logout",
  "data": {}
}
```

**Data:** Empty object

**Causes:**
- User logged out on another device
- Admin revoked session
- Session invalidated server-side
- Force logout command from multi-tab sync

**Usage:**
```javascript
case 'force_logout':
  console.warn('Force logout received');
  
  // Clear session
  localStorage.clear();
  sessionStorage.clear();
  
  // Close WebSocket
  ws.close();
  
  // Redirect to login
  window.location.href = '/login';
  
  // Show message
  showNotification('You have been logged out', 'info');
  break;
```

**UI Flow:**
1. Clear all local data
2. Close WebSocket connection
3. Redirect to login page
4. Optional: Show "Logged out from another device" message

---

## Client-Sent Events

### Send Message

**Used For:** Sending private messages via WebSocket

**Format:**
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

**Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | string | Yes | Always "message" |
| recipientId | integer | Yes | User ID of recipient |
| text | string | Yes | Message content |
| conversationId | integer | No | Existing conversation ID |

**Implementation:**
```javascript
function sendMessage(recipientId, text, conversationId = null) {
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

**Server Processing:**
1. Validate recipient exists
2. Check conversation ownership
3. Store message in database
4. Broadcast to recipient via WebSocket
5. Update conversation last_message

---

## Event Examples

### Complete Message Flow

```javascript
// Setup WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');
const onlineUsers = new Set();

ws.addEventListener('message', (event) => {
  const {event: eventType, data} = JSON.parse(event.data);
  
  switch(eventType) {
    // 1. Connection established
    case 'init':
      console.log('Online users:', data);
      data.forEach(userId => onlineUsers.add(userId));
      renderOnlineList();
      break;
    
    // 2. User comes online
    case 'client_connect':
      onlineUsers.add(data);
      updateUserIndicator(data, true);
      break;
    
    // 3. User goes offline
    case 'client_disconnect':
      onlineUsers.delete(data);
      updateUserIndicator(data, false);
      break;
    
    // 4. Receive message
    case 'new_message':
      displayMessage(data);
      playSound('message.mp3');
      updateConversationPreview(data.conversationId);
      break;
    
    // 5. New post created
    case 'new_post':
      showToast(`${data.author} posted: ${data.title}`);
      break;
    
    // 6. Force logout
    case 'force_logout':
      logout();
      redirectToLogin();
      break;
  }
});

// Send message when user clicks send
function onSendClick() {
  const text = messageInput.value;
  const recipientId = selectedUser.id;
  const conversationId = currentConversation?.id || null;
  
  ws.send(JSON.stringify({
    type: 'message',
    data: {
      recipientId,
      text,
      conversationId
    }
  }));
  
  messageInput.value = '';
}
```

### Typical Session Lifecycle

```
1. User opens browser
   ↓
2. JavaScript connects to WebSocket
   ws = new WebSocket('ws://localhost:8080/ws')
   ↓
3. Server validates session cookie
   ✅ Valid → Connection established
   ❌ Invalid → 401 Unauthorized
   ↓
4. Server sends 'init' event
   {event: 'init', data: [1, 3, 5]}
   ↓
5. User presence updated on frontend
   onlineUsers = {1, 3, 5}
   ✓ Chat sidebar shows online users
   ↓
6. User A sends message to User B
   ws.send({type: 'message', ...})
   ↓
7. Server receives, stores, broadcasts
   If User B online:
     → User B receives 'new_message' event
   If User B offline:
     → Message stored, delivered on reconnect
   ↓
8. Other users go online/offline
   → 'client_connect' / 'client_disconnect' events
   → Real-time presence updated
   ↓
9. User logs out from another tab
   → Server sends 'force_logout' event
   → This tab closes connection & redirects
   ↓
10. User closes browser
    → WebSocket closes
    → Server removes from active users
```

---

## Error Handling

### Connection Errors

**Failed to Connect:**
```javascript
ws.addEventListener('error', (error) => {
  console.error('WebSocket error:', error);
  // Attempt reconnection
  scheduleReconnect();
});
```

**Server Unavailable:**
```
Connection refused → Server not running
→ Show "Server offline" message
→ Retry with exponential backoff
```

**Session Expired:**
```
Connection closed (1000)
→ Redirect to login
→ Clear local session
```

### Message Errors

**Invalid JSON:**
```javascript
try {
  const {event: eventType, data} = JSON.parse(event.data);
} catch (err) {
  console.error('Invalid JSON received:', err);
  // Ignore and wait for next message
}
```

**Unknown Event Type:**
```javascript
default:
  console.warn(`Unknown event type: ${eventType}`);
  // Ignore unknown events for forward compatibility
```

---

## Best Practices

### 1. Always Validate Events
```javascript
if (!eventType || !data) {
  console.warn('Invalid event structure');
  return;
}
```

### 2. Handle Disconnections Gracefully
```javascript
ws.addEventListener('close', () => {
  console.log('Disconnected, attempting to reconnect...');
  setTimeout(reconnect, 3000);
});
```

### 3. Implement Heartbeat (Ping/Pong)
```javascript
// Send ping every 30 seconds
setInterval(() => {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({type: 'ping'}));
  }
}, 30000);
```

### 4. Debounce Real-Time Updates
```javascript
// Don't update every single message
const updates = new Map();

function debounceUpdate(conversationId) {
  if (updates.has(conversationId)) return;
  
  updates.set(conversationId, true);
  updateUI(conversationId);
  
  setTimeout(() => updates.delete(conversationId), 500);
}
```

### 5. Store Message Queue During Disconnect
```javascript
const messageQueue = [];

function sendMessage(msg) {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(msg));
  } else {
    messageQueue.push(msg);
  }
}

ws.addEventListener('open', () => {
  while (messageQueue.length > 0) {
    ws.send(JSON.stringify(messageQueue.shift()));
  }
});
```

### 6. Monitor Connection Health
```javascript
let lastMessageTime = Date.now();
let connectionHealthy = true;

setInterval(() => {
  const elapsed = Date.now() - lastMessageTime;
  
  if (elapsed > 60000) { // 60 seconds
    connectionHealthy = false;
    reconnect();
  }
}, 10000);

ws.addEventListener('message', () => {
  lastMessageTime = Date.now();
  connectionHealthy = true;
});
```

### 7. Security: Validate User Actions
```javascript
// Server-side validation
function SendMessage(w http.ResponseWriter, r *http.Request) {
  // Extract and validate sender from session
  senderId := getUserIDFromCookie(r)
  
  // Validate recipient exists
  validateRecipientExists(receiverId)
  
  // Validate conversation ownership
  validateUserInConversation(senderId, conversationId)
  
  // Only then process message
  storeMessage(...)
}
```

---

## Debugging WebSocket Events

### Enable Logging
```javascript
const originalSend = ws.send;
ws.send = function(data) {
  console.log('📤 Sending:', JSON.parse(data));
  originalSend.call(this, data);
};

ws.addEventListener('message', (event) => {
  console.log('📥 Received:', JSON.parse(event.data));
});
```

### Browser DevTools
```javascript
// In browser console
window.ws.send(JSON.stringify({
  type: 'message',
  data: {
    recipientId: 2,
    text: 'Test message',
    conversationId: 1
  }
}));
```

### Monitor Connection State
```javascript
const states = ['CONNECTING', 'OPEN', 'CLOSING', 'CLOSED'];
console.log('WebSocket state:', states[ws.readyState]);
```

---

## Event Reference Table

| Event | Direction | Type | Frequency | Purpose |
|-------|-----------|------|-----------|---------|
| init | Server→Client | One-time | On connect | Initial online users |
| client_connect | Server→Client | Event | Per connection | User came online |
| client_disconnect | Server→Client | Event | Per disconnection | User went offline |
| new_message | Server→Client | Event | Per message | Incoming message |
| new_post | Server→Client | Event | Per new post | New forum post |
| force_logout | Server→Client | Event | As needed | Force logout user |
| message | Client→Server | Request | Per send | Send private message |

---

## Integration with REST API

WebSocket complements REST API:

| Task | Protocol | Reason |
|------|----------|--------|
| Send message | WebSocket | Real-time delivery |
| Load message history | REST (GET) | Pagination needed |
| Update user profile | REST (POST) | Atomic updates |
| Receive live notifications | WebSocket | Instant delivery |
| Create posts | REST (POST) | File uploads possible |
| Get post list | REST (GET) | Filtering & pagination |


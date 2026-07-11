# System Architecture & Design

## Table of Contents
1. [Application Architecture](#application-architecture)
2. [Component Breakdown](#component-breakdown)
3. [Data Flow](#data-flow)
4. [Request/Response Cycle](#requestresponse-cycle)
5. [WebSocket Communication](#websocket-communication)
6. [Security Architecture](#security-architecture)
7. [Performance Optimizations](#performance-optimizations)

---

## Application Architecture

### Client-Server Model

```
┌─────────────────────────────────────────────────────────────┐
│                         BROWSER (Client)                    │
│                     Single Page Application                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Routes: /feed, /post/:id, /chat, /login, /register │   │
│  │ Components: Post, Comment, Message, Form components│   │
│  │ Storage: localStorage (session tokens, UI state)   │   │
│  └─────────────────────────────────────────────────────┘   │
│                              │                              │
│  ┌──────────────────────────┼────────────────────────────┐  │
│  │ HTTP (REST API)          │    WebSocket (Real-time)   │  │
│  │ • GET /api/posts         │ • Incoming messages        │  │
│  │ • POST /api/comments     │ • User status updates      │  │
│  │ • POST /api/auth         │ • Notifications            │  │
│  └──────────────────────────┼────────────────────────────┘  │
└─────────────────────────────┼─────────────────────────────┘
                              │ HTTP/WSS
                              │
┌─────────────────────────────▼─────────────────────────────┐
│                   SERVER (Backend - Go)                   │
│              HTTP Server + WebSocket Hub                  │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ REST Handlers:                                       │ │
│  │ • Authentication (login, register, logout)          │ │
│  │ • Posts (CRUD, filtering, reactions)                │ │
│  │ • Comments (CRUD, reactions)                        │ │
│  │ • Messages (send, retrieve)                         │ │
│  │ • Users (profile, presence)                         │ │
│  └──────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ WebSocket Hub:                                       │ │
│  │ • Connection management                             │ │
│  │ • Broadcasting events                               │ │
│  │ • User presence tracking                            │ │
│  │ • Message routing                                   │ │
│  └──────────────────────────────────────────────────────┘ │
└─────────────────────────────┬─────────────────────────────┘
                              │ SQL
                              │
┌─────────────────────────────▼─────────────────────────────┐
│               DATABASE (SQLite)                           │
│  • Users & Sessions                                       │
│  • Posts & Comments                                       │
│  • Reactions (Likes/Dislikes)                             │
│  • Conversations & Messages                               │
└─────────────────────────────────────────────────────────┘
```

---

## Component Breakdown

### Frontend Architecture Layers

#### 1. Pages Layer
Each page is a JavaScript module that renders different application states:

```javascript
// pages/feed.js - Main feed
render(data) → Returns HTML string
setup() → Attaches event listeners
```

**Available Pages:**
- `feed.js` - Post feed with filtering
- `post.js` - Single post detail with comments
- `chat.js` - Private messaging interface
- `login.js` - Login form
- `register.js` - Registration form
- `error.js` - Error page display

#### 2. Components Layer
Reusable UI components:

```javascript
Post(post, options) → Returns HTML string
Comment(comment) → Returns HTML string
Header(nickname) → Returns HTML header
```

**Key Components:**
- `Post.js` - Renders post with reactions
- `Comment.js` - Renders comment with actions
- `Conversation.js` - Renders conversation item
- `Message.js` - Renders message bubble
- `LoginForm.js` - Login form component
- `PostCreationForm.js` - Create post modal

#### 3. Services Layer
Core application logic:

```javascript
// services/router.js - Client-side routing
router.navigate(path) → Navigate and render page
router.init() → Initialize routing & auth guards

// services/auth.js - Authentication
isAuthenticated() → Check session validity
logout() → Clear session

// services/websocket.js - Real-time communication
ws.connect() → Establish connection
ws.emit(event, data) → Send event
ws.on(event, callback) → Listen for events

// services/toast.js - Notifications
showToast(message, type) → Display notification
```

#### 4. API Layer
HTTP client functions:

```javascript
// api/posts.js
getPosts(filters) → GET /api/posts
createPost(data) → POST /api/posts/create
PostResolver(id, action) → POST/DELETE /api/posts/:id/:action

// api/auth.js
login(credentials) → POST /api/login
register(data) → POST /api/register
logout() → POST /api/logout

// api/comments.js
CreateComment(data) → POST /api/comments/create
CommentResolver(id, action) → POST/DELETE /api/comments/:id/:action
```

#### 5. Scripts Layer
Event handlers and page initialization:

```javascript
// scripts/_feed.js
setup() → Initialize feed page
setupEvents() → Attach listeners
fetchPosts() → Load posts with filters
renderPosts(posts) → Append posts to DOM

// scripts/_post.js
setup() → Initialize post detail page
setupEventListeners() → Comment actions

// scripts/_chat.js
setup() → Initialize chat interface
sendMessage(data) → Send chat message
reRender(action, userId) → Update UI
```

### Backend Architecture Layers

#### 1. Route Layer (`routes/routes.go`)
Defines all API endpoints with middleware:

```go
http.HandleFunc("/api/login", 
    RateLimit(CheckSessionCookie(Login, false), 2*time.Second))
    
http.HandleFunc("/api/posts",
    RateLimit(CheckSessionCookie(GetPosts, true), 3*time.Second))
```

**Rate Limits:**
- Auth endpoints: 2 seconds
- Post endpoints: 3 seconds  
- Reactions: 250ms
- Messages: 100ms

#### 2. Handler Layer (`handlers/`)
Request processing and business logic:

```go
// handlers/auth.go
func Login(w http.ResponseWriter, r *http.Request)
func Register(w http.ResponseWriter, r *http.Request)
func Logout(w http.ResponseWriter, r *http.Request)

// handlers/post.go
func GetPosts(w http.ResponseWriter, r *http.Request)
func CreatePost(w http.ResponseWriter, r *http.Request)
func GetPostById(w http.ResponseWriter, r *http.Request)
func PostResolver(w http.ResponseWriter, r *http.Request) // like/dislike/delete

// handlers/comment.go
func CreateComment(w http.ResponseWriter, r *http.Request)
func CommentResolver(w http.ResponseWriter, r *http.Request)

// handlers/conversation.go
func GetConversation(w http.ResponseWriter, r *http.Request)
func SendMessage(w http.ResponseWriter, r *http.Request)
```

#### 3. Middleware Layer (`middlewares/`)

```go
// CheckSessionCookie(handler, requiresAuth)
// - Validates session cookie
// - requiresAuth=true: Protects authenticated routes
// - requiresAuth=false: Prevents logged-in users from public pages

// RateLimit(handler, duration)
// - Throttles requests per user per duration
```

#### 4. Model Layer (`models/`)
Data structures:

```go
type User struct {
    Id       int
    Nickname string
    Email    string
    Password string
    FirstName, LastName string
    Age      int
    Gender   string
}

type Post struct {
    Id        int
    UserId    int
    Nickname  string
    Created_at time.Time
    TimeAgo   string
    Title     string
    Text      string
    LikeCount, DislikeCount int
    IsLiked   int // 1:liked, 0:none, -1:disliked
    Comments  []Comment
    Categories []string
    Image     string
}

type Comment struct {
    Id         int
    UserId     int
    PostId     int
    Nickname   string
    Created_at time.Time
    TimeAgo    string
    Text       string
    LikeCount, DislikeCount int
    IsLiked    int
}
```

#### 5. Database Layer (`database/`)

```go
// init.go - Initialize database
func Init(refresh bool) error // Creates tables, loads schema

// Database operations:
db.Query()      // SELECT queries
db.QueryRow()   // Single row queries
db.Exec()       // INSERT/UPDATE/DELETE
```

#### 6. WebSocket Layer (`ws/`)

```go
// handlers/ws.go
func HandlerWs(w http.ResponseWriter, r *http.Request)
// Validates session, upgrades to WebSocket

// ws/client.go
type Client struct {
    id   string
    conn *websocket.Conn
}

func HandleClient(client *Client)
// Reads messages, broadcasts events

// ws/helpers.go
func BroadcastExcept(excludeId, event, data)
func BroadcastToUser(userId, event, data)
func Broadcast(event, data)
```

---

## Data Flow

### 1. Post Creation Flow

```
Frontend                        Backend                     Database
┌─────────┐                     
│ User    │                     
│ Clicks  │                     
│ Create  │                     
└────┬────┘                     
     │ Form Submission          
     ├──→ POST /api/posts/create                           
     │    ├─ Validate request body                         
     │    ├─ Extract userId from session cookie            
     │    ├─ Insert into POSTS table                       
     │                         ├──→ INSERT INTO posts(...)
     │    ├─ Insert categories in POST_CATEGORY            │
     │                         ├──→ INSERT INTO post_category
     │    ├─ Return created post                           │
     │←─── 200 OK + Post data   │                          │
     │ Real-time broadcast via WebSocket                   
     ├─→ ws.emit("new_post")                               
     │    └─→ BroadcastExcept(userId, "new_post", post)    
     │         ├─ All connected clients receive event      
     │         └─ Clients update UI with new post          
     │ Clear form                
     │ Show success toast        
     └─────────────────────────────────────────────────────
```

### 2. Comment Creation & Real-Time Update

```
Frontend                Backend                     Database
┌──────────┐            
│ User     │            
│ Types    │            
│ Comment  │            
└────┬─────┘            
     │ POST /api/comments/create
     ├──→ Validate text & postId
     │    ├─ Check session
     │    ├─ INSERT INTO comments
     │                 ├──→ DB stores comment
     │    ├─ Return comment data
     │←─── 200 OK
     │ Update UI
     │ Append comment to list
     │
     │ Optional: WebSocket broadcast
     │ (real-time if other users viewing)
```

### 3. Like/Dislike Flow

```
Frontend                Backend                 Database
┌──────────┐            
│ User     │            
│ Clicks   │            
│ Like Btn │            
└────┬─────┘            
     │ POST /api/posts/:id/like
     ├──→ Check existing reaction
     │    ├─ SELECT from POST_REACTIONS
     │                 ├──→ Check if exists
     │    ├─ If exists & different: DELETE old
     │                 ├──→ DELETE from reactions
     │    ├─ If new or different: INSERT new
     │                 ├──→ INSERT into reactions
     │    ├─ Return updated counts
     │←─── 200 OK + {likes, dislikes, isLike}
     │ Update UI immediately
     │ Update like count
     │ Toggle button state
```

### 4. Message Flow (WebSocket)

```
User A (Client 1)       WebSocket Hub           User B (Client 2)
┌────────────────┐      ┌──────────────┐       ┌────────────────┐
│ User A types   │      │              │       │ User B         │
│ message        │      │              │       │ waiting...     │
└────────┬───────┘      │              │       └────────────────┘
         │              │              │
         │ ws.send()    │              │
         ├─────────────→ BroadcastToUser
         │              │ (User B)     │
         │              ├─────────────→ emit "new_message"
         │              │              │
         │              │ [Also save]  │
         │              │              │
         │              │ INSERT INTO  │
         │              │ messages(...)│
         │              │              │
         │              │              │ Client receives
         │              │              │ & renders message
         │              │              │
         │              ← ws.on("new_message")
         │              │
         │ [Optional] ← UPDATE offline
         │ INSERT INTO   │
         │ conversations │
```

---

## Request/Response Cycle

### Standard HTTP Request Handling

```
1. Request arrives at main.go
   ↓
2. CORS middleware adds headers
   ↓
3. Max body size middleware checks size
   ↓
4. Router directs to handler
   ↓
5. RateLimit middleware checks request frequency
   ↓
6. CheckSessionCookie middleware validates auth
   ├─ If auth fails & requiresAuth=true → 401 Unauthorized
   ├─ If user logged in & auth=false (login page) → 409 Conflict (redirect)
   └─ Continue if OK
   ↓
7. Handler processes request
   ├─ Parse/validate request body
   ├─ Query/modify database
   ├─ Broadcast WebSocket events if needed
   └─ Return response
   ↓
8. Response sent to client
   ├─ 200 OK + JSON data on success
   ├─ 400+ error status + error message on failure
   └─ Client handles response
```

### Response Format

All endpoints return standardized JSON:

```json
{
  "status_code": 200,
  "message": "success or error description",
  "data": {
    // Actual response data
  }
}
```

---

## WebSocket Communication

### Connection Lifecycle

```
1. Client connects to ws://localhost:8080/ws
   ↓
2. Backend validates session cookie
   ├─ If invalid → HTTP 401
   └─ If valid → Extract userId
   ↓
3. Upgrade to WebSocket connection
   ↓
4. Create Client struct
   ├─ id: userId as string
   └─ conn: WebSocket connection
   ↓
5. Store client in Hub
   ├─ clients[userId] = client
   └─ Broadcast "client_connect" event
   ↓
6. Listen for messages
   ├─ Parse incoming message
   ├─ Route to appropriate handler
   └─ Broadcast or send to specific user
   ↓
7. On disconnect
   ├─ Remove from Hub
   └─ Broadcast "client_disconnect" event
```

### WebSocket Events Reference

#### Server-Sent Events

| Event | Data | Usage |
|-------|------|-------|
| `init` | `[userId1, userId2, ...]` | Initial online users list |
| `client_connect` | `userId` | User came online |
| `client_disconnect` | `userId` | User went offline |
| `new_message` | `{senderId, text, timestamp}` | Incoming message |
| `new_post` | `{postId, title, author}` | New post created |
| `force_logout` | `{}` | Force logout all tabs |

#### Client-Sent Events
```javascript
ws.send(JSON.stringify({
    type: "message",
    data: {
        recipientId: 123,
        text: "Hello!",
        conversationId: 456
    }
}))
```

---

## Security Architecture

### Authentication & Authorization

#### Session-Based Authentication
```
1. User logs in with email/nickname + password
2. Backend validates credentials
3. If valid: Create session
   ├─ Generate UUID
   ├─ Hash password with bcrypt (verification only)
   ├─ Store session in database
   ├─ Set 24-hour expiration
4. Return session_id as HTTP-only secure cookie
5. Client automatically sends cookie with every request
6. Backend validates cookie on protected endpoints
```

#### Protected Routes

| Category | Auth Required | Examples |
|----------|---------------|----------|
| Public (Guest) | No | /login, /register |
| Protected | Yes | /api/posts, /api/messages, /api/comments |
| User-Owned | Yes + ownership | DELETE /api/posts/:id (owner only) |

#### Middleware Flow

```go
CheckSessionCookie(handler, requiresAuth bool):
  ├─ No cookie
  │  ├─ requiresAuth=true → 401 (protected route)
  │  └─ requiresAuth=false → allow (public route)
  ├─ Cookie found
  │  ├─ Query database for session
  │  ├─ Session expired → 401, delete from DB
  │  ├─ Session valid
  │  │  ├─ requiresAuth=true → allow
  │  │  └─ requiresAuth=false → 409 (redirect to home)
  │  └─ DB error → 500
```

### Data Validation

#### Frontend Validation
- Form fields validated before submission
- Categories checked against allowed list
- Post/comment text length limits enforced
- Email format validated

#### Backend Validation
- All inputs re-validated server-side
- SQL injection prevention via parameterized queries
- XSS prevention via output sanitization
- CSRF protection via session-based auth

### Password Security

```
Registration:
  Input password → Validate length (6-20 chars)
                 → Hash with bcrypt
                 → Store hashed password

Login:
  Input password → Hash and compare with stored hash
                 → If match: Create session
                 → If mismatch: Reject
```

---

## Performance Optimizations

### Frontend Optimizations

1. **Lazy Loading of Posts**
   - Initial load: 15 posts
   - Infinite scroll: Load 15 more when near bottom
   - Throttled scroll event (200ms)

2. **Component Memoization**
   - Avoid re-rendering unchanged components
   - Dynamic list updates only modified items

3. **Debouncing & Throttling**
   - Scroll events: 200ms throttle
   - Search/filter: 300ms debounce
   - WebSocket reconnect: exponential backoff

4. **Efficient DOM Updates**
   - Use `insertAdjacentHTML()` for batch inserts
   - Minimize reflows/repaints
   - Event delegation for dynamic content

5. **Asset Optimization**
   - Single HTML file (SPA)
   - Minimized CSS with variables
   - Font loading optimized (Font Awesome via CDN)

### Backend Optimizations

1. **Database Indexing**
   ```sql
   CREATE INDEX idx_username ON users(nickname COLLATE NOCASE);
   -- Speeds up login queries
   ```

2. **Query Optimization**
   - Use DISTINCT for category queries
   - LEFT JOIN for optional data
   - LIMIT/OFFSET for pagination

3. **Connection Pooling**
   - SQLite handles concurrent requests
   - Go's http.Server provides worker pool

4. **Rate Limiting**
   - Per-user throttling
   - Prevents spam and abuse
   - Configurable by endpoint

5. **Caching Strategy**
   - Post reactions cached in client
   - Minimal database round-trips
   - WebSocket for real-time updates

### Network Optimizations

1. **HTTP Compression**
   - JSON responses compressed
   - CSS/JS served compressed

2. **Connection Reuse**
   - Keep-Alive enabled
   - Single WebSocket connection

3. **Request Batching**
   - Categories loaded once
   - Posts loaded in batches

4. **CORS Optimization**
   - Preflight requests minimized
   - Simple requests allowed

---

## Error Handling

### Frontend Error Handling

```javascript
try {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(data.message || "Request failed");
  }
  return data;
} catch (err) {
  showToast(err.message, "error");
  // Log error, retry logic, etc.
}
```

### Backend Error Handling

```go
// Structured error responses
utilities.WriteJSON(w, statusCode, messageString, data)

// Status codes:
// 200 - Success
// 400 - Bad request (validation error)
// 401 - Unauthorized (session invalid)
// 403 - Forbidden (not allowed)
// 404 - Not found
// 409 - Conflict (duplicate, state error)
// 429 - Too many requests (rate limited)
// 500 - Server error
```

---

## Deployment Architecture

### Development
```
Frontend: http://localhost:3000 (Node.js dev server)
Backend: http://localhost:8080 (Go server)
Database: ./forum.db (SQLite file)
```

### Production (Docker)
```
docker-compose up -d
├─ frontend service (port 3000)
├─ backend service (port 8080)
└─ shared volume for database
```

### Environment Configuration
```
Frontend:
  API_URL=http://api.example.com/api
  WS_URL=wss://api.example.com/ws

Backend:
  PORT=8080
  DB_PATH=/data/forum.db
  CORS_ORIGIN=https://app.example.com
  SESSION_TIMEOUT=24h
```


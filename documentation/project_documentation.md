# Real-Time Forum - Project Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Features](#features)
3. [Technology Stack](#technology-stack)
4. [System Architecture](#system-architecture)
5. [Core Features](#core-features)
6. [User Workflows](#user-workflows)
7. [Getting Started](#getting-started)

---

## Project Overview

**Real-Time Forum** is a modern, single-page web application built with JavaScript frontend and Go backend. It provides real-time communication, social features, and an interactive forum experience with WebSocket integration for instant notifications and user presence.

### Key Characteristics
- **Single Page Application (SPA)**: All navigation and page changes handled via JavaScript
- **Real-Time Features**: WebSocket integration for live updates and messaging
- **Session-Based Authentication**: Secure user authentication with cookie-based sessions
- **Responsive Design**: Works seamlessly on desktop and mobile devices
- **Brutalist Design**: Minimalist, editorial-focused UI aesthetic

---

## Features

### 1. Authentication & User Management
- **User Registration**: Complete sign-up with validation
- **User Login**: Flexible login with email or nickname
- **Session Management**: 24-hour server-side sessions with secure HTTP-only cookies
- **User Profiles**: View user information and public profiles
- **Online Status**: Real-time user presence tracking

### 2. Posts & Comments
- **Create Posts**: Users can create posts with title, content, and multiple categories
- **Browse Posts**: Feed displays all posts sorted by creation date (newest first)
- **Filter Posts**: 
  - By category (13+ categories available)
  - Posts created by current user
  - Posts liked by current user
- **Comments**: Users can comment on posts
- **Pagination**: Infinite scroll with dynamic post loading

### 3. Reactions System
- **Like/Dislike Posts**: Users can like or dislike posts
- **Like/Dislike Comments**: Users can rate comments
- **Reaction Tracking**: Real-time reaction count updates
- **User Preference Storage**: Tracks individual user reactions

### 4. Private Messaging
- **Real-Time Chat**: WebSocket-based instant messaging
- **Conversations**: Manage multiple conversations with different users
- **User Presence**: See online/offline status of contacts
- **Message History**: Load past messages with pagination
- **Last Message Display**: Shows preview of last message in conversation list

### 5. Real-Time Notifications
- **WebSocket Events**: Instant notifications for:
  - New posts
  - New messages
  - User online/offline status
  - Force logout alerts (multi-tab sync)
- **Live Updates**: Changes reflected instantly across all client tabs

---

## Technology Stack

### Frontend
- **Language**: JavaScript (ES6+)
- **Architecture**: Single Page Application (SPA)
- **UI Library**: Vanilla JavaScript + CSS
- **Real-Time Communication**: WebSocket API
- **Styling**: Custom CSS with CSS Variables (Brutalist design)
- **State Management**: Client-side JavaScript objects and localStorage

### Backend
- **Language**: Go (Golang)
- **Framework**: Standard `net/http` package
- **Database**: SQLite3
- **Real-Time**: Gorilla WebSocket
- **Authentication**: JWT stored in secure HTTP-only cookies
- **Middleware**: Custom authentication and rate limiting

### DevOps
- **Containerization**: Docker support
- **Port Configuration**: 
  - Frontend: `http://localhost:3000`
  - Backend: `http://localhost:8080`
  - WebSocket: `ws://localhost:8080/ws`

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (JavaScript SPA)               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Pages: Feed, Post, Chat, Login, Register, Error    │  │
│  │ Components: Post, Comment, Conversation, Header    │  │
│  │ Services: Router, Auth, WebSocket, Toast          │  │
│  └──────────────────────────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP/WebSocket
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    Backend (Go)                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Routes: /api/auth, /api/posts, /api/comments       │  │
│  │ Handlers: Auth, Posts, Comments, Messages, Chat    │  │
│  │ Middleware: Auth, Rate Limiting, CORS              │  │
│  │ WebSocket: Real-time messaging & notifications     │  │
│  └──────────────────────────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    Database (SQLite)                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Tables: Users, Sessions, Posts, Comments, etc      │  │
│  │ Schema: Foreign keys, constraints, indexes         │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Frontend Architecture

```
frontend/
├── index.html           # Single entry point
├── pages/               # Page components
│   ├── feed.js         # Main post feed
│   ├── post.js         # Single post detail view
│   ├── chat.js         # Messaging interface
│   ├── login.js        # Login form
│   ├── register.js     # Registration form
│   └── error.js        # Error pages
├── components/         # Reusable UI components
│   ├── Post.js        # Post display
│   ├── Comment.js     # Comment display
│   ├── Header.js      # Navigation header
│   └── ...
├── services/          # Core services
│   ├── router.js      # Client-side routing
│   ├── auth.js        # Authentication
│   ├── websocket.js   # WebSocket connection
│   └── toast.js       # Notifications
├── api/              # API client functions
│   ├── posts.js      # Post endpoints
│   ├── auth.js       # Auth endpoints
│   ├── comments.js   # Comment endpoints
│   └── conversations.js  # Messaging endpoints
└── styles/           # Stylesheets
```

### Backend Architecture

```
backend/
├── main.go           # Server entry point, middleware setup
├── routes/
│   └── routes.go     # Route definitions & rate limiting
├── handlers/         # HTTP request handlers
│   ├── auth.go       # Login, register, logout
│   ├── posts.go      # Post CRUD operations
│   ├── comments.go   # Comment operations
│   ├── reaction.go   # Like/dislike logic
│   ├── conversation.go  # Chat management
│   ├── ws.go         # WebSocket initialization
│   └── users.go      # User profile endpoints
├── middlewares/      # Middleware functions
│   ├── auth.go       # Session validation
│   └── rate_limit.go # Rate limiting
├── models/           # Data structures
│   ├── user.go
│   ├── post.go
│   └── comment.go
├── ws/              # WebSocket handling
│   ├── client.go    # Client connection management
│   └── helpers.go   # Broadcasting utilities
├── database/        # Database operations
│   ├── init.go      # Database initialization
│   ├── schema.sql   # Database schema
│   └── seed.go      # Sample data
└── utilities/       # Helper functions
```

---

## Core Features

### 1. Authentication Flow
1. User registers with email, nickname, password, and profile info
2. Backend validates and hashes password with bcrypt
3. Session created in database (24-hour expiration)
4. Session ID sent as secure HTTP-only cookie
5. Protected endpoints verify session validity
6. Logout clears session and cookie

### 2. Post Lifecycle
1. User creates post with title, content, and categories
2. Post stored in database with timestamp
3. Post displayed in feed sorted by creation date
4. Users can filter posts by:
   - Category selection
   - Own posts
   - Liked posts
5. Users can like/dislike, comment, or delete (if owner)
6. Real-time updates broadcast to all connected clients

### 3. Real-Time Messaging
1. User initiates conversation with another user
2. WebSocket connection established
3. Messages sent instantly via WebSocket
4. Messages stored in database for history
5. Conversation sorted by last message timestamp
6. New users appear in alphabetical order

### 4. WebSocket Events
| Event | Direction | Purpose |
|-------|-----------|---------|
| `init` | Server→Client | Send list of online users |
| `client_connect` | Server→Client | Notify user came online |
| `client_disconnect` | Server→Client | Notify user went offline |
| `new_message` | Server→Client | Incoming message in chat |
| `new_post` | Server→Client | New post created |
| `force_logout` | Server→Client | Multi-tab logout sync |

---

## User Workflows

### Registration & Login
```
Guest → Register Form → Validation → Account Created → Login
Login → Session Created → Cookie Set → Redirected Home → Feed
```

### Creating & Sharing Posts
```
Feed → Create Post Button → Form Modal → Enter Title/Content/Categories
Submit → Validation → Post Stored → Real-Time Update → Visible in Feed
```

### Interacting with Posts
```
Feed → See Post → Like/Dislike (instant count update)
                → Click Post → View Comments → Add Comment
                → Delete (if owner) → Post Removed
```

### Messaging Flow
```
Chat → Select User → Load Conversation History → Type Message
Send → Real-Time Update → Message Appears → Recipient Notified
```

---

## Getting Started

### Prerequisites
- Node.js (for frontend dev server)
- Go 1.16+ (for backend)
- SQLite3
- Docker (optional)

### Backend Setup
```bash
cd backend
go mod download
go run main.go          # Run with fresh database
go run main.go refresh  # Reset database with seed data
```

### Frontend Setup
```bash
cd frontend
npm install
npm start               # Runs on http://localhost:3000
```

### Environment Variables
Backend (.env):
```
DB_PATH=./forum.db
SESSION_SECRET=your-secret-key
CORS_ORIGIN=http://localhost:3000
```

### API Base URLs
- REST API: `http://localhost:8080/api`
- WebSocket: `ws://localhost:8080/ws`
- Frontend: `http://localhost:3000`

---

## Rate Limiting

| Endpoint | Limit | Purpose |
|----------|-------|---------|
| POST /api/login | 1 per 2s | Prevent brute force |
| POST /api/register | 1 per 2s | Prevent spam registration |
| GET /api/posts | 1 per 3s | Reduce server load |
| POST /api/posts/create | 1 per 3s | Prevent spam posts |
| POST/DELETE /api/posts/:id/:action | 1 per 250ms | Throttle reactions |
| POST /api/comments/create | 1 per 250ms | Throttle comments |
| POST /api/messages | 1 per 100ms | Allow rapid messaging |

---

## Key Design Decisions

1. **Single Page Application**: Eliminates page reloads for seamless UX
2. **WebSocket Integration**: Enables real-time features without polling
3. **Session-Based Auth**: Simpler than JWT, secure with HTTP-only cookies
4. **SQLite**: Lightweight, file-based, perfect for smaller deployments
5. **Vanilla JavaScript**: No framework overhead, full control over UI
6. **Brutalist CSS**: Fast-loading, accessible, editorial aesthetic

---

## Deployment

### Docker Deployment
```bash
docker-compose up -d
```

### Production Considerations
- Use HTTPS/WSS in production
- Implement proper CORS configuration
- Set secure cookie flags
- Enable database backups
- Monitor WebSocket connections
- Implement proper error logging

---

## Support & Documentation

For detailed information on specific components:
- [Architecture Documentation](./archi_doc.md)
- [Database Schema](./db.md)
- [API Documentation](./api_doc/)
  - [Authentication](./api_doc/authentification.md)
  - [Posts](./api_doc/posts.md)
  - [Comments](./api_doc/comments.md)
  - [Likes/Reactions](./api_doc/likes.md)
  - [Messaging](./api_doc/notification.md)
  - [WebSocket Events](./api_doc/websocket-events.md)
  - [Profile](./api_doc/profile.md)


# real-time-forum

Welcome to the Real-Time Forum Documentation! This is your comprehensive guide to understanding and working with the project.

## Quick Navigation

### 📚 For Everyone
- **[Project Documentation](./documentation/project_documentation.md)** - Start here! Overview of features, tech stack, and getting started

### 🏗️ For Developers
- **[Architecture Documentation](./documentation/archi_doc.md)** - System design, component breakdown, data flow, WebSocket communication
- **[Database Schema](./documentation/db.md)** - Complete database design, tables, relationships, and query patterns

### 🔌 For API Integration
- **[Authentication](./documentation/api_doc/authentification.md)** - Login, register, sessions, and security
- **[Posts API](./documentation/api_doc/posts.md)** - Create, retrieve, filter, and manage posts
- **[Comments API](./documentation/api_doc/comments.md)** - Comment operations and reactions
- **[Likes & Reactions](./documentation/api_doc/likes.md)** - Like/dislike system for posts and comments
- **[Real-Time Messaging](./documentation/api_doc/notification.md)** - WebSocket events, messaging, conversations
- **[User Profiles](./documentation/api_doc/profile.md)** - User information and presence tracking
- **[Admin API](./documentation/api_doc/admin.md)** - Admin operations (future)

---

## Documentation Structure

```
documentation/
├── README.md (this file)
├── project_documentation.md      ← START HERE
├── archi_doc.md                  ← System architecture
├── db.md                         ← Database schema
└── api_doc/
    ├── authentification.md
    ├── posts.md
    ├── comments.md
    ├── likes.md
    ├── notification.md
    ├── profile.md
    ├── admin.md
    └── real-time-forum.postman_collection.json
```

---

## Getting Started

### 1. **New to the Project?**
   → Read [project_documentation.md](./documentation/project_documentation.md)
   - Features overview
   - Technology stack
   - System architecture at high level
   - Getting started guide

### 2. **Want to Understand the Architecture?**
   → Read [archi_doc.md](./documentation/archi_doc.md)
   - Frontend/backend layer breakdown
   - Component interactions
   - Data flow diagrams
   - Security architecture
   - Performance optimizations

### 3. **Need Database Information?**
   → Read [db.md](./documentation/db.md)
   - All tables and relationships
   - Schema diagrams
   - Query patterns
   - Data integrity rules

### 4. **Building API Integration?**
   → Browse [api_doc/](./documentation/api_doc/)
   - Each endpoint documented separately
   - Request/response examples
   - Error handling
   - Rate limits

---

## Key Sections by Topic

### 🔐 Authentication & Security
- **Where:** [archi_doc.md > Security Architecture](./documentation/archi_doc.md#security-architecture)
- **Where:** [authentification.md](./documentation/api_doc/authentification.md)
- **Topics:** Session management, password hashing, middleware validation

### 📝 Posts & Content
- **Where:** [api_doc/posts.md](./documentation/api_doc/posts.md)
- **Where:** [db.md > POST_CATEGORY](./documentation/db.md#5-post_category)
- **Topics:** Creating, filtering, categorization

### 💬 Comments & Reactions
- **Where:** [api_doc/comments.md](./documentation/api_doc/comments.md)
- **Where:** [api_doc/likes.md](./documentation/api_doc/likes.md)
- **Where:** [db.md > COMMENT_REACTIONS](./documentation/db.md#8-comment_reactions)

### 💬 Messaging & Real-Time
- **Where:** [api_doc/notification.md](./documentation/api_doc/notification.md)
- **Where:** [archi_doc.md > WebSocket Communication](./documentation/archi_doc.md#websocket-communication)
- **Where:** [db.md > CONVERSATIONS & MESSAGES](./documentation/db.md)

### 👤 Users & Profiles
- **Where:** [api_doc/profile.md](./documentation/api_doc/profile.md)
- **Where:** [db.md > USERS](./documentation/db.md#1-users)

### 🏗️ Frontend Components
- **Where:** [archi_doc.md > Frontend Architecture](./documentation/archi_doc.md#frontend-architecture-layers)
- **Pages:** Feed, Post detail, Chat, Auth
- **Components:** Post, Comment, Message, Header

### 🗄️ Backend Handlers
- **Where:** [archi_doc.md > Backend Architecture](./documentation/archi_doc.md#backend-architecture-layers)
- **Handlers:** Auth, Posts, Comments, Messages, WebSocket

---

## Code Examples Quick Links

### Frontend
```javascript
// Routing
import { router } from './services/router.js';
await router.navigate('/feed');

// API Calls
const { getPosts, CreatePost } = await import('./api/posts.js');
const posts = await getPosts({limit: 15, offset: 0});

// WebSocket
import { ws } from './services/websocket.js';
ws.on('new_message', handleNewMessage);
```

### Backend
```go
// Routes
http.HandleFunc("/api/posts", CheckSessionCookie(GetPosts, true))

// Handlers
func GetPosts(w http.ResponseWriter, r *http.Request)

// Database
database.Database.Query("SELECT ...")
database.Database.Exec("INSERT INTO ...")
```

---

## API Endpoints Summary

### Authentication
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/register` | Create new account |
| POST | `/api/login` | Authenticate user |
| POST | `/api/logout` | End session |
| GET | `/api/me` | Get current user |

### Posts
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/api/posts` | List posts (with filters) |
| POST | `/api/posts/create` | Create post |
| GET | `/api/posts/{id}` | Get single post |
| POST | `/api/posts/{id}/like` | Like post |
| POST | `/api/posts/{id}/dislike` | Dislike post |
| DELETE | `/api/posts/{id}` | Delete post |

### Comments
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/comments/create` | Create comment |
| POST | `/api/comments/{id}/like` | Like comment |
| POST | `/api/comments/{id}/dislike` | Dislike comment |
| DELETE | `/api/comments/{id}` | Delete comment |

### Messaging
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/api/conversations` | List conversations |
| GET | `/api/conversation/{id}` | Get conversation messages |
| POST | `/api/messages` | Send message |
| WS | `/ws` | Real-time WebSocket |

### Users
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/api/users/{id}` | Get user profile |

---

## Rate Limits Reference

| Endpoint | Limit | Purpose |
|----------|-------|---------|
| /api/login | 1/2s | Prevent brute force |
| /api/register | 1/2s | Prevent spam |
| /api/posts | 1/3s | Load management |
| /api/posts/create | 1/3s | Prevent spam posts |
| /api/posts/:id/:action | 1/250ms | Throttle reactions |
| /api/comments/create | 1/250ms | Throttle comments |
| /api/messages | 1/100ms | Allow messaging |

---

## Database Tables Quick Reference

| Table | Purpose |
|-------|---------|
| USERS | User accounts |
| SESSIONS | Active user sessions |
| POSTS | User posts |
| CATEGORY | Post categories |
| POST_CATEGORY | Post-category mapping |
| COMMENTS | Post comments |
| POST_REACTIONS | Post likes/dislikes |
| COMMENT_REACTIONS | Comment likes/dislikes |
| CONVERSATIONS | User conversations |
| MESSAGES | Private messages |

---

## Architecture Diagrams

### High-Level System
```
┌─────────────────────────┐
│   Browser (SPA)         │
│  JavaScript + WebSocket │
└────────────┬────────────┘
             │ HTTP/WS
┌────────────▼────────────┐
│   Go Backend Server     │
│  REST API + WebSocket   │
└────────────┬────────────┘
             │ SQL
┌────────────▼────────────┐
│  SQLite Database        │
│  Relational Data        │
└─────────────────────────┘
```

### Frontend Layers
```
Pages (Feed, Post, Chat)
         ↓
Components (Post, Comment, Message)
         ↓
Services (Router, Auth, WebSocket)
         ↓
API Layer (HTTP clients)
```

### Backend Layers
```
Routes (Endpoints + Middleware)
   ↓
Handlers (Business Logic)
   ↓
Models (Data Structures)
   ↓
Database (SQLite)
```

---

## Deployment

### Development
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- WebSocket: `ws://localhost:8080/ws`

### Docker
```bash
docker-compose up -d
```

---

## Troubleshooting

### Common Issues
See specific API docs for error handling:
- Auth errors → [authentification.md](./documentation/api_doc/authentification.md)
- Post errors → [posts.md](./documentation/api_doc/posts.md)
- WebSocket errors → [notification.md](./documentation/api_doc/notification.md)

### Performance Tips
See [archi_doc.md > Performance Optimizations](./documentation/archi_doc.md#performance-optimizations)

### Security Concerns
See [archi_doc.md > Security Architecture](./documentation/archi_doc.md#security-architecture)

---

## Contributing to Documentation

When updating features:
1. Update relevant API doc file
2. Update architecture doc if design changes
3. Update database doc if schema changes
4. Update project overview if features change

### Documentation Checklist
- [ ] Endpoint documented
- [ ] Request/response examples provided
- [ ] Error cases covered
- [ ] Rate limits specified
- [ ] Related concepts linked

---

## Contact & Support

For questions about:
- **Features** → See project_documentation.md
- **Architecture** → See archi_doc.md
- **Database** → See db.md
- **Specific API** → See api_doc/*.md

---

**Last Updated:** 2026-06-20
**Documentation Version:** 2.0
**Project Version:** Latest

## Authors

- [herrabba](https://github.com/hamzaerrhh)
- [thakkou](https://github.com/thakkou)
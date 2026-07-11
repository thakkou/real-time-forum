# Database Schema & Design

## Table of Contents
1. [Database Overview](#database-overview)
2. [Schema Diagram](#schema-diagram)
3. [Table Definitions](#table-definitions)
4. [Relationships](#relationships)
5. [Constraints & Indexes](#constraints--indexes)
6. [Query Patterns](#query-patterns)
7. [Data Integrity](#data-integrity)

---

## Database Overview

### Technology
- **DBMS**: SQLite3
- **Type**: Relational Database
- **File-Based**: Single `.db` file (portable)
- **Connection**: No server overhead

### Key Characteristics
- Supports concurrent reads
- Transactions for data integrity
- Foreign key constraints
- Index support for query optimization
- AUTOINCREMENT for primary keys

---

## Schema Diagram

```
┌──────────────────┐
│     USERS        │
│──────────────────│
│ id (PK)          │
│ nickname (UNIQUE)│
│ firstname        │
│ lastname         │
│ age              │
│ gender           │
│ email (UNIQUE)   │
│ password         │
│ last_seen        │
└────┬─────────────┘
     │ (1:N)
     ├────────────────────────────┬─────────────────────┬──────────────┐
     │                            │                     │              │
     ▼                            ▼                     ▼              ▼
┌──────────────┐        ┌──────────────┐    ┌──────────────┐  ┌──────────────┐
│  SESSIONS    │        │    POSTS     │    │  COMMENTS    │  │ CONVERSATION │
│──────────────│        │──────────────│    │──────────────│  │──────────────│
│ id (PK/UUID) │        │ id (PK)      │    │ id (PK)      │  │ id (PK)      │
│ expires_at   │        │ user_id (FK) │    │ user_id (FK) │  │ user1_id (FK)│
│ user_id (FK) │        │ created_at   │    │ post_id (FK) │  │ user2_id (FK)│
└──────────────┘        │ title        │    │ created_at   │  │ created_at   │
                        │ text         │    │ text         │  │ last_message_at
                        │ image        │    └──────────────┘  └──────────────┘
                        └─────────┬────┘           ▲                 │
                                  │ (1:N)         │ (1:N)            │ (1:N)
                                  │               │                  │
                        ┌─────────┴───────┐      │                  ▼
                        │                 │      │         ┌──────────────┐
                        ▼                 ▼      │         │  MESSAGES    │
                ┌───────────────────┐    │      │         │──────────────│
                │   POST_CATEGORY   │    │      │         │ id (PK)      │
                │───────────────────│    │      │         │ sender_id(FK)│
                │ post_id (FK/PK)   │    │      │         │ conv_id (FK) │
                │ category_id(FK/PK)│    │      │         │ created_at   │
                └─────────┬─────────┘    │      │         │ text         │
                          │ (N:M)        │      │         │ is_read      │
                          │              │      │         └──────────────┘
                          ▼              │      │
                ┌───────────────────┐    │      │
                │    CATEGORY       │    │      │
                │───────────────────│    │      │
                │ id (PK)           │    │      │
                │ name (UNIQUE)     │    │      │
                └───────────────────┘    │      │
                                         │      │
                ┌────────────────────────┘      │
                │                               │
                ▼                               ▼
        ┌──────────────────┐      ┌──────────────────────┐
        │  POST_REACTIONS  │      │ COMMENT_REACTIONS    │
        │──────────────────│      │──────────────────────│
        │ user_id (FK/PK)  │      │ user_id (FK/PK)      │
        │ post_id (FK/PK)  │      │ comment_id (FK/PK)   │
        │ is_like (1/-1)   │      │ is_like (1/-1)       │
        └──────────────────┘      └──────────────────────┘
```

---

## Table Definitions

### 1. USERS

Stores user account information.

```sql
CREATE TABLE USERS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT NOT NULL UNIQUE,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    age INTEGER NOT NULL,
    gender TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT,
    last_seen DATETIME
);

CREATE INDEX idx_username ON users(nickname COLLATE NOCASE);
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique user ID |
| nickname | TEXT | NOT NULL, UNIQUE | Display name (case-insensitive search) |
| firstname | TEXT | NOT NULL | User's first name |
| lastname | TEXT | NOT NULL | User's last name |
| age | INTEGER | NOT NULL | User age |
| gender | TEXT | NOT NULL | Gender (male/female/other) |
| email | TEXT | NOT NULL, UNIQUE | Email address for login |
| password | TEXT | - | Bcrypt hashed password |
| last_seen | DATETIME | - | Last activity timestamp |

**Example:**
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
  "last_seen": "2026-06-20T10:30:00"
}
```

---

### 2. SESSIONS

Tracks active user sessions for authentication.

```sql
CREATE TABLE SESSIONS (
    id TEXT PRIMARY KEY UNIQUE,
    expires_at DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| id | TEXT | PRIMARY KEY, UNIQUE | UUID session token |
| expires_at | DATETIME | NOT NULL | Session expiration time |
| user_id | INTEGER | NOT NULL, FK | Reference to USERS.id |

**Behavior:**
- Expires after 24 hours
- Cascade delete when user deleted
- Returned as HTTP-only cookie to client

**Example:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2026-06-21T10:00:00",
  "user_id": 1
}
```

---

### 3. POSTS

Stores user-created posts.

```sql
CREATE TABLE POSTS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    title TEXT,
    text TEXT,
    image TEXT,
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE
);
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique post ID |
| user_id | INTEGER | NOT NULL, FK | Author of post |
| created_at | DATETIME | NOT NULL | Creation timestamp |
| title | TEXT | - | Post title (max 255 chars) |
| text | TEXT | - | Post content (max 1000 chars) |
| image | TEXT | - | Image URL/path |

**Example:**
```json
{
  "id": 4,
  "user_id": 1,
  "created_at": "2026-06-18T13:51:55.859",
  "title": "My First Go Post",
  "text": "This is a test post created from Postman",
  "image": "/uploads/post-image.jpg"
}
```

---

### 4. CATEGORY

List of available post categories.

```sql
CREATE TABLE CATEGORY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

INSERT OR IGNORE INTO CATEGORY (name) VALUES 
('General'), ('Lifestyle'), ('Health & Fitness'), ('Travel'),
('Food & Cooking'), ('Education'), ('Business'), ('Finance'),
('Entertainment'), ('Sports'), ('Personal Dev'), ('Culture'), ('News');
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Category ID |
| name | TEXT | NOT NULL, UNIQUE | Category name |

**Available Categories:** (13 total)
- General
- Lifestyle
- Health & Fitness
- Travel
- Food & Cooking
- Education
- Business
- Finance
- Entertainment
- Sports
- Personal Dev
- Culture
- News

---

### 5. POST_CATEGORY

Maps posts to categories (many-to-many relationship).

```sql
CREATE TABLE POST_CATEGORY (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES POSTS(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES CATEGORY(id) ON DELETE CASCADE
);
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| post_id | INTEGER | PRIMARY KEY, FK | Reference to POSTS.id |
| category_id | INTEGER | PRIMARY KEY, FK | Reference to CATEGORY.id |

**Purpose:** Associates posts with multiple categories.

**Example:**
```
Post #4 is associated with:
├─ category_id 1 (General)
├─ category_id 6 (Education)
└─ category_id 2 (Lifestyle)
```

---

### 6. COMMENTS

Stores comments on posts.

```sql
CREATE TABLE COMMENTS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    text TEXT,
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES POSTS (id) ON DELETE CASCADE
);
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Comment ID |
| user_id | INTEGER | NOT NULL, FK | Comment author |
| post_id | INTEGER | NOT NULL, FK | Parent post |
| created_at | DATETIME | NOT NULL | Creation timestamp |
| text | TEXT | - | Comment text (max 1000 chars) |

**Example:**
```json
{
  "id": 1,
  "user_id": 1,
  "post_id": 4,
  "created_at": "2026-06-18T16:56:57.973",
  "text": "Great post! Very informative."
}
```

---

### 7. POST_REACTIONS

Tracks likes and dislikes on posts.

```sql
CREATE TABLE POST_REACTIONS (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    is_like INTEGER NOT NULL DEFAULT 1 CHECK (is_like IN (-1, 1)),
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES POSTS (id) ON DELETE CASCADE
);
```

**Columns:**
| Column | Type | Constraint | Description |
|--------|------|-----------|------------|
| user_id | INTEGER | NOT NULL, FK | User reacting |
| post_id | INTEGER | NOT NULL, FK | Post being reacted to |
| is_like | INTEGER | NOT NULL, CHECK | 1 for like, -1 for dislike |

**Note:** Composite primary key (user_id, post_id) ensures one reaction per user per post.

**Example:**
```
User #1 likes Post #4:
  {user_id: 1, post_id: 4, is_like: 1}

User #2 dislikes Post #4:
  {user_id: 2, post_id: 4, is_like: -1}

User #1 changes to dislike:
  {user_id: 1, post_id: 4, is_like: -1} (updated)
```

---

### 8. COMMENT_REACTIONS

Tracks likes and dislikes on comments.

```sql
CREATE TABLE COMMENT_REACTIONS (
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    is_like INTEGER NOT NULL DEFAULT 1 CHECK (is_like IN (-1, 1)),
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES COMMENTS (id) ON DELETE CASCADE
);
```

**Columns:** Same structure as POST_REACTIONS but for comments.

---

### 9. CONVERSATIONS

Tracks conversations between users.

```sql
-- Note: Check backend for exact schema
CREATE TABLE CONVERSATIONS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user1_id INTEGER NOT NULL,
    user2_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    last_message TEXT,
    last_message_at DATETIME,
    FOREIGN KEY (user1_id) REFERENCES USERS (id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES USERS (id) ON DELETE CASCADE
);
```

**Purpose:** Tracks direct conversations between two users.

---

### 10. MESSAGES

Stores private messages in conversations.

```sql
-- Note: Check backend for exact schema
CREATE TABLE MESSAGES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    text TEXT NOT NULL,
    is_read INTEGER DEFAULT 0,
    FOREIGN KEY (conversation_id) REFERENCES CONVERSATIONS (id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES USERS (id) ON DELETE CASCADE
);
```

**Purpose:** Stores individual messages within conversations.

---

## Relationships

### One-to-Many Relationships

#### Users → Sessions
```
One user can have multiple sessions
└─ A user logged in on different devices/tabs
```

#### Users → Posts
```
One user can create many posts
└─ Each post has one author
```

#### Users → Comments
```
One user can write many comments
└─ Each comment has one author
```

#### Posts → Comments
```
One post can have many comments
└─ Comments are tied to specific posts
```

#### Posts → Post Reactions
```
One post can have many reactions
└─ Multiple users can like/dislike same post
```

#### Comments → Comment Reactions
```
One comment can have many reactions
└─ Multiple users can like/dislike same comment
```

#### Conversations → Messages
```
One conversation can have many messages
└─ Messages stored chronologically
```

### Many-to-Many Relationships

#### Posts ↔ Categories (via POST_CATEGORY)
```
One post can have many categories
One category can have many posts

Example:
Post #1: [Education, Lifestyle, General]
Post #2: [General, Business]
Post #3: [Education, Business, Finance]

Category "Education": [Post #1, Post #3]
```

---

## Constraints & Indexes

### Primary Key Constraints
- Ensures uniqueness and fast lookup
- AUTOINCREMENT for sequential IDs
- UUID for session tokens

### Unique Constraints
```sql
-- Users
UNIQUE (nickname)  -- Case-insensitive search
UNIQUE (email)     -- Email login support

-- Category
UNIQUE (name)      -- Prevent duplicate categories

-- Sessions
UNIQUE (id)        -- Token uniqueness
```

### Foreign Key Constraints
```
ON DELETE CASCADE: When parent deleted, children deleted too

Example: Delete user
├─ Sessions deleted
├─ Posts deleted
│  ├─ Comments deleted
│  ├─ Post_Reactions deleted
│  └─ Post_Category entries deleted
├─ Comments deleted (direct)
├─ Comment_Reactions deleted
└─ Messages deleted
```

### Check Constraints
```sql
-- POST_REACTIONS & COMMENT_REACTIONS
CHECK (is_like IN (-1, 1))
-- Ensures only valid reaction values: 1=like, -1=dislike
```

### Indexes

```sql
-- Username search optimization (case-insensitive)
CREATE INDEX idx_username ON users(nickname COLLATE NOCASE);
```

**Why Indexes Matter:**
- Login queries search by nickname/email (frequent)
- Index speeds up LIKE queries
- COLLATE NOCASE makes search case-insensitive

---

## Query Patterns

### 1. User Authentication

**Login Query:**
```sql
SELECT id, nickname, password
FROM users
WHERE email = ? OR nickname = ?
LIMIT 1;
```
Uses: nickname index for fast lookup

**Session Validation:**
```sql
SELECT expires_at FROM sessions WHERE id = ?;
```

---

### 2. Post Retrieval

**Get All Posts (Latest First):**
```sql
SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text
FROM posts p
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?;
```

**Get Posts by Category:**
```sql
SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text
FROM posts p
LEFT JOIN post_category pc ON p.id = pc.post_id
LEFT JOIN category c ON pc.category_id = c.id
WHERE c.name IN (?, ?, ...)
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?;
```

**Get User's Own Posts:**
```sql
SELECT id, user_id, created_at, title, text
FROM posts
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
```

**Get Posts User Liked:**
```sql
SELECT DISTINCT p.id, p.user_id, p.created_at, p.title, p.text
FROM posts p
JOIN post_reactions pr ON p.id = pr.post_id
WHERE pr.user_id = ? AND pr.is_like = 1
ORDER BY p.created_at DESC
LIMIT ? OFFSET ?;
```

---

### 3. Post Enrichment

**Get Post Categories:**
```sql
SELECT c.name
FROM category c
JOIN post_category pc ON c.id = pc.category_id
WHERE pc.post_id = ?
ORDER BY c.name;
```

**Get Post Reactions:**
```sql
-- Like count
SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND is_like = 1;

-- Dislike count
SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND is_like = -1;

-- User's reaction
SELECT is_like FROM post_reactions WHERE user_id = ? AND post_id = ?;
```

**Get Post Comments:**
```sql
SELECT id, user_id, post_id, created_at, text
FROM comments
WHERE post_id = ?
ORDER BY created_at DESC
LIMIT ?;
```

---

### 4. User Presence (Chat)

**Get Online Users:**
```sql
-- From WebSocket Hub (in-memory)
SELECT userId FROM hub.clients
```

**Get Conversations (Sorted by Last Message):**
```sql
SELECT c.id, c.last_message, c.last_message_at, u.nickname
FROM conversations c
JOIN users u ON (
    (c.user1_id = ? AND c.user2_id = u.id) OR
    (c.user2_id = ? AND c.user1_id = u.id)
)
WHERE c.user1_id = ? OR c.user2_id = ?
ORDER BY c.last_message_at DESC
LIMIT ? OFFSET ?;
```

**Get Messages in Conversation:**
```sql
SELECT id, sender_id, created_at, text, is_read
FROM messages
WHERE conversation_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
```

---

## Data Integrity

### Cascade Delete Strategy

**Deleting a User triggers:**
```
User deleted
├─ All SESSIONS deleted
├─ All POSTS deleted (which cascades to)
│  ├─ COMMENTS deleted
│  ├─ POST_REACTIONS deleted
│  └─ POST_CATEGORY entries deleted
├─ All COMMENTS (direct) deleted
├─ All COMMENT_REACTIONS deleted
├─ All CONVERSATIONS deleted (which cascades to)
│  └─ MESSAGES deleted
└─ All MESSAGES sent by user deleted
```

**Deleting a Post triggers:**
```
Post deleted
├─ COMMENTS deleted
├─ POST_REACTIONS deleted
└─ POST_CATEGORY entries deleted
```

**Deleting a Comment triggers:**
```
Comment deleted
└─ COMMENT_REACTIONS deleted
```

### Transaction Safety

**Critical Operations use Transactions:**
```go
tx, err := database.Database.Begin()
defer tx.Rollback()

// Multiple operations
tx.Exec("INSERT INTO posts...")
tx.Exec("INSERT INTO post_category...")

tx.Commit() // All succeed or all fail
```

---

## Data Size Estimates

### Typical Data Volumes

| Table | Typical Rows | Notes |
|-------|------------|-------|
| USERS | 100-1000 | One per user |
| SESSIONS | 10-100 | 1-24 hours lifespan |
| POSTS | 1000-10000 | Grows with user activity |
| COMMENTS | 5000-50000 | Multiple per post |
| POST_REACTIONS | 5000-50000 | Each user reacts to posts |
| COMMENT_REACTIONS | 2000-20000 | Less frequent |
| CONVERSATIONS | 100-1000 | Sparse for active users |
| MESSAGES | 1000-100000 | Depends on chat volume |

### Storage Example
```
100 users × 500 posts × 100 KB/post = ~5 GB
(Realistic: < 1 GB for small-medium deployments)
```

---

## Optimization Tips

### 1. Query Optimization
- Always use LIMIT/OFFSET for pagination
- Use indices on frequently filtered columns
- Avoid SELECT * - specify needed columns

### 2. Index Strategy
```sql
-- Add indices on:
-- Frequently filtered columns
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);

-- Foreign keys
CREATE INDEX idx_comments_post_id ON comments(post_id);

-- Reactions lookups
CREATE INDEX idx_post_reactions_user ON post_reactions(user_id);
```

### 3. Maintenance
- Regular database backups
- Vacuum database periodically (SQLite)
- Monitor query performance
- Archive old messages if needed


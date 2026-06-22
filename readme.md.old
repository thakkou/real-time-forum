# 01Forum

A full-stack web forum built with Go, SQLite, and Docker. Users can register, create posts with categories, comment, and like or dislike content.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
  - [Run with Docker](#run-with-docker)
  - [Run Locally](#run-locally)
- [Usage](#usage)
- [Database](#database)
- [API Routes](#api-routes)
- [Authors](#authors)

---

## Features

- **Authentication** — Register and login with email or username. Passwords are encrypted with `bcrypt`. Sessions use UUID cookies with expiry.
- **Posts** — Registered users can create posts with a title, text, and one or more categories.
- **Comments** — Registered users can comment on any post.
- **Reactions** — Registered users can like or dislike posts and comments. Counts are visible to all users.
- **Filtering** — Filter posts by category, by posts you created, or by posts you liked.
- **Rate Limiting** — All POST routes are rate-limited per IP to prevent spam.
- **Error Handling** — HTTP 400, 401, 403, 404, 405, 429, and 500 errors are all handled with a dedicated error page.

---

## Project Structure

```
├── database/
│   ├── init.go         # Opens SQLite DB and runs schema
│   └── schema.sql      # All CREATE TABLE statements + seed categories
├── forum-api/
│   ├── comment.go      # Comment queries
│   ├── post.go         # Post queries and filters
│   ├── reaction.go     # Like/dislike logic
│   └── session.go      # Session deletion
├── handlers/
│   ├── api.go          # CreatePost, CreateComment, PostResolver, CommentResolver
│   ├── error.go        # HandleError
│   ├── forum.go        # Main forum page + filtering
│   ├── login.go        # Login handler
│   ├── logout.go       # Logout handler
│   ├── register.go     # Register handler
│   ├── statichandler.go
│   └── template.go     # RenderTemplate helper
├── helper/
│   └── GetUserId.go    # Get user ID from session cookie
├── middlewares/
│   ├── auth.go         # Session cookie validation middleware
│   └── Ratelimit.go    # Rate limiting middleware
├── routing/
│   └── rountig.go      # All route registrations
├── static/
│   ├── script.js       # Like/dislike fetch calls
│   └── style.css
├── templates/
│   ├── error.html
│   ├── index.html      # Main forum page
│   ├── login.html
│   └── register.html
├── main.go
├── Dockerfile
└── docker.sh
```

---

## Requirements

- [Docker](https://docs.docker.com/get-docker/) — recommended
- Or: Go 1.21+, GCC (for sqlite3 CGo)

---

## Getting Started

### Run with Docker

```bash
# 1. Build the image
docker build -t forum .

# 2. Run the container
docker run -p 8080:8080 forum
```

Or use the provided script:

```bash
bash docker.sh
```

Then open your browser at [http://localhost:8080](http://localhost:8080)

---

### Run Locally

```bash
# Install dependencies
go mod download

# Run the server
go run .
```

> Requires GCC installed for the `go-sqlite3` driver (CGo).

---

## Usage

| Action                                  | Who can do it         |
| --------------------------------------- | --------------------- |
| View posts and comments                 | Everyone              |
| See like/dislike counts                 | Everyone              |
| Register / Login                        | Everyone              |
| Create a post                           | Registered users only |
| Create a comment                        | Registered users only |
| Like or dislike a post/comment          | Registered users only |
| Delete own post or comment              | Registered users only |
| Filter by category                      | Everyone              |
| Filter by "my posts" or "posts I liked" | Registered users only |

### OAuth

Create a .env file and set provider's client id and secret:

* Google:
  1. create project & app (Google developer console)
  2. create client (id & secret)
  3. set google provider's variables: GOOGLE_CLIENT_ID & GOOGLE_CLIENT_SECRET

* Github:
  1. Create an OAuth app (Settings > developer settings > oauth apps)
  2. Generate a new client secret
  3. set github provider's variables: GITHUB_CLIENT_ID & GITHUB_CLIENT_SECRET

---

## Database

SQLite database stored at `./database/forum.db`.

Tables: `USERS`, `SESSIONS`, `POSTS`, `CATEGORY`, `POST_CATEGORY`, `COMMENTS`, `POST_REACTIONS`, `COMMENT_REACTIONS`, `rate_limits`

Seed categories loaded automatically on startup: General, Lifestyle, Health & Fitness, Travel, Food & Cooking, Education, Business, Finance, Entertainment, Sports, Personal Dev, Culture, News.

---

## API Routes

| Method | Route                        | Description         | Auth required |
| ------ | ---------------------------- | ------------------- | ------------- |
| GET    | `/`                          | Forum home page     | No            |
| GET    | `/login`                     | Login page          | No            |
| POST   | `/login`                     | Submit login        | No            |
| GET    | `/register`                  | Register page       | No            |
| POST   | `/register`                  | Submit registration | No            |
| POST   | `/logout`                    | Logout              | Yes           |
| POST   | `/posts/create`              | Create a post       | Yes           |
| POST   | `/comments/create`           | Create a comment    | Yes           |
| POST   | `/api/posts/{id}/like`       | Like a post         | Yes           |
| POST   | `/api/posts/{id}/dislike`    | Dislike a post      | Yes           |
| POST   | `/api/posts/{id}/delete`     | Delete a post       | Yes           |
| POST   | `/api/comments/{id}/like`    | Like a comment      | Yes           |
| POST   | `/api/comments/{id}/dislike` | Dislike a comment   | Yes           |
| POST   | `/api/comments/{id}/delete`  | Delete a comment    | Yes           |

---

## Allowed Packages

- Standard Go library
- `github.com/mattn/go-sqlite3` — SQLite driver
- `golang.org/x/crypto/bcrypt` — Password hashing
- `github.com/google/uuid` — Session IDs

---

## Authors

- thakkou - [Github](https://github.com/thakkou)
- halhyane - [Github](https://github.com/Houssam-Alhyane)
- erezzoug - [Github](https://github.com/elmehdi-rezoug)
- herraba - [Github](https://github.com/hamzaerrhh)
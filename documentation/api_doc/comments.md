# Comments API Documentation

## Table of Contents
1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Endpoints](#endpoints)
4. [Error Handling](#error-handling)

---

## Overview

The Comments API allows authenticated users to:
- Create comments on posts
- Like/dislike comments
- Delete their own comments
- Retrieve comments for a post

---

## Authentication

All comment endpoints require a valid session.

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

### Create Comment

Creates a new comment on a post.

#### Endpoint
```http
POST /api/comments/create
```

#### Rate Limit
1 request per 250 milliseconds

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| postId | integer/string | Yes | ID of the post |
| text | string | Yes | Comment text (max 1000 chars) |

#### Example Request

```json
{
  "postId": 4,
  "text": "Great post! Very informative."
}
```

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "comment created",
  "data": {
    "id": 1,
    "userId": 1,
    "postId": 4,
    "nickname": "john_doe",
    "createdAt": "2026-06-18T16:56:57.973",
    "text": "Great post! Very informative.",
    "likeCount": 0,
    "dislikeCount": 0,
    "isLiked": 0
  }
}
```

#### Error Responses

**Missing or Empty Text:**
```json
{
  "status_code": 400,
  "message": "Comment cannot be empty"
}
```

**Invalid Post ID:**
```json
{
  "status_code": 400,
  "message": "Invalid post ID"
}
```

**Text Exceeds Max Length:**
```json
{
  "status_code": 400,
  "message": "Comment cannot exceed 1000 characters"
}
```

**Post Not Found:**
```json
{
  "status_code": 404,
  "message": "Post not found"
}
```

**Unauthorized (Invalid Session):**
```json
{
  "status_code": 401,
  "message": "Invalid or expired session"
}
```

---

### Like Comment

Likes a comment.

#### Endpoint
```http
POST /api/comments/{id}/like
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Comment ID |

#### Rate Limit
1 request per 500 milliseconds

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "liked",
  "data": {
    "likes": 1,
    "dislikes": 0,
    "isLike": 1
  }
}
```

#### Behavior
- First like: Creates new reaction
- Like already exists: Removes reaction (neutral)
- Dislike exists: Replaces with like

---

### Dislike Comment

Dislikes a comment.

#### Endpoint
```http
POST /api/comments/{id}/dislike
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Comment ID |

#### Rate Limit
1 request per 500 milliseconds

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "disliked",
  "data": {
    "likes": 0,
    "dislikes": 1,
    "isLike": -1
  }
}
```

#### Behavior
- First dislike: Creates new reaction
- Dislike already exists: Removes reaction (neutral)
- Like exists: Replaces with dislike

---

### Delete Comment

Deletes a comment (owner only).

#### Endpoint
```http
DELETE /api/comments/{id}/delete
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Comment ID |

#### Rate Limit
1 request per 500 milliseconds

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "deleted",
  "data": {}
}
```

#### Error Responses

**Comment Not Found:**
```json
{
  "status_code": 404,
  "message": "Comment not found"
}
```

**Not Authorized (Not Comment Owner):**
```json
{
  "status_code": 403,
  "message": "Unauthorized"
}
```

---

## Resolver Pattern

Comments use a flexible resolver pattern for actions:

```http
POST /api/comments/{id}/like      # Like comment
POST /api/comments/{id}/dislike   # Dislike comment
DELETE /api/comments/{id}/delete  # Delete comment
```

---

## Error Handling

### Standard Error Response Format

```json
{
  "status_code": <HTTP_STATUS>,
  "message": "Error description",
  "data": null
}
```

### Common HTTP Status Codes

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Bad request (validation error) |
| 401 | Unauthorized (invalid session) |
| 403 | Forbidden (not authorized) |
| 404 | Not found |
| 429 | Too many requests (rate limited) |
| 500 | Server error |

---

## Reaction Logic

### Like/Dislike State Machine

```
┌─────────────┐
│   NEUTRAL   │
│ (isLike=0)  │
└──────┬──────┘
       │
   ╔═══╩═══╗
   ║       ║
   ▼       ▼
┌──────┐  ┌────────┐
│ LIKE │  │DISLIKE │
│  (1) │  │  (-1)  │
└──────┘  └────────┘
   ▲       ▲
   │       │
   ╠═══════╣
   │       │
┌──────────────┐
│ Toggle Again │
│(Remove Reaction)
└──────────────┘

Example Sequence:
1. User clicks Like → isLike = 1
2. User clicks Like again → isLike = 0 (neutral)
3. User clicks Dislike → isLike = -1
4. User clicks Like → isLike = 1
```

---

## Example Flow

### Creating and Reacting to a Comment

```
1. Create Comment
   POST /api/comments/create
   {postId: 4, text: "Nice!"}
   ← 200 OK {id: 1}

2. Like the Comment
   POST /api/comments/1/like
   ← 200 OK {likes: 1, isLike: 1}

3. Change to Dislike
   POST /api/comments/1/dislike
   ← 200 OK {dislikes: 1, isLike: -1}

4. Remove Reaction
   POST /api/comments/1/dislike (again)
   ← 200 OK {dislikes: 0, isLike: 0}

5. Delete Comment
   DELETE /api/comments/1
   ← 200 OK
```


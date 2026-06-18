# Posts Service

The Posts Service handles post creation, retrieval, filtering, reactions, and deletion.

---

# Overview

Posts allow authenticated users to publish content categorized under one or more categories.

Features include:

* Create posts
* Delete own posts
* Retrieve all posts
* Filter posts by categories
* Retrieve posts created by the authenticated user
* Retrieve posts liked by the authenticated user
* Like / dislike posts

---

# Authentication

Most post actions require a valid session.

### Required Cookie

| Name       | Type   | Required |
| ---------- | ------ | -------- |
| session_id | Cookie | Yes      |

Example:

```http
Cookie: session_id=550e8400-e29b-41d4-a716-446655440000
```

---

# Endpoints

## Create Post

Creates a new post.

### Endpoint

```http
POST /api/posts/create
```

### Authentication

Requires a valid `session_id` cookie.

### Request Body

| Field      | Type         | Required |
| ---------- | ------------ | -------- |
| title      | string       | Yes      |
| text       | string       | Yes      |
| categories | string[]     | Yes      |
| image      | string (URL) | No       |

### Example Request

```json
{
  "title": "My First Go Post",
  "text": "This is a test post created from Postman using JSON.",
  "categories": [
    "General",
    "Lifestyle",
    "Education"
  ],
  "image": "https://example.com/image.jpg"
}
```

### Validation Rules

| Field      | Rule                         |
| ---------- | ---------------------------- |
| title      | 1 - 255 characters           |
| text       | 1 - 1000 characters          |
| categories | At least 1 category required |
| image      | Optional valid image URL     |

---

### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "post created successfully",
  "data": {
    "categories": [
      "General",
      "Lifestyle",
      "Education"
    ],
    "post_id": 4,
    "text": "This is a test post created from Postman using JSON.",
    "title": "My First Go Post",
    "user_id": 1
  }
}
```

---

### Error Responses

#### Missing Title or Text

```json
{
  "status_code": 400,
  "message": "Title and text cannot be empty"
}
```

#### Missing Categories

```json
{
  "status_code": 400,
  "message": "At least one category must be selected"
}
```

#### Invalid Category

```json
{
  "status_code": 400,
  "message": "Invalid category: Education"
}
```

#### Invalid Session

```json
{
  "status_code": 401,
  "message": "Invalid or expired session"
}
```

#### Internal Server Error

```json
{
  "status_code": 500,
  "message": "Could not create post"
}
```

---

## Get Posts

Returns posts ordered by creation date (newest first).

### Endpoint

```http
GET /api/posts/getPosts
```

### Query Parameters

| Parameter      | Type     | Required | Description                          |
| -------------- | -------- | -------- | ------------------------------------ |
| categories     | string[] | No       | Filter by categories                 |
| my-liked-posts | boolean  | No       | Return posts liked by current user   |
| my-creat-posts | boolean  | No       | Return posts created by current user |

---

### Example Requests

#### Get All Posts

```http
GET /api/posts/getPosts
```

#### Filter By Categories

```http
GET /api/posts/getPosts?categories=Education&categories=Lifestyle
```

#### Get My Posts

```http
GET /api/posts/getPosts?my-creat-posts=true
```

#### Get Posts I Liked

```http
GET /api/posts/getPosts?my-liked-posts=true
```

#### Combined Filters

```http
GET /api/posts/getPosts?my-liked-posts=true&categories=Education
```

---

### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "posts fetched successfully",
  "data": [
    {
      "Id": 4,
      "UserId": 1,
      "Nickname": "john_doe",
      "Created_at": "2026-06-18T13:51:55.859009543+01:00",
      "TimeAgo": "47 minutes ago",
      "Title": "My First Go Post",
      "Text": "This is a test post created from Postman using JSON.",
      "LikeCount": 0,
      "DislikeCount": 0,
      "IsLiked": 0,
      "Comments": null,
      "Categories": [
        "Education",
        "General",
        "Lifestyle"
      ],
      "Image": ""
    }
  ]
}
```

---

### Error Responses

#### Invalid Request

```json
{
  "status_code": 400,
  "message": "Bad request"
}
```

#### Internal Server Error

```json
{
  "status_code": 500,
  "message": "failed to get posts"
}
```

---

## Delete Post

Deletes a post owned by the authenticated user.

### Endpoint

```http
DELETE /api/posts/{id}/delete
```

### Authentication

Requires a valid `session_id` cookie.

### URL Parameters

| Parameter | Type    | Required |
| --------- | ------- | -------- |
| id        | integer | Yes      |

---

### Example Request

```http
DELETE /api/posts/1/delete
```

---

### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "post deleted successfully",
  "data": {
    "deleted": true,
    "post_id": 1
  }
}
```

---

### Error Responses

#### Invalid Post ID

```json
{
  "status_code": 400,
  "message": "Invalid post ID"
}
```

#### Not Logged In

```json
{
  "status_code": 401,
  "message": "Not logged in"
}
```

#### Invalid Session

```json
{
  "status_code": 401,
  "message": "Invalid session"
}
```

#### Post Not Found

```json
{
  "status_code": 403,
  "message": "post not found"
}
```

#### Not Your Post

```json
{
  "status_code": 403,
  "message": "not your post"
}
```

---

## Like Post

Adds a like reaction to a post.

### Endpoint

```http
POST /api/posts/{id}/like
```

### Authentication

Requires a valid `session_id` cookie.

### Success Response

```json
{
  "status_code": 200,
  "message": "post liked successfully"
}
```

---

## Dislike Post

Adds a dislike reaction to a post.

### Endpoint

```http
POST /api/posts/{id}/dislike
```

### Authentication

Requires a valid `session_id` cookie.

### Success Response

```json
{
  "status_code": 200,
  "message": "post disliked successfully"
}
```

---

# Standard Response Format

All API responses follow the same structure:

```json
{
  "status_code": 200,
  "message": "success message",
  "data": {}
}
```

| Field       | Description                         |
| ----------- | ----------------------------------- |
| status_code | HTTP status code                    |
| message     | Response message                    |
| data        | Additional response data (optional) |

---

# Post Flow

```text
Login
   ↓
Create Post
   ↓
Post Stored
   ↓
Retrieve Posts
   ↓
Like / Dislike
   ↓
Delete Post
```

---


### Allowed Categories

Posts must contain at least one category selected from the following list:

| Available Categories |
| -------------------- |
| General              |
| Lifestyle            |
| Health & Fitness     |
| Travel               |
| Food & Cooking       |
| Education            |
| Business             |
| Finance              |
| Entertainment        |
| Sports               |
| Personal Dev         |
| Culture              |
| News                 |

### Example

```json
{
  "title": "My First Go Post",
  "text": "This is a test post created from Postman using JSON.",
  "categories": [
    "General",
    "Education",
    "Lifestyle"
  ]
}
```

### Invalid Category Example

```json
{
  "title": "My First Go Post",
  "text": "This is a test post created from Postman using JSON.",
  "categories": [
    "Gaming"
  ]
}
```

### Error Response

```json
{
  "status_code": 400,
  "message": "Invalid category: Gaming"
}
```




# Security Notes

* Only authenticated users can create posts.
* Only post owners can delete their posts.
* Session authentication is required for protected endpoints.
* Categories are validated before association.
* User IDs are derived from the active session and cannot be supplied by clients.
* Posts are returned ordered by newest first.

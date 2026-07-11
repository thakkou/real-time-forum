# Likes & Reactions API Documentation

## Table of Contents
1. [Overview](#overview)
2. [Like/Dislike System](#likedislike-system)
3. [Post Reactions](#post-reactions)
4. [Error Handling](#error-handling)

---

## Overview

The Likes & Reactions API allows authenticated users to express their opinions on posts and comments through like and dislike reactions.

### Key Features
- Toggle like/dislike reactions
- Real-time reaction count updates
- One reaction per user per item
- Track user's own reaction state

---

## Like/Dislike System

### Reaction States

Each user can have one of three states per post/comment:

| State | Value | Meaning |
|-------|-------|---------|
| Neutral | 0 | No reaction |
| Like | 1 | User likes item |
| Dislike | -1 | User dislikes item |

### Reaction Toggle Behavior

```
User Perspective:
┌─────────────────────────┐
│ User clicks Like Button │
└────────┬────────────────┘
         │
    ┌────▼─────┐
    │ State?    │
    └────┬─────┘
    ┌────┴────────┬──────────┬──────────┐
    │             │          │          │
    ▼             ▼          ▼          ▼
┌────────┐  ┌────────┐ ┌─────────┐ ┌────────┐
│Neutral │  │  Like  │ │Dislike  │ │Like    │
│   ↓    │  │   ↓    │ │   ↓     │ │   ↓    │
│ Like   │  │Neutral │ │Like     │ │Neutral │
└────────┘  └────────┘ └─────────┘ └────────┘
```

---

## Post Reactions

### Like Post

Likes a post. Clicking like again removes the reaction.

#### Endpoint
```http
POST /api/posts/{id}/like
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Post ID |

#### Rate Limit
1 request per 250 milliseconds

#### Request
No body required (uses session cookie for user identification)

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "liked",
  "data": {
    "likes": 5,
    "dislikes": 2,
    "isLike": 1
  }
}
```

#### Response Data

| Field | Type | Description |
|-------|------|-------------|
| likes | integer | Total likes on post |
| dislikes | integer | Total dislikes on post |
| isLike | integer | User's current reaction (1, -1, or 0) |

#### Behavior

**First Like:**
- State changes from Neutral (0) to Like (1)
- Like count incremented

**Second Click (Toggle Off):**
- State changes from Like (1) to Neutral (0)
- Like count decremented

**Change from Dislike:**
- State changes from Dislike (-1) to Like (1)
- Like count incremented
- Dislike count decremented

#### Example Sequence

```
Initial: {likes: 4, dislikes: 1, isLike: 0}

Click Like:
POST /api/posts/4/like
← {likes: 5, dislikes: 1, isLike: 1}

Click Like Again:
POST /api/posts/4/like
← {likes: 4, dislikes: 1, isLike: 0}
```

---

### Dislike Post

Dislikes a post. Clicking dislike again removes the reaction.

#### Endpoint
```http
POST /api/posts/{id}/dislike
```

#### URL Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Post ID |

#### Rate Limit
1 request per 250 milliseconds

#### Success Response

**Status:** `200 OK`

```json
{
  "status_code": 200,
  "message": "disliked",
  "data": {
    "likes": 4,
    "dislikes": 2,
    "isLike": -1
  }
}
```

#### Behavior

**First Dislike:**
- State changes from Neutral (0) to Dislike (-1)
- Dislike count incremented

**Second Click (Toggle Off):**
- State changes from Dislike (-1) to Neutral (0)
- Dislike count decremented

**Change from Like:**
- State changes from Like (1) to Dislike (-1)
- Dislike count incremented
- Like count decremented

---

### Comment Reactions

Comment reactions work identically to post reactions:

**Like Comment:**
```http
POST /api/comments/{id}/like
```

**Dislike Comment:**
```http
POST /api/comments/{id}/dislike
```

All behavior and response formats are the same as post reactions.

---

## Frontend Integration

### Reaction Display & Updating

The frontend displays reactions in real-time:

```javascript
// Example: Like button click
const likeBtn = document.querySelector('.like-btn');
likeBtn.addEventListener('click', async () => {
  try {
    const res = await PostResolver({id: postId, type: 'like'});
    
    if (res?.message === 'liked') {
      // Update UI with new counts and state
      likeCount.innerText = res.data.likes;
      dislikeCount.innerText = res.data.dislikes;
      
      // Update button visual state
      likeBtn.classList.toggle('active');
      dislikeBtn.classList.remove('active');
    }
  } catch (err) {
    showToast(err.message, 'error');
  }
});
```

### Visual Feedback

**Active Reaction Button:**
```css
.like-btn.active {
  border-color: #4ade80;  /* Green for like */
  color: #4ade80;
}

.dislike-btn.active {
  border-color: #ff4545;  /* Red for dislike */
  color: #ff4545;
}
```

---

## Reaction Counting Logic

### Database Storage

```sql
POST_REACTIONS table:
user_id (PK) | post_id (PK) | is_like
     1       |      4       |   1     ← User 1 likes
     2       |      4       |   1     ← User 2 likes
     3       |      4       |  -1     ← User 3 dislikes
```

### Count Calculation

```sql
-- Get like count
SELECT COUNT(*) FROM post_reactions 
WHERE post_id = ? AND is_like = 1;

-- Get dislike count
SELECT COUNT(*) FROM post_reactions 
WHERE post_id = ? AND is_like = -1;

-- Get user's reaction
SELECT is_like FROM post_reactions 
WHERE user_id = ? AND post_id = ?;
```

---

## Error Handling

### Error Response Format

```json
{
  "status_code": <STATUS>,
  "message": "Error description",
  "data": null
}
```

### Common Errors

#### Post Not Found
```json
{
  "status_code": 404,
  "message": "Post not found"
}
```

#### Unauthorized (Invalid Session)
```json
{
  "status_code": 401,
  "message": "Invalid or expired session"
}
```

#### Rate Limited
```json
{
  "status_code": 429,
  "message": "Too many requests"
}
```

#### Server Error
```json
{
  "status_code": 500,
  "message": "Could not update reaction"
}
```

---

## Real-Time Behavior

### Scenario: Multiple Users

```
User A's Screen          Database              User B's Screen
Post: Likes: 3          
      Dislikes: 1       ┌─────────────┐       Post: Likes: 3
                        │ POST_ID: 4  │       Dislikes: 1
User A clicks Like      │ Likes: 3    │
  → Send request        │ Dislikes: 1 │
    POST /api/posts/4/like               User B viewing same post

    ← Response:                           (Sees old counts)
    {likes: 4, isLike: 1}

    Updates UI:                           
    Likes: 4                              (User B's screen
    Button becomes active                  not automatically
                                           updated - needs refresh
                                           or WebSocket)
```

**Note:** Currently, like/dislike updates are shown to the reacting user immediately, but other users see updates when:
1. They refresh the page
2. They navigate back to feed
3. WebSocket broadcast sends notification (if implemented)

---

## Best Practices

### Frontend
1. **Optimistic Updates**: Update UI immediately, revert on error
2. **Loading States**: Show loading indicator during request
3. **Error Handling**: Display error toast on failure
4. **Visual Feedback**: Clear active state indication

### Backend
1. **Idempotency**: Multiple identical requests have same effect
2. **Rate Limiting**: Prevent spam reactions
3. **Data Validation**: Verify post/comment exists before updating
4. **Transaction Safety**: Atomic reaction updates

### Performance
1. **Batch Updates**: In future, batch multiple reactions
2. **Caching**: Cache reaction counts with invalidation
3. **Query Optimization**: Use indexed columns
4. **Connection Pooling**: Reuse database connections


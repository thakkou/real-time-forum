-- USERS
CREATE TABLE IF NOT EXISTS USERS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- SESSIONS
CREATE TABLE IF NOT EXISTS SESSIONS (
    id TEXT PRIMARY KEY UNIQUE, -- uuid
    expires_at DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- POSTS
CREATE TABLE IF NOT EXISTS POSTS (
    id         INTEGER  NOT NULL UNIQUE,
    user_id    INTEGER  NOT NULL,
    created_at DATETIME NOT NULL, -- ~ 1GB
    title      TEXT     NULL    ,
    text       TEXT     NULL    ,
    PRIMARY KEY (id AUTOINCREMENT),
    FOREIGN KEY (user_id) REFERENCES USERS (id)
);
---CATEGORY 
CREATE TABLE IF NOT EXISTS CATEGORY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);
---category post 
CREATE TABLE IF NOT EXISTS POST_CATEGORY (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,

    PRIMARY KEY (post_id, category_id),

    FOREIGN KEY (post_id) REFERENCES POSTS(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES CATEGORY(id) ON DELETE CASCADE
);

-- COMMENTS
CREATE TABLE IF NOT EXISTS COMMENTS (
    id         INTEGER  NOT NULL UNIQUE,
    user_id    INTEGER  NOT NULL,
    post_id    INTEGER  NOT NULL,
    created_at DATETIME NOT NULL,
    text       TEXT     NULL    ,
    PRIMARY KEY (id AUTOINCREMENT),
    FOREIGN KEY (user_id) REFERENCES USERS (id),
    FOREIGN KEY (post_id) REFERENCES POSTS (id)
);

-- POST REACTIONS
CREATE TABLE IF NOT EXISTS POST_REACTIONS (
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  is_like INTEGER NOT NULL DEFAULT 1,
  FOREIGN KEY (user_id) REFERENCES USERS (id),
  FOREIGN KEY (post_id) REFERENCES POSTS (id)
);
-- reactions are unique by combination of both user_id and post_id

-- COMMENT REACTIONS
CREATE TABLE IF NOT EXISTS COMMENT_REACTIONS (
  user_id INTEGER NOT NULL,
  comment_id INTEGER NOT NULL,
  is_like INTEGER NOT NULL DEFAULT 1,
  FOREIGN KEY (user_id) REFERENCES USERS (id),
  FOREIGN KEY (comment_id) REFERENCES COMMENTS (id)
);

--Rate Limits
CREATE TABLE IF NOT EXISTS rate_limits (
    ip TEXT NOT NULL,
    route TEXT NOT NULL,
    last_request DATETIME NOT NULL,
    PRIMARY KEY (ip, route)
);
INSERT OR IGNORE INTO CATEGORY (name) VALUES 
('General'),
('Lifestyle'),
('Health & Fitness'),
('Travel'),
('Food & Cooking'),
('Education'),
('Business'),
('Finance'),
('Entertainment'),
('Sports'),
('Personal Dev'),
('Culture'),
('News');
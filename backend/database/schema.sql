-- USERS
CREATE TABLE IF NOT EXISTS USERS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT NOT NULL UNIQUE,
 

    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    age INTEGER NOT NULL,
    gender TEXT NOT NULL,

    email TEXT NOT NULL UNIQUE,
    password TEXT ,
    last_seen DATETIME
);

ALTER TABLE USERS ADD COLUMN last_seen DATETIME;

CREATE INDEX IF NOT EXISTS idx_username  ON users(nickname COLLATE NOCASE);
-- SESSIONS
CREATE TABLE IF NOT EXISTS SESSIONS (
    id TEXT PRIMARY KEY UNIQUE, -- uuid
    expires_at DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- POSTS
CREATE TABLE IF NOT EXISTS POSTS (
    id         INTEGER  NOT NULL UNIQUE,
    user_id    INTEGER  NOT NULL,
    created_at DATETIME NOT NULL, -- ~ 1GB
    title      TEXT     NULL    ,
    text       TEXT     NULL    ,
    image      TEXT     NULL, 
    PRIMARY KEY (id AUTOINCREMENT),
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE
);

---CATEGORY 
CREATE TABLE IF NOT EXISTS CATEGORY (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
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
    FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES POSTS (id) ON DELETE CASCADE 
);

-- POST REACTIONS
CREATE TABLE IF NOT EXISTS POST_REACTIONS (
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  is_like INTEGER NOT NULL DEFAULT 1 CHECK (is_like IN (-1, 1)),
  FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES POSTS (id) ON DELETE CASCADE 
);
-- reactions are unique by combination of both user_id and post_id

-- COMMENT REACTIONS
CREATE TABLE IF NOT EXISTS COMMENT_REACTIONS (
  user_id INTEGER NOT NULL,
  comment_id INTEGER NOT NULL,
  is_like INTEGER NOT NULL DEFAULT 1 CHECK (is_like IN (-1, 1)), -- 1 for like / -1 for dislike
  FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE ,
  FOREIGN KEY (comment_id) REFERENCES COMMENTS (id) ON DELETE CASCADE 
);

--Rate Limits
CREATE TABLE IF NOT EXISTS rate_limits (
    ip TEXT NOT NULL,
    route TEXT NOT NULL,
    last_request DATETIME NOT NULL,
    PRIMARY KEY (ip, route)
);
--this for private messages
CREATE TABLE IF NOT EXISTS MESSAGES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    conversation_id INTEGER NOT NULL,

    sender_id INTEGER NOT NULL,

    text TEXT NOT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    is_read INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (conversation_id) REFERENCES CONVERSATIONS(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES USERS(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS CONVERSATIONS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    user1_id INTEGER NOT NULL,
    user2_id INTEGER NOT NULL,

    last_message TEXT,
    last_message_at DATETIME,

    user1_last_read_message_id INTEGER,
    user2_last_read_message_id INTEGER,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user1_id, user2_id),

    FOREIGN KEY (user1_id) REFERENCES USERS(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES USERS(id) ON DELETE CASCADE
);


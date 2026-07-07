package database

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Nickname  string
	FirstName string
	LastName  string
	Age       int
	Gender    string
	Email     string
}

func RefreshAndSeed(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// --------------------------------------------------
	// USERS
	// --------------------------------------------------

	users := []seedUser{
		{"john_doe", "John", "Doe", 25, "male", "john_doe@example.com"},
		{"john_doe1", "John", "Doe", 26, "male", "john_doe1@example1.com"},
		{"john_doe2", "John", "Doe", 27, "male", "john_doe2@example2.com"},
		{"john_doe3", "John", "Doe", 28, "male", "john_doe3@example3.com"},
		{"john_doe4", "John", "Doe", 29, "male", "john_doe4@example4.com"},
		{"jane_doe", "Jane", "Doe", 24, "female", "jane_doe@example.com"},
		{"alice_smith", "Alice", "Smith", 30, "female", "alice@example.com"},
		{"bob_jones", "Bob", "Jones", 32, "male", "bob@example.com"},
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
	INSERT INTO USERS
	(nickname, firstname, lastname, age, gender, email, password, last_seen)
	VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, u := range users {
		_, err := stmt.Exec(
			u.Nickname,
			u.FirstName,
			u.LastName,
			u.Age,
			u.Gender,
			u.Email,
			string(passwordHash),
		)
		if err != nil {
			return err
		}
	}

	// --------------------------------------------------
	// SESSIONS
	// --------------------------------------------------

	_, err = tx.Exec(`
	INSERT INTO SESSIONS (id, expires_at, user_id)
	VALUES
	('sess_1', datetime('now','+1 day'), 1),
	('sess_2', datetime('now','+1 day'), 2),
	('sess_3', datetime('now','+1 day'), 3)
	`)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// POSTS
	// --------------------------------------------------

	_, err = tx.Exec(`
INSERT INTO POSTS (user_id, created_at, title, text, image)
VALUES
(1, datetime('now'), 'Hello World', 'My first post', NULL),
(2, datetime('now'), 'Travel vibes', 'I love Morocco!', NULL),
(3, datetime('now'), 'Fitness update', 'Gym every day 💪', NULL),
(4, datetime('now'), 'Food time', 'Tagine is amazing', NULL),
(5, datetime('now'), 'Tech talk', 'SQLite is awesome', NULL),

(1, datetime('now'), 'Morning Coffee', 'Coffee before coding ☕', NULL),
(2, datetime('now'), 'Beach Day', 'Spent the afternoon at the beach.', NULL),
(3, datetime('now'), 'Learning Go', 'Concurrency is really interesting.', NULL),
(4, datetime('now'), 'Movie Night', 'Watched an amazing sci-fi movie.', NULL),
(5, datetime('now'), 'Football Match', 'What a great game yesterday!', NULL),

(6, datetime('now'), 'Weekend Plans', 'Thinking about hiking.', NULL),
(7, datetime('now'), 'Photography', 'Captured a beautiful sunset.', NULL),
(8, datetime('now'), 'Cooking', 'Made homemade pizza today.', NULL),
(1, datetime('now'), 'Music Time', 'Listening to my favorite playlist.', NULL),
(2, datetime('now'), 'Book Review', 'Just finished a fantastic novel.', NULL),

(3, datetime('now'), 'Coding Tips', 'Always write clean code.', NULL),
(4, datetime('now'), 'Gaming', 'Reached a new high score!', NULL),
(5, datetime('now'), 'Nature Walk', 'Fresh air feels amazing.', NULL),
(6, datetime('now'), 'Study Session', 'Preparing for tomorrow''s exam.', NULL),
(7, datetime('now'), 'Photography Tips', 'Golden hour is the best.', NULL),

(8, datetime('now'), 'Fitness Goal', 'Ran 10 kilometers today.', NULL),
(1, datetime('now'), 'Healthy Food', 'Salads can actually taste great.', NULL),
(2, datetime('now'), 'Road Trip', 'Planning a drive across the country.', NULL),
(3, datetime('now'), 'Programming', 'Interfaces make code flexible.', NULL),
(4, datetime('now'), 'Cats', 'My cat slept all day.', NULL),

(5, datetime('now'), 'Dogs', 'Took my dog to the park.', NULL),
(6, datetime('now'), 'Music Festival', 'Canot wait for the weekend.', NULL),
(7, datetime('now'), 'Learning SQL', 'Joins finally make sense!', NULL),
(8, datetime('now'), 'Basketball', 'Played with friends tonight.', NULL),
(1, datetime('now'), 'Gardening', 'Planted some tomatoes.', NULL),

(2, datetime('now'), 'Cloud Computing', 'Trying out Docker today.', NULL),
(3, datetime('now'), 'Web Development', 'CSS can be tricky.', NULL),
(4, datetime('now'), 'Travel Diary', 'Visited Chefchaouen.', NULL),
(5, datetime('now'), 'Photography Gear', 'Bought a new lens.', NULL),
(6, datetime('now'), 'Morning Run', 'Started the day with exercise.', NULL),

(7, datetime('now'), 'AI Discussion', 'Machine learning is fascinating.', NULL),
(8, datetime('now'), 'Chess', 'Won three games in a row.', NULL),
(1, datetime('now'), 'Rust vs Go', 'Both have their strengths.', NULL),
(2, datetime('now'), 'Open Source', 'Made my first contribution.', NULL),
(3, datetime('now'), 'Reading', 'Trying to read more books.', NULL),

(4, datetime('now'), 'Family Dinner', 'Had a wonderful evening.', NULL),
(5, datetime('now'), 'UI Design', 'Simple designs are often the best.', NULL),
(6, datetime('now'), 'Linux', 'Switched my laptop to Linux.', NULL),
(7, datetime('now'), 'Backend Development', 'Working with APIs today.', NULL),
(8, datetime('now'), 'Weekend Relax', 'Time to recharge!', NULL)
`)
	if err != nil {
		fmt.Println("post creation err")
		return err
	}

	// --------------------------------------------------
	// POST CATEGORY
	// --------------------------------------------------

	_, err = tx.Exec(`
INSERT INTO POST_CATEGORY (post_id, category_id)
VALUES
(1,1),
(2,4),
(3,3),
(4,5),
(5,7),

(6,2),
(7,4),
(8,6),
(9,3),
(10,5),

(11,1),
(12,7),
(13,6),
(14,2),
(15,4),

(16,3),
(17,5),
(18,1),
(19,2),
(20,7),

(21,3),
(22,6),
(23,4),
(24,2),
(25,1),

(26,5),
(27,6),
(28,3),
(29,7),
(30,2),

(31,4),
(32,6),
(33,1),
(34,5),
(35,7),

(36,3),
(37,2),
(38,4),
(39,6),
(40,1),

(41,5),
(42,7),
(43,3),
(44,2),
(45,4)
`)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// COMMENTS
	// --------------------------------------------------

	_, err = tx.Exec(`
	INSERT INTO COMMENTS (user_id, post_id, created_at, text)
	VALUES
	(2,1,datetime('now'),'Nice post!'),
	(3,1,datetime('now'),'Welcome 👋'),
	(1,2,datetime('now'),'Thanks!'),
	(4,3,datetime('now'),'Keep going!'),
	(5,4,datetime('now'),'Yummy 😋')
	`)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// POST REACTIONS
	// --------------------------------------------------

	_, err = tx.Exec(`
    INSERT OR IGNORE INTO POST_REACTIONS (user_id, post_id, is_like)
    VALUES
    (2,1,1),
    (3,1,1),
    (4,1,-1),
    (1,2,1),
    (5,3,1)
    `)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// COMMENT REACTIONS
	// --------------------------------------------------

	_, err = tx.Exec(`
	INSERT INTO COMMENT_REACTIONS (user_id, comment_id, is_like)
	VALUES
	(1,1,1),
	(3,1,1),
	(2,2,1),
	(4,3,-1)
	`)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// CONVERSATIONS
	// --------------------------------------------------

	_, err = tx.Exec(`
	INSERT INTO CONVERSATIONS
	(user1_id, user2_id, last_message, last_message_at)
	VALUES
	(1,2,'Hey!',datetime('now')),
	(2,3,'What''s up?',datetime('now')),
	(3,4,'Hello 👋',datetime('now'))
	`)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// MESSAGES
	// --------------------------------------------------

	_, err = tx.Exec(`
	INSERT INTO MESSAGES
	(conversation_id, sender_id, text, created_at, is_read)
	VALUES
	(1,1,'Hey John1!',datetime('now'),1),
	(1,2,'Hey John!',datetime('now'),1),
	(2,2,'How are you?',datetime('now'),0),
	(2,3,'Good you?',datetime('now'),0),
	(3,3,'Hello Bob!',datetime('now'),1)
	`)
	if err != nil {
		return err
	}

	// --------------------------------------------------
	// RATE LIMITS
	// --------------------------------------------------

	_, err = tx.Exec(`
	INSERT INTO rate_limits (ip, route, last_request)
	VALUES
	('127.0.0.1','/login',datetime('now')),
	('127.0.0.1','/posts',datetime('now'))
	`)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	return nil
}

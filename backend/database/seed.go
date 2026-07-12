package database

import (
	"database/sql"
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

		{"emma_wilson", "Emma", "Wilson", 22, "female", "emma@example.com"},
		{"liam_brown", "Liam", "Brown", 31, "male", "liam@example.com"},
		{"olivia_taylor", "Olivia", "Taylor", 27, "female", "olivia@example.com"},
		{"noah_miller", "Noah", "Miller", 35, "male", "noah@example.com"},
		{"ava_davis", "Ava", "Davis", 26, "female", "ava@example.com"},
		{"ethan_white", "Ethan", "White", 29, "male", "ethan@example.com"},
		{"mia_clark", "Mia", "Clark", 24, "female", "mia@example.com"},
		{"lucas_hall", "Lucas", "Hall", 33, "male", "lucas@example.com"},
		{"sophia_lewis", "Sophia", "Lewis", 28, "female", "sophia@example.com"},
		{"james_walker", "James", "Walker", 36, "male", "james@example.com"},
		{"charlotte_young", "Charlotte", "Young", 23, "female", "charlotte@example.com"},
		{"henry_king", "Henry", "King", 34, "male", "henry@example.com"},
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		// fmt.Println("err users")
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
(6, datetime('now'), 'Weekend Plans', 'Going hiking this weekend.', NULL),
(7, datetime('now'), 'Morning Coffee', 'Nothing beats coffee at sunrise.', NULL),
(8, datetime('now'), 'Coding', 'Go is becoming my favorite language!', NULL),
(9, datetime('now'), 'Photography', 'Captured an amazing sunset today.', NULL),
(10, datetime('now'), 'Football', 'Great match last night!', NULL),
(11, datetime('now'), 'Reading', 'Finished a fantastic novel.', NULL),
(12, datetime('now'), 'Movies', 'Any sci-fi recommendations?', NULL),
(13, datetime('now'), 'Cooking', 'Homemade pizza tonight!', NULL),
(14, datetime('now'), 'Gaming', 'Reached Diamond rank today!', NULL),
(15, datetime('now'), 'Music', 'Listening to some classic rock.', NULL),
(16, datetime('now'), 'Nature', 'Beautiful walk in the forest.', NULL),
(17, datetime('now'), 'Programming', 'Building a chat application in Go.', NULL),
(18, datetime('now'), 'Travel', 'Dreaming about visiting Japan.', NULL),
(19, datetime('now'), 'Pets', 'My cat finally likes me 😂', NULL),
(20, datetime('now'), 'Life Update', 'Started learning Docker today.', NULL),
(1, datetime('now'), 'Go Tips', 'Interfaces are more powerful than they look.', NULL),
(3, datetime('now'), 'Workout', 'Leg day is always the hardest.', NULL),
(5, datetime('now'), 'Database', 'SQLite is perfect for small projects.', NULL),
(7, datetime('now'), 'Breakfast', 'Pancakes and coffee!', NULL),
(9, datetime('now'), 'Photography Gear', 'Thinking about buying a new lens.', NULL),
(12, datetime('now'), 'Series', 'Watching a new mystery show.', NULL),
(15, datetime('now'), 'Concert', 'Canot wait for next weekend!', NULL),
(18, datetime('now'), 'Adventure', 'Road trips are the best.', NULL),
(20, datetime('now'), 'Learning', 'Today I practiced SQL joins.', NULL)
`)
	if err != nil {
		// fmt.Println("err posts")

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
(7,5),
(8,7),
(9,6),
(10,8),
(11,9),
(12,10),
(13,5),
(14,7),
(15,11),
(16,6),
(17,7),
(18,4),
(19,12),
(20,2),
(21,7),
(22,3),
(23,7),
(24,5),
(25,6),
(26,10),
(27,11),
(28,4),
(29,2);
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
	INSERT INTO POST_REACTIONS (user_id, post_id, is_like)
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

package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func Init(refresh bool) error {
	var err error

	if refresh {
		// Close existing connection if any
		if Database != nil {
			Database.Close()
		}

		// Delete the database
		_ = os.Remove("./database/forum.db")
	}

	Database, err = sql.Open("sqlite3", "./database/forum.db?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("can't open/create forum.db: %v", err)
	}

	if err := Database.Ping(); err != nil {
		return fmt.Errorf("can't connect to database: %v", err)
	}

	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return fmt.Errorf("can't read schema: %v", err)
	}

	if _, err := Database.Exec(string(schema)); err != nil {
		return fmt.Errorf("schema execution failed: %v", err)
	}

	if refresh {
		if err := RefreshAndSeed(Database); err != nil {
			return fmt.Errorf("seed failed: %v", err)
		}
	}

	return nil
}

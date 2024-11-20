package storage

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/VsProger/snippetbox/pkg/config"
)

func NewSqlite(config config.Config) (*sql.DB, error) {
	log.Printf("Initializing database with driver: %s, DSN: %s", config.Driver, config.DSN)
	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}
	if err = CreateTables(db, config); err != nil {
		log.Printf("Failed to create tables: %v", err)
		return nil, err
	}
	log.Println("Database successfully initialized")
	return db, nil
}

func CreateTables(db *sql.DB, config config.Config) error {
	file, err := os.ReadFile(config.Database)
	if err != nil {
		return err
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		request = strings.TrimSpace(request)
		if request == "" {
			continue
		}

		_, err := db.Exec(request)
		if err != nil {
			log.Printf("Error executing SQL statement: %s, error: %v", request, err)
			return err
		}
	}
	log.Printf("Successfully executing SQL statements")

	return nil
}

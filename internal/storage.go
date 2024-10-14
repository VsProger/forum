package internal

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/VsProger/snippetbox/pkg/config"
)

func NewSqlite(config config.Config) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.Dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err = CreateTables(db, config); err != nil {
		return nil, err
	}
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

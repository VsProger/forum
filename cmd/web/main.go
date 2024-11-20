package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/VsProger/snippetbox/internal"
	"github.com/VsProger/snippetbox/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := internal.NewSqlite(*cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)

	templateCache, err := handlers.newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("starting server", slog.String("addr", cfg.Port))

	err = http.ListenAndServe(cfg.Port, handlers.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

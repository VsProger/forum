package server

import (
	"fmt"
	"net/http"

	"github.com/VsProger/snippetbox/internal/handlers"
	"github.com/VsProger/snippetbox/internal/repository"
	"github.com/VsProger/snippetbox/internal/service"
	"github.com/VsProger/snippetbox/internal/storage"
	"github.com/VsProger/snippetbox/logger"
	"github.com/VsProger/snippetbox/pkg/config"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{cfg: cfg}
}

func (app *App) Run() error {
	logger := logger.NewLogger()
	db, err := storage.NewSqlite(app.cfg)
	if err != nil {
		return err
	}
	logger.Info("Database successfully connected!")

	repo := repository.NewRepo(db)

	logger.Info("Repository working...")

	service := service.NewService(repo)

	logger.Info("Service working...")
	service.PostService.CreateCategory("IT")
	service.PostService.CreateCategory("Economy")
	service.PostService.CreateCategory("Medicine")
	service.PostService.CreateCategory("Other")

	handler := handlers.NewHandler(service)

	logger.Info("Handler working...")

	logger.Info("Server successfully started!")
	fmt.Printf("Server running on http://localhost%v\n", app.cfg.Port)

	return http.ListenAndServe(app.cfg.Port, handler.Router())
}

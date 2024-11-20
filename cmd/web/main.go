package main

import (
	"log"

	"github.com/VsProger/snippetbox/internal/server"
	"github.com/VsProger/snippetbox/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	app := server.NewApp(*cfg)

	if err := app.Run(); err != nil {
		log.Fatalln(err)
		return
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/VsProger/snippetbox/internal/server"
	"github.com/VsProger/snippetbox/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(cfg)

	app := server.NewApp(*cfg)

	if err := app.Run(); err != nil {
		log.Fatalln(err)
		return
	}
}

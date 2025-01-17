package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Driver   string `json:"Driver"`
	DSN      string `json:"DSN"`
	Database string `json:"Database"`
}

func NewConfig() (*Config, error) {
	configFile := "pkg/config/config.json"
	data, err := os.ReadFile(configFile)

	fmt.Print(string(data), "\n")
	if err != nil {
		log.Printf("ERROR: Read file in config encountered problem: %v", err)
		log.Fatal(err)
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("ERROR: Unmarshalling in config encountered problem: %v", err)
		log.Fatal(err)
		return nil, err
	}

	return &config, nil
}

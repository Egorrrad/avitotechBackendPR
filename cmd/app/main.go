package main

import (
	"log"

	"github.com/Egorrrad/avitotechBackendPR/config"
	"github.com/Egorrrad/avitotechBackendPR/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}

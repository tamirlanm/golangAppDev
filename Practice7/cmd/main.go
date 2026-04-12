package main

import (
	"Practice7/config"
	"Practice7/internal/app"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Config error:", err)
	}
	app.Run(cfg)
}

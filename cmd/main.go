package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/CactusBros/smaila/config"
	"github.com/CactusBros/smaila/handler"
	"github.com/joho/godotenv"
)

func main() {
	envPath := flag.String("env", "", "Path to .env file (optional)")
	flag.Parse()

	if *envPath == "" {
		*envPath = os.Getenv("CONFIG_PATH")
	}

	if *envPath != "" {
		if err := godotenv.Load(*envPath); err != nil {
			slog.Warn("Warning: failed to load .env file", "error", err)
		}
	}

	cfg := config.MustReadConfigFromEnv()

	handler.Run(cfg)
}

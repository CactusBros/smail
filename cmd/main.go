package main

import (
	"flag"
	"os"

	"github.com/CactusBros/smaila/config"
	"github.com/CactusBros/smaila/handler"
	"github.com/joho/godotenv"
)

var envPath = flag.String("env", ".env", "path to the environment variable file")

func main() {
	flag.Parse()
	if v := os.Getenv("CONFIG_PATH"); len(v) != 0 {
		*envPath = v
	}
	cfg := MustInitConfig()

	handler.Run(cfg)
}

func MustInitConfig() config.Config {
	err := godotenv.Load(*envPath)
	if err != nil {
		panic(err)
	}
	return config.MustReadConfigFromEnv()
}

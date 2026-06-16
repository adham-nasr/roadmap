package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"ETL/internal/app"
	"ETL/internal/config"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PipelineTimeout)
	defer cancel()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

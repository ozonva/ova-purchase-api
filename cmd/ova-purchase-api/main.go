package main

import (
	"fmt"
	"github.com/joho/godotenv"
	db2 "github.com/ozonva/ova-purchase-api/internal/db"
	"github.com/ozonva/ova-purchase-api/internal/server"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msgf("Failed to load env file %v", err)
	}
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := db2.NewDB(url)

	if err != nil {
		log.Fatal().Msgf("Failed to create connection to db %v", err)
	}

	if err := goose.Up(db.DB, "db/migration"); err != nil {
		log.Fatal().Msgf("Failed to run migrations %v", err)
	}

	server := server.NewServer(db, 81, 8181)
	server.Run()
}

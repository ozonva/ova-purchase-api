package main

import (
	"github.com/joho/godotenv"
	app2 "github.com/ozonva/ova-purchase-api/internal/app"
	"github.com/ozonva/ova-purchase-api/internal/config"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load env file")
	}
	configuration, err := config.LoadConfiguration("config/application.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	control := make(chan struct{}, 1)
	control <- struct{}{}

	app := app2.NewApp(configuration)

	for {
		select {
		case <-control:
			app.Start()
		case <-quit:
			app.Stop()
			return
		}
	}
}

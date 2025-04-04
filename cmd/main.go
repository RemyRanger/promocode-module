package main

import (
	"APIs/internal/app"
	"APIs/internal/common/config"
	"APIs/internal/common/server"
	"APIs/internal/common/telemetry"
	"os"

	"github.com/rs/zerolog/log"
)

const app_name = "promocode-api"
const app_config_filename = "config"

func main() {
	configFileName := app_config_filename
	if len(os.Args) < 2 {
		log.Info().Msgf("Config file name not specified, using default value %s (usage: go run main.go <configFileName>)", app_config_filename)
	} else {
		configFileName = os.Args[1]
	}

	// Init openTelemetry.
	otelShutdown := telemetry.SetupOTelSDK(app_name)

	// Init config
	app_config := config.LoadConfig(configFileName)

	// Run server
	_, handler := app.NewApp(app_name, app_config)
	if err := server.Run(app_name, app_config, handler, otelShutdown); err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}

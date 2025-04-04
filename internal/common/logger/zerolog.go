package logger

import (
	"APIs/internal/common/config"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const log_name = "log"

// New : initialize logger
func NewZerolog(app_config config.Config) {
	// Initialize Zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if app_config.Server.Env == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Str("service", log_name).Msg("Env is DEV, pretty logging enabled")
	}

	// Set Log level from config file
	switch logsLevel := app_config.Logs.Level; logsLevel {
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Str("service", log_name).Msg("Log level is INFO")
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Str("service", log_name).Msg("Log level is DEBUG")
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Info().Str("service", log_name).Msg("Log level is TRACE")
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Str("service", log_name).Msg("Config empty - Log level default is INFO")
	}
}

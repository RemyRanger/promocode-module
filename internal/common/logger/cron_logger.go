package logger

import (
	"github.com/rs/zerolog/log"
)

const cron_name = "cron"

type CronLogger struct {
}

func NewCronLogger() *CronLogger {
	return &CronLogger{}
}

func (l *CronLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Info().Str("service", cron_name).Msgf(msg, keysAndValues...)
}

func (l *CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Error().Err(err).Str("service", cron_name).Msgf(msg, keysAndValues...)
}

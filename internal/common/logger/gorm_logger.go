package logger

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const orm_name = "orm"

// MysqlLogger : Mysql custom Logger
type GormLogger struct {
}

func NewGormLogger() *GormLogger {
	return &GormLogger{}
}

// LogMode log mode
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...any) {
	log.Info().Str("service", orm_name).Msgf(msg, data...)
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	log.Warn().Str("service", orm_name).Msgf(msg, data...)
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...any) {
	log.Error().Str("service", orm_name).Msgf(msg, data...)
}

// Trace print sql message
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	traceStr := "%s\n" + "[%.3fms] " + "[rows:%v]" + " %s"

	sql, rows := fc()
	if err != nil {
		if rows == -1 {
			log.Debug().Str("service", orm_name).Msgf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Debug().Str("service", orm_name).Msgf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	} else {
		if rows == -1 {
			log.Trace().Str("service", orm_name).Msgf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Trace().Str("service", orm_name).Msgf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

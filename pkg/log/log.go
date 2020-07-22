package log

import (
	"go.uber.org/zap"
	"log"
	"time"
)

var logger *zap.Logger

func init() {
	SetupDev()
}

func SetupDev() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger = logger.WithOptions(zap.AddCallerSkip(1))
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Err(err error) zap.Field {
	return zap.String("err", err.Error())
}

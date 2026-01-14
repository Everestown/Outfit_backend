package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	red    = "\x1b[31m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	blue   = "\x1b[34m"
	reset  = "\x1b[0m"
)

func colorize(color string, s string) string {
	return color + s + reset
}

type Logger struct {
	*zap.Logger
}

func New(level string) *Logger {
	lvl, _ := zapcore.ParseLevel(level)

	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(lvl),
	)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)

	return &Logger{logger}
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func Success(msg string, fields ...zap.Field) {
	zap.L().Info(colorize(green, msg), fields...)
}

func Info(msg string, fields ...zap.Field) {
	zap.L().Info(colorize(blue, msg), fields...)
}

func Warning(msg string, fields ...zap.Field) {
	zap.L().Warn(colorize(yellow, msg), fields...)
}

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(colorize(red, msg), fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zap.L().Fatal(colorize(red, msg), fields...)
}

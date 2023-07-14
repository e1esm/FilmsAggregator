package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var Logger zap.Logger

func init() {
	Logger = newZapLogger()
}

func newZapLogger() zap.Logger {
	config := zap.NewDevelopmentEncoderConfig()

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)

	logFile, err := os.OpenFile("logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Couldn't have open/create logs.json file: %v", err)
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, defaultLogLevel))

	return *zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

package loger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Loger *zap.Logger
var ModelName = "[Loger]"

func InitLog() {
	ConsoleLogConfig := zap.NewProductionEncoderConfig()
	ConsoleLogConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	ConsoleLogConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	var ConsoleLogerEcoder zapcore.Encoder
	ConsoleLogerEcoder = zapcore.NewConsoleEncoder(ConsoleLogConfig)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(wd+"/log", 0775)
	if err != nil {
		panic(err)
	}
	File, err := os.Create(wd + "/log/" + time.Now().Format("2006-01-02_15_04_05") + ".log")
	if err != nil {
		panic(err)
	}
	FileConfig := zap.NewProductionEncoderConfig()
	FileConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	FileEcoder := zapcore.NewConsoleEncoder(FileConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(ConsoleLogerEcoder, os.Stdout, zap.InfoLevel),
		zapcore.NewCore(FileEcoder, File, zap.InfoLevel),
	)
	Loger = zap.New(core)
	Loger.Info(ModelName + "OK")
}

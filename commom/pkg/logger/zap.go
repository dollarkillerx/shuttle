package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.dev/google/common/pkg/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zp *zap.Logger
var logger *zap.SugaredLogger

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args)
}

func InitLogger(conf conf.LoggerConfig) {
	writeSyncer := getLogWriter(conf)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	zp = zap.New(core, zap.AddCallerSkip(1), zap.WithCaller(true))
	logger = zp.Sugar()
}

func GetSugaredLogger() *zap.SugaredLogger {
	return logger
}

func GetLogger() *zap.Logger {
	return zp
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(conf conf.LoggerConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.Filename, // 日志文件的位置
		MaxSize:    10,            // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,             // 保留旧文件的最大个数
		MaxAge:     30,            // 保留旧文件的最大天数
		Compress:   false,         // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}

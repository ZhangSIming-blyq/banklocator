package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	MyLogger   = InitLogger()
	DefaultLog = MyLogger.Sugar()
)

/*
InitLogger initializes logger for the program, example:
1. Output to console with readable format as well as log's level above "Info" var log = logger.DefaultLog
2. Output to console with readable format as well as log's level above "Debug" var log = logger.InitLogger("", "console", "debug").Sugar()
3. Output to file with readable format as well as log's level above "Debug" var log = logger.InitLogger("", "file", "debug").Sugar()
4. Output to console with json format as well as log's level above "Debug" var log = logger.InitLogger("json", "console", "debug").Sugar()
available variables:
format: json, ""
logType: file, json, ""
priority: debug, info, error, ""
*/
func InitLogger(logArgs ...string) *zap.Logger {
	var logger *zap.Logger
	var coreArr []zapcore.Core
	var format, logType, priority string

	// get the parameters
	switch {
	case len(logArgs) >= 3:
		format = logArgs[0]
		logType = logArgs[1]
		priority = logArgs[2]
	case len(logArgs) == 2:
		format = logArgs[0]
		logType = logArgs[1]
	case len(logArgs) == 1:
		format = logArgs[0]
	}

	// get encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // time format
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // use different color for various log levels
	// uncomment next line to show full path of the code
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	// NewJSONEncoder() for json，NewConsoleEncoder() for normal
	if format == "" {
		format = "normal"
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// log levels
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.InfoLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.InfoLevel && lev >= zap.DebugLevel
	})
	if logType == "" {
		logType = "console"
	}
	// writeSyncer for debug file
	debugFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/debug.log", // will create if not exist
		MaxSize:    128,               // max size for log file, unit:MB
		MaxBackups: 3,                 // max backup's count
		MaxAge:     10,                // max reserved days for log file
		Compress:   false,             // whether to compress or not
	})
	debugFileCore := zapcore.NewCore(encoder, os.Stdout, debugPriority)
	if logType == "file" {
		debugFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(debugFileWriteSyncer, zapcore.AddSync(os.Stdout)), debugPriority)
	}
	// writeSyncer for info file
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log",
		MaxSize:    128,
		MaxBackups: 3,
		MaxAge:     10,
		Compress:   false,
	})
	infoFileCore := zapcore.NewCore(encoder, os.Stdout, infoPriority)
	if logType == "file" {
		infoFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), infoPriority)
	}
	// writeSyncer for error file
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    128,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   false,
	})
	errorFileCore := zapcore.NewCore(encoder, os.Stdout, errorPriority)
	if logType == "file" {
		errorFileCore = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), errorPriority)
	}

	switch priority {
	case "":
		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	case "info":
		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	case "error":
		coreArr = append(coreArr, errorFileCore)
	case "debug":
		coreArr = append(coreArr, debugFileCore)
		coreArr = append(coreArr, infoFileCore)
		coreArr = append(coreArr, errorFileCore)
	}
	logger = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller() is to show the line number
	return logger
}

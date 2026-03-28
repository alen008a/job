package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"siteVideoJob/config"
)

var (
	ZapLog *zap.SugaredLogger // 简易版日志文件
	// Logger *zap.Logger // 这个日志强大一些, 目前还用不到

	logLevel = zap.NewAtomicLevel()
)

// InitLog 初始化日志文件
func getLogWriter(logPath, logType string) zapcore.WriteSyncer {
	if logType == "file" {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    128, // MB
			LocalTime:  true,
			Compress:   true,
			MaxBackups: 8, // 最多保留 n 个备份
		})
		return w
	} else {
		return zapcore.Lock(os.Stdout) //标准输出
	}
}

func getEncoder(logType string) zapcore.Encoder {
	var encoder zapcore.Encoder
	if logType == "file" {
		c := zap.NewProductionEncoderConfig()
		c.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(c)
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}
	return encoder
}

// InitLog 初始化日志文件
func InitLog() zapcore.WriteSyncer {
	logConf := config.GetConfig().Logger
	/**
	logConf := config.Logger{
		LogPath:  "",
		LogLevel: "INFO",
		LogType:  "",
	}**/

	loglevel := zapcore.InfoLevel
	switch logConf.LogLevel {
	case "INFO":
		loglevel = zapcore.InfoLevel
	case "ERROR":
		loglevel = zapcore.ErrorLevel
	}
	setLevel(loglevel)

	encoder := getEncoder(logConf.LogType)                   //获取编码方式
	writer := getLogWriter(logConf.LogPath, logConf.LogType) //获取writer
	core := zapcore.NewCore(encoder, writer, loglevel)
	ZapLog = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return writer
}

func setLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}

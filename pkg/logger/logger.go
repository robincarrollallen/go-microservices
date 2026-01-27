package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// colorizedEncodeLevel 自定义的Level编码函数，支持颜色输出
func colorizedEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorCode, resetCode string

	switch level {
	case zapcore.DebugLevel:
		colorCode = "\033[36m" // 青色
	case zapcore.InfoLevel:
		colorCode = "\033[32m" // 绿色
	case zapcore.WarnLevel:
		colorCode = "\033[33m" // 黄色
	case zapcore.ErrorLevel:
		colorCode = "\033[31m" // 红色
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		colorCode = "\033[35m" // 品红色
	default:
		colorCode = ""
	}

	resetCode = "\033[0m" // 重置颜色

	if colorCode != "" {
		enc.AppendString(colorCode + level.CapitalString() + resetCode)
	} else {
		enc.AppendString(level.CapitalString())
	}
}

func Init(level string) {
	// 使用开发环境配置以支持控制台输出和颜色
	cfg := zap.NewDevelopmentConfig()

	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = colorizedEncodeLevel // 添加颜色编码器
	cfg.Encoding = "console" // 改用控制台编码器而不是JSON
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	if level == "debug" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	l, err := cfg.Build(
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	log = l
}

func L() *zap.Logger {
	return log
}

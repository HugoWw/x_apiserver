package clog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type clog struct {
	LogLevel zap.AtomicLevel
	Logger   *zap.SugaredLogger
}

var customTimeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
}

var customLevelEncoder = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

var customCallerEncoder = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

var customCallName = func(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "]")
}

func initZapSugLogger(atom zap.AtomicLevel) *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:          "msg",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "logger",
		CallerKey:           "caller",
		FunctionKey:         zapcore.OmitKey,
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         customLevelEncoder,
		EncodeTime:          customTimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        customCallerEncoder,
		EncodeName:          customCallName,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    " - ",
	}

	zapConfig := zap.Config{
		Level:             atom,
		DisableCaller:     false,
		DisableStacktrace: true,
		Development:       false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := zapConfig.Build()
	return logger.Sugar()
}

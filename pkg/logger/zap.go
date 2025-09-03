package logger

import (
	zapLogger "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewZapLogger(cfg *Config) log.Logger {
	if cfg == nil {
		return nil
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    int(cfg.MaxSize),
		MaxBackups: int(cfg.MaxBackups),
		MaxAge:     int(cfg.MaxAge),
	}
	// 同时输出到控制台和文件
	var syncers []zapcore.WriteSyncer
	syncers = append(syncers, zapcore.AddSync(lumberJackLogger))
	if cfg.Stdout {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(
		syncers...,
	)

	var lvl = new(zapcore.Level)
	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil
	}

	core := zapcore.NewCore(jsonEncoder, writeSyncer, lvl)
	logger := zap.New(core)
	wrapped := log.With(zapLogger.NewLogger(logger))
	return wrapped
}

package logger

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

func NewLogger(config *Config) log.Logger {
	switch config.LoggerType {
	default:
		fallthrough
	case Std:
		return NewStdLogger()
	case Zap:
		return NewZapLogger(config)
	case Fluent:
		return nil
	case Logrus:
		return nil
	case Aliyun:
		return nil
	case Tencent:
		return nil
	}
}

// NewLoggerProvider 创建一个新的日志记录器提供者
func NewLoggerProvider(serviceInfo *ServiceInfo, opt ...Option) log.Logger {
	config := defaultConfig
	for _, o := range opt {
		o(config)
	}
	l := NewLogger(config)

	return log.With(
		l,
		"service.id", serviceInfo.Id,
		"service.name", serviceInfo.Name,
		"service.version", serviceInfo.Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", ifElse(config.SimpleTrace, GetTraceId(), tracing.TraceID()),
		"span_id", tracing.SpanID(),
	)
}

func ifElse(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

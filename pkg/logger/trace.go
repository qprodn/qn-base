package logger

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/uuid"
)

type TraceIdKey struct {
}

func SimpleTraceIdProvider() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			v7, err := uuid.NewV7()
			ctx = context.WithValue(ctx, TraceIdKey{}, v7)
			return handler(ctx, req)
		}
	}
}

func GetTraceId() log.Valuer {
	return func(ctx context.Context) any {
		return ctx.Value(TraceIdKey{})
	}
}

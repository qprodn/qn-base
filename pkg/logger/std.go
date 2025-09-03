package logger

import (
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

// NewStdLogger 创建一个新的日志记录器 - Kratos内置，控制台输出
func NewStdLogger() log.Logger {
	l := log.NewStdLogger(os.Stdout)
	return l
}

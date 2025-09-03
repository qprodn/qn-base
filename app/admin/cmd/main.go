package main

import (
	"flag"
	"os"
	"qn-base/pkg/logger"

	"qn-base/app/admin/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	// 初始化配置
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 初始化日志
	loggerProvider := logger.NewLoggerProvider(
		&logger.ServiceInfo{
			Id:      "",
			Name:    Name,
			Version: Version,
		},
		logger.WithLoggerType(logger.Zap),
		logger.WithFile(bc.Log.Filename),
		logger.WithLevel(bc.Log.Level),
		logger.WithMaxAge(bc.Log.MaxAge),
		logger.WithMaxSize(bc.Log.MaxSize),
		logger.WithMaxBackups(bc.Log.MaxBackups),
		logger.WithStdout(bc.Log.Stdout),
		logger.WithSimpleTrace(bc.Env.Active == "dev"),
	)

	// 初始化服务
	app, cleanup, err := wireApp(&bc, loggerProvider)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

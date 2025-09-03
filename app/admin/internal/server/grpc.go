package server

import (
	adminV1 "qn-base/api/gen/go/admin/v1"
	"qn-base/app/admin/internal/conf"
	"qn-base/app/admin/internal/service/systemuser"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, userService *systemuser.UserService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(c.Server.Grpc.Timeout)))
	}
	srv := grpc.NewServer(opts...)
	adminV1.RegisterUserServer(srv, userService)
	return srv
}

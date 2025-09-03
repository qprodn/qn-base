//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"qn-base/app/admin/internal/biz"
	"qn-base/app/admin/internal/conf"
	datainit "qn-base/app/admin/internal/data/wire"
	"qn-base/app/admin/internal/server"
	"qn-base/app/admin/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, datainit.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

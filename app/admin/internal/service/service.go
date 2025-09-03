package service

import (
	"qn-base/app/admin/internal/service/systemuser"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(systemuser.NewUserService)

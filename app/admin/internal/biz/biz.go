package biz

import (
	"context"
	"qn-base/app/admin/internal/biz/systemuser"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(systemuser.NewUserUsecase)

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}

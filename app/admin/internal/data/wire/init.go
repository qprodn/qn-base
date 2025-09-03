package wire

import (
	"github.com/google/wire"
	"qn-base/app/admin/internal/data/data"
	"qn-base/app/admin/internal/data/db"
	"qn-base/app/admin/internal/data/idgen"
	"qn-base/app/admin/internal/data/systemuser"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(data.NewData, data.NewTransaction, systemuser.NewSystemUserRepo, db.NewDB, idgen.NewIDGenerator)

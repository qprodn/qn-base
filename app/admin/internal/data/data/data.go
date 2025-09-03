package data

import (
	"qn-base/app/admin/internal/biz"
	"qn-base/app/admin/internal/data/ent"

	"github.com/go-kratos/kratos/v2/log"
)

// Data .
type Data struct {
	db *ent.Database
}

// NewTransaction .
func NewTransaction(data *Data) biz.Transaction {
	return data.db
}

// Database returns the ent database client.
func (d *Data) Database() *ent.Database {
	return d.db
}

// NewData .
func NewData(logger log.Logger, db *ent.Database) (*Data, func(), error) {
	helper := log.NewHelper(logger)
	d := &Data{
		db: db,
	}

	return d, func() {
		helper.Info("message", "closing the data resources")
		d.db.Close()
	}, nil
}

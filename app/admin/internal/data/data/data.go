package data

import (
	"qn-base/app/admin/internal/biz"
	"qn-base/app/admin/internal/data/ent"

	"github.com/go-kratos/kratos/v2/log"
)

// Data .
type Data struct {
	DB *ent.Database
}

// NewTransaction .
func NewTransaction(data *Data) biz.Transaction {
	return data.DB
}

// Database returns the ent database client.
func (d *Data) Database() *ent.Database {
	return d.DB
}

// NewData .
func NewData(logger log.Logger, db *ent.Database) (*Data, func(), error) {
	helper := log.NewHelper(logger)
	d := &Data{
		DB: db,
	}

	return d, func() {
		helper.Info("message", "closing the data resources")
		d.DB.Close()
	}, nil
}

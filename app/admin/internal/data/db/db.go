package db

import (
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"qn-base/app/admin/internal/conf"
	"qn-base/app/admin/internal/data/ent"
	"time"
)

func NewDB(c *conf.Bootstrap, logger log.Logger) *ent.Database {
	helper := log.NewHelper(logger)
	drv, err := sql.Open(
		c.Data.Database.Driver,
		c.Data.Database.Source,
	)

	if err != nil {
		helper.Errorf("failed opening connection to db: %v", err)
		panic(err)
	}

	// 连接池中最多保留的空闲连接数量
	drv.DB().SetMaxIdleConns(int(c.Data.Database.MaxIdleConns))
	// 连接池在同一时间打开连接的最大数量
	drv.DB().SetMaxOpenConns(int(c.Data.Database.MaxOpenConns))
	// 连接可重用的最大时间长度
	drv.DB().SetConnMaxLifetime(time.Duration(c.Data.Database.ConnMaxLifetime))

	return ent.NewDatabase(ent.Driver(drv))

}

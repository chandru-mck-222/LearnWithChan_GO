package db

import (
	"sync"

	"chanlearning/pkg/service/common/envvar"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	err error
	mu  sync.Mutex
	DB  *gorm.DB
)

func DBinit() *gorm.DB {
	if DB == nil {
		mu.Lock()
		defer mu.Unlock()
		if DB == nil {
			dsn := os.Getenv(envvar.MYSQLCLIENT)
			if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err == nil {
				return DB
			} else {
				panic(err)
			}
		} else {
			return DB
		}
	} else {
		return DB
	}
}

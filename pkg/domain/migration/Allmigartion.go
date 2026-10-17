package migration

import (
	"chanlearning/pkg/domain/db"
	"time"

	"gorm.io/gorm"
)

func All() []Migration {
	return []Migration{
		Products,
	}
}

func Migratersetup() {
	all := All()
	Always := Filter(all, func(migration Migration) bool {
		return !migration.IsOnetime
	})
	OneTime := Filter(all, func(migration Migration) bool {
		return migration.IsOnetime
	})
	for _, d := range Always {
		d.Up(db.DBinit())
	}

	for _, d := range OneTime {
		if err := db.DBinit().Where("name = ?", d.Name).First(&OnetimeTable{}).Error; err != nil {
			if err := gorm.ErrRecordNotFound; err != nil {
				d.Up(db.DBinit())
				db.DB.Create(&OnetimeTable{
					Name:   d.Name,
					Active: time.Now(),
				})
			}
		}
	}
}

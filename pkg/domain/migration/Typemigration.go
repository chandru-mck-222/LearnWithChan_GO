package migration

import (
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	Name      string
	Up        func(*gorm.DB)
	Down      func(*gorm.DB)
	IsOnetime bool
}

func Filter[T any](migration []T, test func(T) bool) []T {
	var Res = []T{}
	for _, data := range migration {
		if test(data) {
			Res = append(Res, data)

		}
	}
	return Res
}

type OnetimeTable struct {
	Name   string
	Active time.Time
}

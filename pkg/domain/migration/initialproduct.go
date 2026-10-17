package migration

import (
	"chanlearning/pkg/domain/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

var Products = Migration{
	Name: "products",
	Up: func(tx *gorm.DB) {
		tx.AutoMigrate(&models.ProductsDetails{})
		var SeedingProduct = []models.ProductsDetails{
			{
				ProductId:    "chan123",
				Name:         "mobile",
				ProductImage: "",
				Description:  "iphone 15",
				Price:        85000.21,
			}, {
				ProductId:    "chan125",
				Name:         "T-shirt",
				ProductImage: "",
				Description:  "loosely shirt style manner",
				Price:        4300,
			}, {
				ProductId:    "chan126",
				Name:         "Monitors",
				ProductImage: "",
				Description:  "lenovo legion 27*inches curve display",
				Price:        20000,
			},
		}
		for _, Val := range SeedingProduct {
			var exists = &models.ProductsDetails{}
			err := tx.Where("product_id=?", Val.ProductId).First(exists).Error
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("failed to check product %s: %v",
					Val.ProductId, err)
				continue
			}

			if err := tx.Create(&Val).Error; err != nil {
				log.Printf("failed to insertdata %v", Val)
			} else {
				continue
			}

		}
	},
	Down: func(tx *gorm.DB) {
		tx.Migrator().DropTable(&models.ProductsDetails{})
	},
	IsOnetime: true,
}

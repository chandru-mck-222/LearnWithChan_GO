package models

import (
	"chanlearning/pkg/domain/db"
	"time"

	"gorm.io/gorm"
)

type ProductsDetails struct {
	Id           int    `json:"id"`
	ProductId    string `json:"product_id"`
	Name         string `json:"name"`
	ProductImage string
	Description  string    `json:"description"`
	Price        float64   `json:"prices"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"update_at"`
	Deleted_at   time.Time `json:"Deleted_at"`
}

var SeedingProduct = []ProductsDetails{
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

type Products struct {
	ProductId   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"prices"`
}

type ProductDef struct {
	ProductId   string    `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"prices"`
	Created_at  time.Time `json:"created_at"`
}

func Database() *gorm.DB {
	return db.DBinit()
}

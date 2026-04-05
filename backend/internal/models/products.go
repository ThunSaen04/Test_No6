package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductCode string `gorm:"size:60;uniqueIndex:idx_product_code_active,where:deleted_at IS NULL;not null" json:"product_code"`
	BarcodePath string `gorm:"size:255;not null" json:"barcode_path"`
}

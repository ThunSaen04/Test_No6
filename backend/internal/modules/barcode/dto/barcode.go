package dto

import "time"

type CreateBarcodeRequest struct {
	ProductCode string `json:"product_code"`
	ItemCode    string `json:"item_code"`
	ItemsCode   string `json:"items_code"`
}

type BarcodeResponse struct {
	ID          uint      `json:"id"`
	ProductCode string    `json:"product_code"`
	BarcodePath string    `json:"barcode_path"`
	BarcodeURL  string    `json:"barcode_url"`
	CreatedAt   time.Time `json:"created_at"`
}

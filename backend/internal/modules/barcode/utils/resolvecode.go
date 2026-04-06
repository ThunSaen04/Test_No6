package utils

import (
	"strings"

	"no6/backend/internal/modules/barcode/dto"
)

func ResolveCode(r dto.CreateBarcodeRequest) string {
	if code := strings.TrimSpace(r.ProductCode); code != "" {
		return code
	}
	return strings.TrimSpace(r.ItemsCode)
}

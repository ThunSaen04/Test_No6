package utils

import (
	"strings"

	"no6/backend/internal/modules/products/dto"
)

func ResolveCode(r dto.CreateBarcodeRequest) string {
	if code := strings.TrimSpace(r.ProductCode); code != "" {
		return code
	}
	return strings.TrimSpace(r.ItemsCode)
}

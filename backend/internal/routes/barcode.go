package routes

import (
	"no6/backend/internal/config"
	productHandler "no6/backend/internal/modules/barcode/handler"
	productRepository "no6/backend/internal/modules/barcode/repository"
	productService "no6/backend/internal/modules/barcode/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setupBarcodeRoutes(router fiber.Router, db *gorm.DB, cfg config.Config) {
	repo := productRepository.NewProductRepository(db)
	svc := productService.NewBarcodeService(repo, cfg)
	hdl := productHandler.NewBarcodeHandler(svc, cfg)

	barcode := router.Group("/barcodes")
	barcode.Get("/", hdl.List)
	barcode.Post("/", hdl.Create)
	barcode.Delete("/:id", hdl.Delete)
}

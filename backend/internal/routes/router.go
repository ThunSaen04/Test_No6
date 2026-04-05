package routes

import (
	"no6/backend/internal/config"
	productHandler "no6/backend/internal/modules/products/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, handler *productHandler.BarcodeHandler, cfg config.Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "barcode-backend",
		})
	})

	app.Static("/images", cfg.BarcodeDir)

	api := app.Group("/api")
	barcodes := api.Group("/barcodes")

	barcodes.Get("", handler.List)
	barcodes.Post("", handler.Create)
	barcodes.Delete("/:id", handler.Delete)
}

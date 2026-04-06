package routes

import (
	"no6/backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg config.Config) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "barcode-backend",
		})
	})

	app.Static("/images", cfg.BarcodeDir)

	api := app.Group("/api")

	setupBarcodeRoutes(api, db, cfg)
}

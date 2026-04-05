package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"no6/backend/internal/config"
	"no6/backend/internal/database"
	productHandler "no6/backend/internal/modules/products/handler"
	productRepository "no6/backend/internal/modules/products/repository"
	productService "no6/backend/internal/modules/products/service"
	"no6/backend/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if err := os.MkdirAll(cfg.BarcodeDir, 0o755); err != nil {
		log.Fatalf("create barcode directory failed: %v", err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	barcodeRepo := productRepository.NewProductRepository(db)
	barcodeService := productService.NewBarcodeService(barcodeRepo, cfg)
	barcodeHandler := productHandler.NewBarcodeHandler(barcodeService, cfg)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.Setup(app, barcodeHandler, cfg)

	addr := ":" + cfg.ServerPort
	log.Printf("server running on %s", addr)
	log.Fatal(app.Listen(addr))
}

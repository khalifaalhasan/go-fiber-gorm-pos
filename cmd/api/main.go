package main

import (
	"go-fiber-pos/internal/config"
	"go-fiber-pos/internal/routes" // Import package routes yang baru dibuat
	"go-fiber-pos/pkg/validator"
	"go-fiber-pos/pkg/logger"



	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Inisialisasi Logger Custom & Validator
	logger.InitLogger()
	validator.InitValidator()

	// 2. Load .env
	err := godotenv.Load()
	if err != nil {
		logger.Log.Warn("File .env tidak ditemukan, menggunakan variabel OS")
	}

	// 3. Konek GORM ke Database
	config.ConnectDatabase()

	// 4. Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Bangga Punya Web - POS API",
	})

	app.Use(fiberlog.New())
	app.Use(recover.New())

	// 5. Panggil Setup Routes 
	routes.SetupRoutes(app)

	// 6. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	logger.Log.Infof("Starting server on port %s", port)
	err = app.Listen(":" + port)
	if err != nil {
		logger.Log.Fatalf("Error starting server: %v", err)
	}
}
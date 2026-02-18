package main

import (
	"go-fiber-pos/config"
	"go-fiber-pos/routes" // Import package routes yang baru dibuat
	"go-fiber-pos/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Inisialisasi Logger Custom & Validator
	utils.InitLogger()
	utils.InitValidator()

	// 2. Load .env
	err := godotenv.Load()
	if err != nil {
		utils.Log.Warn("File .env tidak ditemukan, menggunakan variabel OS")
	}

	// 3. Konek GORM ke Database
	config.ConnectDatabase()

	// 4. Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Bangga Punya Web - POS API",
	})

	// Middleware bawaan Fiber untuk log request dan anti-crash
	app.Use(logger.New())
	app.Use(recover.New())

	// 5. Panggil Setup Routes (Magic-nya ada di sini)
	routes.SetupRoutes(app)

	// 6. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	utils.Log.Infof("Starting server on port %s", port)
	err = app.Listen(":" + port)
	if err != nil {
		utils.Log.Fatalf("Error starting server: %v", err)
	}
}
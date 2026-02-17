package main

import (
	"os"
	"go-fiber-pos/config"
	"go-fiber-pos/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Inisialisasi Logger Custom
	utils.InitLogger()

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

	// Middleware bawaan Fiber untuk log request dan anti-crash (recover)
	app.Use(logger.New())
	app.Use(recover.New())

	// Route Test Ping
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong! Server Fiber berjalan secepat kilat 🚀",
		})
	})

	// 5. Jalankan Server
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
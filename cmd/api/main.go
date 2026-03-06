package main

import (
	"go-fiber-pos/internal/config"
	"go-fiber-pos/internal/routes" // Import package routes yang baru dibuat
	"go-fiber-pos/pkg/logger"
	"go-fiber-pos/pkg/validator"

	"os"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

// @title Bangga Punya Web - POS API
// @version 1.0
// @description API Dokumentasi untuk Sistem POS (Point of Sale) Bangga Punya Web.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @query.collection.format multi

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your-jwt-token>" to authenticate.
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
	
	// 4. Konek ke Redis
	config.ConnectRedis()

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
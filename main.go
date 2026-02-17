package main

import (
	"go-fiber-pos/config"
	"go-fiber-pos/controller"
	"go-fiber-pos/middleware"
	"go-fiber-pos/repository"
	"go-fiber-pos/service"
	"go-fiber-pos/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Inisialisasi Logger Custom

	// Middleware bawaan Fiber untuk log request dan anti-crash (recover)


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

	app.Use(logger.New())
	app.Use(recover.New())

	// Auth
	authRepo := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

	// category
	categoryRepo := repository.NewCategoryRepository(config.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)

	// product
	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	api := app.Group("/api")
	

	// auth
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	api.Use(middleware.Protected())

	// category	
	api.Post("/categories", categoryController.Create)
	api.Get("/categories", categoryController.GetAll)

	// product
	api.Post("/products", productController.Create)
	api.Get("/products", productController.GetAll)

	

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
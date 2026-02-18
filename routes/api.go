package routes

import (
	"go-fiber-pos/config"
	"go-fiber-pos/controller"
	"go-fiber-pos/middleware"
	"go-fiber-pos/repository"
	"go-fiber-pos/service"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes merakit semua dependency dan mendaftarkan endpoint
func SetupRoutes(app *fiber.App) {
	// ==========================================
	// 1. DEPENDENCY INJECTION (Setup Variabel)
	// ==========================================
	
	// Auth
	authRepo := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

	// Category
	categoryRepo := repository.NewCategoryRepository(config.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)
	
	// Product
	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	// public
	publicCategoryController := controller.NewPublicCategoryController(categoryService)
	publicProductController := controller.NewPublicProductController(productService)

	// 💡 TODO Besok: Variabel Controller Public untuk Pelanggan (Di-comment dulu)
	// publicCategoryController := controller.NewPublicCategoryController(categoryService)
	// publicProductController := controller.NewPublicProductController(productService)

	// ==========================================
	// 2. ROUTING
	// ==========================================
	api := app.Group("/api")

	// Route Test Ping
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong! Server Fiber berjalan secepat kilat 🚀",
		})
	})

	// A. AUTHENTICATION (Bebas Akses)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

// B. PUBLIC ROUTE (Katalog Menu untuk Pelanggan / Scan QR)
	publicGroup := api.Group("/public")
	
	// Endpoint Get All Menu (Hapus parameter :store_id)
	publicGroup.Get("/menu/categories", publicCategoryController.GetAllMenu)
	publicGroup.Get("/menu/products", publicProductController.GetAllMenu)


	// C. ADMIN / DASHBOARD ROUTE (Wajib Token JWT)
	// Pasang Middleware JWT khusus untuk grup admin
	adminGroup := api.Group("/admin", middleware.Protected()) 

	// Category Management (Kasir)
	// localhost:8080/api/admin/

	adminGroup.Post("/categories", categoryController.Create)  // Kasir tambah kategori

	adminGroup.Post("/products", productController.Create)     // Kasir tambah produk
}
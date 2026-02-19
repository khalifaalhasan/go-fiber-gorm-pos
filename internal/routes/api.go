package routes

import (
	"go-fiber-pos/internal/config"
	"go-fiber-pos/internal/middleware"
	"go-fiber-pos/internal/modules/auth"
	"go-fiber-pos/internal/modules/category"
	"go-fiber-pos/internal/modules/product"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes merakit semua dependency dan mendaftarkan endpoint
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Route Test Ping
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong! Server Fiber berjalan secepat kilat 🚀",
		})
	})

	// ==========================================
	// DEFINISI GRUP RUTE UTAMA
	// ==========================================
	
	// A. Bebas Akses (Auth)
	authGroup := api.Group("/auth")
	
	// B. Public Route (Katalog Pelanggan / QR)
	publicGroup := api.Group("/public")
	
	// C. Admin Route (Wajib Token JWT)
	adminGroup := api.Group("/admin", middleware.Protected())

	// ==========================================
	// PENDAFTARAN RUTE PER MODUL (DDD STYLE)
	// ==========================================
	
	// Lempar grup dan koneksi DB ke masing-masing modul
	auth.SetupRoutes(authGroup, config.DB)
	category.SetupRoutes(adminGroup, publicGroup, config.DB)
	product.SetupRoutes(adminGroup, publicGroup, config.DB)
}
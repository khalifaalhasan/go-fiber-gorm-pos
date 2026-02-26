package routes

import (
	"go-fiber-pos/internal/config"
	"go-fiber-pos/internal/infrastructure/provider"
	"go-fiber-pos/internal/middleware"
	"go-fiber-pos/internal/modules/auth"
	"go-fiber-pos/internal/modules/category"
	"go-fiber-pos/internal/modules/inventory"
	"go-fiber-pos/internal/modules/order"
	"go-fiber-pos/internal/modules/payment"
	"go-fiber-pos/internal/modules/product"
	"go-fiber-pos/internal/modules/store"
	"go-fiber-pos/internal/modules/voucher"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes merakit semua dependency dan mendaftarkan endpoint
func SetupRoutes(app *fiber.App) {
	v := validator.New()

	// Inisialisasi Payment Gateway Adapter (PORT & ADAPTER pattern)
	// Mudah diganti dengan adapter lain tanpa mengubah service layer
	midtransAdapter := provider.NewMidtransAdapter()

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

	// D. Webhook Route (PUBLIC — tidak butuh JWT, diakses oleh payment gateway)
	webhookGroup := api.Group("/webhook")

	// ==========================================
	// PENDAFTARAN RUTE PER MODUL (DDD STYLE)
	// ==========================================

	// Existing modules
	auth.SetupRoutes(authGroup, config.DB)
	category.SetupRoutes(adminGroup, publicGroup, config.DB)
	product.SetupRoutes(adminGroup, publicGroup, config.DB, v)

	invRepo := inventory.NewInventoryRepository(config.DB)
	invService := inventory.NewInventoryService(invRepo, v)

	// New modules
	store.SetupRoutes(adminGroup, config.DB, v)
	voucher.SetupRoutes(adminGroup, config.DB, v)
	inventory.SetupRoutes(adminGroup, config.DB, v)
	order.SetupRoutes(adminGroup, config.DB, v)
	payment.SetupRoutes(adminGroup, webhookGroup, config.DB, v, midtransAdapter, invService)
}
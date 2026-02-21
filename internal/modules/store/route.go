package store

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(adminGroup fiber.Router, db *gorm.DB, v *validator.Validate) {
	repo := NewStoreRepository(db)
	service := NewStoreService(repo, v)
	ctrl := NewStoreController(service, v)

	// Admin-only endpoints
	adminGroup.Get("/store-profile", ctrl.GetProfile)
	adminGroup.Put("/store-profile", ctrl.UpdateProfile)
}

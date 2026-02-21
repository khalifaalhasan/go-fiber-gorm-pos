package store

import (
	"errors"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type StoreController struct {
	service StoreService
	v       *validator.Validate
}

func NewStoreController(service StoreService, v *validator.Validate) *StoreController {
	return &StoreController{service: service, v: v}
}

func (ctrl *StoreController) GetProfile(c *fiber.Ctx) error {
	profile, err := ctrl.service.GetProfile()
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Profil toko belum dikonfigurasi",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": profile})
}

func (ctrl *StoreController) UpdateProfile(c *fiber.Ctx) error {
	var req UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	profile, err := ctrl.service.UpdateProfile(req)
	if err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validasi gagal", "details": valErr.Error()})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profil toko berhasil diperbarui",
		"data":    profile,
	})
}

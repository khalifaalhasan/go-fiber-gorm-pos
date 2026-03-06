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

// GetProfile godoc
// @Summary      Get store profile
// @Description  Retrieve the store profile information (Admin only).
// @Tags         store
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/stores/profile [get]
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

// UpdateProfile godoc
// @Summary      Update store profile
// @Description  Update the store profile information (Admin only).
// @Tags         store
// @Accept       json
// @Produce      json
// @Param        request body UpdateStoreRequest true "Update Store Request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/stores/profile [put]
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

package controller

import (
	"go-fiber-pos/model"

	"github.com/gofiber/fiber/v2"
)

type PublicCategoryController struct {
	service model.CategoryService
}

func NewPublicCategoryController(service model.CategoryService) *PublicCategoryController {
	return &PublicCategoryController{service: service}
}

// GetAllMenu mengambil semua kategori untuk tampilan pelanggan
func (ctrl *PublicCategoryController) GetAllMenu(c *fiber.Ctx) error {
	// Lihat! Betapa bersihnya kodingan ini tanpa c.Params("store_id")
	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data kategori menu"})
	}

	// Mapping ke DTO agar format JSON rapi
	res := model.ToCategoryResponseList(categories)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
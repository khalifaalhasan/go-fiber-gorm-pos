package controller

import (
	"go-fiber-pos/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryController struct {
	service model.CategoryService
}

func NewCategoryController(service model.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (ctrl *CategoryController) Create(c *fiber.Ctx) error {
	// Ambil store_id dari token JWT (Hasil kerja Middleware)
	storeIDContext := c.Locals("store_id").(uuid.UUID)

	var req model.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	// res
	category, err := ctrl.service.CreateCategory(storeIDContext, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Bungkus pakai kardus DTO
	res := model.ToCategoryResponse(category)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Kategori berhasil dibuat",
		"data":    res,
	})
}

func (ctrl *CategoryController) GetAll(c *fiber.Ctx) error {
	storeIDContext := c.Locals("store_id").(uuid.UUID)

	// res
	categories, err := ctrl.service.GetCategoriesByStore(storeIDContext)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

	// Bungkus array-nya pakai kardus DTO
	res := model.ToCategoryResponseList(categories)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
package category

import (
	model "go-fiber-pos/internal/core"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	service model.CategoryService
}

func NewCategoryController(service model.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (ctrl *CategoryController) Create(c *fiber.Ctx) error {
	var req model.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	category, err := ctrl.service.CreateCategory(req)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": category,
		"message": "Kategori berhasil dibuat",
	})
}


func (ctrl *CategoryController) GetAll(c *fiber.Ctx) error {
	// Panggil service (Bersih tanpa storeIDContext)
	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	res := model.ToCategoryResponseList(categories)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
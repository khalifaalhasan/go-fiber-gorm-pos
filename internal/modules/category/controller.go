package category

import (
	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	service CategoryService
}

func NewCategoryController(service CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (ctrl *CategoryController) Create(c *fiber.Ctx) error {
	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	category, err := ctrl.service.CreateCategory(req)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Kategori berhasil dibuat",
		"data": category,
	})
}


func (ctrl *CategoryController) GetAll(c *fiber.Ctx) error {
	// Panggil service (Bersih tanpa storeIDContext)
	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	res := ToCategoryResponseList(categories)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
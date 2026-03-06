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

// Create godoc
// @Summary      Create a new category
// @Description  Create a new category for products (Admin only).
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request body CreateCategoryRequest true "Create Category Request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/categories [post]
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


// GetAll godoc
// @Summary      Get all categories
// @Description  Retrieve a list of all product categories (Admin only).
// @Tags         categories
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/categories [get]
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
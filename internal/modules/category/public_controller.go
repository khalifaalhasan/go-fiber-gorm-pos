package category

import (
	"github.com/gofiber/fiber/v2"
)

type PublicCategoryController struct {
	service CategoryService
}

func NewPublicCategoryController(service CategoryService) *PublicCategoryController {
	return &PublicCategoryController{service: service}
}


// GetAllMenu godoc
// @Summary      Get all categories for public menu
// @Description  Retrieve a list of all product categories for the public customer menu.
// @Tags         public
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /public/categories [get]
func (ctrl *PublicCategoryController) GetAllMenu(c *fiber.Ctx) error {

	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}


	res := ToCategoryResponseList(categories)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
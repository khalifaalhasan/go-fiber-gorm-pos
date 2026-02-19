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
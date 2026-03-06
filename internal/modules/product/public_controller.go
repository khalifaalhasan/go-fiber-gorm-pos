package product

import (
	"github.com/gofiber/fiber/v2"
)

type PublicProductController struct {
	service ProductService
}

// 1. Return type diperbaiki menjadi *PublicProductController
func NewPublicProductController(service ProductService) *PublicProductController {
	return &PublicProductController{service: service}
}

// 2. Receiver diperbaiki menjadi *PublicProductController
// GetAllMenu godoc
// @Summary      Get all products for public menu
// @Description  Retrieve a list of all products for the public customer menu.
// @Tags         public
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /public/products [get]
func (ctrl *PublicProductController) GetAllMenu(c *fiber.Ctx) error {
	products, err := ctrl.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Panggil fungsi mapper yang baru
	res := ToProductResponseList(products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
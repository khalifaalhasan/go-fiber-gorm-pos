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
package controller

import (
	"go-fiber-pos/model"

	"github.com/gofiber/fiber/v2"
)

type PublicProductController struct{
	service model.ProductService
}

func NewPublicProductController(service model.ProductService)*ProductController{
	return &ProductController{service: service}
}

func (ctrl *ProductController) GetAllMenu(c *fiber.Ctx) error {
	products, err := ctrl.service.GetAllProducts()
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	} 

	res := model.ToProductResponseList(products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
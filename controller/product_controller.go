package controller

import (
	"go-fiber-pos/model"
	"github.com/gofiber/fiber/v2"

)

type ProductController struct {
	service model.ProductService
}

func NewProductController(service model.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (ctrl *ProductController) Create(c *fiber.Ctx) error {
	// Ambil store_id dari token JWT (Hasil kerja Middleware)
	

	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	product, err := ctrl.service.CreateProduct(req)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produk berhasil dibuat",
		"data":    product,
	})
}

func (ctrl *ProductController) GetAll(c *fiber.Ctx) error {
	
	categories, err := ctrl.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": categories,
	})
}
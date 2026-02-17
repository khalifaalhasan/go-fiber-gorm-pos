package controller

import (
	"go-fiber-pos/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController struct {
	service model.ProductService
}

func NewProductController(service model.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (ctrl *ProductController) Create(c *fiber.Ctx) error {
	// Ambil store_id dari token JWT (Hasil kerja Middleware)
	storeIDContext := c.Locals("store_id").(uuid.UUID)

	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	Product, err := ctrl.service.CreateProduct(storeIDContext, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produk berhasil dibuat",
		"data":    Product,
	})
}

func (ctrl *ProductController) GetAll(c *fiber.Ctx) error {
	storeIDContext := c.Locals("store_id").(uuid.UUID)

	categories, err := ctrl.service.GetProductsByStore(storeIDContext)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": categories,
	})
}
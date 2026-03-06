package product

import (
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	service ProductService
}

func NewProductController(service ProductService) *ProductController {
	return &ProductController{service: service}
}

// Create godoc
// @Summary      Create a new product
// @Description  Create a new product with details (Admin only).
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body CreateProductRequest true "Create Product Request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products [post]
func (ctrl *ProductController) Create(c *fiber.Ctx) error {

	var req CreateProductRequest
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

// GetAll godoc
// @Summary      Get all products
// @Description  Retrieve a list of all products (Admin only).
// @Tags         products
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products [get]
func (ctrl *ProductController) GetAll(c *fiber.Ctx) error {
	
	products, err := ctrl.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	res := ToProductResponseList(products)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res,
	})
}
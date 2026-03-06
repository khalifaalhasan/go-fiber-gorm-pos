package order

import (
	"errors"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderController struct {
	service OrderService
}

func NewOrderController(service OrderService) *OrderController {
	return &OrderController{service: service}
}

// Checkout godoc
// @Summary      Create a new order (Checkout)
// @Description  Create a new order with items and apply voucher if valid (Admin only).
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body CheckoutRequest true "Checkout Request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/orders/checkout [post]
func (ctrl *OrderController) Checkout(c *fiber.Ctx) error {
	var req CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	order, err := ctrl.service.Checkout(req)
	if err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validasi gagal", "details": valErr.Error()})
		}
		// Petakan sentinel errors ke HTTP status yang tepat
		if errors.Is(err, core.ErrInsufficientStock) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, core.ErrVoucherInvalid) || errors.Is(err, core.ErrVoucherMinOrder) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, core.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order berhasil dibuat",
		"data":    order,
	})
}

// GetAll godoc
// @Summary      Get all orders
// @Description  Retrieve a list of all orders (Admin only).
// @Tags         orders
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/orders [get]
func (ctrl *OrderController) GetAll(c *fiber.Ctx) error {
	orders, err := ctrl.service.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": orders})
}

// GetByID godoc
// @Summary      Get order by ID
// @Description  Retrieve the details of a specific order by its ID (Admin only).
// @Tags         orders
// @Param        id   path      string  true  "Order UUID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/orders/{id} [get]
func (ctrl *OrderController) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID order tidak valid"})
	}

	order, err := ctrl.service.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": order})
}

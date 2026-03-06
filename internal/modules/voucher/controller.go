package voucher

import (
	"errors"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VoucherController struct {
	service VoucherService
}

func NewVoucherController(service VoucherService) *VoucherController {
	return &VoucherController{service: service}
}

// Create godoc
// @Summary      Create a new voucher
// @Description  Create a new discount voucher (Admin only).
// @Tags         vouchers
// @Accept       json
// @Produce      json
// @Param        request body CreateVoucherRequest true "Create Voucher Request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/vouchers [post]
func (ctrl *VoucherController) Create(c *fiber.Ctx) error {
	var req CreateVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	voucher, err := ctrl.service.CreateVoucher(req)
	if err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validasi gagal", "details": valErr.Error()})
		}
		if errors.Is(err, core.ErrAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		if errors.Is(err, core.ErrVoucherInvalid) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Voucher berhasil dibuat",
		"data":    voucher,
	})
}

// GetAll godoc
// @Summary      Get all vouchers
// @Description  Retrieve a list of all active/inactive vouchers (Admin only).
// @Tags         vouchers
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/vouchers [get]
func (ctrl *VoucherController) GetAll(c *fiber.Ctx) error {
	vouchers, err := ctrl.service.GetAllVouchers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": vouchers})
}

// Delete godoc
// @Summary      Delete a voucher
// @Description  Delete a voucher by its ID (Admin only).
// @Tags         vouchers
// @Param        id   path      string  true  "Voucher UUID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/vouchers/{id} [delete]
func (ctrl *VoucherController) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID voucher tidak valid"})
	}

	if err := ctrl.service.DeleteVoucher(id); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Voucher berhasil dihapus"})
}

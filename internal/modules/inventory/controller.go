package inventory

import (
	"errors"

	"go-fiber-pos/internal/core"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryController struct {
	service InventoryService
}

func NewInventoryController(service InventoryService) *InventoryController {
	return &InventoryController{service: service}
}

func (c *InventoryController) GetStockByProductID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("productId")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID produk tidak valid"})
	}

	res, err := c.service.GetStockByProductID(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Stok tidak ditemukan"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil stok"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil stok",
		"data":    res,
	})
}

func (c *InventoryController) AdjustStock(ctx *fiber.Ctx) error {
	var req AdjustStockRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := c.service.AdjustStock(ctx.Context(), req); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Stok tidak ditemukan"})
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Gagal menyesuaikan stok: " + err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menyesuaikan stok",
	})
}

func (c *InventoryController) GetMovements(ctx *fiber.Ctx) error {
	idParam := ctx.Params("productId")
	productID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID produk tidak valid"})
	}

	res, err := c.service.GetMovements(ctx.Context(), productID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Stok tidak ditemukan"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil riwayat pergerakan stok"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil riwayat pergerakan stok",
		"data":    res,
	})
}

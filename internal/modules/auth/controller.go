package auth

import (
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service AuthService
}

func NewAuthController(service AuthService) *AuthController {
	return &AuthController{service: service}
}

// Register godoc
// @Summary      Register a new admin
// @Description  Create a new admin account for the POS system.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register Request"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := ctrl.service.Register(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil, silakan login",
	})
}

// Login godoc
// @Summary      Login as admin
// @Description  Authenticate admin and return JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Request"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	token, err := ctrl.service.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
	})
}
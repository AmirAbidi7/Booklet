package controllers

import (
	"log/slog"

	"backend/internal/dtos"
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
	logger      *slog.Logger
}

func NewAuthController(authService *services.AuthService, logger *slog.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      logger,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req dtos.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := ac.authService.Register(&req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		ac.logger.Error("Registration failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process registration",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := ac.authService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (ac *AuthController) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	user, err := ac.authService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}

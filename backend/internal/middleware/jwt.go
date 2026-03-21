package middleware

import (
	"strings"

	"backend/internal/auth"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		tokenString := tokenParts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store user information in context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(int)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role information not found",
			})
		}

		// Assuming Admin role = 1 based on your models
		if role != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		return c.Next()
	}
}

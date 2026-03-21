package server

import (
	"backend/internal/controllers"
	"backend/internal/middleware"
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	authService := services.NewAuthService(s.DB)
	authController := controllers.NewAuthController(authService, s.Logger)

	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	auth := s.App.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	api := s.App.Group("/api", middleware.JWTProtected())
	api.Get("/me", authController.Me)

	admin := s.App.Group("/admin", middleware.JWTProtected(), middleware.RequireAdmin())
	admin.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Admin only route"})
	})
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

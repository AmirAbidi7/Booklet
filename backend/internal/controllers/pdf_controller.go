package controllers

import (
	"log/slog"
	"net/http"

	"backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type PdfController struct {
	PdfService *services.PdfService
	Logger     *slog.Logger
}

func NewPdfController(pdfService *services.PdfService, Logger *slog.Logger) *PdfController {
	return &PdfController{
		pdfService,
		Logger,
	}
}

func (pc PdfController) displayPdfs(c fiber.Ctx) error {
	userID, exists := c.Locals("user_id").(uint)

	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}
	pdfs, err := pc.PdfService.GetPdfs(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"List": pdfs,
	})
}

package dtos

import "backend/internal/models"

type PdfRequest struct {
	Title string `json:"title" validate:"required,min=4,max=255"`
	Pages uint   `json:"pages" validate:"omitempty,min=1"`
}

type PdfResponse struct {
	ID          uint
	Title       string
	Pages       uint
	currentPage uint
}

func (pdfR PdfRequest) toPdf() (*models.Pdf, *models.UserPdf) {
	pdf := models.Pdf{}
	pdf.Pages = pdfR.Pages

	userPdf := models.UserPdf{}
	userPdf.Title = pdfR.Title

	return &pdf, &userPdf
}

func pdfToResponse(pdf models.Pdf, userPdf models.UserPdf) *PdfResponse {
	return &PdfResponse{
		pdf.ID,
		userPdf.Title,
		pdf.Pages,
		userPdf.CurrentPage,
	}
}

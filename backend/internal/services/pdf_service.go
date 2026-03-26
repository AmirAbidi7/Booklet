package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"backend/internal/models"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"gorm.io/gorm"
)

type PdfService struct {
	DB     *gorm.DB
	Client *azblob.Client
}

func NewPdfService(DB *gorm.DB, Client *azblob.Client) *PdfService {
	return &PdfService{
		DB,
		Client,
	}
}

func (ps PdfService) AddPdf(pdf models.Pdf, file *multipart.FileHeader, userID uint) error {
	ctx := context.Background()
	contentType := file.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return fmt.Errorf("invalid content type: %s. Expected application/pdf", contentType)
	}
	if file.Size >= 50*1024*1024 {
		return fmt.Errorf("file size cannot surpass 50MB")
	}
	openedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer openedFile.Close()

	hasher := sha256.New()
	if _, err = io.Copy(hasher, openedFile); err != nil {
		return err
	}
	pdf.Checksum = hex.EncodeToString(hasher.Sum(nil))

	if _, err := openedFile.Seek(0, 0); err != nil {
		return err
	}

	return ps.DB.Transaction(func(tx *gorm.DB) error {
		existing, err := gorm.G[models.Pdf](tx).Where("checksum = ?", pdf.Checksum).First(ctx)
		if err != nil {

			if err = gorm.G[models.Pdf](tx).Create(ctx, &pdf); err != nil {
				return err
			}

			containerName := "pdf-storage"
			_, err = ps.Client.UploadStream(ctx, containerName,
				strconv.FormatInt(int64(pdf.ID), 10),
				openedFile, nil)
			if err != nil {
				return err
			}
		} else {
			pdf.ID = existing.ID
		}

		userPdf := models.UserPdf{
			PdfID:       pdf.ID,
			UserID:      userID,
			CurrentPage: 1,
		}

		return gorm.G[models.UserPdf](tx).Create(ctx, &userPdf)
	})
}

func (ps PdfService) DeletePdf(pdfID uint, userID uint) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		ctx := context.Background()
		_, err := gorm.G[models.UserPdf](tx).Where("pdf_id = ? AND user_id = ? ", pdfID, userID).Delete(ctx)
		if err != nil {
			return err
		}
		count, err := gorm.G[models.UserPdf](tx).Where("pdf_id = ?", pdfID).Count(ctx, "pdf_id")
		if err != nil {
			return err
		}
		if count == 0 {
			_, err := gorm.G[models.Pdf](tx).Where("id = ? ", pdfID).Delete(ctx)
			if err != nil {
				return err
			}
			_, err = ps.Client.DeleteBlob(ctx, "pdf-storage", strconv.FormatInt(int64(pdfID), 10), nil)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (ps PdfService) UpdateProgress(userPdf models.UserPdf) error {
	ctx := context.Background()
	_, err := gorm.G[models.UserPdf](ps.DB).Where("user_id = ? AND pdf_id = ? ", userPdf.UserID, userPdf.PdfID).First(ctx)
	if err != nil {
		return err
	}

	return ps.DB.Save(&userPdf).Error
}

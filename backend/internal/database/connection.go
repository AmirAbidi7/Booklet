package database

import (
	"os"

	"backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetConn() (*gorm.DB, error) {
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.SetupJoinTable(&models.Pdf{}, "Users", &models.UserPdf{})
	if err != nil {
		return nil, err
	}

	err = db.SetupJoinTable(&models.User{}, "Pdfs", &models.UserPdf{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(models.User{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(models.Pdf{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(models.UserPdf{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

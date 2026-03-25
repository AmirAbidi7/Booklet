package server

import (
	"fmt"
	"log/slog"
	"os"

	"backend/internal/database"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type FiberServer struct {
	*fiber.App
	*gorm.DB
	*slog.Logger
}

func New() *FiberServer {
	DB, err := database.GetConn()
	if err != nil {
		fmt.Printf("error connecting to the database: %v", err)
		os.Exit(1)
	}
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "backend",
			AppName:      "backend",
		}),
		DB:     DB,
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	return server
}

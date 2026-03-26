package server

import (
	"fmt"
	"log/slog"
	"os"

	"backend/internal/database"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type FiberServer struct {
	*fiber.App
	*gorm.DB
	*slog.Logger
	Client *azblob.Client
}

func New() *FiberServer {
	DB, err := database.GetConn()
	if err != nil {
		fmt.Printf("error connecting to the database: %v", err)
		os.Exit(1)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Printf("error authenticated to azure: %v", err)
		os.Exit(1)
	}

	accountURI := "https://bookletstorage.blob.core.windows.net/"
	client, err := azblob.NewClient(accountURI, cred, nil)
	if err != nil {
		fmt.Printf("error authenticated to azure: %v", err)
		os.Exit(1)

	}
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "backend",
			AppName:      "backend",
		}),
		DB:     DB,
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Client: client,
	}

	return server
}

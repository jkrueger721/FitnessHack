package server

import (
	"github.com/gofiber/fiber/v2"

	"fitness-hack/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "fitness-hack",
			AppName:      "fitness-hack",
		}),

		db: database.New(),
	}

	return server
}

package main

import (
	"github.com/antiloger/nhostel-go/database"
	"github.com/antiloger/nhostel-go/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Connect()

	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024,
	})

	app.Static("/uploads", "./uploads")

	// middlewares

	app.Use(logger.New())
	app.Use(cors.New())

	routers.RegisterRoutes(app)

	app.Listen(":3000")
}

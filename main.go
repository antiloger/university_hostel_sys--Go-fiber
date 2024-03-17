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

	app := fiber.New()

	// middlewares
	app.Use(logger.New())
	app.Use(cors.New())

	routers.RegisterRoutes(app)

	app.Listen(":3000")
}

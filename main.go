package main

import (
	"github.com/antiloger/nhostel-go/config"
	"github.com/antiloger/nhostel-go/handler"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", handler.Hello)
	app.Post("/login", handler.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Jwt_Secret)},
	}))
	app.Get("/atest", handler.Authtest)
	app.Post("/users", handler.Insertuser)
	app.Get("/allusers", handler.Getusers)
	app.Get("/home", handler.HomeLoad)
	app.Get("/homeloc", handler.HomeLoadloc)

	app.Listen(":3000")
}

package handler

import (
	"time"

	"github.com/antiloger/nhostel-go/config"
	"github.com/antiloger/nhostel-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *fiber.Ctx) error {
	loginreq := new(models.LoginRequest)
	if err := c.BodyParser(loginreq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if loginreq.Email == "test" || loginreq.Password == "test" {
		u := models.User{
			ID:       1,
			Email:    loginreq.Email,
			Password: loginreq.Password,
			Role:     "admin",
		}

		day := time.Hour * 24

		claims := jwt.MapClaims{
			"id":    u.ID,
			"email": u.Email,
			"role":  u.Role,
			"exp":   time.Now().Add(day * 1).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(config.Jwt_Secret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(models.LoginResponse{
			Token: t,
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Password required",
		})
	}
}

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// Home & Search Handler

func HomeLoad(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func HomeLoadloc(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Search(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Hosteldetails(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// user: student handler

func Studentsignup(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// user: hostel owner handler

func Hostelownersignup(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Hostelregister(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Hostelupdate(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Hosteldelete(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

// user: admin handler

func Adminregister(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func AdminLogin(c *fiber.Ctx) error {
	return c.SendString("admin login")
}

func Hostelapprovetable(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Hostelowapprovetable(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Studentapprovetable(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

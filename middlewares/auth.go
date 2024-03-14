package middlewares

import (
	"github.com/antiloger/nhostel-go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RoleMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Request().Header.Peek("Authorization")
		if token == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(string(token), &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Jwt_Secret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		if !parsedToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		userRole := claims["role"].(string)
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}

		return c.Next()
	}
}

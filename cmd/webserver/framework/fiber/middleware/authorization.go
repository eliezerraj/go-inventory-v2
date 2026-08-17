package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AuthorizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}
		return c.Next()
	}
}

func HeaderMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Accepts("application/json")

		securityHeaders := map[string]string{
			"X-XSS-Protection":          "1; mode=block",
			"X-Content-Type-Options":    "nosniff",
			"X-Download-Options":        "noopen",
			"Strict-Transport-Security": "max-age=5184000",
			"X-Frame-Options":           "SAMEORIGIN",
			"X-DNS-Prefetch-Control":    "off",
		}

		corsHeaders := map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
			"Access-Control-Allow-Methods":     "POST, GET, PUT, DELETE",
		}

		for key, value := range securityHeaders {
			c.Set(key, value)
		}
		for key, value := range corsHeaders {
			c.Set(key, value)
		}

		return c.Next()
	}
}
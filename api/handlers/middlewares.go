package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
)

func (s *Server) initializeMiddlewares() {
	s.Router.Use(cors.New())
	// JWT Middleware
	s.Router.Use("/api", jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		},
		SigningKey: []byte(os.Getenv("API_JWT_SECRET")),
	}))
}

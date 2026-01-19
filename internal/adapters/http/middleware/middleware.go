package middleware

import (
	"spsc-loaneasy/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Setup configures all middlewares for the application
func Setup(app *fiber.App, cfg *config.Config) {
	// Recover middleware - catches panics
	app.Use(recover.New())

	// Logger middleware
	if cfg.IsDev() {
		app.Use(logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
		}))
	} else {
		app.Use(logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		}))
	}

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))
}

// CustomErrorHandler handles errors globally
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}

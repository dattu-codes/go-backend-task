package middleware

import (
	"time"

	"go-backend-task/internal/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger returns a Fiber middleware that logs HTTP requests with latency.
// It integrates with our global Zap logger and records method, path, status code, IP, request ID, and execution duration.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Execute next handlers in the stack
		err := c.Next()

		duration := time.Since(start)
		reqID, _ := c.Locals("requestId").(string)

		// Extract correct HTTP status code (handling standard router responses and explicit Fiber errors)
		status := c.Response().StatusCode()
		if err != nil {
			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		// Log structured info
		logger.Log.Info("HTTP Request Processed",
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		)

		return err
	}
}

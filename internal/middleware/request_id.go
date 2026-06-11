package middleware

import (
	"crypto/rand"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// RequestID returns a Fiber middleware that injects or propagates a Request ID.
// It stores the Request ID in the locals context and returns it in the X-Request-ID response header.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			var err error
			reqID, err = generateUUID()
			if err != nil {
				// Fallback if random generator fails
				reqID = "unknown-request-id"
			}
		}

		// Store in locals so subsequent middlewares (like Logger) and handlers can access it.
		c.Locals("requestId", reqID)

		// Set response header
		c.Set("X-Request-ID", reqID)

		return c.Next()
	}
}

// generateUUID creates a cryptographically secure RFC 4122 compliant UUID v4.
// This is a great interview talking point as it avoids external packages for UUID generation.
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Set the version to 4 (random)
	b[6] = (b[6] & 0x0f) | 0x40
	// Set the variant to RFC 4122
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

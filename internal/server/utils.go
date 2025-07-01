package server

import (
	"github.com/gofiber/fiber/v2"
)

// LogError logs an error with CloudWatch-compatible structured logging
func LogError(s *FiberServer, level, message string, err error, c *fiber.Ctx, metadata map[string]interface{}) {
	s.logError(level, message, err, c, metadata)
}

// LogDatabaseError logs database operation errors
func LogDatabaseError(s *FiberServer, operation string, err error, c *fiber.Ctx) {
	LogError(s, "ERROR", "Database operation failed", err, c, map[string]interface{}{
		"operation": operation,
		"component": "database",
	})
}

// LogCacheError logs cache operation errors
func LogCacheError(s *FiberServer, operation string, err error, c *fiber.Ctx) {
	LogError(s, "WARN", "Cache operation failed", err, c, map[string]interface{}{
		"operation": operation,
		"component": "cache",
	})
}

// LogAuthError logs authentication/authorization errors
func LogAuthError(s *FiberServer, message string, err error, c *fiber.Ctx) {
	LogError(s, "WARN", message, err, c, map[string]interface{}{
		"component": "authentication",
	})
}

// LogValidationError logs request validation errors
func LogValidationError(s *FiberServer, field string, err error, c *fiber.Ctx) {
	LogError(s, "INFO", "Validation error", err, c, map[string]interface{}{
		"component": "validation",
		"field":     field,
	})
}

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"fitness-hack/internal/database"
)

type FiberServer struct {
	*fiber.App
	db    database.Service
	cache *redis.Client
}

// CloudWatchLogEntry represents a structured log entry for AWS CloudWatch
type CloudWatchLogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Error      string                 `json:"error,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	UserID     string                 `json:"user_id,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Path       string                 `json:"path,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Latency    string                 `json:"latency,omitempty"`
	IP         string                 `json:"ip,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Stack      []string               `json:"stack,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// logError logs structured errors for CloudWatch
func (s *FiberServer) logError(level, message string, err error, c *fiber.Ctx, metadata map[string]interface{}) {
	entry := CloudWatchLogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Metadata:  metadata,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	if c != nil {
		entry.RequestID = c.Get("X-Request-ID")
		entry.Method = c.Method()
		entry.Path = c.Path()
		entry.IP = c.IP()
		entry.UserAgent = c.Get("User-Agent")

		// Extract user ID from JWT if available
		if userID, err := getUserIDFromJWT(c); err == nil {
			entry.UserID = userID
		}
	}

	// Add stack trace for errors
	if level == "ERROR" || level == "FATAL" {
		entry.Stack = getStackTrace()
	}

	// Output as JSON
	if logData, err := json.Marshal(entry); err == nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(logData))
	}
}

// getStackTrace returns the current stack trace
func getStackTrace() []string {
	var stack []string
	for i := 1; i < 10; i++ {
		if pc, file, line, ok := runtime.Caller(i); ok {
			fn := runtime.FuncForPC(pc)
			stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
		}
	}
	return stack
}

// errorHandler middleware for catching and logging errors
func (s *FiberServer) errorHandler(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	latency := time.Since(start)

	// Log errors
	if err != nil {
		s.logError("ERROR", "Request failed", err, c, map[string]interface{}{
			"latency": latency.String(),
		})
		return err
	}

	// Log 4xx and 5xx status codes
	status := c.Response().StatusCode()
	if status >= 400 {
		s.logError("WARN", fmt.Sprintf("HTTP %d", status), nil, c, map[string]interface{}{
			"status_code": status,
			"latency":     latency.String(),
		})
	}

	return nil
}

func New() *FiberServer {
	// Redis config from env
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if dbInt, err := strconv.Atoi(dbStr); err == nil {
			redisDB = dbInt
		}
	}
	cache := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "fitness-hack",
			AppName:      "fitness-hack",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				// We'll set up the error handler after server creation
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
			},
		}),
		db:    database.New(),
		cache: cache,
	}

	// Add error logging middleware first
	server.App.Use(server.errorHandler)

	// Add request logging middleware
	server.App.Use(logger.New(logger.Config{
		Format:     "${time} | ${method} | ${path} | ${status} | ${latency} | ${ip} | ${userAgent}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stdout,
	}))

	return server
}

// getUserIDFromJWT extracts the user_id from the JWT claims in the Fiber context
func getUserIDFromJWT(c *fiber.Ctx) (string, error) {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return "", errors.New("invalid or missing JWT token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid JWT claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}
	return userID, nil
}

// SetCache sets a value in Redis with expiration (in seconds)
func (s *FiberServer) SetCache(ctx context.Context, key string, value string, expiration time.Duration) error {
	return s.cache.Set(ctx, key, value, expiration).Err()
}

// GetCache gets a value from Redis
func (s *FiberServer) GetCache(ctx context.Context, key string) (string, error) {
	return s.cache.Get(ctx, key).Result()
}

// DeleteCache deletes a key from Redis
func (s *FiberServer) DeleteCache(ctx context.Context, key string) error {
	return s.cache.Del(ctx, key).Err()
}

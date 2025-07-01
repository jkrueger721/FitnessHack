package server

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	// Health and basic routes
	s.App.Get("/", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)

	// API v1 group
	api := s.App.Group("/api/v1")

	// Public routes (no JWT required)
	api.Post("/auth/login", s.loginUser)
	api.Post("/users", s.createUser)

	// JWT Middleware for all other /api/v1 routes
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		},
	}))

	// Protected Users routes
	users := api.Group("/users")
	users.Get("/", s.listUsers)
	users.Get("/:id", s.getUser)
	users.Put("/:id", s.updateUser)
	users.Delete("/:id", s.deleteUser)

	// Workouts routes
	workouts := api.Group("/workouts")
	workouts.Post("/", s.createWorkout)
	workouts.Get("/", s.listWorkouts)
	workouts.Get("/:id", s.getWorkout)
	workouts.Put("/:id", s.updateWorkout)
	workouts.Delete("/:id", s.deleteWorkout)

	// Exercises routes
	exercises := api.Group("/exercises")
	exercises.Post("/", s.createExercise)
	exercises.Get("/", s.listExercises)
	exercises.Get("/:id", s.getExercise)
	exercises.Put("/:id", s.updateExercise)
	exercises.Delete("/:id", s.deleteExercise)

	// Workout exercises routes
	workoutExercises := api.Group("/workout-exercises")
	workoutExercises.Post("/", s.createWorkoutExercise)
	workoutExercises.Get("/", s.listWorkoutExercises)
	workoutExercises.Get("/:id", s.getWorkoutExercise)
	workoutExercises.Put("/:id", s.updateWorkoutExercise)
	workoutExercises.Delete("/:id", s.deleteWorkoutExercise)

	// Workout sessions routes
	workoutSessions := api.Group("/workout-sessions")
	workoutSessions.Post("/", s.createWorkoutSession)
	workoutSessions.Get("/", s.listWorkoutSessions)
	workoutSessions.Get("/:id", s.getWorkoutSession)
	workoutSessions.Put("/:id", s.updateWorkoutSession)
	workoutSessions.Delete("/:id", s.deleteWorkoutSession)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

// Helper function to get pagination parameters
func getPaginationParams(c *fiber.Ctx) (limit, offset int) {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, _ = strconv.Atoi(limitStr)
	offset, _ = strconv.Atoi(offsetStr)

	// Set reasonable defaults and limits
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

// Helper function to create error response
func errorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// Helper function to create success response
func successResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"data": data,
	})
}

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fitness-hack/internal/database"

	"github.com/gofiber/fiber/v2"
)

// Cache key helpers
func workoutCacheKey(id string) string {
	return fmt.Sprintf("workout:%s", id)
}

func workoutsListCacheKey(limit, offset int) string {
	return fmt.Sprintf("workouts:list:%d:%d", limit, offset)
}

// Helper to convert database workout to response model
func workoutToResponse(workout *database.Workouts) database.WorkoutResponse {
	return database.WorkoutResponse{
		ID:              workout.Id,
		UserID:          workout.User_id,
		Name:            workout.Name,
		Description:     workout.Description,
		DurationMinutes: workout.Duration_minutes,
		CreatedAt:       workout.Created_at,
		UpdatedAt:       workout.Updated_at,
	}
}

// Workouts handlers
func (s *FiberServer) createWorkout(c *fiber.Ctx) error {
	var req database.CreateWorkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id").(string)

	// Create database workout
	workout := database.Workouts{
		User_id:          userID,
		Name:             req.Name,
		Description:      req.Description,
		Duration_minutes: req.DurationMinutes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdWorkout, err := s.db.CreateWorkout(ctx, &workout)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create workout: "+err.Error())
	}

	// Invalidate workouts list cache
	s.cache.Del(ctx, "workouts:list:*")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": workoutToResponse(createdWorkout),
	})
}

func (s *FiberServer) getWorkout(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := workoutCacheKey(id)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var workout database.Workouts
		if json.Unmarshal([]byte(cachedData), &workout) == nil {
			return successResponse(c, workoutToResponse(&workout))
		}
	}

	// Get from database
	workout, err := s.db.GetWorkoutByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Workout not found")
	}

	// Cache the workout data
	if workoutData, err := json.Marshal(workout); err == nil {
		s.SetCache(ctx, cacheKey, string(workoutData), 10*time.Minute)
	}

	return successResponse(c, workoutToResponse(workout))
}

func (s *FiberServer) listWorkouts(c *fiber.Ctx) error {
	limit, offset := getPaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := workoutsListCacheKey(limit, offset)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var workouts []database.Workouts
		if json.Unmarshal([]byte(cachedData), &workouts) == nil {
			// Convert to response models
			responses := make([]database.WorkoutResponse, len(workouts))
			for i, workout := range workouts {
				responses[i] = workoutToResponse(&workout)
			}
			return successResponse(c, responses)
		}
	}

	// Get from database
	workouts, err := s.db.ListWorkouts(ctx, limit, offset)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to fetch workouts: "+err.Error())
	}

	// Cache the workouts data
	if workoutsData, err := json.Marshal(workouts); err == nil {
		s.SetCache(ctx, cacheKey, string(workoutsData), 10*time.Minute)
	}

	// Convert to response models
	responses := make([]database.WorkoutResponse, len(workouts))
	for i, workout := range workouts {
		responses[i] = workoutToResponse(&workout)
	}

	return successResponse(c, responses)
}

func (s *FiberServer) updateWorkout(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout ID is required")
	}

	var req database.UpdateWorkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing workout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingWorkout, err := s.db.GetWorkoutByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Workout not found")
	}

	// Update fields if provided
	if req.Name != nil {
		existingWorkout.Name = *req.Name
	}
	if req.Description != nil {
		existingWorkout.Description = *req.Description
	}
	if req.DurationMinutes != nil {
		existingWorkout.Duration_minutes = *req.DurationMinutes
	}
	existingWorkout.Updated_at = time.Now()

	updatedWorkout, err := s.db.UpdateWorkout(ctx, existingWorkout)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to update workout: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, workoutCacheKey(id))
	s.cache.Del(ctx, "workouts:list:*")

	return successResponse(c, workoutToResponse(updatedWorkout))
}

func (s *FiberServer) deleteWorkout(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.DeleteWorkout(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to delete workout: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, workoutCacheKey(id))
	s.cache.Del(ctx, "workouts:list:*")

	return c.Status(fiber.StatusNoContent).Send(nil)
}

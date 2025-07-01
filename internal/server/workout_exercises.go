package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fitness-hack/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// Cache key helpers
func workoutExerciseCacheKey(id string) string {
	return fmt.Sprintf("workout_exercise:%s", id)
}

func workoutExercisesListCacheKey(limit, offset int) string {
	return fmt.Sprintf("workout_exercises:list:%d:%d", limit, offset)
}

// Helper to convert database workout exercise to response model
func workoutExerciseToResponse(we *database.Workout_exercises) database.WorkoutExerciseResponse {
	// Handle type assertions safely
	var workoutId, exerciseId string
	if we.Workout_id != "" {
		workoutId = we.Workout_id
	}
	if we.Exercise_id != "" {
		exerciseId = we.Exercise_id
	}

	return database.WorkoutExerciseResponse{
		ID:              we.Id,
		WorkoutID:       workoutId,
		ExerciseID:      exerciseId,
		Sets:            we.Sets,
		Reps:            we.Reps,
		WeightKg:        we.Weight_kg.InexactFloat64(),
		DurationSeconds: we.Duration_seconds,
		OrderIndex:      we.Order_index,
		RestSeconds:     we.Rest_seconds,
		Notes:           we.Notes,
		CreatedAt:       we.Created_at,
	}
}

// Workout exercises handlers
func (s *FiberServer) createWorkoutExercise(c *fiber.Ctx) error {
	var req database.CreateWorkoutExerciseRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Create database workout exercise
	workoutExercise := database.Workout_exercises{
		Workout_id:       req.WorkoutID,
		Exercise_id:      req.ExerciseID,
		Sets:             req.Sets,
		Reps:             req.Reps,
		Weight_kg:        decimal.NewFromFloat(req.WeightKg),
		Duration_seconds: req.DurationSeconds,
		Order_index:      req.OrderIndex,
		Rest_seconds:     req.RestSeconds,
		Notes:            req.Notes,
		Created_at:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdWorkoutExercise, err := s.db.CreateWorkoutExercise(ctx, &workoutExercise)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create workout exercise: "+err.Error())
	}

	// Invalidate workout exercises list cache
	s.cache.Del(ctx, "workout_exercises:list:*")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": workoutExerciseToResponse(createdWorkoutExercise),
	})
}

func (s *FiberServer) getWorkoutExercise(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout exercise ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := workoutExerciseCacheKey(id)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var workoutExercise database.Workout_exercises
		if json.Unmarshal([]byte(cachedData), &workoutExercise) == nil {
			return successResponse(c, workoutExerciseToResponse(&workoutExercise))
		}
	}

	// Get from database
	workoutExercise, err := s.db.GetWorkoutExerciseByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Workout exercise not found")
	}

	// Cache the workout exercise data
	if workoutExerciseData, err := json.Marshal(workoutExercise); err == nil {
		s.SetCache(ctx, cacheKey, string(workoutExerciseData), 10*time.Minute)
	}

	return successResponse(c, workoutExerciseToResponse(workoutExercise))
}

func (s *FiberServer) listWorkoutExercises(c *fiber.Ctx) error {
	limit, offset := getPaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := workoutExercisesListCacheKey(limit, offset)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var workoutExercises []database.Workout_exercises
		if json.Unmarshal([]byte(cachedData), &workoutExercises) == nil {
			// Convert to response models
			responses := make([]database.WorkoutExerciseResponse, len(workoutExercises))
			for i, we := range workoutExercises {
				responses[i] = workoutExerciseToResponse(&we)
			}
			return successResponse(c, responses)
		}
	}

	// Get from database
	workoutExercises, err := s.db.ListWorkoutExercises(ctx, limit, offset)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to fetch workout exercises: "+err.Error())
	}

	// Cache the workout exercises data
	if workoutExercisesData, err := json.Marshal(workoutExercises); err == nil {
		s.SetCache(ctx, cacheKey, string(workoutExercisesData), 10*time.Minute)
	}

	// Convert to response models
	responses := make([]database.WorkoutExerciseResponse, len(workoutExercises))
	for i, we := range workoutExercises {
		responses[i] = workoutExerciseToResponse(&we)
	}

	return successResponse(c, responses)
}

func (s *FiberServer) updateWorkoutExercise(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout exercise ID is required")
	}

	var req database.UpdateWorkoutExerciseRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing workout exercise
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingWorkoutExercise, err := s.db.GetWorkoutExerciseByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Workout exercise not found")
	}

	// Update fields if provided
	if req.WorkoutID != nil {
		existingWorkoutExercise.Workout_id = *req.WorkoutID
	}
	if req.ExerciseID != nil {
		existingWorkoutExercise.Exercise_id = *req.ExerciseID
	}
	if req.Sets != nil {
		existingWorkoutExercise.Sets = *req.Sets
	}
	if req.Reps != nil {
		existingWorkoutExercise.Reps = *req.Reps
	}
	if req.WeightKg != nil {
		existingWorkoutExercise.Weight_kg = decimal.NewFromFloat(*req.WeightKg)
	}
	if req.DurationSeconds != nil {
		existingWorkoutExercise.Duration_seconds = *req.DurationSeconds
	}
	if req.OrderIndex != nil {
		existingWorkoutExercise.Order_index = *req.OrderIndex
	}
	if req.RestSeconds != nil {
		existingWorkoutExercise.Rest_seconds = *req.RestSeconds
	}
	if req.Notes != nil {
		existingWorkoutExercise.Notes = *req.Notes
	}

	updatedWorkoutExercise, err := s.db.UpdateWorkoutExercise(ctx, existingWorkoutExercise)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to update workout exercise: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, workoutExerciseCacheKey(id))
	s.cache.Del(ctx, "workout_exercises:list:*")

	return successResponse(c, workoutExerciseToResponse(updatedWorkoutExercise))
}

func (s *FiberServer) deleteWorkoutExercise(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout exercise ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.DeleteWorkoutExercise(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to delete workout exercise: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, workoutExerciseCacheKey(id))
	s.cache.Del(ctx, "workout_exercises:list:*")

	return c.Status(fiber.StatusNoContent).Send(nil)
}

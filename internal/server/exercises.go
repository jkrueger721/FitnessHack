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
func exerciseCacheKey(id string) string {
	return fmt.Sprintf("exercise:%s", id)
}

func exercisesListCacheKey(limit, offset int) string {
	return fmt.Sprintf("exercises:list:%d:%d", limit, offset)
}

// Helper to convert database exercise to response model
func exerciseToResponse(exercise *database.Exercises) database.ExerciseResponse {
	return database.ExerciseResponse{
		Id:               exercise.Id,
		Name:             exercise.Name,
		Description:      exercise.Description,
		Muscle_group:     exercise.Muscle_group,
		Equipment:        exercise.Equipment,
		Difficulty_level: exercise.Difficulty_level,
		Instructions:     exercise.Instructions,
		Created_at:       exercise.Created_at,
		Updated_at:       exercise.Updated_at,
	}
}

// Exercises handlers
func (s *FiberServer) createExercise(c *fiber.Ctx) error {
	var req database.CreateExerciseRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Create database exercise
	exercise := database.Exercises{
		Name:             req.Name,
		Description:      req.Description,
		Muscle_group:     req.Muscle_group,
		Equipment:        req.Equipment,
		Difficulty_level: req.Difficulty_level,
		Instructions:     req.Instructions,
		Created_at:       time.Now(),
		Updated_at:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdExercise, err := s.db.CreateExercise(ctx, &exercise)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create exercise: "+err.Error())
	}

	// Invalidate exercises list cache
	s.cache.Del(ctx, "exercises:list:*")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": exerciseToResponse(createdExercise),
	})
}

func (s *FiberServer) getExercise(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Exercise ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := exerciseCacheKey(id)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var exercise database.Exercises
		if json.Unmarshal([]byte(cachedData), &exercise) == nil {
			return successResponse(c, exerciseToResponse(&exercise))
		}
	}

	// Get from database
	exercise, err := s.db.GetExerciseByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Exercise not found")
	}

	// Cache the exercise data
	if exerciseData, err := json.Marshal(exercise); err == nil {
		s.SetCache(ctx, cacheKey, string(exerciseData), 10*time.Minute)
	}

	return successResponse(c, exerciseToResponse(exercise))
}

func (s *FiberServer) listExercises(c *fiber.Ctx) error {
	limit, offset := getPaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := exercisesListCacheKey(limit, offset)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var exercises []database.Exercises
		if json.Unmarshal([]byte(cachedData), &exercises) == nil {
			// Convert to response models
			responses := make([]database.ExerciseResponse, len(exercises))
			for i, exercise := range exercises {
				responses[i] = exerciseToResponse(&exercise)
			}
			return successResponse(c, responses)
		}
	}

	// Get from database
	exercises, err := s.db.ListExercises(ctx, limit, offset)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to fetch exercises: "+err.Error())
	}

	// Cache the exercises data
	if exercisesData, err := json.Marshal(exercises); err == nil {
		s.SetCache(ctx, cacheKey, string(exercisesData), 10*time.Minute)
	}

	// Convert to response models
	responses := make([]database.ExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = exerciseToResponse(&exercise)
	}

	return successResponse(c, responses)
}

func (s *FiberServer) updateExercise(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Exercise ID is required")
	}

	var req database.UpdateExerciseRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing exercise
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingExercise, err := s.db.GetExerciseByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Exercise not found")
	}

	// Update fields if provided
	if req.Name != nil {
		existingExercise.Name = *req.Name
	}
	if req.Description != nil {
		existingExercise.Description = req.Description
	}
	if req.Muscle_group != nil {
		existingExercise.Muscle_group = req.Muscle_group
	}
	if req.Equipment != nil {
		existingExercise.Equipment = req.Equipment
	}
	if req.Difficulty_level != nil {
		existingExercise.Difficulty_level = req.Difficulty_level
	}
	if req.Instructions != nil {
		existingExercise.Instructions = req.Instructions
	}
	existingExercise.Updated_at = time.Now()

	updatedExercise, err := s.db.UpdateExercise(ctx, existingExercise)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to update exercise: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, exerciseCacheKey(id))
	s.cache.Del(ctx, "exercises:list:*")

	return successResponse(c, exerciseToResponse(updatedExercise))
}

func (s *FiberServer) deleteExercise(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Exercise ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.DeleteExercise(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to delete exercise: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, exerciseCacheKey(id))
	s.cache.Del(ctx, "exercises:list:*")

	return c.Status(fiber.StatusNoContent).Send(nil)
}

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
func workoutSessionCacheKey(id string) string {
	return fmt.Sprintf("workout_session:%s", id)
}

func workoutSessionsListCacheKey(limit, offset int) string {
	return fmt.Sprintf("workout_sessions:list:%d:%d", limit, offset)
}

// Helper to convert database workout session to response model
func workoutSessionToResponse(ws *database.Workout_sessions) database.WorkoutSessionResponse {
	return database.WorkoutSessionResponse{
		ID:              ws.Id,
		UserID:          ws.User_id,
		WorkoutID:       ws.Workout_id,
		Name:            ws.Name.(string),
		StartedAt:       ws.Started_at,
		CompletedAt:     &ws.Completed_at,
		DurationMinutes: ws.Duration_minutes,
		Notes:           ws.Notes,
		CreatedAt:       ws.Created_at,
		UpdatedAt:       ws.Updated_at,
	}
}

// Workout sessions handlers
func (s *FiberServer) createWorkoutSession(c *fiber.Ctx) error {
	var req database.CreateWorkoutSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get user ID from JWT token
	userID := c.Locals("user_id").(string)

	// Set default started_at if not provided
	startedAt := time.Now()
	if req.StartedAt != nil {
		startedAt = *req.StartedAt
	}

	// Create database workout session
	workoutSession := database.Workout_sessions{
		User_id:          userID,
		Workout_id:       req.WorkoutID,
		Name:             req.Name,
		Started_at:       startedAt,
		Completed_at:     *req.CompletedAt,
		Duration_minutes: req.DurationMinutes,
		Notes:            req.Notes,
		Created_at:       time.Now(),
		Updated_at:       time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdWorkoutSession, err := s.db.CreateWorkoutSession(ctx, &workoutSession)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create workout session: "+err.Error())
	}

	// Invalidate workout sessions list cache
	s.cache.Del(ctx, "workout_sessions:list:*")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": workoutSessionToResponse(createdWorkoutSession),
	})
}

func (s *FiberServer) getWorkoutSession(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout session ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := workoutSessionCacheKey(id)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var workoutSession database.Workout_sessions
		if json.Unmarshal([]byte(cachedData), &workoutSession) == nil {
			return successResponse(c, workoutSessionToResponse(&workoutSession))
		}
	}

	// Get from database
	workoutSession, err := s.db.GetWorkoutSessionByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Workout session not found")
	}

	// Cache the workout session data
	if workoutSessionData, err := json.Marshal(workoutSession); err == nil {
		s.SetCache(ctx, cacheKey, string(workoutSessionData), 10*time.Minute)
	}

	return successResponse(c, workoutSessionToResponse(workoutSession))
}

func (s *FiberServer) listWorkoutSessions(c *fiber.Ctx) error {
	limit, offset := getPaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := workoutSessionsListCacheKey(limit, offset)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var workoutSessions []database.Workout_sessions
		if json.Unmarshal([]byte(cachedData), &workoutSessions) == nil {
			// Convert to response models
			responses := make([]database.WorkoutSessionResponse, len(workoutSessions))
			for i, ws := range workoutSessions {
				responses[i] = workoutSessionToResponse(&ws)
			}
			return successResponse(c, responses)
		}
	}

	// Get from database
	workoutSessions, err := s.db.ListWorkoutSessions(ctx, limit, offset)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to fetch workout sessions: "+err.Error())
	}

	// Cache the workout sessions data
	if workoutSessionsData, err := json.Marshal(workoutSessions); err == nil {
		s.SetCache(ctx, cacheKey, string(workoutSessionsData), 10*time.Minute)
	}

	// Convert to response models
	responses := make([]database.WorkoutSessionResponse, len(workoutSessions))
	for i, ws := range workoutSessions {
		responses[i] = workoutSessionToResponse(&ws)
	}

	return successResponse(c, responses)
}

func (s *FiberServer) updateWorkoutSession(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout session ID is required")
	}

	var req database.UpdateWorkoutSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing workout session
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingWorkoutSession, err := s.db.GetWorkoutSessionByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Workout session not found")
	}

	// Update fields if provided
	if req.WorkoutID != nil {
		existingWorkoutSession.Workout_id = *req.WorkoutID
	}
	if req.Name != nil {
		existingWorkoutSession.Name = *req.Name
	}
	if req.StartedAt != nil {
		existingWorkoutSession.Started_at = *req.StartedAt
	}
	if req.CompletedAt != nil {
		existingWorkoutSession.Completed_at = *req.CompletedAt
	}
	if req.DurationMinutes != nil {
		existingWorkoutSession.Duration_minutes = *req.DurationMinutes
	}
	if req.Notes != nil {
		existingWorkoutSession.Notes = *req.Notes
	}
	existingWorkoutSession.Updated_at = time.Now()

	updatedWorkoutSession, err := s.db.UpdateWorkoutSession(ctx, existingWorkoutSession)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to update workout session: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, workoutSessionCacheKey(id))
	s.cache.Del(ctx, "workout_sessions:list:*")

	return successResponse(c, workoutSessionToResponse(updatedWorkoutSession))
}

func (s *FiberServer) deleteWorkoutSession(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Workout session ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.DeleteWorkoutSession(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to delete workout session: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, workoutSessionCacheKey(id))
	s.cache.Del(ctx, "workout_sessions:list:*")

	return c.Status(fiber.StatusNoContent).Send(nil)
}

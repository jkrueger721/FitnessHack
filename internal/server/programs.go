package server

import (
	"time"

	"fitness-hack/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ProgramResponse represents the response structure for programs
type ProgramResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	UserID        string    `json:"userId"`
	DurationWeeks *int      `json:"durationWeeks,omitempty"`
	Difficulty    *string   `json:"difficulty,omitempty"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// CreateProgramRequest represents the request structure for creating programs
type CreateProgramRequest struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	DurationWeeks *int    `json:"durationWeeks,omitempty"`
	Difficulty    *string `json:"difficulty,omitempty"`
}

// UpdateProgramRequest represents the request structure for updating programs
type UpdateProgramRequest struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	DurationWeeks *int    `json:"durationWeeks,omitempty"`
	Difficulty    *string `json:"difficulty,omitempty"`
	IsActive      *bool   `json:"isActive,omitempty"`
}

// convertProgramToResponse converts a database Programs to ProgramResponse
func convertProgramToResponse(program *database.Programs) *ProgramResponse {
	// Handle type assertions safely
	var name string
	if program.Name != nil {
		if str, ok := program.Name.(string); ok {
			name = str
		}
	}

	var difficulty *string
	if program.Difficulty != nil {
		if str, ok := program.Difficulty.(string); ok {
			difficulty = &str
		}
	}

	return &ProgramResponse{
		ID:            program.Id,
		Name:          name,
		Description:   &program.Description,
		UserID:        program.User_id,
		DurationWeeks: &program.Duration_weeks,
		Difficulty:    difficulty,
		IsActive:      program.Is_active,
		CreatedAt:     program.Created_at,
		UpdatedAt:     program.Updated_at,
	}
}

// convertRequestToProgram converts a CreateProgramRequest to database Programs
func convertRequestToProgram(req *CreateProgramRequest, userID string) *database.Programs {
	now := time.Now()

	// Convert optional fields
	var description string
	if req.Description != nil {
		description = *req.Description
	}

	var durationWeeks int
	if req.DurationWeeks != nil {
		durationWeeks = *req.DurationWeeks
	}

	var difficulty interface{}
	if req.Difficulty != nil {
		difficulty = *req.Difficulty
	}

	return &database.Programs{
		Id:             uuid.New().String(),
		Name:           req.Name,
		Description:    description,
		User_id:        userID,
		Duration_weeks: durationWeeks,
		Difficulty:     difficulty,
		Is_active:      true,
		Created_at:     now,
		Updated_at:     now,
	}
}

// createProgram handles POST /api/programs
func (s *FiberServer) createProgram(c *fiber.Ctx) error {
	var req CreateProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// TODO: Get user ID from authentication context
	// For now, using a placeholder user ID
	userID := "placeholder-user-id"

	program := convertRequestToProgram(&req, userID)

	createdProgram, err := s.db.CreateProgram(c.Context(), program)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create program")
	}

	response := convertProgramToResponse(createdProgram)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// getProgram handles GET /api/programs/{id}
func (s *FiberServer) getProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	program, err := s.db.GetProgramByID(c.Context(), id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Program not found")
	}

	response := convertProgramToResponse(program)
	return c.JSON(response)
}

// listPrograms handles GET /api/programs
func (s *FiberServer) listPrograms(c *fiber.Ctx) error {
	limit, offset := getPaginationParams(c)

	programs, err := s.db.ListPrograms(c.Context(), limit, offset)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to list programs")
	}

	responses := make([]*ProgramResponse, len(programs))
	for i, program := range programs {
		responses[i] = convertProgramToResponse(&program)
	}

	return c.JSON(responses)
}

// updateProgram handles PUT /api/programs/{id}
func (s *FiberServer) updateProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing program
	existingProgram, err := s.db.GetProgramByID(c.Context(), id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "Program not found")
	}

	// Update fields if provided
	if req.Name != nil {
		existingProgram.Name = *req.Name
	}
	if req.Description != nil {
		existingProgram.Description = *req.Description
	}
	if req.DurationWeeks != nil {
		existingProgram.Duration_weeks = *req.DurationWeeks
	}
	if req.Difficulty != nil {
		existingProgram.Difficulty = *req.Difficulty
	}
	if req.IsActive != nil {
		existingProgram.Is_active = *req.IsActive
	}
	existingProgram.Updated_at = time.Now()

	updatedProgram, err := s.db.UpdateProgram(c.Context(), existingProgram)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to update program")
	}

	response := convertProgramToResponse(updatedProgram)
	return c.JSON(response)
}

// deleteProgram handles DELETE /api/programs/{id}
func (s *FiberServer) deleteProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	err := s.db.DeleteProgram(c.Context(), id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to delete program")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

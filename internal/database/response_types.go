package database

import "time"

// UserResponse represents the response structure for users
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateUserRequest represents the request structure for creating users
type CreateUserRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// UpdateUserRequest represents the request structure for updating users
type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty"`
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
}

// LoginRequest represents the request structure for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response structure for user login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// WorkoutResponse represents the response structure for workouts
type WorkoutResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"durationMinutes"`
	ProgramID       string    `json:"programId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CreateWorkoutRequest represents the request structure for creating workouts
type CreateWorkoutRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
	ProgramID       string `json:"programId"`
}

// UpdateWorkoutRequest represents the request structure for updating workouts
type UpdateWorkoutRequest struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	DurationMinutes *int    `json:"durationMinutes,omitempty"`
	ProgramID       *string `json:"programId,omitempty"`
}

// ExerciseResponse represents the response structure for exercises
type ExerciseResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	MuscleGroup     string    `json:"muscleGroup"`
	Equipment       string    `json:"equipment"`
	DifficultyLevel string    `json:"difficultyLevel"`
	Instructions    string    `json:"instructions"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CreateExerciseRequest represents the request structure for creating exercises
type CreateExerciseRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	MuscleGroup     string `json:"muscleGroup"`
	Equipment       string `json:"equipment"`
	DifficultyLevel string `json:"difficultyLevel"`
	Instructions    string `json:"instructions"`
}

// UpdateExerciseRequest represents the request structure for updating exercises
type UpdateExerciseRequest struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	MuscleGroup     *string `json:"muscleGroup,omitempty"`
	Equipment       *string `json:"equipment,omitempty"`
	DifficultyLevel *string `json:"difficultyLevel,omitempty"`
	Instructions    *string `json:"instructions,omitempty"`
}

// WorkoutExerciseResponse represents the response structure for workout exercises
type WorkoutExerciseResponse struct {
	ID              string    `json:"id"`
	WorkoutID       string    `json:"workoutId"`
	ExerciseID      string    `json:"exerciseId"`
	Sets            int       `json:"sets"`
	Reps            int       `json:"reps"`
	WeightKg        float64   `json:"weightKg"`
	DurationSeconds int       `json:"durationSeconds"`
	OrderIndex      int       `json:"orderIndex"`
	RestSeconds     int       `json:"restSeconds"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"createdAt"`
}

// CreateWorkoutExerciseRequest represents the request structure for creating workout exercises
type CreateWorkoutExerciseRequest struct {
	WorkoutID       string  `json:"workoutId"`
	ExerciseID      string  `json:"exerciseId"`
	Sets            int     `json:"sets"`
	Reps            int     `json:"reps"`
	WeightKg        float64 `json:"weightKg"`
	DurationSeconds int     `json:"durationSeconds"`
	OrderIndex      int     `json:"orderIndex"`
	RestSeconds     int     `json:"restSeconds"`
	Notes           string  `json:"notes"`
}

// UpdateWorkoutExerciseRequest represents the request structure for updating workout exercises
type UpdateWorkoutExerciseRequest struct {
	WorkoutID       *string  `json:"workoutId,omitempty"`
	ExerciseID      *string  `json:"exerciseId,omitempty"`
	Sets            *int     `json:"sets,omitempty"`
	Reps            *int     `json:"reps,omitempty"`
	WeightKg        *float64 `json:"weightKg,omitempty"`
	DurationSeconds *int     `json:"durationSeconds,omitempty"`
	OrderIndex      *int     `json:"orderIndex,omitempty"`
	RestSeconds     *int     `json:"restSeconds,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
}

// WorkoutSessionResponse represents the response structure for workout sessions
type WorkoutSessionResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"userId"`
	WorkoutID       string     `json:"workoutId"`
	Name            string     `json:"name"`
	StartedAt       time.Time  `json:"startedAt"`
	CompletedAt     *time.Time `json:"completedAt,omitempty"`
	DurationMinutes int        `json:"durationMinutes"`
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// CreateWorkoutSessionRequest represents the request structure for creating workout sessions
type CreateWorkoutSessionRequest struct {
	WorkoutID       string     `json:"workoutId"`
	Name            string     `json:"name"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	CompletedAt     *time.Time `json:"completedAt,omitempty"`
	DurationMinutes int        `json:"durationMinutes"`
	Notes           string     `json:"notes"`
}

// UpdateWorkoutSessionRequest represents the request structure for updating workout sessions
type UpdateWorkoutSessionRequest struct {
	WorkoutID       *string    `json:"workoutId,omitempty"`
	Name            *string    `json:"name,omitempty"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	CompletedAt     *time.Time `json:"completedAt,omitempty"`
	DurationMinutes *int       `json:"durationMinutes,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
}

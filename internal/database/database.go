package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetDB returns the underlying sqlx.DB instance for direct access
	GetDB() *sqlx.DB

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (*sqlx.Tx, error)

	// PingContext pings the database with context
	PingContext(ctx context.Context) error

	// Stats returns database statistics
	Stats() map[string]interface{}

	// --- USERS CRUD ---
	CreateUser(ctx context.Context, user *Users) (*Users, error)
	GetUserByID(ctx context.Context, id string) (*Users, error)
	GetUserByEmail(ctx context.Context, email string) (*Users, error)
	ListUsers(ctx context.Context, limit, offset int) ([]Users, error)
	UpdateUser(ctx context.Context, user *Users) (*Users, error)
	DeleteUser(ctx context.Context, id string) error

	// --- WORKOUTS CRUD ---
	CreateWorkout(ctx context.Context, workout *Workouts) (*Workouts, error)
	GetWorkoutByID(ctx context.Context, id string) (*Workouts, error)
	ListWorkouts(ctx context.Context, limit, offset int) ([]Workouts, error)
	UpdateWorkout(ctx context.Context, workout *Workouts) (*Workouts, error)
	DeleteWorkout(ctx context.Context, id string) error

	// --- EXERCISES CRUD ---
	CreateExercise(ctx context.Context, exercise *Exercises) (*Exercises, error)
	GetExerciseByID(ctx context.Context, id string) (*Exercises, error)
	ListExercises(ctx context.Context, limit, offset int) ([]Exercises, error)
	UpdateExercise(ctx context.Context, exercise *Exercises) (*Exercises, error)
	DeleteExercise(ctx context.Context, id string) error

	// --- WORKOUT_EXERCISES CRUD ---
	CreateWorkoutExercise(ctx context.Context, we *Workout_exercises) (*Workout_exercises, error)
	GetWorkoutExerciseByID(ctx context.Context, id string) (*Workout_exercises, error)
	ListWorkoutExercises(ctx context.Context, limit, offset int) ([]Workout_exercises, error)
	UpdateWorkoutExercise(ctx context.Context, we *Workout_exercises) (*Workout_exercises, error)
	DeleteWorkoutExercise(ctx context.Context, id string) error

	// --- WORKOUT_SESSIONS CRUD ---
	CreateWorkoutSession(ctx context.Context, ws *Workout_sessions) (*Workout_sessions, error)
	GetWorkoutSessionByID(ctx context.Context, id string) (*Workout_sessions, error)
	ListWorkoutSessions(ctx context.Context, limit, offset int) ([]Workout_sessions, error)
	UpdateWorkoutSession(ctx context.Context, ws *Workout_sessions) (*Workout_sessions, error)
	DeleteWorkoutSession(ctx context.Context, id string) error

	// --- PROGRAMS CRUD ---
	CreateProgram(ctx context.Context, program *Programs) (*Programs, error)
	GetProgramByID(ctx context.Context, id string) (*Programs, error)
	ListPrograms(ctx context.Context, limit, offset int) ([]Programs, error)
	UpdateProgram(ctx context.Context, program *Programs) (*Programs, error)
	DeleteProgram(ctx context.Context, id string) error
}

type service struct {
	db *sqlx.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

// Config holds database configuration
type Config struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns default database configuration
func DefaultConfig() *Config {
	return &Config{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// New creates a new database service instance with default configuration
func New() Service {
	return NewWithConfig(DefaultConfig())
}

// NewWithConfig creates a new database service instance with custom configuration
func NewWithConfig(config *Config) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		username, password, host, port, database, schema)

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	dbInstance = &service{
		db: db,
	}

	log.Printf("Successfully connected to database: %s", database)
	return dbInstance
}

// GetDB returns the underlying sqlx.DB instance for direct access
func (s *service) GetDB() *sqlx.DB {
	return s.db
}

// BeginTx starts a new transaction
func (s *service) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, nil)
}

// PingContext pings the database with context
func (s *service) PingContext(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Stats returns database statistics
func (s *service) Stats() map[string]interface{} {
	dbStats := s.db.Stats()
	return map[string]interface{}{
		"open_connections":    dbStats.OpenConnections,
		"in_use":              dbStats.InUse,
		"idle":                dbStats.Idle,
		"wait_count":          dbStats.WaitCount,
		"wait_duration":       dbStats.WaitDuration,
		"max_idle_closed":     dbStats.MaxIdleClosed,
		"max_lifetime_closed": dbStats.MaxLifetimeClosed,
	}
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("Database health check failed: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "Database is healthy"

	// Get database stats
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 20 {
		stats["message"] = "Database is experiencing high connection load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "Database has high wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections being closed, consider adjusting pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections closed due to max lifetime, consider increasing lifetime."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnecting from database: %s", database)
	return s.db.Close()
}

func (s *service) CreateUser(ctx context.Context, user *Users) (*Users, error) {
	query := `INSERT INTO users (email, username, password_hash, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *`

	// Handle type assertions for interface{} fields
	var email, username, passwordHash, firstName, lastName string

	if user.Email != nil {
		if str, ok := user.Email.(string); ok {
			email = str
		}
	}
	if user.Username != nil {
		if str, ok := user.Username.(string); ok {
			username = str
		}
	}
	if user.Password_hash != nil {
		if str, ok := user.Password_hash.(string); ok {
			passwordHash = str
		}
	}
	if user.First_name != nil {
		if str, ok := user.First_name.(string); ok {
			firstName = str
		}
	}
	if user.Last_name != nil {
		if str, ok := user.Last_name.(string); ok {
			lastName = str
		}
	}

	// Log the values being inserted for debugging
	fmt.Printf("DEBUG: Inserting user with values: email=%s, username=%s, passwordHash=%s, firstName=%s, lastName=%s\n",
		email, username, passwordHash, firstName, lastName)

	row := s.db.QueryRowContext(ctx, query, email, username, passwordHash, firstName, lastName, user.Created_at, user.Updated_at)

	var created Users
	err := row.Scan(&created.Id, &created.Email, &created.Username, &created.Password_hash, &created.First_name, &created.Last_name, &created.Created_at, &created.Updated_at)
	if err != nil {
		fmt.Printf("DEBUG: Error scanning result: %v\n", err)
		return nil, fmt.Errorf("failed to scan user result: %w", err)
	}

	return &created, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*Users, error) {
	var user Users
	query := `SELECT * FROM users WHERE id = $1`
	err := s.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*Users, error) {
	var user Users
	query := `SELECT * FROM users WHERE email = $1`
	err := s.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]Users, error) {
	var users []Users
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}

func (s *service) UpdateUser(ctx context.Context, user *Users) (*Users, error) {
	query := `UPDATE users SET email=:email, username=:username, password_hash=:password_hash, first_name=:first_name, last_name=:last_name, updated_at=:updated_at WHERE id=:id RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var updated Users
		if err := row.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, fmt.Errorf("failed to update user")
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// --- WORKOUTS CRUD ---
func (s *service) CreateWorkout(ctx context.Context, workout *Workouts) (*Workouts, error) {
	query := `INSERT INTO workouts (id, user_id, name, description, duration_minutes, program_id, created_at, updated_at)
		VALUES (:id, :user_id, :name, :description, :duration_minutes, :program_id, :created_at, :updated_at)
		RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, workout)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var created Workouts
		if err := row.StructScan(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}
	return nil, fmt.Errorf("failed to insert workout")
}

func (s *service) GetWorkoutByID(ctx context.Context, id string) (*Workouts, error) {
	var workout Workouts
	query := `SELECT * FROM workouts WHERE id = $1`
	err := s.db.GetContext(ctx, &workout, query, id)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (s *service) ListWorkouts(ctx context.Context, limit, offset int) ([]Workouts, error) {
	var workouts []Workouts
	query := `SELECT * FROM workouts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &workouts, query, limit, offset)
	return workouts, err
}

func (s *service) UpdateWorkout(ctx context.Context, workout *Workouts) (*Workouts, error) {
	query := `UPDATE workouts SET user_id=:user_id, name=:name, description=:description, duration_minutes=:duration_minutes, program_id=:program_id, updated_at=:updated_at WHERE id=:id RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, workout)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var updated Workouts
		if err := row.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, fmt.Errorf("failed to update workout")
}

func (s *service) DeleteWorkout(ctx context.Context, id string) error {
	query := `DELETE FROM workouts WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// --- EXERCISES CRUD ---
func (s *service) CreateExercise(ctx context.Context, exercise *Exercises) (*Exercises, error) {
	query := `INSERT INTO exercises (id, name, description, muscle_group, equipment, difficulty_level, instructions, created_at, updated_at)
		VALUES (:id, :name, :description, :muscle_group, :equipment, :difficulty_level, :instructions, :created_at, :updated_at)
		RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, exercise)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var created Exercises
		if err := row.StructScan(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}
	return nil, fmt.Errorf("failed to insert exercise")
}

func (s *service) GetExerciseByID(ctx context.Context, id string) (*Exercises, error) {
	var exercise Exercises
	query := `SELECT * FROM exercises WHERE id = $1`
	err := s.db.GetContext(ctx, &exercise, query, id)
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (s *service) ListExercises(ctx context.Context, limit, offset int) ([]Exercises, error) {
	var exercises []Exercises
	query := `SELECT * FROM exercises ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &exercises, query, limit, offset)
	return exercises, err
}

func (s *service) UpdateExercise(ctx context.Context, exercise *Exercises) (*Exercises, error) {
	query := `UPDATE exercises SET name=:name, description=:description, muscle_group=:muscle_group, equipment=:equipment, difficulty_level=:difficulty_level, instructions=:instructions, updated_at=:updated_at WHERE id=:id RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, exercise)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var updated Exercises
		if err := row.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, fmt.Errorf("failed to update exercise")
}

func (s *service) DeleteExercise(ctx context.Context, id string) error {
	query := `DELETE FROM exercises WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// --- WORKOUT_EXERCISES CRUD ---
func (s *service) CreateWorkoutExercise(ctx context.Context, we *Workout_exercises) (*Workout_exercises, error) {
	query := `INSERT INTO workout_exercises (id, workout_id, exercise_id, sets, reps, weight_kg, duration_seconds, order_index, rest_seconds, notes, created_at)
		VALUES (:id, :workout_id, :exercise_id, :sets, :reps, :weight_kg, :duration_seconds, :order_index, :rest_seconds, :notes, :created_at)
		RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, we)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var created Workout_exercises
		if err := row.StructScan(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}
	return nil, fmt.Errorf("failed to insert workout_exercise")
}

func (s *service) GetWorkoutExerciseByID(ctx context.Context, id string) (*Workout_exercises, error) {
	var we Workout_exercises
	query := `SELECT * FROM workout_exercises WHERE id = $1`
	err := s.db.GetContext(ctx, &we, query, id)
	if err != nil {
		return nil, err
	}
	return &we, nil
}

func (s *service) ListWorkoutExercises(ctx context.Context, limit, offset int) ([]Workout_exercises, error) {
	var wes []Workout_exercises
	query := `SELECT * FROM workout_exercises ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &wes, query, limit, offset)
	return wes, err
}

func (s *service) UpdateWorkoutExercise(ctx context.Context, we *Workout_exercises) (*Workout_exercises, error) {
	query := `UPDATE workout_exercises SET workout_id=:workout_id, exercise_id=:exercise_id, sets=:sets, reps=:reps, weight_kg=:weight_kg, duration_seconds=:duration_seconds, order_index=:order_index, rest_seconds=:rest_seconds, notes=:notes WHERE id=:id RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, we)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var updated Workout_exercises
		if err := row.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, fmt.Errorf("failed to update workout_exercise")
}

func (s *service) DeleteWorkoutExercise(ctx context.Context, id string) error {
	query := `DELETE FROM workout_exercises WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// --- WORKOUT_SESSIONS CRUD ---
func (s *service) CreateWorkoutSession(ctx context.Context, ws *Workout_sessions) (*Workout_sessions, error) {
	query := `INSERT INTO workout_sessions (id, user_id, workout_id, name, started_at, completed_at, duration_minutes, notes, created_at, updated_at)
		VALUES (:id, :user_id, :workout_id, :name, :started_at, :completed_at, :duration_minutes, :notes, :created_at, :updated_at)
		RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, ws)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var created Workout_sessions
		if err := row.StructScan(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}
	return nil, fmt.Errorf("failed to insert workout_session")
}

func (s *service) GetWorkoutSessionByID(ctx context.Context, id string) (*Workout_sessions, error) {
	var ws Workout_sessions
	query := `SELECT * FROM workout_sessions WHERE id = $1`
	err := s.db.GetContext(ctx, &ws, query, id)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

func (s *service) ListWorkoutSessions(ctx context.Context, limit, offset int) ([]Workout_sessions, error) {
	var wss []Workout_sessions
	query := `SELECT * FROM workout_sessions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &wss, query, limit, offset)
	return wss, err
}

func (s *service) UpdateWorkoutSession(ctx context.Context, ws *Workout_sessions) (*Workout_sessions, error) {
	query := `UPDATE workout_sessions SET user_id=:user_id, workout_id=:workout_id, name=:name, started_at=:started_at, completed_at=:completed_at, duration_minutes=:duration_minutes, notes=:notes, updated_at=:updated_at WHERE id=:id RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, ws)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var updated Workout_sessions
		if err := row.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, fmt.Errorf("failed to update workout_session")
}

func (s *service) DeleteWorkoutSession(ctx context.Context, id string) error {
	query := `DELETE FROM workout_sessions WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// --- PROGRAMS CRUD ---
func (s *service) CreateProgram(ctx context.Context, program *Programs) (*Programs, error) {
	query := `INSERT INTO programs (id, name, description, user_id, duration_weeks, difficulty, is_active, created_at, updated_at)
		VALUES (:id, :name, :description, :user_id, :duration_weeks, :difficulty, :is_active, :created_at, :updated_at)
		RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, program)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var created Programs
		if err := row.StructScan(&created); err != nil {
			return nil, err
		}
		return &created, nil
	}
	return nil, fmt.Errorf("failed to insert program")
}

func (s *service) GetProgramByID(ctx context.Context, id string) (*Programs, error) {
	var program Programs
	query := `SELECT * FROM programs WHERE id = $1`
	err := s.db.GetContext(ctx, &program, query, id)
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (s *service) ListPrograms(ctx context.Context, limit, offset int) ([]Programs, error) {
	var programs []Programs
	query := `SELECT * FROM programs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := s.db.SelectContext(ctx, &programs, query, limit, offset)
	return programs, err
}

func (s *service) UpdateProgram(ctx context.Context, program *Programs) (*Programs, error) {
	query := `UPDATE programs SET name=:name, description=:description, user_id=:user_id, duration_weeks=:duration_weeks, difficulty=:difficulty, is_active=:is_active, updated_at=:updated_at WHERE id=:id RETURNING *`
	row, err := s.db.NamedQueryContext(ctx, query, program)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		var updated Programs
		if err := row.StructScan(&updated); err != nil {
			return nil, err
		}
		return &updated, nil
	}
	return nil, fmt.Errorf("failed to update program")
}

func (s *service) DeleteProgram(ctx context.Context, id string) error {
	query := `DELETE FROM programs WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

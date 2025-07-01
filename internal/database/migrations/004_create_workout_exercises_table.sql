-- Migration: 004_create_workout_exercises_table
-- Description: Creates the junction table for workout-exercise relationships
-- Date: 2024-01-01

CREATE TABLE IF NOT EXISTS workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    sets INTEGER DEFAULT 1,
    reps INTEGER,
    weight_kg DECIMAL(5,2),
    duration_seconds INTEGER,
    order_index INTEGER DEFAULT 0,
    rest_seconds INTEGER DEFAULT 60,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(workout_id, exercise_id, order_index)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_workout_exercises_workout_id ON workout_exercises(workout_id);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_exercise_id ON workout_exercises(exercise_id);
CREATE INDEX IF NOT EXISTS idx_workout_exercises_order ON workout_exercises(workout_id, order_index);

-- Add comments for documentation
COMMENT ON TABLE workout_exercises IS 'Junction table linking workouts to exercises with specific parameters';
COMMENT ON COLUMN workout_exercises.id IS 'Unique identifier for the workout-exercise relationship';
COMMENT ON COLUMN workout_exercises.workout_id IS 'Reference to the workout';
COMMENT ON COLUMN workout_exercises.exercise_id IS 'Reference to the exercise';
COMMENT ON COLUMN workout_exercises.sets IS 'Number of sets for this exercise in the workout';
COMMENT ON COLUMN workout_exercises.reps IS 'Number of repetitions per set';
COMMENT ON COLUMN workout_exercises.weight_kg IS 'Weight to use for this exercise (in kg)';
COMMENT ON COLUMN workout_exercises.duration_seconds IS 'Duration for time-based exercises (in seconds)';
COMMENT ON COLUMN workout_exercises.order_index IS 'Order of this exercise within the workout';
COMMENT ON COLUMN workout_exercises.rest_seconds IS 'Rest time after this exercise (in seconds)';
COMMENT ON COLUMN workout_exercises.notes IS 'Additional notes for this exercise in the workout'; 
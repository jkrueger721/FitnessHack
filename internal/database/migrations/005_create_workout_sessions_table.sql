-- Migration: 005_create_workout_sessions_table
-- Description: Creates the workout sessions table for tracking completed workouts
-- Date: 2024-01-01

CREATE TABLE IF NOT EXISTS workout_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id UUID REFERENCES workouts(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_workout_sessions_user_id ON workout_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_workout_id ON workout_sessions(workout_id);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_started_at ON workout_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_workout_sessions_completed_at ON workout_sessions(completed_at);

-- Add comments for documentation
COMMENT ON TABLE workout_sessions IS 'Tracks completed workout sessions by users';
COMMENT ON COLUMN workout_sessions.id IS 'Unique identifier for the workout session';
COMMENT ON COLUMN workout_sessions.user_id IS 'Reference to the user who performed the workout';
COMMENT ON COLUMN workout_sessions.workout_id IS 'Reference to the workout plan (optional for custom workouts)';
COMMENT ON COLUMN workout_sessions.name IS 'Name of the workout session';
COMMENT ON COLUMN workout_sessions.started_at IS 'When the workout session started';
COMMENT ON COLUMN workout_sessions.completed_at IS 'When the workout session was completed';
COMMENT ON COLUMN workout_sessions.duration_minutes IS 'Actual duration of the workout session';
COMMENT ON COLUMN workout_sessions.notes IS 'User notes about the workout session'; 
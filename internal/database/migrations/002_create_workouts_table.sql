-- Migration: 002_create_workouts_table
-- Description: Creates the workouts table linked to users
-- Date: 2024-01-01

CREATE TABLE IF NOT EXISTS workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_created_at ON workouts(created_at);

-- Add comments for documentation
COMMENT ON TABLE workouts IS 'Stores workout plans created by users';
COMMENT ON COLUMN workouts.id IS 'Unique identifier for the workout';
COMMENT ON COLUMN workouts.user_id IS 'Reference to the user who created this workout';
COMMENT ON COLUMN workouts.name IS 'Name of the workout';
COMMENT ON COLUMN workouts.description IS 'Detailed description of the workout';
COMMENT ON COLUMN workouts.duration_minutes IS 'Estimated duration of the workout in minutes'; 
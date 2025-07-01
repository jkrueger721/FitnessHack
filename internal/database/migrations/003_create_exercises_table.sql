-- Migration: 003_create_exercises_table
-- Description: Creates the exercises table with exercise metadata
-- Date: 2024-01-01

CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    muscle_group VARCHAR(100),
    equipment VARCHAR(100),
    difficulty_level VARCHAR(50),
    instructions TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_exercises_name ON exercises(name);
CREATE INDEX IF NOT EXISTS idx_exercises_muscle_group ON exercises(muscle_group);
CREATE INDEX IF NOT EXISTS idx_exercises_equipment ON exercises(equipment);

-- Add comments for documentation
COMMENT ON TABLE exercises IS 'Stores exercise definitions and metadata';
COMMENT ON COLUMN exercises.id IS 'Unique identifier for the exercise';
COMMENT ON COLUMN exercises.name IS 'Name of the exercise';
COMMENT ON COLUMN exercises.description IS 'Detailed description of the exercise';
COMMENT ON COLUMN exercises.muscle_group IS 'Primary muscle group targeted by this exercise';
COMMENT ON COLUMN exercises.equipment IS 'Equipment required for this exercise';
COMMENT ON COLUMN exercises.difficulty_level IS 'Difficulty level (beginner, intermediate, advanced)';
COMMENT ON COLUMN exercises.instructions IS 'Step-by-step instructions for performing the exercise'; 
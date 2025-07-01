-- Migration: 006_add_programs_table.sql
-- Description: add programs table
-- Date: 2025-06-30

-- Create programs table
CREATE TABLE IF NOT EXISTS programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    duration_weeks INTEGER,
    difficulty VARCHAR(50) CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_programs_user_id ON programs(user_id);
CREATE INDEX IF NOT EXISTS idx_programs_is_active ON programs(is_active);

-- Add program_id column to workouts table
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS program_id UUID REFERENCES programs(id) ON DELETE SET NULL;

-- Add index for the new foreign key
CREATE INDEX IF NOT EXISTS idx_workouts_program_id ON workouts(program_id);

-- Add comments for documentation
COMMENT ON TABLE programs IS 'Fitness programs that contain multiple workouts';
COMMENT ON COLUMN programs.name IS 'Name of the fitness program';
COMMENT ON COLUMN programs.description IS 'Description of the program goals and structure';
COMMENT ON COLUMN programs.duration_weeks IS 'Expected duration of the program in weeks';
COMMENT ON COLUMN programs.difficulty IS 'Difficulty level of the program';
COMMENT ON COLUMN programs.is_active IS 'Whether the program is currently active';

COMMENT ON COLUMN workouts.program_id IS 'Reference to the program this workout belongs to';

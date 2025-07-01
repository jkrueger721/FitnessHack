-- Migration: 001_create_users_table
-- Description: Creates the users table with authentication fields
-- Date: 2024-01-01

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Add comments for documentation
COMMENT ON TABLE users IS 'Stores user account information and authentication data';
COMMENT ON COLUMN users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN users.email IS 'User email address (unique)';
COMMENT ON COLUMN users.username IS 'User username (unique)';
COMMENT ON COLUMN users.password_hash IS 'Hashed password for authentication'; 
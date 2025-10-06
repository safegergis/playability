-- Migration: Add Verifications Table and User Verification
-- Created: 2024-10-05
-- Description: Adds verifications table for email verification and password reset, and adds verified column to users table

-- Add verified column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS verified BOOLEAN NOT NULL DEFAULT FALSE;

-- Create verifications table
CREATE TABLE IF NOT EXISTS verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    type INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_verifications_email_type ON verifications(email, type);
CREATE INDEX IF NOT EXISTS idx_verifications_expires_at ON verifications(expires_at);
CREATE INDEX IF NOT EXISTS idx_users_verified ON users(verified);

-- Add comment to describe the type column
COMMENT ON COLUMN verifications.type IS '1 = Email Confirmation, 2 = Password Reset';

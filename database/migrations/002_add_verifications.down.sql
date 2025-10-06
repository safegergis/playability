-- Migration Rollback: Remove Verifications Table and User Verification
-- Created: 2024-10-05
-- Description: Removes verifications table and verified column from users table

-- Drop indexes
DROP INDEX IF EXISTS idx_users_verified;
DROP INDEX IF EXISTS idx_verifications_expires_at;
DROP INDEX IF EXISTS idx_verifications_email_type;

-- Drop verifications table
DROP TABLE IF EXISTS verifications;

-- Remove verified column from users table
ALTER TABLE users DROP COLUMN IF EXISTS verified;

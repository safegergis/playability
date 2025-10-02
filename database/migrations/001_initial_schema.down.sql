-- Migration Rollback: Initial Schema
-- Description: Drops all tables and types created in the initial schema migration

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_games_updated_at ON games;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_games_name;
DROP INDEX IF EXISTS idx_reports_game_id_created_at;
DROP INDEX IF EXISTS idx_reports_user_id;
DROP INDEX IF EXISTS idx_reports_game_user;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS feature_support;

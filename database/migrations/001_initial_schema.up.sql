-- Migration: Initial Schema
-- Created: 2024-10-01
-- Description: Creates initial database schema with users, games, and reports tables

-- Create custom enum type for feature support levels
CREATE TYPE feature_support AS ENUM ('false', 'unknown', 'limited', 'true');

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hash VARCHAR(255) NOT NULL,
    num_reports INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create games table
CREATE TABLE IF NOT EXISTS games (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    summary TEXT NOT NULL,
    cover_art VARCHAR(255) NOT NULL,
    platforms JSONB NOT NULL,
    closed_captions feature_support NOT NULL,
    color_blind feature_support NOT NULL,
    full_controller_support feature_support NOT NULL,
    controller_remapping feature_support NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create reports table
CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    game_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    closed_captions feature_support NOT NULL,
    color_blind feature_support NOT NULL,
    full_controller_support feature_support NOT NULL,
    controller_remapping feature_support NOT NULL,
    score VARCHAR(255) NOT NULL,
    report TEXT
);

-- Add unique constraints
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);

-- Add foreign key constraints
ALTER TABLE reports ADD CONSTRAINT fk_reports_game
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;

ALTER TABLE reports ADD CONSTRAINT fk_reports_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add check constraints
ALTER TABLE users ADD CONSTRAINT check_num_reports_positive
    CHECK (num_reports >= 0);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_games_name ON games(name);
CREATE INDEX idx_reports_game_id_created_at ON reports(game_id, created_at DESC);
CREATE INDEX idx_reports_user_id ON reports(user_id);
CREATE INDEX idx_reports_game_user ON reports(game_id, user_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_games_updated_at BEFORE UPDATE ON games
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

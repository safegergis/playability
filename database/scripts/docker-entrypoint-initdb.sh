#!/bin/bash
set -e

echo "======================================"
echo "Docker PostgreSQL Initialization"
echo "======================================"

echo "Running database initialization scripts..."
echo ""

# PostgreSQL runs initialization scripts from /docker-entrypoint-initdb.d/
# This script is executed automatically when the container starts for the first time

echo "Step 1: Creating custom types..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TYPE IF NOT EXISTS feature_support AS ENUM ('false', 'unknown', 'limited', 'true');
EOSQL
echo "✓ Custom types created"

echo ""
echo "Step 2: Creating tables..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
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
EOSQL
echo "✓ Tables created"

echo ""
echo "Step 3: Adding constraints..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Add unique constraints
    DO \$\$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'unique_email'
        ) THEN
            ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'unique_username'
        ) THEN
            ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);
        END IF;
    END \$\$;

    -- Add foreign key constraints
    DO \$\$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'fk_reports_game'
        ) THEN
            ALTER TABLE reports ADD CONSTRAINT fk_reports_game
            FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'fk_reports_user'
        ) THEN
            ALTER TABLE reports ADD CONSTRAINT fk_reports_user
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        END IF;
    END \$\$;

    -- Add check constraints
    DO \$\$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint WHERE conname = 'check_num_reports_positive'
        ) THEN
            ALTER TABLE users ADD CONSTRAINT check_num_reports_positive
            CHECK (num_reports >= 0);
        END IF;
    END \$\$;
EOSQL
echo "✓ Constraints added"

echo ""
echo "Step 4: Creating indexes..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
    CREATE INDEX IF NOT EXISTS idx_games_name ON games(name);
    CREATE INDEX IF NOT EXISTS idx_reports_game_id_created_at ON reports(game_id, created_at DESC);
    CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);
    CREATE INDEX IF NOT EXISTS idx_reports_game_user ON reports(game_id, user_id);
EOSQL
echo "✓ Indexes created"

echo ""
echo "Step 5: Creating triggers..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create function to update updated_at timestamp
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS \$func\$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    \$func\$ language 'plpgsql';

    -- Create triggers
    DROP TRIGGER IF EXISTS update_users_updated_at ON users;
    CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

    DROP TRIGGER IF EXISTS update_games_updated_at ON games;
    CREATE TRIGGER update_games_updated_at BEFORE UPDATE ON games
        FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
EOSQL
echo "✓ Triggers created"

echo ""
echo "Step 6: Creating migration tracking table..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS schema_migrations (
        id SERIAL PRIMARY KEY,
        version VARCHAR(255) NOT NULL UNIQUE,
        applied_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    -- Mark initial migration as applied
    INSERT INTO schema_migrations (version)
    VALUES ('001_initial_schema')
    ON CONFLICT (version) DO NOTHING;
EOSQL
echo "✓ Migration tracking ready"

echo ""
echo "======================================"
echo "Database initialization complete!"
echo "======================================"
echo ""
echo "Tables created:"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "\dt"

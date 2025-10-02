#!/bin/bash
set -e

echo "======================================"
echo "Playability Database Initialization"
echo "======================================"

if [ -f backend/.env ]; then
    export $(cat backend/.env | grep -v '^#' | xargs)
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-playability_user}"
DB_NAME="${DB_NAME:-playability}"

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo ""

check_postgres() {
    echo "Checking PostgreSQL connection..."
    if ! PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c '\q' 2>/dev/null; then
        echo "ERROR: Cannot connect to PostgreSQL"
        echo "Please ensure PostgreSQL is running and credentials are correct"
        exit 1
    fi
    echo "✓ PostgreSQL connection successful"
}

create_database() {
    echo ""
    echo "Checking if database exists..."

    if PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
        echo "⚠ Database '$DB_NAME' already exists"
        read -p "Do you want to drop and recreate it? (yes/no): " -r
        if [[ $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
            echo "Dropping database '$DB_NAME'..."
            PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
            echo "Creating database '$DB_NAME'..."
            PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME;"
            echo "✓ Database recreated"
        else
            echo "Using existing database"
        fi
    else
        echo "Creating database '$DB_NAME'..."
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME;"
        echo "✓ Database created"
    fi
}

run_migrations() {
    echo ""
    echo "Running database migrations..."

    SQL_DIR="$(dirname "$0")/.."

    SQL_FILES=(
        "feature_support.sql"
        "users.sql"
        "games.sql"
        "reports.sql"
    )

    for sql_file in "${SQL_FILES[@]}"; do
        file_path="$SQL_DIR/$sql_file"
        if [ -f "$file_path" ]; then
            echo "Executing: $sql_file"
            PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file_path"
            echo "✓ $sql_file executed"
        else
            echo "⚠ Warning: $file_path not found, skipping..."
        fi
    done
}

add_constraints() {
    echo ""
    echo "Adding constraints and indexes..."

    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << 'SQL'
-- Add unique constraints
DO $$
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
END $$;

-- Add foreign key constraints
DO $$
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
END $$;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_games_name ON games(name);
CREATE INDEX IF NOT EXISTS idx_reports_game_id_created_at ON reports(game_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);
CREATE INDEX IF NOT EXISTS idx_reports_game_user ON reports(game_id, user_id);

SQL

    echo "✓ Constraints and indexes added"
}

verify_schema() {
    echo ""
    echo "Verifying schema..."

    TABLE_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';")

    echo "Tables created: $TABLE_COUNT"

    if [ "$TABLE_COUNT" -ge 3 ]; then
        echo "✓ Schema verification passed"
    else
        echo "⚠ Warning: Expected at least 3 tables, found $TABLE_COUNT"
    fi

    echo ""
    echo "Database tables:"
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\dt"
}

main() {
    check_postgres
    create_database
    run_migrations
    add_constraints
    verify_schema

    echo ""
    echo "======================================"
    echo "Database initialization complete!"
    echo "======================================"
}

main "$@"

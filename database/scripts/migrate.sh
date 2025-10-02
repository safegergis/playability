#!/bin/bash
set -e

echo "======================================"
echo "Playability Database Migration Tool"
echo "======================================"

if [ -f backend/.env ]; then
    export $(cat backend/.env | grep -v '^#' | xargs)
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-playability_user}"
DB_NAME="${DB_NAME:-playability}"

MIGRATIONS_DIR="$(dirname "$0")/../migrations"
MIGRATION_TABLE="schema_migrations"

echo "Database: $DB_NAME"
echo "Migrations directory: $MIGRATIONS_DIR"
echo ""

check_postgres() {
    if ! PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; then
        echo "ERROR: Cannot connect to database '$DB_NAME'"
        echo "Please ensure the database exists and credentials are correct"
        exit 1
    fi
}

create_migration_table() {
    echo "Checking migration table..."

    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << 'SQL'
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    version VARCHAR(255) NOT NULL UNIQUE,
    applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);
SQL

    echo "✓ Migration table ready"
}

get_current_version() {
    VERSION=$(PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c \
        "SELECT version FROM $MIGRATION_TABLE ORDER BY version DESC LIMIT 1;" 2>/dev/null | xargs)

    if [ -z "$VERSION" ]; then
        echo "0"
    else
        echo "$VERSION"
    fi
}

get_pending_migrations() {
    CURRENT_VERSION=$(get_current_version)

    MIGRATIONS=()
    for file in "$MIGRATIONS_DIR"/*.up.sql; do
        if [ -f "$file" ]; then
            filename=$(basename "$file")
            version="${filename%%.up.sql}"

            if [[ "$version" > "$CURRENT_VERSION" ]] || [[ "$CURRENT_VERSION" == "0" ]]; then
                MIGRATIONS+=("$version")
            fi
        fi
    done

    printf '%s\n' "${MIGRATIONS[@]}" | sort
}

migrate_up() {
    echo ""
    echo "Running migrations..."

    CURRENT_VERSION=$(get_current_version)
    echo "Current version: $CURRENT_VERSION"

    PENDING=$(get_pending_migrations)

    if [ -z "$PENDING" ]; then
        echo "✓ Database is up to date"
        return
    fi

    echo ""
    echo "Pending migrations:"
    echo "$PENDING"
    echo ""

    for version in $PENDING; do
        MIGRATION_FILE="$MIGRATIONS_DIR/${version}.up.sql"

        if [ ! -f "$MIGRATION_FILE" ]; then
            echo "✗ Migration file not found: $MIGRATION_FILE"
            exit 1
        fi

        echo "Applying migration: $version"

        PGPASSWORD=$DB_PASSWORD psql \
            -h "$DB_HOST" \
            -p "$DB_PORT" \
            -U "$DB_USER" \
            -d "$DB_NAME" \
            -f "$MIGRATION_FILE"

        PGPASSWORD=$DB_PASSWORD psql \
            -h "$DB_HOST" \
            -p "$DB_PORT" \
            -U "$DB_USER" \
            -d "$DB_NAME" \
            -c "INSERT INTO $MIGRATION_TABLE (version) VALUES ('$version');"

        echo "✓ Applied: $version"
    done

    echo ""
    echo "✓ All migrations applied successfully"
}

migrate_down() {
    echo ""
    echo "Rolling back last migration..."

    CURRENT_VERSION=$(get_current_version)

    if [ "$CURRENT_VERSION" == "0" ]; then
        echo "✓ No migrations to roll back"
        return
    fi

    echo "Current version: $CURRENT_VERSION"

    MIGRATION_FILE="$MIGRATIONS_DIR/${CURRENT_VERSION}.down.sql"

    if [ ! -f "$MIGRATION_FILE" ]; then
        echo "✗ Rollback file not found: $MIGRATION_FILE"
        exit 1
    fi

    echo "Rolling back: $CURRENT_VERSION"

    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -f "$MIGRATION_FILE"

    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -c "DELETE FROM $MIGRATION_TABLE WHERE version = '$CURRENT_VERSION';"

    echo "✓ Rolled back: $CURRENT_VERSION"
}

migration_status() {
    echo ""
    echo "Migration Status"
    echo "----------------------------------------"

    CURRENT_VERSION=$(get_current_version)
    echo "Current version: $CURRENT_VERSION"

    echo ""
    echo "Applied migrations:"
    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -c "SELECT version, applied_at FROM $MIGRATION_TABLE ORDER BY version;"

    echo ""
    PENDING=$(get_pending_migrations)
    if [ -z "$PENDING" ]; then
        echo "Pending migrations: None"
    else
        echo "Pending migrations:"
        echo "$PENDING"
    fi
}

create_migration() {
    if [ -z "$1" ]; then
        echo "Usage: $0 create <migration_name>"
        exit 1
    fi

    TIMESTAMP=$(date +%Y%m%d%H%M%S)
    NAME=$(echo "$1" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')
    VERSION="${TIMESTAMP}"

    UP_FILE="$MIGRATIONS_DIR/${VERSION}_${NAME}.up.sql"
    DOWN_FILE="$MIGRATIONS_DIR/${VERSION}_${NAME}.down.sql"

    cat > "$UP_FILE" << EOF
-- Migration: $NAME
-- Created: $(date +%Y-%m-%d)
-- Description: TODO

-- Add your migration SQL here

EOF

    cat > "$DOWN_FILE" << EOF
-- Migration Rollback: $NAME
-- Description: TODO

-- Add your rollback SQL here

EOF

    echo "Created migration files:"
    echo "  $UP_FILE"
    echo "  $DOWN_FILE"
}

main() {
    case "${1:-up}" in
        up)
            check_postgres
            create_migration_table
            migrate_up
            ;;
        down)
            check_postgres
            create_migration_table
            migrate_down
            ;;
        status)
            check_postgres
            create_migration_table
            migration_status
            ;;
        create)
            create_migration "$2"
            ;;
        *)
            echo "Usage: $0 {up|down|status|create}"
            echo ""
            echo "  up      - Apply pending migrations"
            echo "  down    - Rollback last migration"
            echo "  status  - Show migration status"
            echo "  create  - Create new migration files"
            exit 1
            ;;
    esac

    echo ""
    echo "======================================"
}

main "$@"

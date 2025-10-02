#!/bin/bash
set -e

echo "======================================"
echo "Playability Database Restore"
echo "======================================"

if [ -f backend/.env ]; then
    export $(cat backend/.env | grep -v '^#' | xargs)
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-playability_user}"
DB_NAME="${DB_NAME:-playability}"

BACKUP_DIR="${BACKUP_DIR:-./backups}"

if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file>"
    echo ""
    echo "Available backups:"
    ls -1 "$BACKUP_DIR" | grep playability_ || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
    BACKUP_FILE="$BACKUP_DIR/$1"
fi

if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "Database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"
echo ""

decompress_if_needed() {
    if [[ "$BACKUP_FILE" == *.gz ]]; then
        echo "Decompressing backup file..."
        DECOMPRESSED_FILE="${BACKUP_FILE%.gz}"
        gunzip -c "$BACKUP_FILE" > "$DECOMPRESSED_FILE"
        BACKUP_FILE="$DECOMPRESSED_FILE"
        echo "✓ Backup decompressed"
    fi
}

confirm_restore() {
    echo "⚠️  WARNING ⚠️"
    echo "This will DROP and recreate the database '$DB_NAME'"
    echo "All existing data will be lost!"
    echo ""
    read -p "Are you sure you want to continue? (yes/no): " -r

    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        echo "Restore cancelled"
        exit 0
    fi
}

drop_database() {
    echo ""
    echo "Dropping existing database..."

    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d postgres \
        -c "DROP DATABASE IF EXISTS $DB_NAME;"

    echo "✓ Database dropped"
}

create_database() {
    echo "Creating database..."

    PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d postgres \
        -c "CREATE DATABASE $DB_NAME;"

    echo "✓ Database created"
}

restore_backup() {
    echo ""
    echo "Restoring backup..."

    if [[ "$BACKUP_FILE" == *.dump ]]; then
        PGPASSWORD=$DB_PASSWORD pg_restore \
            -h "$DB_HOST" \
            -p "$DB_PORT" \
            -U "$DB_USER" \
            -d "$DB_NAME" \
            --no-owner \
            --no-acl \
            "$BACKUP_FILE"
    else
        PGPASSWORD=$DB_PASSWORD psql \
            -h "$DB_HOST" \
            -p "$DB_PORT" \
            -U "$DB_USER" \
            -d "$DB_NAME" \
            -f "$BACKUP_FILE"
    fi

    echo "✓ Backup restored"
}

verify_restore() {
    echo ""
    echo "Verifying restore..."

    TABLE_COUNT=$(PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE';")

    echo "Tables restored: $TABLE_COUNT"

    if [ "$TABLE_COUNT" -ge 3 ]; then
        echo "✓ Restore verification passed"
    else
        echo "⚠ Warning: Expected at least 3 tables, found $TABLE_COUNT"
    fi

    echo ""
    USERS_COUNT=$(PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -c "SELECT COUNT(*) FROM users;" 2>/dev/null || echo "0")

    GAMES_COUNT=$(PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -c "SELECT COUNT(*) FROM games;" 2>/dev/null || echo "0")

    REPORTS_COUNT=$(PGPASSWORD=$DB_PASSWORD psql \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -t \
        -c "SELECT COUNT(*) FROM reports;" 2>/dev/null || echo "0")

    echo "Data restored:"
    echo "  Users: $USERS_COUNT"
    echo "  Games: $GAMES_COUNT"
    echo "  Reports: $REPORTS_COUNT"
}

cleanup_temp_files() {
    if [[ "$1" == *".decompressed" ]]; then
        rm -f "$1"
    fi
}

main() {
    decompress_if_needed
    confirm_restore
    drop_database
    create_database
    restore_backup
    verify_restore
    cleanup_temp_files "$BACKUP_FILE"

    echo ""
    echo "======================================"
    echo "Restore complete!"
    echo "======================================"
}

main "$@"

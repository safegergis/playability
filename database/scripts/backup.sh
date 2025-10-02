#!/bin/bash
set -e

echo "======================================"
echo "Playability Database Backup"
echo "======================================"

if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-playability_user}"
DB_NAME="${DB_NAME:-playability}"

BACKUP_DIR="${BACKUP_DIR:-./backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS="${RETENTION_DAYS:-30}"

mkdir -p "$BACKUP_DIR"

echo "Database: $DB_NAME"
echo "Backup directory: $BACKUP_DIR"
echo "Timestamp: $TIMESTAMP"
echo ""

backup_full() {
    echo "Creating full database backup..."
    BACKUP_FILE="$BACKUP_DIR/playability_full_$TIMESTAMP.dump"

    PGPASSWORD=$DB_PASSWORD pg_dump \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        -F c \
        -f "$BACKUP_FILE"

    echo "✓ Full backup created: $BACKUP_FILE"
    echo "Size: $(du -h "$BACKUP_FILE" | cut -f1)"
}

backup_schema() {
    echo ""
    echo "Creating schema-only backup..."
    SCHEMA_FILE="$BACKUP_DIR/playability_schema_$TIMESTAMP.sql"

    PGPASSWORD=$DB_PASSWORD pg_dump \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --schema-only \
        -f "$SCHEMA_FILE"

    echo "✓ Schema backup created: $SCHEMA_FILE"
}

backup_data() {
    echo ""
    echo "Creating data-only backup..."
    DATA_FILE="$BACKUP_DIR/playability_data_$TIMESTAMP.sql"

    PGPASSWORD=$DB_PASSWORD pg_dump \
        -h "$DB_HOST" \
        -p "$DB_PORT" \
        -U "$DB_USER" \
        -d "$DB_NAME" \
        --data-only \
        -f "$DATA_FILE"

    echo "✓ Data backup created: $DATA_FILE"
}

cleanup_old_backups() {
    echo ""
    echo "Cleaning up backups older than $RETENTION_DAYS days..."

    DELETED_COUNT=$(find "$BACKUP_DIR" -name "playability_*.dump" -type f -mtime +$RETENTION_DAYS -delete -print | wc -l)
    DELETED_COUNT=$((DELETED_COUNT + $(find "$BACKUP_DIR" -name "playability_*.sql" -type f -mtime +$RETENTION_DAYS -delete -print | wc -l)))

    if [ "$DELETED_COUNT" -gt 0 ]; then
        echo "✓ Deleted $DELETED_COUNT old backup files"
    else
        echo "✓ No old backups to delete"
    fi
}

compress_backup() {
    echo ""
    echo "Compressing backups..."

    for file in "$BACKUP_DIR"/playability_*_$TIMESTAMP.*; do
        if [ -f "$file" ]; then
            gzip -f "$file"
            echo "✓ Compressed: $file.gz"
        fi
    done
}

verify_backup() {
    echo ""
    echo "Verifying backup integrity..."

    BACKUP_FILE="$BACKUP_DIR/playability_full_$TIMESTAMP.dump.gz"

    if [ -f "$BACKUP_FILE" ]; then
        if gunzip -t "$BACKUP_FILE" 2>/dev/null; then
            echo "✓ Backup integrity verified"
        else
            echo "✗ Backup verification failed!"
            exit 1
        fi
    fi
}

list_backups() {
    echo ""
    echo "Current backups:"
    echo "----------------------------------------"
    ls -lh "$BACKUP_DIR" | grep playability_ || echo "No backups found"
    echo "----------------------------------------"

    TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
    echo "Total backup size: $TOTAL_SIZE"
}

main() {
    case "${1:-full}" in
        full)
            backup_full
            backup_schema
            compress_backup
            verify_backup
            cleanup_old_backups
            list_backups
            ;;
        schema)
            backup_schema
            compress_backup
            ;;
        data)
            backup_data
            compress_backup
            ;;
        list)
            list_backups
            ;;
        *)
            echo "Usage: $0 {full|schema|data|list}"
            echo ""
            echo "  full   - Full database backup (default)"
            echo "  schema - Schema-only backup"
            echo "  data   - Data-only backup"
            echo "  list   - List all backups"
            exit 1
            ;;
    esac

    echo ""
    echo "======================================"
    echo "Backup complete!"
    echo "======================================"
}

main "$@"

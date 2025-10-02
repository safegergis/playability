# Database Scripts and Utilities

This directory contains SQL schema files and utility scripts for managing the Playability database.

## Quick Start

### Option 1: Docker (Recommended)

The database will be automatically initialized when using Docker Compose:

```bash
docker-compose up -d postgres
```

SQL files are executed automatically on first startup.

### Option 2: Manual Setup

1. Create `.env` file in `playability-backend/`:
   ```bash
   cp playability-backend/.env.example playability-backend/.env
   # Edit .env with your database credentials
   ```

2. Initialize the database:
   ```bash
   ./database/scripts/init.sh
   ```

## Directory Structure

```
database/
├── README.md                    # This file
├── feature_support.sql          # Custom ENUM type definition
├── users.sql                    # Users table schema
├── games.sql                    # Games table schema
├── reports.sql                  # Reports table schema
├── migrations/                  # Database migration files
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
└── scripts/                     # Utility scripts
    ├── init.sh                  # Initialize database
    ├── migrate.sh               # Run migrations
    ├── backup.sh                # Backup database
    └── restore.sh               # Restore from backup
```

## Available Scripts

### 1. Initialize Database (`init.sh`)

Creates the database and runs all schema files with constraints and indexes.

```bash
./database/scripts/init.sh
```

**What it does:**
- Checks PostgreSQL connection
- Creates database (prompts to drop if exists)
- Runs all SQL schema files
- Adds foreign keys and constraints
- Creates performance indexes
- Verifies schema

**Environment variables:**
- `DB_HOST` (default: localhost)
- `DB_PORT` (default: 5432)
- `DB_USER` (default: playability_user)
- `DB_PASSWORD` (required)
- `DB_NAME` (default: playability)

---

### 2. Migration Tool (`migrate.sh`)

Manages database schema migrations with version tracking.

```bash
# Apply all pending migrations
./database/scripts/migrate.sh up

# Rollback last migration
./database/scripts/migrate.sh down

# Check migration status
./database/scripts/migrate.sh status

# Create new migration
./database/scripts/migrate.sh create add_user_roles
```

**Features:**
- Tracks applied migrations in `schema_migrations` table
- Supports up (apply) and down (rollback) migrations
- Creates timestamped migration files
- Shows pending and applied migrations

**Migration file naming:**
- Format: `{version}_{name}.{up|down}.sql`
- Example: `001_initial_schema.up.sql`

---

### 3. Backup Database (`backup.sh`)

Creates compressed database backups with automatic cleanup.

```bash
# Full backup (schema + data)
./database/scripts/backup.sh full

# Schema only
./database/scripts/backup.sh schema

# Data only
./database/scripts/backup.sh data

# List all backups
./database/scripts/backup.sh list
```

**Features:**
- Creates timestamped backups
- Compresses backups with gzip
- Verifies backup integrity
- Automatically deletes backups older than 30 days
- Shows backup sizes and total storage

**Environment variables:**
- `BACKUP_DIR` (default: ./backups)
- `RETENTION_DAYS` (default: 30)

**Backup formats:**
- Full: `.dump.gz` (PostgreSQL custom format, compressed)
- Schema: `.sql.gz` (SQL format, compressed)
- Data: `.sql.gz` (SQL format, compressed)

---

### 4. Restore Database (`restore.sh`)

Restores database from a backup file.

```bash
# Restore from backup
./database/scripts/restore.sh backups/playability_full_20241001_120000.dump.gz

# Or just provide the filename (will search in backup directory)
./database/scripts/restore.sh playability_full_20241001_120000.dump.gz

# List available backups if no file provided
./database/scripts/restore.sh
```

**Features:**
- Automatically decompresses `.gz` files
- Drops and recreates database
- Restores from `.dump` or `.sql` files
- Verifies restore by counting tables and records
- Safety confirmation prompt

**⚠️ Warning:**
This script will **DROP** the existing database and all data. Always backup first!

---

## Schema Files

### Execution Order

When using Docker or the init script, SQL files are executed in this order:

1. `feature_support.sql` - Creates ENUM type
2. `users.sql` - Creates users table
3. `games.sql` - Creates games table
4. `reports.sql` - Creates reports table

Then constraints and indexes are added programmatically.

### Manual Execution

If you need to run schema files manually:

```bash
export PGPASSWORD=your_password

# Create database
psql -h localhost -U playability_user -d postgres -c "CREATE DATABASE playability;"

# Run schema files
psql -h localhost -U playability_user -d playability -f database/feature_support.sql
psql -h localhost -U playability_user -d playability -f database/users.sql
psql -h localhost -U playability_user -d playability -f database/games.sql
psql -h localhost -U playability_user -d playability -f database/reports.sql
```

---

## Production Deployment

### First Deployment

1. Set production environment variables in `.env`:
   ```bash
   DB_HOST=your-db-host.rds.amazonaws.com
   DB_SSLMODE=require
   # ... other vars
   ```

2. Run initialization:
   ```bash
   ./database/scripts/init.sh
   ```

   Or use migrations:
   ```bash
   ./database/scripts/migrate.sh up
   ```

### Subsequent Deployments

1. Create migration for schema changes:
   ```bash
   ./database/scripts/migrate.sh create add_new_column
   ```

2. Edit migration files in `database/migrations/`

3. Test in staging environment

4. Apply to production:
   ```bash
   ./database/scripts/migrate.sh up
   ```

---

## Automated Backups (Production)

### Daily Backup Cron Job

Add to crontab (`crontab -e`):

```bash
# Daily backup at 2 AM
0 2 * * * cd /path/to/playability && ./database/scripts/backup.sh full >> /var/log/playability_backup.log 2>&1
```

### Weekly Backup to S3 (Example)

```bash
#!/bin/bash
./database/scripts/backup.sh full
aws s3 sync ./backups/ s3://your-bucket/playability-backups/
```

---

## Troubleshooting

### Connection refused
```
ERROR: Cannot connect to PostgreSQL
```

**Solutions:**
- Check if PostgreSQL is running: `systemctl status postgresql`
- Verify credentials in `.env` file
- Check network/firewall settings
- For Docker: ensure container is running

### Permission denied
```
ERROR: permission denied for database
```

**Solutions:**
- Verify database user has proper permissions
- For Docker: user is created automatically
- For manual setup: grant permissions:
  ```sql
  GRANT ALL PRIVILEGES ON DATABASE playability TO playability_user;
  ```

### Migration already applied
```
ERROR: duplicate key value violates unique constraint
```

**Solutions:**
- Check migration status: `./database/scripts/migrate.sh status`
- If migration failed midway, manually fix and re-run
- For stuck migrations, manually delete from `schema_migrations` table

### Backup file not found
```
Error: Backup file not found
```

**Solutions:**
- List available backups: `./database/scripts/backup.sh list`
- Provide full path to backup file
- Check `BACKUP_DIR` location

---

## Best Practices

### Development
1. Use Docker Compose for local development
2. Test migrations locally before production
3. Always create down migration for rollback
4. Keep schema files in sync with migrations

### Production
1. Always backup before schema changes
2. Test restores regularly
3. Monitor backup success/failures
4. Keep at least 30 days of backups
5. Store backups off-site (S3, etc.)
6. Use SSL connections (`DB_SSLMODE=require`)
7. Limit database user permissions
8. Enable connection pooling

### Migrations
1. One migration per logical change
2. Make migrations idempotent when possible
3. Test up and down migrations
4. Document breaking changes
5. Use descriptive migration names
6. Never edit applied migrations

---

## Additional Resources

- [DATABASE.md](../DATABASE.md) - Complete database documentation
- [DEPLOYMENT.md](../DEPLOYMENT.md) - Deployment checklist
- PostgreSQL documentation: https://www.postgresql.org/docs/

---

## Support

For issues or questions:
1. Check DATABASE.md for detailed schema information
2. Review troubleshooting section above
3. Check application logs in `playability-backend/`
4. Verify environment variables are set correctly

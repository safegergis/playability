# Database Documentation

This directory contains the PostgreSQL database schema, migrations, and management scripts for Playability.

## Quick Start

The database initializes automatically when you start Docker Compose for the first time:

```bash
docker-compose up -d
```

The initialization script creates all tables, indexes, constraints, and custom types automatically.

## Database Schema

### Custom Types

**`feature_support`** - ENUM for accessibility feature support levels:
- `'false'` - Feature not supported
- `'unknown'` - Support status unknown
- `'limited'` - Partial support
- `'true'` - Full support

### Tables

#### `users`
Stores user account information and authentication data.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-incrementing user ID |
| `username` | VARCHAR(255) | NOT NULL | User's display name |
| `email` | VARCHAR(255) | NOT NULL | User's email address |
| `hash` | VARCHAR(255) | NOT NULL | Bcrypt password hash |
| `num_reports` | INTEGER | NOT NULL, DEFAULT 0 | Number of reports submitted |

#### `games`
Stores game information and accessibility features.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Game ID (matches IGDB ID) |
| `name` | VARCHAR(255) | NOT NULL | Game title |
| `summary` | TEXT | NOT NULL | Game description |
| `cover_art` | VARCHAR(255) | NOT NULL | Cover image URL |
| `platforms` | JSONB | NOT NULL | Array of platform objects |
| `closed_captions` | feature_support | NOT NULL | Caption support level |
| `color_blind` | feature_support | NOT NULL | Colorblind mode support |
| `full_controller_support` | feature_support | NOT NULL | Controller support level |
| `controller_remapping` | feature_support | NOT NULL | Remapping support level |

#### `reports`
Stores user-submitted accessibility reports.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-incrementing report ID |
| `created_at` | TIMESTAMP | NOT NULL, DEFAULT now() | Report creation time |
| `game_id` | INTEGER | NOT NULL, FK | Reference to games table |
| `user_id` | INTEGER | NOT NULL, FK | Reference to users table |
| `closed_captions` | feature_support | NOT NULL | User's caption assessment |
| `color_blind` | feature_support | NOT NULL | User's colorblind assessment |
| `full_controller_support` | feature_support | NOT NULL | User's controller assessment |
| `controller_remapping` | feature_support | NOT NULL | User's remapping assessment |
| `score` | VARCHAR(255) | NOT NULL | Overall accessibility score (0-100) |
| `report` | VARCHAR(255) | NULL | Optional text feedback |

### Relationships

```
users (1) ──────< (N) reports (N) >────── (1) games
```

- One user can submit many reports
- One game can have many reports
- Each report belongs to one user and one game

## Directory Structure

```
database/
├── feature_support.sql      # ENUM type definition
├── users.sql               # Users table schema
├── games.sql               # Games table schema
├── reports.sql             # Reports table schema
├── migrations/             # Database migration files
│   ├── 001_initial_schema.up.sql
│   ├── 001_initial_schema.down.sql
│   ├── 002_add_verifications.up.sql
│   └── 002_add_verifications.down.sql
├── scripts/                # Database management scripts
│   ├── docker-entrypoint-initdb.sh
│   ├── backup.sh
│   └── restore.sh
└── README.md              # This file
```

## Common Operations

### Database Shell Access

```bash
docker-compose exec postgres psql -U playability_user -d playability
```

Once connected, you can run SQL queries:
```sql
SELECT * FROM users;
SELECT * FROM games WHERE id = 1;
SELECT * FROM reports WHERE game_id = 1;
```

### Backup Database

```bash
docker-compose --profile backup run --rm db-backup
```

Backups are saved to `./backups/` with timestamp:
```
backups/playability_full_YYYYMMDD_HHMMSS.dump.gz
```

### Restore from Backup

```bash
./database/scripts/restore.sh backups/playability_full_YYYYMMDD_HHMMSS.dump.gz
```

### Running Migrations

```bash
# Apply all pending migrations
docker-compose run --rm db-migrate /migrate.sh up

# Check migration status
docker-compose run --rm db-migrate /migrate.sh status

# Rollback last migration
docker-compose run --rm db-migrate /migrate.sh down
```

### Reset Database (Development Only)

```bash
# Warning: This deletes all data
docker-compose down -v
docker-compose up -d
```

## Initialization Process

When the PostgreSQL container starts for the first time, it automatically runs:

**`scripts/docker-entrypoint-initdb.sh`**

This script:
1. Creates the `feature_support` ENUM type
2. Creates `users`, `games`, and `reports` tables
3. Adds foreign key constraints
4. Creates indexes for performance
5. Sets up triggers and defaults

## Scripts

### `docker-entrypoint-initdb.sh`
Automatic initialization script that runs on first container startup. Creates all database objects in the correct order.

**When it runs:**
- Automatically on first `docker-compose up`
- Only runs if database is empty

**What it creates:**
- Custom ENUM types
- All tables with constraints
- Foreign key relationships
- Performance indexes
- Default values and triggers

### `backup.sh`
Creates compressed database backups.

**Usage:**
```bash
docker-compose --profile backup run --rm db-backup
```

**Output:**
- Location: `./backups/`
- Format: `playability_full_YYYYMMDD_HHMMSS.dump.gz`
- Type: PostgreSQL custom format with gzip compression

### `restore.sh`
Restores database from backup file.

**Usage:**
```bash
./database/scripts/restore.sh backups/playability_full_YYYYMMDD_HHMMSS.dump.gz
```

**Process:**
1. Drops existing database
2. Creates fresh database
3. Restores from backup file
4. Validates restoration

## Migrations

Database migrations are stored in `migrations/` directory and managed using golang-migrate.

### Migration Files

Migrations come in pairs:
- `XXX_name.up.sql` - Apply the migration
- `XXX_name.down.sql` - Rollback the migration

Example:
```
001_initial_schema.up.sql     # Creates initial tables
001_initial_schema.down.sql   # Drops initial tables
```

### Migration Commands

```bash
# Apply all pending migrations
docker-compose run --rm db-migrate /migrate.sh up

# Apply next migration
docker-compose run --rm db-migrate /migrate.sh up 1

# Rollback last migration
docker-compose run --rm db-migrate /migrate.sh down

# Check current version
docker-compose run --rm db-migrate /migrate.sh version

# View migration status
docker-compose run --rm db-migrate /migrate.sh status
```

### Creating New Migrations

For development, modify `scripts/docker-entrypoint-initdb.sh` and restart with fresh database.

For production schema changes:
1. Create migration files in `migrations/` directory
2. Follow naming convention: `XXX_description.up.sql` and `XXX_description.down.sql`
3. Test locally before applying to production
4. Always create both up and down migrations

## Data Types

### JSONB: `platforms` Column

The `games.platforms` column stores platform data as JSONB:

```json
[
  {"id": 6, "name": "PC"},
  {"id": 48, "name": "PlayStation 4"},
  {"id": 49, "name": "Xbox One"}
]
```

Query examples:
```sql
-- Find games on PC
SELECT * FROM games WHERE platforms @> '[{"id": 6}]';

-- Get all platform names
SELECT jsonb_array_elements(platforms)->>'name' AS platform
FROM games WHERE id = 1;
```

## Performance Considerations

### Indexes

The initialization script creates indexes on:
- Primary keys (automatic)
- Foreign keys (`reports.game_id`, `reports.user_id`)
- Frequently queried columns

### Connection Pooling

Configure in application code:
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Query Optimization

- All queries use parameterized statements
- JSONB columns use GIN indexes where needed
- Foreign key indexes improve join performance

## Security

### SQL Injection Protection
All application queries use parameterized statements:
```go
db.QueryRow("SELECT * FROM users WHERE email = $1", email)
```

### Password Security
- Passwords are hashed using bcrypt
- Only hashes are stored in `users.hash`
- Cost factor: 10 (default)

### Connection Security

**Development:** SSL disabled (`sslmode=disable`)
**Production:** SSL required (`sslmode=require`)

## Environment Variables

Database connection configured via:

| Variable | Example | Description |
|----------|---------|-------------|
| `DB_HOST` | `postgres` | Database hostname |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `playability_user` | Database username |
| `DB_PASSWORD` | `***` | Database password |
| `DB_NAME` | `playability` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |

For PostgreSQL container:
- `POSTGRES_DB` - Must match `DB_NAME`
- `POSTGRES_USER` - Must match `DB_USER`
- `POSTGRES_PASSWORD` - Must match `DB_PASSWORD`

## Troubleshooting

### Database won't initialize

```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Verify initialization script ran
docker-compose exec postgres ls -la /docker-entrypoint-initdb.d/
```

### Connection refused

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Verify connection settings in .env
cat .env | grep DB_
```

### Migration stuck

```bash
# Check migration status
docker-compose run --rm db-migrate /migrate.sh status

# View schema_migrations table
docker-compose exec postgres psql -U playability_user -d playability \
  -c "SELECT * FROM schema_migrations;"
```

### Fresh start needed

```bash
# Delete everything and start over
docker-compose down -v
docker-compose up -d
```

## Best Practices

1. **Always backup before schema changes**
   ```bash
   docker-compose --profile backup run --rm db-backup
   ```

2. **Test migrations locally first**
   ```bash
   docker-compose down -v
   docker-compose up -d
   # Test new migration
   ```

3. **Use transactions for multiple operations**
   ```go
   tx, _ := db.Begin()
   defer tx.Rollback()
   // ... operations ...
   tx.Commit()
   ```

4. **Close database connections properly**
   ```go
   rows, _ := db.Query(sql)
   defer rows.Close()
   ```

5. **Monitor slow queries**
   ```sql
   -- Enable in PostgreSQL config
   log_min_duration_statement = 1000
   ```

## Monitoring

### Useful Queries

**Check table sizes:**
```sql
SELECT
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

**Check connection count:**
```sql
SELECT count(*) FROM pg_stat_activity;
```

**Find slow queries:**
```sql
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

# Playability Database Documentation

## Overview

The Playability database is a PostgreSQL database designed to store game accessibility information, user accounts, and user-submitted accessibility reports. The schema supports tracking detailed accessibility features for video games and allowing users to contribute their own accessibility assessments.

---

## Database Schema

### Custom Types

#### `feature_support` ENUM

Represents the level of support for an accessibility feature.

**Values:**
- `'false'` - Feature is not supported
- `'unknown'` - Support status is unknown or not verified
- `'limited'` - Feature is partially supported with limitations
- `'true'` - Feature is fully supported

**File:** `database/feature_support.sql`

---

### Tables

#### 1. `users` Table

Stores user account information for authentication and tracking contributions.

**Columns:**

| Column Name   | Data Type          | Constraints                 | Description                          |
|---------------|--------------------|-----------------------------|--------------------------------------|
| `id`          | SERIAL (INTEGER)   | PRIMARY KEY, NOT NULL       | Auto-incrementing user identifier    |
| `username`    | VARCHAR(255)       | NOT NULL, UNIQUE*           | User's display name                  |
| `email`       | VARCHAR(255)       | NOT NULL, UNIQUE*           | User's email address for login       |
| `hash`        | VARCHAR(255)       | NOT NULL                    | Bcrypt hashed password               |
| `num_reports` | INTEGER            | NOT NULL, DEFAULT 0         | Count of reports submitted by user   |

**Indexes:**
- Primary key index on `id`
- Should have unique index on `email` (enforced in application)
- Should have unique index on `username` (enforced in application)

**Security Notes:**
- Passwords are stored using bcrypt hashing (see `auth/user.go`)
- Email addresses should be stored in lowercase for consistent lookups
- The application checks for duplicate emails/usernames before insertion

**File:** `database/users.sql`

---

#### 2. `games` Table

Stores game information and accessibility feature data sourced from IGDB and PCGamingWiki APIs.

**Columns:**

| Column Name                | Data Type          | Constraints              | Description                                    |
|----------------------------|--------------------|--------------------------|------------------------------------------------|
| `id`                       | SERIAL (INTEGER)   | PRIMARY KEY, NOT NULL    | Game identifier (matches IGDB ID)              |
| `name`                     | VARCHAR(255)       | NOT NULL                 | Game title                                     |
| `summary`                  | TEXT               | NOT NULL                 | Game description/summary                       |
| `cover_art`                | VARCHAR(255)       | NOT NULL                 | URL to game cover image                        |
| `platforms`                | JSONB              | NOT NULL                 | Array of platform objects (id, name, etc.)     |
| `closed_captions`          | feature_support    | NOT NULL                 | Closed caption/subtitle support level          |
| `color_blind`              | feature_support    | NOT NULL                 | Colorblind mode support level                  |
| `full_controller_support`  | feature_support    | NOT NULL                 | Full controller support level                  |
| `controller_remapping`     | feature_support    | NOT NULL                 | Controller remapping support level             |

**Indexes:**
- Primary key index on `id`
- Recommended: Index on `name` for search performance

**JSONB Column (`platforms`):**
The `platforms` column stores an array of platform objects in JSONB format:
```json
[
  {"id": 6, "name": "PC"},
  {"id": 48, "name": "PlayStation 4"}
]
```

**Data Source:**
- Game metadata from IGDB API
- Accessibility features from PCGamingWiki API

**File:** `database/games.sql`

---

#### 3. `reports` Table

Stores user-submitted accessibility reports for games, allowing crowd-sourced verification of accessibility features.

**Columns:**

| Column Name                | Data Type          | Constraints              | Description                                    |
|----------------------------|--------------------|--------------------------|------------------------------------------------|
| `id`                       | SERIAL (INTEGER)   | PRIMARY KEY, NOT NULL    | Auto-incrementing report identifier            |
| `created_at`               | TIMESTAMP          | NOT NULL, DEFAULT now()  | Timestamp when report was created              |
| `game_id`                  | INTEGER            | NOT NULL, FK → games(id) | Reference to the game being reported on        |
| `user_id`                  | INTEGER            | NOT NULL, FK → users(id) | Reference to the user who created the report   |
| `closed_captions`          | feature_support    | NOT NULL                 | User's assessment of closed caption support    |
| `color_blind`              | feature_support    | NOT NULL                 | User's assessment of colorblind mode support   |
| `full_controller_support`  | feature_support    | NOT NULL                 | User's assessment of controller support        |
| `controller_remapping`     | feature_support    | NOT NULL                 | User's assessment of remapping support         |
| `score`                    | VARCHAR(255)       | NOT NULL                 | Overall accessibility score (0-100)            |
| `report`                   | VARCHAR(255)       | NULL                     | Optional text feedback/comments                |

**Indexes:**
- Primary key index on `id`
- Recommended: Index on `(game_id, created_at)` for efficient report retrieval
- Recommended: Index on `user_id` for user report lookups
- Recommended: Index on `(game_id, user_id)` for duplicate detection

**Foreign Keys (Recommended):**
- `game_id` → `games(id)` ON DELETE CASCADE
- `user_id` → `users(id)` ON DELETE CASCADE

**Business Rules (Enforced in Application):**
- One report per user per game (enforced in `db/reports.go`)
- All reports pass through AI moderation before storage (see `pkg/ai/moderation.go`)

**File:** `database/reports.sql`

---

## Relationships

### Entity Relationship Diagram

```
┌─────────────┐
│    users    │
│─────────────│
│ id (PK)     │
│ username    │
│ email       │
│ hash        │
│ num_reports │
└──────┬──────┘
       │
       │ 1:N
       │
       ▼
┌──────────────┐         ┌─────────────┐
│   reports    │   N:1   │    games    │
│──────────────│◄────────│─────────────│
│ id (PK)      │         │ id (PK)     │
│ created_at   │         │ name        │
│ game_id (FK) │─────────│ summary     │
│ user_id (FK) │         │ cover_art   │
│ ...features  │         │ platforms   │
│ score        │         │ ...features │
│ report       │         └─────────────┘
└──────────────┘
```

### Relationship Details

1. **users → reports** (One-to-Many)
   - One user can create many reports
   - Each report belongs to exactly one user
   - Constraint: One report per user per game

2. **games → reports** (One-to-Many)
   - One game can have many reports
   - Each report is about exactly one game
   - Reports provide crowd-sourced accessibility verification

---

## Database Operations

### Database Layer (`playability-backend/db/`)

The database layer uses Go's `database/sql` package with prepared statements for security.

#### Files:

- **`database.go`** - Database connection initialization
  - `InitDB()` - Establishes PostgreSQL connection using environment variables
  - Connection string format: `host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`

- **`user.go`** - User management operations
  - `InsertUser(user)` - Create new user with duplicate checking
  - `CheckUser(email, password)` - Authenticate user and return user ID
  - `QueryUser(userID)` - Retrieve user information

- **`game.go`** - Game data operations
  - `InsertGame(body)` - Insert game with accessibility features
  - `QueryGame(id)` - Retrieve game by ID with JSONB platform parsing

- **`reports.go`** - Report management operations
  - `InsertReport(report)` - Create new report (checks for duplicates)
  - `QueryReportCards(gameID)` - Get all report cards for a game
  - `QueryAccessibilityScores(gameID)` - Get all scores for score calculation

- **`features.go`** - Feature aggregation operations
  - `QueryFeatureReports(gameID)` - Get all feature assessments for consensus calculation

---

## Queries and Performance

### Common Query Patterns

#### 1. User Authentication
```sql
SELECT hash FROM users WHERE email = $1;
SELECT id FROM users WHERE email = $1;
```

#### 2. Game Retrieval
```sql
SELECT * FROM games WHERE id = $1;
```

#### 3. Report Submission
```sql
-- Check for existing report
SELECT id FROM reports WHERE game_id = $1 AND user_id = $2;

-- Insert new report
INSERT INTO reports (game_id, user_id, closed_captions, color_blind,
                     full_controller_support, controller_remapping, score, report)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
```

#### 4. Report Retrieval
```sql
-- Get report cards for a game
SELECT id, created_at, game_id, user_id, score, report
FROM reports
WHERE game_id = $1
ORDER BY created_at DESC;

-- Get feature reports for consensus
SELECT closed_captions, color_blind, full_controller_support, controller_remapping
FROM reports
WHERE game_id = $1;
```

### Performance Recommendations

#### Critical Indexes to Add:

```sql
-- Users table
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);

-- Games table
CREATE INDEX idx_games_name ON games(name);
CREATE INDEX idx_games_name_trgm ON games USING gin(name gin_trgm_ops); -- For fuzzy search

-- Reports table
CREATE INDEX idx_reports_game_id_created_at ON reports(game_id, created_at DESC);
CREATE INDEX idx_reports_user_id ON reports(user_id);
CREATE INDEX idx_reports_game_user ON reports(game_id, user_id); -- For duplicate checks
```

#### Query Optimization Tips:

1. **Enable pg_trgm extension** for fuzzy game name search:
   ```sql
   CREATE EXTENSION pg_trgm;
   ```

2. **Connection Pooling** - Configure in `database.go`:
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **Prepared Statements** - Already implemented in all queries

4. **ANALYZE Tables** - Run periodically:
   ```sql
   ANALYZE users;
   ANALYZE games;
   ANALYZE reports;
   ```

---

## Data Integrity and Constraints

### Current Constraints

#### Defined in SQL:
- Primary keys on all tables
- NOT NULL constraints on essential fields
- DEFAULT values for `created_at` and `num_reports`

#### Enforced in Application Layer:
- **Unique email/username** (`db/user.go` lines 24-45)
- **One report per user per game** (`db/reports.go` lines 20-32)
- **Foreign key relationships** (should be added to schema)
- **Password complexity** (should be added)
- **Email format validation** (should be added)

### Missing Constraints (Recommended):

```sql
-- Add foreign keys
ALTER TABLE reports
ADD CONSTRAINT fk_reports_game
FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE;

ALTER TABLE reports
ADD CONSTRAINT fk_reports_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add unique constraints
ALTER TABLE users
ADD CONSTRAINT unique_email UNIQUE (email);

ALTER TABLE users
ADD CONSTRAINT unique_username UNIQUE (username);

-- Add check constraints
ALTER TABLE reports
ADD CONSTRAINT check_score_range
CHECK (score::integer >= 0 AND score::integer <= 100);

ALTER TABLE users
ADD CONSTRAINT check_num_reports_positive
CHECK (num_reports >= 0);
```

---

## Security Considerations

### 1. SQL Injection Prevention
✅ All queries use parameterized statements (`$1`, `$2`, etc.)
✅ No string concatenation in SQL queries

### 2. Password Security
✅ Bcrypt hashing with default cost (10)
✅ Passwords never stored in plain text
✅ Hash stored in `users.hash` column (VARCHAR(255))

### 3. Data Validation

**Current Validation:**
- Email/username uniqueness checked before insertion
- Duplicate report detection
- AI moderation on user-submitted reports

**Missing Validation:**
- Email format validation
- Password complexity requirements
- Input length validation
- XSS protection for text fields

### 4. Connection Security

**Current:**
- SSL mode disabled in development (`sslmode=disable`)

**Production Requirements:**
- Enable SSL/TLS: `sslmode=require`
- Use SSL certificates
- Restrict database network access
- Use connection pooling with limits

---

## Database Initialization

### Current Setup

The database can be initialized manually by running SQL files in order:

```bash
psql -U playability_user -d playability -f database/feature_support.sql
psql -U playability_user -d playability -f database/users.sql
psql -U playability_user -d playability -f database/games.sql
psql -U playability_user -d playability -f database/reports.sql
```

### Docker Setup

When using Docker Compose, SQL files in `/database` are automatically executed on first startup via the `docker-entrypoint-initdb.d` mechanism (see `docker-compose.yml`).

**Execution Order:**
Files are executed in alphabetical order. Current naming doesn't guarantee correct order.

**Recommended Renaming:**
```
001_feature_support.sql
002_users.sql
003_games.sql
004_reports.sql
```

---

## Backup and Recovery

### Backup Strategy (Recommended)

#### 1. Full Database Backup
```bash
pg_dump -U playability_user -d playability -F c -f backup_$(date +%Y%m%d_%H%M%S).dump
```

#### 2. Schema-Only Backup
```bash
pg_dump -U playability_user -d playability --schema-only -f schema_backup.sql
```

#### 3. Data-Only Backup
```bash
pg_dump -U playability_user -d playability --data-only -f data_backup.sql
```

#### 4. Table-Specific Backup
```bash
pg_dump -U playability_user -d playability -t users -t reports -f user_data_backup.sql
```

### Restore Procedures

#### From Custom Format Dump:
```bash
pg_restore -U playability_user -d playability -c backup.dump
```

#### From SQL File:
```bash
psql -U playability_user -d playability -f backup.sql
```

### Automated Backup (Production)

See `database/scripts/backup.sh` for automated daily backup script.

---

## Migration Strategy

### Current State
- No migration system in place
- Schema changes must be applied manually
- Risk of schema drift between environments

### Recommended Migration Tools

#### Option 1: golang-migrate
```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir database/migrations -seq create_indexes

# Run migrations
migrate -path database/migrations -database "postgres://localhost/playability?sslmode=disable" up
```

#### Option 2: goose
```bash
# Install
go install github.com/pressly/goose/v3/cmd/goose@latest

# Create migration
goose -dir database/migrations create add_foreign_keys sql

# Run migrations
goose -dir database/migrations postgres "user=playability_user dbname=playability" up
```

---

## Environment Variables

Database connection configuration via environment variables:

| Variable       | Description                  | Example                    | Required |
|----------------|------------------------------|----------------------------|----------|
| `DB_HOST`      | Database host                | `localhost`, `postgres`    | Yes      |
| `DB_PORT`      | Database port                | `5432`                     | Yes      |
| `DB_USER`      | Database username            | `playability_user`         | Yes      |
| `DB_PASSWORD`  | Database password            | `secure_password_123`      | Yes      |
| `DB_NAME`      | Database name                | `playability`              | Yes      |
| `DB_SSLMODE`   | SSL mode                     | `disable`, `require`       | No       |

**Default SSL Mode:** `disable` (hardcoded in `db/database.go`)

**Production Recommendation:** Make `DB_SSLMODE` configurable via environment variable.

---

## Schema Version History

### v1.0 (Current)
- Initial schema with `users`, `games`, `reports` tables
- Custom `feature_support` ENUM type
- Basic constraints (primary keys, NOT NULL)
- JSONB for platform data

### Pending Changes (v1.1)
- Add foreign key constraints
- Add unique constraints on users.email and users.username
- Add indexes for performance
- Add check constraints for data validation
- Make sslmode configurable
- Rename SQL files with numeric prefixes

---

## Known Issues and Technical Debt

### 1. Missing Foreign Keys
**Issue:** Foreign key relationships not enforced at database level
**Impact:** Risk of orphaned records
**Priority:** High
**Fix:** Add foreign key constraints (see recommended constraints)

### 2. No Unique Constraints on Email/Username
**Issue:** Uniqueness only enforced in application layer
**Impact:** Race condition could allow duplicates
**Priority:** High
**Fix:** Add unique indexes/constraints

### 3. Hardcoded SSL Mode
**Issue:** SSL mode hardcoded to `disable` in `database.go:30`
**Impact:** Cannot enable SSL without code changes
**Priority:** Medium
**Fix:** Read from `DB_SSLMODE` environment variable

### 4. No Migration System
**Issue:** No automated schema versioning or migration
**Impact:** Difficult to track and apply schema changes
**Priority:** Medium
**Fix:** Implement golang-migrate or goose

### 5. Missing Indexes
**Issue:** No indexes on foreign keys or frequently queried columns
**Impact:** Poor query performance as data grows
**Priority:** Medium
**Fix:** Add recommended indexes

### 6. Score Column Type
**Issue:** `reports.score` is VARCHAR(255) but stored as integer
**Impact:** Type mismatch, unnecessary storage
**Priority:** Low
**Fix:** Change to INTEGER or SMALLINT

### 7. Report Column Size
**Issue:** `reports.report` limited to 255 characters
**Impact:** Users cannot leave detailed feedback
**Priority:** Low
**Fix:** Change to TEXT type

### 8. No Created/Updated Timestamps on Users/Games
**Issue:** Only reports table has timestamp tracking
**Impact:** Cannot track when users registered or games were added
**Priority:** Low
**Fix:** Add `created_at` and `updated_at` columns

### 9. Connection Pool Not Configured
**Issue:** No connection pool limits set in `database.go`
**Impact:** Potential resource exhaustion
**Priority:** Medium
**Fix:** Add `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`

### 10. SQL File Execution Order
**Issue:** Files in `database/` don't have numeric prefixes
**Impact:** Unpredictable execution order in Docker init
**Priority:** Low
**Fix:** Rename files to `001_`, `002_`, etc.

---

## Best Practices for Development

### 1. Always Use Transactions for Multi-Statement Operations
```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// ... perform operations ...

return tx.Commit()
```

### 2. Check for sql.ErrNoRows
```go
err := db.QueryRow(query, id).Scan(&result)
if err == sql.ErrNoRows {
    // Handle not found case
    return nil, fmt.Errorf("not found")
}
```

### 3. Always Close Rows
```go
rows, err := db.Query(query)
if err != nil {
    return err
}
defer rows.Close()
```

### 4. Use Context for Timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := db.QueryRowContext(ctx, query, id).Scan(&result)
```

### 5. Log Errors, Don't Panic
```go
// Bad
log.Fatal("Error:", err)

// Good
log.Println("Error:", err)
return err
```

---

## Testing

### Test Database Setup

Create a separate test database:
```sql
CREATE DATABASE playability_test;
```

Set environment variables for tests:
```bash
export DB_NAME=playability_test
go test ./db/...
```

### Recommended Test Coverage

1. **User Operations**
   - User registration with duplicate email/username
   - User login with correct/incorrect credentials
   - User retrieval

2. **Game Operations**
   - Game insertion and retrieval
   - JSONB platform parsing

3. **Report Operations**
   - Report submission
   - Duplicate report prevention
   - Report retrieval by game
   - Score aggregation

4. **Edge Cases**
   - NULL values
   - Empty strings
   - SQL injection attempts
   - Unicode characters

---

## Monitoring and Maintenance

### Database Metrics to Monitor

1. **Connection Pool**
   - Active connections
   - Idle connections
   - Connection wait time

2. **Query Performance**
   - Slow query log
   - Query execution times
   - Index usage statistics

3. **Storage**
   - Database size
   - Table sizes
   - Index sizes

4. **Locks and Deadlocks**
   - Lock wait times
   - Deadlock frequency

### Useful PostgreSQL Queries

#### Table Sizes:
```sql
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### Index Usage:
```sql
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;
```

#### Slow Queries (requires pg_stat_statements):
```sql
SELECT
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

---

## Appendix: Full Schema SQL

See migration scripts in `database/scripts/` for complete schema creation with all recommended constraints and indexes.

---

## Questions or Issues?

For database-related questions or issues, please refer to:
- This documentation
- DEPLOYMENT.md for deployment-specific database setup
- playability-backend/db/ for implementation details

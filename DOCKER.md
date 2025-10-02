# Docker Setup and Usage Guide

This guide covers all Docker-related operations for the Playability application.

## Quick Start

### First Time Setup

1. **Copy environment files:**
   ```bash
   cp .env.example .env
   cp playability-backend/.env.example playability-backend/.env
   ```

2. **Edit environment files with your credentials:**
   ```bash
   # Generate a secure JWT secret
   openssl rand -base64 32

   # Edit playability-backend/.env and add:
   # - DB_PASSWORD
   # - JWT_SECRET (from command above)
   # - IGDB_ACCESS_TOKEN
   # - CLAUDE_API_KEY
   ```

3. **Start the application:**
   ```bash
   make dev
   ```

   Or without make:
   ```bash
   docker-compose up -d
   ```

4. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

## Using the Makefile

The Makefile provides convenient shortcuts for common operations.

### View All Commands
```bash
make help
```

### Development Commands

```bash
# Start all services
make up

# View logs (all services)
make logs

# View logs for specific service
make logs SVC=backend
make logs SVC=frontend
make logs SVC=postgres

# Stop all services
make down

# Restart all services
make restart

# View running containers
make ps
```

### Database Commands

```bash
# Initialize database (first time setup)
make db-init

# Run migrations
make db-migrate

# Check migration status
make db-migrate-status

# Rollback last migration
make db-migrate-down

# Create backup
make db-backup

# Restore from backup
make db-restore

# Open PostgreSQL shell
make db-shell

# View database logs
make db-logs
```

### Build and Cleanup

```bash
# Rebuild Docker images
make build

# Clean everything (WARNING: deletes data)
make clean
```

## Docker Compose Services

### Main Services

#### **postgres**
- PostgreSQL 16 database
- Automatically initializes schema on first run
- Data persisted in `postgres_data` volume
- Backups stored in `./backups` directory

#### **backend**
- Go API server
- Connects to postgres
- Exposes port 8080
- Auto-restarts on failure

#### **frontend**
- Nuxt 3 application
- Connects to backend
- Exposes port 3000
- Auto-restarts on failure

### Utility Services (Profiles)

#### **db-backup** (profile: backup)
- Creates database backups on-demand
- Usage: `docker-compose --profile backup run --rm db-backup`
- Or: `make db-backup`

#### **db-migrate** (profile: tools)
- Runs database migrations
- Usage: `docker-compose run --rm db-migrate /migrate.sh up`
- Or: `make db-migrate`

## Database Initialization

### Automatic Initialization

When you start the postgres container for the first time, it automatically:

1. Creates the database
2. Creates custom types (feature_support enum)
3. Creates all tables (users, games, reports)
4. Adds foreign key constraints
5. Adds unique constraints
6. Creates indexes for performance
7. Sets up triggers for updated_at columns
8. Creates migration tracking table

**Script:** `database/scripts/docker-entrypoint-initdb.sh`

### Manual Initialization

If you need to reinitialize the database:

```bash
# Stop and remove the database volume
docker-compose down
docker volume rm playability_postgres_data

# Start again (will reinitialize)
docker-compose up -d postgres
```

Or use the init script:
```bash
make db-init
```

## Running Migrations

### Apply Migrations

```bash
# Using make
make db-migrate

# Using docker-compose
docker-compose run --rm db-migrate /migrate.sh up

# Check status
make db-migrate-status
```

### Create New Migration

```bash
# On host machine
./database/scripts/migrate.sh create add_user_roles

# Edit the created files:
# database/migrations/{timestamp}_add_user_roles.up.sql
# database/migrations/{timestamp}_add_user_roles.down.sql

# Apply the migration
make db-migrate
```

### Rollback Migration

```bash
# Rollback last migration
make db-migrate-down

# Or with docker-compose
docker-compose run --rm db-migrate /migrate.sh down
```

## Database Backups

### Create Backup

```bash
# Using make (recommended)
make db-backup

# Using docker-compose
docker-compose --profile backup run --rm db-backup

# Backups are saved to ./backups/ directory
```

### Restore from Backup

```bash
# Using make (interactive)
make db-restore

# Or manually
./database/scripts/restore.sh backups/playability_full_20241001_120000.dump.gz
```

### Automated Backups

For production, set up a cron job:

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * cd /path/to/playability && make db-backup >> /var/log/playability-backup.log 2>&1
```

Or use the GitHub Actions workflow:
- File: `.github/workflows/backup.yml`
- Runs daily at 2 AM UTC
- Stores backups as GitHub Artifacts (30 day retention)

## Development Workflow

### Starting Fresh

```bash
# Clean everything
make clean

# Start from scratch
make dev

# Check logs
make logs
```

### Making Code Changes

**Frontend changes:**
- Rebuild: `docker-compose up -d --build frontend`
- Or: `make build && make up`

**Backend changes:**
- Rebuild: `docker-compose up -d --build backend`
- Or: `make build && make up`

### Debugging

```bash
# View logs for all services
make logs

# View logs for specific service
make logs SVC=backend

# Open shell in container
docker-compose exec backend sh
docker-compose exec frontend sh

# Open database shell
make db-shell
```

## Production Deployment

### Using Production Compose File

```bash
# Start production services
make prod-up

# Or manually
docker-compose -f docker-compose.production.yml up -d

# Create production backup
make prod-backup

# Stop production services
make prod-down
```

### Production Differences

The production compose file (`docker-compose.production.yml`) includes:

- **SSL/TLS configuration** for database
- **Resource limits** (CPU/memory)
- **Multiple replicas** for frontend/backend
- **Nginx reverse proxy** (requires configuration)
- **Always restart policy**
- **Stricter security settings**

### Required Production Setup

Before using production mode:

1. **Create nginx configuration:**
   ```bash
   mkdir -p nginx
   # Add nginx.conf (see DEPLOYMENT.md)
   ```

2. **Set up SSL certificates:**
   ```bash
   mkdir -p nginx/ssl
   # Add SSL certificates
   ```

3. **Configure environment variables:**
   - Set strong JWT_SECRET
   - Set ALLOWED_ORIGINS to your domain
   - Set DB_SSLMODE=require
   - Configure all API keys

4. **Test locally:**
   ```bash
   docker-compose -f docker-compose.production.yml up
   ```

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs backend

# Check status
make ps

# Restart specific service
docker-compose restart backend
```

### Database connection issues

```bash
# Check if postgres is healthy
docker-compose ps postgres

# View postgres logs
make db-logs

# Test connection
docker-compose exec postgres psql -U playability_user -d playability -c "SELECT 1;"
```

### Port already in use

```bash
# Find process using port 3000, 8080, or 5432
lsof -i :3000
lsof -i :8080
lsof -i :5432

# Kill the process or change port in docker-compose.yml
```

### Out of disk space

```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove everything (WARNING: deletes data)
docker system prune -a --volumes
```

### Database not initializing

```bash
# Remove volume and recreate
docker-compose down
docker volume rm playability_postgres_data
docker-compose up -d postgres

# Check initialization logs
docker-compose logs postgres
```

### Migration stuck

```bash
# Check migration status
make db-migrate-status

# If stuck, manually check
docker-compose exec postgres psql -U playability_user -d playability -c "SELECT * FROM schema_migrations;"

# Fix manually if needed
docker-compose exec postgres psql -U playability_user -d playability
```

## Environment Variables

### Frontend (.env)

```bash
# IGDB API
IGDB_CLIENT_SECRET=your_secret_here

# Backend URL (change for production)
NUXT_PUBLIC_API_URL=http://localhost:8080
```

### Backend (playability-backend/.env)

```bash
# Database
DB_HOST=postgres  # Use 'postgres' for Docker, 'localhost' for local dev
DB_PORT=5432
DB_USER=playability_user
DB_PASSWORD=your_secure_password
DB_NAME=playability
DB_SSLMODE=disable  # Use 'require' for production

# Security
JWT_SECRET=your_32_char_secret_here

# APIs
IGDB_ACCESS_TOKEN=your_token
CLAUDE_API_KEY=your_key

# App Config
PORT=8080
ENVIRONMENT=development
ALLOWED_ORIGINS=http://localhost:3000
```

## Health Checks

All services have health checks configured:

```bash
# Check health status
docker-compose ps

# Healthy services show: "healthy"
# Unhealthy services show: "unhealthy"
```

**Health check endpoints:**
- Backend: `http://localhost:8080/health` (needs to be implemented)
- Frontend: `http://localhost:3000/api/health` (needs to be implemented)
- Database: `pg_isready` command

## Volumes

### Named Volumes

- **postgres_data**: Stores all database data
  - Location: Docker manages this
  - Backup: Use `make db-backup`
  - Remove: `docker volume rm playability_postgres_data`

### Bind Mounts

- **./backups**: Database backups
- **./database/scripts**: Database scripts
- **./database/migrations**: Migration files

## Networking

All services communicate via the `playability-network` bridge network:

```bash
# Inspect network
docker network inspect playability_playability-network

# Services can reach each other by service name:
# - backend → postgres
# - frontend → backend
# - db-migrate → postgres
```

## Best Practices

### Development

1. **Always use make commands** for consistency
2. **Check logs** when something fails: `make logs SVC=<service>`
3. **Create backups** before major changes: `make db-backup`
4. **Test migrations** before applying: `make db-migrate-status`

### Production

1. **Use production compose file**: `docker-compose.production.yml`
2. **Enable SSL** for database connections
3. **Set resource limits** to prevent resource exhaustion
4. **Use secrets management** (not .env files)
5. **Monitor health checks** and set up alerts
6. **Automate backups** with cron or CI/CD
7. **Test restore procedure** regularly

### Security

1. **Never commit .env files**
2. **Use strong passwords** for DB_PASSWORD and JWT_SECRET
3. **Rotate secrets** regularly
4. **Limit CORS origins** in production
5. **Enable SSL/TLS** for all connections
6. **Run containers as non-root** (already configured)
7. **Scan images** for vulnerabilities

## Additional Resources

- [DATABASE.md](DATABASE.md) - Database documentation
- [DEPLOYMENT.md](DEPLOYMENT.md) - Full deployment checklist
- [database/README.md](database/README.md) - Database scripts guide
- [Docker Compose Docs](https://docs.docker.com/compose/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)

## Support

For issues:
1. Check logs: `make logs`
2. Check this troubleshooting section
3. Review container status: `make ps`
4. Check environment variables are set correctly

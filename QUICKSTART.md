# Playability Quick Start

## First Time Setup

1. **Create your environment file:**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys
   ```

2. **Start everything:**
   ```bash
   make dev
   ```

That's it! Database creates itself automatically.

## Daily Development

```bash
make dev        # Start all services
make logs       # View all logs
make down       # Stop everything
```

## Common Commands

| What you want to do | Command |
|-------------------|---------|
| Start development | `make dev` |
| View logs | `make logs` |
| View specific service logs | `make logs SVC=backend` |
| Stop everything | `make down` |
| Restart everything | `make restart` |
| Access database | `make db-shell` |
| Backup database | `make db-backup` |
| Restore database | `make db-restore` |
| Start fresh (delete all data) | `docker-compose down -v && make dev` |

## Services & Ports

- **Frontend:** http://localhost:3000 (with hot reload)
- **Backend:** http://localhost:8080
- **Database:** localhost:5432

## Environment Variables You Need

In your `.env` file:

```bash
# Database (already set for Docker)
DB_HOST=postgres
DB_PORT=5432
DB_USER=playability_user
DB_PASSWORD=your_secure_password
DB_NAME=playability

# Backend
JWT_SECRET=your_jwt_secret_here_minimum_32_characters

# External APIs (get these from their websites)
IGDB_ACCESS_TOKEN=your_token_here
CLAUDE_API_KEY=your_key_here
```

## Troubleshooting

**"Port already in use"**
```bash
make down
# Wait 5 seconds
make dev
```

**"Database connection failed"**
```bash
docker-compose down -v  # Delete old database
make dev                # Start fresh
```

**"Can't see my code changes"**
- Frontend: Should auto-reload (check logs: `make logs SVC=frontend-dev`)
- Backend: Rebuild with `docker-compose build backend && make restart`

**"I messed up the database"**
```bash
docker-compose down -v  # Delete database
make dev                # Start fresh
```

## Need More Help?

- Database info: `database/README.md`
- Full docs: `CLAUDE.md`
- Docker info: `DOCKER.md`

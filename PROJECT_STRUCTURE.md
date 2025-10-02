# Playability Project Structure

This document describes the organization of the Playability codebase.

## Directory Structure

```
playability/
├── frontend/                    # Nuxt 3 Frontend Application
│   ├── assets/                 # Static assets (CSS, images)
│   ├── components/             # Vue components
│   │   └── ui/                # ShadCN UI components
│   ├── layouts/               # Nuxt layouts
│   ├── lib/                   # Utility libraries
│   ├── pages/                 # Nuxt pages (routes)
│   ├── public/                # Public static files
│   ├── server/                # Nuxt server API routes
│   │   └── api/              # API route handlers
│   ├── types/                 # TypeScript type definitions
│   ├── .dockerignore          # Docker ignore file for frontend
│   ├── .env.example           # Environment variable template
│   ├── Dockerfile             # Frontend Docker build configuration
│   ├── components.json        # ShadCN component configuration
│   ├── eslint.config.mjs      # ESLint configuration
│   ├── nuxt.config.ts         # Nuxt configuration
│   ├── package.json           # Frontend dependencies
│   ├── pnpm-lock.yaml         # Lock file for pnpm
│   ├── tailwind.config.js     # Tailwind CSS configuration
│   └── tsconfig.json          # TypeScript configuration
│
├── backend/                    # Go Backend API
│   ├── auth/                  # Authentication utilities
│   │   └── user.go           # JWT and password hashing
│   ├── cmd/                   # Application entry points
│   │   └── api/
│   │       └── main.go       # Main API server
│   ├── db/                    # Database layer
│   │   ├── database.go       # Database connection
│   │   ├── game.go           # Game queries
│   │   ├── user.go           # User queries
│   │   ├── reports.go        # Report queries
│   │   └── features.go       # Feature aggregation
│   ├── handlers/              # HTTP request handlers
│   ├── pkg/                   # Utility packages
│   │   ├── ai/               # AI moderation services
│   │   ├── calc/             # Score calculation
│   │   └── fetch/            # External API integrations
│   ├── types/                 # Type definitions
│   │   └── types.go          # Go struct definitions
│   ├── .dockerignore          # Docker ignore file for backend
│   ├── .env.example           # Environment variable template
│   ├── Dockerfile             # Backend Docker build configuration
│   ├── go.mod                 # Go module definition
│   └── go.sum                 # Go dependency checksums
│
├── database/                   # Database Schema and Scripts
│   ├── migrations/            # Database migrations
│   │   ├── 001_initial_schema.up.sql
│   │   └── 001_initial_schema.down.sql
│   ├── scripts/               # Database utility scripts
│   │   ├── docker-entrypoint-initdb.sh  # Auto-init for Docker
│   │   ├── init.sh                      # Manual initialization
│   │   ├── migrate.sh                   # Migration runner
│   │   ├── backup.sh                    # Backup utility
│   │   └── restore.sh                   # Restore utility
│   ├── feature_support.sql    # ENUM type definition
│   ├── games.sql              # Games table schema
│   ├── reports.sql            # Reports table schema
│   ├── users.sql              # Users table schema
│   └── README.md              # Database scripts documentation
│
├── .github/                    # GitHub Configuration
│   └── workflows/             # GitHub Actions workflows
│       └── backup.yml         # Automated backup workflow
│
├── backups/                    # Database backups (gitignored)
│
├── docker-compose.yml          # Development environment setup
├── docker-compose.production.yml  # Production environment setup
├── Dockerfile                  # (Legacy - moved to frontend/)
├── Makefile                    # Convenient development commands
│
├── CLAUDE.md                   # Instructions for Claude Code
├── DATABASE.md                 # Complete database documentation
├── DEPLOYMENT.md               # Deployment checklist
├── DOCKER.md                   # Docker usage guide
├── PROJECT_STRUCTURE.md        # This file
├── README.md                   # Project overview
│
└── .gitignore                  # Git ignore configuration
```

## Component Responsibilities

### Frontend (`frontend/`)

The frontend is a **Nuxt 3** application with **Vue.js** and **TypeScript**.

**Key Technologies:**
- Nuxt 3 (SSR framework)
- Vue 3 (UI framework)
- ShadCN Vue (UI components)
- Tailwind CSS (styling)
- Vee Validate (form validation)
- Nuxt Icon (icons)
- Embla Carousel (carousels)

**Key Files:**
- `nuxt.config.ts` - Main configuration
- `server/api/` - Server-side API routes (proxy to backend)
- `pages/` - Application routes
- `components/ui/` - Reusable UI components

**Environment Variables:**
- `IGDB_CLIENT_SECRET` - IGDB API key
- `NUXT_PUBLIC_API_URL` - Backend API URL

---

### Backend (`backend/`)

The backend is a **Go REST API** with the **Chi** router.

**Key Technologies:**
- Go 1.23.1
- Chi router
- PostgreSQL driver (lib/pq)
- JWT authentication
- Bcrypt password hashing
- Claude AI integration
- IGDB and PCGamingWiki APIs

**Architecture:**
- Dependency injection pattern (`Env` struct)
- Layered architecture (handlers → db → types)
- Middleware for auth, logging, CORS

**Key Files:**
- `cmd/api/main.go` - Application entry point
- `handlers/` - HTTP request handlers
- `db/` - Database access layer
- `auth/` - Authentication utilities
- `pkg/` - External integrations

**Environment Variables:**
- Database: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- Security: `JWT_SECRET`
- APIs: `IGDB_ACCESS_TOKEN`, `CLAUDE_API_KEY`
- Config: `PORT`, `ENVIRONMENT`, `ALLOWED_ORIGINS`

---

### Database (`database/`)

PostgreSQL database with migration system and utility scripts.

**Schema Files:**
- `feature_support.sql` - Custom ENUM type
- `users.sql` - User accounts table
- `games.sql` - Game information table
- `reports.sql` - User-submitted reports table

**Migrations:**
- Located in `migrations/` directory
- Managed by `migrate.sh` script
- Version tracked in `schema_migrations` table

**Scripts:**
- `init.sh` - Initialize database from scratch
- `migrate.sh` - Apply/rollback migrations
- `backup.sh` - Create backups
- `restore.sh` - Restore from backup
- `docker-entrypoint-initdb.sh` - Auto-init for Docker

---

## Data Flow

### User Registration Flow
```
Frontend           Backend              Database
   |                  |                    |
   |-- POST /user/register -------------->|
   |                  |                    |
   |                  |-- Hash password -->|
   |                  |                    |
   |                  |-- INSERT user ---->|
   |                  |<--- user_id -------|
   |                  |                    |
   |<--- JWT token ---|                    |
```

### Game Data Flow
```
Frontend           Backend              IGDB/PCGaming     Database
   |                  |                      |              |
   |-- GET /search -->|                      |              |
   |                  |-- Search games ----->|              |
   |                  |<--- Game data -------|              |
   |                  |                                     |
   |                  |-- INSERT/UPDATE game -------------->|
   |<--- Games -------|                                     |
```

### Report Submission Flow
```
Frontend           Backend              Claude AI         Database
   |                  |                      |              |
   |-- POST /user/report (JWT) ------------>|              |
   |                  |                      |              |
   |                  |-- Moderate text ---->|              |
   |                  |<--- Violations ------|              |
   |                  |                                     |
   |                  |-- If approved: INSERT report ------>|
   |<--- Success -----|                                     |
```

## Configuration Files

### Root Level

**docker-compose.yml**
- Development environment
- Services: postgres, backend, frontend, db-backup, db-migrate
- Networks: playability-network
- Volumes: postgres_data, backups

**docker-compose.production.yml**
- Production environment
- Additional services: nginx reverse proxy
- Resource limits and replicas
- SSL/TLS configuration

**Makefile**
- Convenient commands for development
- Database operations
- Docker management
- See: `make help`

### Frontend Configuration

**nuxt.config.ts**
- Nuxt modules and plugins
- Tailwind CSS integration
- Runtime configuration
- ShadCN component setup

**tailwind.config.js**
- Custom color palette
- Typography settings
- Animation extensions

**eslint.config.mjs**
- Linting rules
- Uses Nuxt preset

### Backend Configuration

**go.mod**
- Go module definition
- Dependencies with versions

**main.go**
- Server setup and routing
- Middleware mounting
- Database initialization

## Development Workflow

### Starting Development

1. **Clone and setup:**
   ```bash
   git clone <repo>
   cd playability
   ```

2. **Configure environment:**
   ```bash
   cp frontend/.env.example frontend/.env
   cp backend/.env.example backend/.env
   # Edit backend/.env with credentials
   ```

3. **Start with Docker:**
   ```bash
   make dev
   ```

   Or manually:
   ```bash
   docker-compose up -d
   ```

4. **Access services:**
   - Frontend: http://localhost:3000
   - Backend: http://localhost:8080
   - Database: localhost:5432

### Making Changes

**Frontend changes:**
```bash
# Edit files in frontend/
# Hot reload is automatic with Nuxt

# Rebuild if needed
docker-compose up -d --build frontend
```

**Backend changes:**
```bash
# Edit files in backend/
# Rebuild required
docker-compose up -d --build backend
```

**Database changes:**
```bash
# Create migration
./database/scripts/migrate.sh create your_change_name

# Edit migration files in database/migrations/

# Apply migration
make db-migrate
```

## Testing

### Frontend Testing
```bash
# Run tests (configure first)
docker-compose exec frontend pnpm test

# Run linter
docker-compose exec frontend pnpm lint
make lint

# Type checking
docker-compose exec frontend pnpm typecheck
```

### Backend Testing
```bash
# Run tests
docker-compose exec backend go test ./...

# Run linter (requires golangci-lint)
docker-compose exec backend golangci-lint run
```

## Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for complete deployment checklist.

### Quick Production Deploy

1. **Configure production environment:**
   ```bash
   # Set all production environment variables in backend/.env
   # Ensure JWT_SECRET is strong (openssl rand -base64 32)
   # Set DB_SSLMODE=require
   ```

2. **Build and deploy:**
   ```bash
   make prod-up
   ```

3. **Verify deployment:**
   ```bash
   docker-compose -f docker-compose.production.yml ps
   docker-compose -f docker-compose.production.yml logs
   ```

## Useful Commands

### Development
```bash
make dev              # Start all services
make logs             # View all logs
make logs SVC=backend # View specific service
make restart          # Restart all services
make down             # Stop all services
```

### Database
```bash
make db-shell         # PostgreSQL shell
make db-backup        # Create backup
make db-migrate       # Run migrations
make db-migrate-status # Check migration status
```

### Building
```bash
make build            # Rebuild all images
make clean            # Stop and remove volumes
```

## Environment Variables Reference

### Frontend (.env)
```bash
IGDB_CLIENT_SECRET=xxx
NUXT_PUBLIC_API_URL=http://localhost:8080  # Change for production
```

### Backend (.env)
```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=playability_user
DB_PASSWORD=xxx
DB_NAME=playability
DB_SSLMODE=disable  # Use 'require' in production

# Security
JWT_SECRET=xxx  # Use: openssl rand -base64 32

# APIs
IGDB_ACCESS_TOKEN=xxx
CLAUDE_API_KEY=xxx

# Application
PORT=8080
ENVIRONMENT=development
ALLOWED_ORIGINS=http://localhost:3000
FRONTEND_URL=http://localhost:3000
```

## Additional Documentation

- [CLAUDE.md](CLAUDE.md) - Instructions for Claude Code
- [DATABASE.md](DATABASE.md) - Complete database documentation
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment checklist
- [DOCKER.md](DOCKER.md) - Docker usage guide
- [database/README.md](database/README.md) - Database scripts guide
- [README.md](README.md) - Project overview

## Contributing

When adding new features:

1. **Frontend**: Add to `frontend/` directory
2. **Backend**: Add to appropriate `backend/` subdirectory
3. **Database**: Create migration in `database/migrations/`
4. **Documentation**: Update relevant .md files
5. **Docker**: Update docker-compose.yml if adding services

## Questions?

For help with:
- **Frontend issues**: Check nuxt.config.ts, package.json
- **Backend issues**: Check cmd/api/main.go, db/database.go
- **Database issues**: See DATABASE.md, database/README.md
- **Docker issues**: See DOCKER.md, docker-compose.yml
- **Deployment**: See DEPLOYMENT.md

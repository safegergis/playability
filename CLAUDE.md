# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Playability is a full-stack web application focused on game accessibility. It provides comprehensive accessibility information for video games, allowing users to submit reports and rate accessibility features to help gamers with disabilities make informed decisions.

## Architecture

### Project Structure
```
playability/
├── frontend/       # Nuxt 3 frontend application
├── backend/        # Go REST API
└── database/       # Database schema and scripts
```

### Frontend (Nuxt 3 + Vue.js)
**Location:** `frontend/`
- **Framework**: Nuxt 3 with Vue.js for SSR and improved SEO
- **UI Components**: ShadCN Vue component library with Radix Vue primitives
- **Styling**: Tailwind CSS with custom configuration
- **Form Validation**: Vee Validate with Yup/Zod schema validation
- **Icons**: Nuxt Icon module with Lucide Vue Next
- **Carousel**: Embla Carousel Vue for game showcases

### Backend (Go REST API)
**Location:** `backend/`

- **Framework**: Chi router with dependency injection pattern
- **Database**: PostgreSQL with `database/sql` package
- **Authentication**: JWT-based with bcrypt password hashing
- **AI Integration**: Claude AI for content moderation
- **External APIs**: IGDB and PCGamingWiki for game data
- **Port**: Runs on localhost:8080

#### Backend Structure
```
backend/
├── auth/           # Authentication utilities
├── cmd/api/        # Application entry point (main.go)
├── db/             # Database layer (games, users, reports)
├── handlers/       # HTTP request handlers
├── pkg/            # Utility packages
│   ├── ai/         # AI moderation services
│   ├── calc/       # Score calculation utilities
│   └── fetch/      # External API integrations
└── types/          # Type definitions
```

### Database Schema
- **games**: Game information and accessibility features (closed_captions, color_blind, full_controller_support, controller_remapping)
- **users**: User management with authentication (id, username, email, hash, num_reports)
- **reports**: User-submitted accessibility feedback linked to games

### API Architecture

#### Public Endpoints
- `GET /search` - Game search
- `GET /games?id={}` - Game details
- `GET /featured` - Featured games
- `GET /reports/cards/{game}` - Report cards for game
- `GET /reports/features/{game}` - Feature statistics
- `GET /reports/score/{game}` - Accessibility score

#### User Endpoints
- `POST /user/login` - User authentication
- `POST /user/register` - User registration
- `GET /user/{id}` - User profile

#### Protected Endpoints (JWT Required)
- `POST /user/report` - Submit accessibility report

## Development Commands

### Quick Start (Docker)
```bash
# Start all services
make dev

# View logs
make logs

# Stop services
make down
```

### Frontend (pnpm)
```bash
# Navigate to frontend
cd frontend

# Development server
pnpm dev

# Build for production
pnpm build

# Generate static site
pnpm generate

# Preview production build
pnpm preview

# Lint code (ESLint with Nuxt preset)
pnpm lint

# Type checking
pnpm typecheck
```

### Backend (Go)
```bash
# Navigate to backend directory
cd backend

# Download dependencies
go mod download

# Run development server
go run cmd/api/main.go

# Build binary
go build cmd/api/main.go
```

### Database
```bash
# Initialize database
make db-init

# Run migrations
make db-migrate

# Create backup
make db-backup

# Open database shell
make db-shell
```

## Configuration

### Frontend Environment Variables
- `IGDB_CLIENT_SECRET`: Required for IGDB API integration

### Backend Environment Variables
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Database connection
- `JWT_SECRET` - JWT signing key
- `IGDB_ACCESS_TOKEN` - IGDB API authentication
- `CLAUDE_API_KEY` - Anthropic Claude API key

### Key Configuration Files
- `nuxt.config.ts`: Main Nuxt configuration with modules and runtime config
- `tailwind.config.js`: Tailwind CSS customization
- `components.json`: ShadCN component configuration
- `eslint.config.mjs`: ESLint configuration using Nuxt preset
- `backend/.env`: Backend environment variables
- `docker-compose.yml`: Development Docker setup
- `Makefile`: Convenient development commands

## External Dependencies

### Frontend to Backend Communication
The Nuxt frontend communicates with the Go backend through:
- Server API routes in `/server/api/` that proxy requests
- JWT tokens stored as HTTP-only cookies
- Backend running on `localhost:8080`

### Game Data Sources
- **IGDB API**: Primary source for video game information
- **PCGamingWiki API**: Additional accessibility information for PC games

### AI Services
- **Claude AI**: Content moderation for user reports with violation detection

## Key Backend Components

### Authentication & Security
- JWT authentication with HS256 signing
- 24-hour token expiration
- bcrypt password hashing
- AI-powered content moderation
- CORS configuration and input validation

### External Integrations
- **IGDB**: Game search, details, cover art, Steam ID resolution
- **PCGamingWiki**: Accessibility feature detection
- **Claude AI**: Content moderation categories (hate, harassment, violence, etc.)

### Score Calculation
- Average calculation of user-submitted accessibility scores
- Feature consensus determination with statistical analysis
- Support level categorization (true/limited/false)

## Development Notes

- Frontend and backend are separate applications
- Backend uses dependency injection pattern with `Env` struct
- ShadCN components configured with no prefix in `./components/ui`
- All UI components follow accessibility best practices
- Database uses parameterized queries for SQL injection protection
- AI moderation runs on all user-submitted content
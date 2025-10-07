# Playability

> Making video games accessible and inclusive for everyone

[![Live Site](https://img.shields.io/badge/Live-playablty.com-blue?style=for-the-badge)](https://playablty.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Nuxt](https://img.shields.io/badge/Nuxt-4-00DC82?style=for-the-badge&logo=nuxt.js)](https://nuxt.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)

---

## Mission Statement

At Playability, we believe that everyone should have the opportunity to enjoy video games, regardless of their physical or cognitive abilities. We strive to break down barriers in gaming by providing accurate, user-driven accessibility information and advocating for more inclusive game design. Our goal is to create a world where no gamer is left behind due to lack of accessibility features.

**Visit us live at: [playablty.com](https://playablty.com)**

---

## Screenshots

<img width="1509" alt="Home Page" src="https://github.com/user-attachments/assets/53403101-f4bd-42d2-898f-a7ed39e86505">
<img width="1509" alt="Game Details" src="https://github.com/user-attachments/assets/d08bfbbf-42cd-4ab4-b2c9-dab104d785d6">
<img width="1509" alt="Search Results" src="https://github.com/user-attachments/assets/cf2e78bf-6ea1-4869-9c77-dc2d3bc41690">

---

## Features

### User-Facing Features

- **User Registration & Authentication:** Secure JWT-based authentication with bcrypt password hashing
- **Comprehensive Game Database:** Access thousands of games with detailed information sourced from IGDB and PCGamingWiki
- **Detailed Game Pages:** View comprehensive game information including:
  - Cover art and screenshots
  - Supported platforms
  - Release date and developer information
  - Genre and tags
  - Accessibility features with user-submitted ratings
- **User-Driven Accessibility Feedback:**
  - Submit detailed reports on game accessibility features
  - Rate the effectiveness of existing accessibility options
  - Provide comments and tips for other users
- **Accessibility Information:** Track features such as:
  - Closed captions and subtitle options
  - Colorblind modes and visual assistance settings
  - Full controller support and button remapping capabilities
- **Accessibility Scoring System:** Aggregate user feedback to provide overall accessibility scores for games
- **Responsive Design:** Mobile-friendly interface using Tailwind CSS, accessible across all devices

### Technical Features

- **RESTful API:** Well-structured API endpoints for seamless frontend-backend communication
- **JWT Authentication:** Secure, token-based authentication system for protected routes
- **AI Content Moderation:** Automated moderation of user submissions to ensure quality and safety
- **Docker Support:** Complete containerization for easy deployment and development
- **Database Migrations:** Version-controlled schema management
- **WCAG Compliance:** Platform itself follows accessibility best practices

---

## Tech Stack

### Frontend
- **[Nuxt 3](https://nuxt.com/)** - Vue.js framework with SSR for improved SEO
- **[Vue.js 3](https://vuejs.org/)** - Progressive JavaScript framework
- **[Tailwind CSS 4](https://tailwindcss.com/)** - Utility-first CSS framework
- **[ShadCN Vue](https://www.shadcn-vue.com/)** - Accessible UI component library built on Radix Vue
- **[Vee Validate](https://vee-validate.logaretm.com/v4/)** - Form validation with Yup/Zod schemas
- **[Embla Carousel](https://www.embla-carousel.com/)** - Accessible carousel component
- **[Lucide Vue](https://lucide.dev/)** - Icon library
- **[Pinia](https://pinia.vuejs.org/)** - State management for Vue

### Backend
- **[Go 1.23](https://golang.org/)** - High-performance programming language
- **[Chi Router](https://github.com/go-chi/chi)** - Lightweight and flexible HTTP router
- **[PostgreSQL 16](https://www.postgresql.org/)** - Robust relational database
- **[JWT Auth](https://jwt.io/)** - Secure authentication with go-chi/jwtauth
- **[Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Password hashing

### External APIs
- **[IGDB API](https://api-docs.igdb.com/)** - Comprehensive video game database
- **[PCGamingWiki API](https://www.pcgamingwiki.com/wiki/PCGamingWiki:API)** - PC game accessibility information

### DevOps
- **Docker & Docker Compose** - Containerization and orchestration
- **GitHub Actions** - CI/CD pipeline

---

## Project Structure

```
playability/
├── frontend/              # Nuxt 3 frontend application
│   ├── components/       # Vue components (including ShadCN UI)
│   ├── pages/           # Nuxt pages (routes)
│   ├── server/          # Server API routes
│   └── types/           # TypeScript definitions
├── backend/              # Go REST API
│   ├── auth/            # JWT & password hashing
│   ├── cmd/api/         # Application entry point
│   ├── db/              # Database layer
│   ├── handlers/        # HTTP request handlers
│   ├── pkg/             # Utilities (AI, calc, fetch)
│   └── types/           # Go type definitions
├── database/             # Database schema & scripts
│   ├── migrations/      # Database migrations
│   └── scripts/         # Backup, restore, init scripts
├── docker-compose.yml    # Development environment
└── Makefile             # Development commands
```

---

## Quick Start

### Prerequisites
- **Docker** and **Docker Compose** installed
- API keys for IGDB API and AI moderation service

### Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd playability
   ```

2. **Create environment file:**
   ```bash
   cp .env.example .env
   ```

3. **Generate JWT secret:**
   ```bash
   openssl rand -base64 32
   ```

4. **Edit `.env` with your credentials:**
   ```bash
   # Required values:
   DB_PASSWORD=your_secure_password
   JWT_SECRET=<paste generated secret>
   IGDB_ACCESS_TOKEN=your_igdb_token
   IGDB_CLIENT_SECRET=your_igdb_secret

   # Match DB_PASSWORD with POSTGRES_PASSWORD
   POSTGRES_PASSWORD=your_secure_password
   ```

5. **Start the application:**
   ```bash
   docker-compose up -d
   ```

6. **Access the application:**
   - **Frontend:** http://localhost:3000
   - **Backend API:** http://localhost:8080
   - **Database:** localhost:5432

The database initializes automatically on first run.

---

## Development Commands

### Docker Operations

```bash
# Start all services
docker-compose up -d

# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop all services
docker-compose down

# Restart all services
docker-compose restart

# View running containers
docker-compose ps

# Rebuild specific service
docker-compose up -d --build backend

# Stop and remove everything (including data)
docker-compose down -v
```

### Database Operations

```bash
# Open PostgreSQL shell
docker-compose exec postgres psql -U playability_user -d playability

# Create database backup
docker-compose --profile backup run --rm db-backup

# Run migrations
docker-compose run --rm db-migrate /migrate.sh up

# Check migration status
docker-compose run --rm db-migrate /migrate.sh status

# Rollback migration
docker-compose run --rm db-migrate /migrate.sh down
```

### Frontend Development

```bash
cd frontend

# Development server
pnpm dev

# Production build
pnpm build

# Preview production build
pnpm preview

# Lint code
pnpm lint

# Type checking
pnpm typecheck
```

### Backend Development

```bash
cd backend

# Run development server
go run cmd/api/main.go

# Build binary
go build cmd/api/main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

---

## Database

### Schema

The database uses PostgreSQL with three main tables:

- **`users`** - User accounts and authentication
- **`games`** - Game information with accessibility features
- **`reports`** - User-submitted accessibility reports

Custom types:
- **`feature_support`** - ENUM: `'false'`, `'unknown'`, `'limited'`, `'true'`

### Migrations

```bash
# Create new migration
./database/scripts/migrate.sh create migration_name

# Apply migrations
docker-compose run --rm db-migrate /migrate.sh up

# Check status
docker-compose run --rm db-migrate /migrate.sh status

# Rollback last migration
docker-compose run --rm db-migrate /migrate.sh down
```

### Backups

```bash
# Create backup
docker-compose --profile backup run --rm db-backup

# Backups are saved to ./backups/ directory

# Restore from backup
./database/scripts/restore.sh backups/playability_full_YYYYMMDD_HHMMSS.dump.gz
```

---

## API Endpoints

### Public Endpoints
- `GET /search` - Search games
- `GET /games?id={id}` - Get game details
- `GET /featured` - Get featured games
- `GET /reports/cards/{game}` - Get report cards for game
- `GET /reports/features/{game}` - Get feature statistics
- `GET /reports/score/{game}` - Get accessibility score

### User Endpoints
- `POST /user/login` - User authentication
- `POST /user/register` - User registration
- `GET /user/{id}` - Get user profile

### Protected Endpoints (JWT Required)
- `POST /user/report` - Submit accessibility report

---

## Security Features

- **JWT Authentication** with HS256 signing and 24-hour expiration
- **Bcrypt Password Hashing** with default cost factor (10)
- **AI Content Moderation** for all user-submitted reports
- **Parameterized SQL Queries** for injection protection
- **CORS Configuration** with environment-based origins
- **HTTP-Only Cookies** for secure token storage

---

## Testing

### Frontend Tests
```bash
cd frontend
pnpm test          # Run tests
pnpm lint          # Lint code
pnpm typecheck     # Type checking
```

### Backend Tests
```bash
cd backend
go test ./...                    # Run all tests
go test -v ./db/                 # Test database layer
go test -cover ./handlers/       # Test with coverage
```

---

## Contributing

We welcome contributions! Here's how to get started:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- **Frontend:** Follow Vue.js and Nuxt best practices
- **Backend:** Follow Go conventions and use dependency injection
- **Database:** Always create migrations for schema changes
- **Testing:** Write tests for new features
- **Documentation:** Keep code well-documented

---

## Roadmap

- Enhanced search with filters and sorting
- User profiles and contribution tracking
- Game recommendations based on accessibility needs
- Community forums and discussions
- Mobile native apps (iOS/Android)
- API documentation with Swagger/OpenAPI
- Advanced analytics and reporting
- Email notifications for game updates
- Multi-language support
- Accessibility audit tools for developers

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Acknowledgments

- **IGDB** for providing comprehensive game data
- **PCGamingWiki** for accessibility information
- **ShadCN** for the UI component library
- **The gaming accessibility community** for inspiring this project

---

## Contact & Support

- **Live Site:** [playablty.com](https://playablty.com)
- **Issues:** Please use the GitHub issue tracker
- **Discussions:** Use GitHub Discussions for questions and ideas

---

<div align="center">

**Made with care for accessible gaming**

[Visit Playability](https://playablty.com) | [Report Bug](https://github.com/yourusername/playability/issues) | [Request Feature](https://github.com/yourusername/playability/issues)

</div>

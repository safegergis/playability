# Quick Setup Guide

## Environment Configuration

All environment variables are managed from a **single `.env` file** at the project root.

Docker Compose reads variables from `.env` and passes them to containers using the `environment:` section with `${VAR}` syntax.

### Setup Steps

1. **Copy the example file:**
   ```bash
   cp .env.example .env
   ```

2. **Generate a secure JWT secret:**
   ```bash
   openssl rand -base64 32
   ```

3. **Edit `.env` and fill in your values:**
   ```bash
   # Required values:
   DB_PASSWORD=your_secure_password
   JWT_SECRET=<paste the generated secret from step 2>
   IGDB_ACCESS_TOKEN=your_token
   IGDB_CLIENT_SECRET=your_secret
   CLAUDE_API_KEY=your_key
   
   # Also update POSTGRES_PASSWORD to match DB_PASSWORD
   POSTGRES_PASSWORD=your_secure_password
   ```

4. **Start the application:**
   ```bash
   make dev
   ```
   
   Or manually:
   ```bash
   docker-compose up -d
   ```

### What This `.env` File Controls

- ✅ **PostgreSQL Database** - DB credentials and settings
- ✅ **Backend API** - All Go API configuration
- ✅ **Frontend** - Nuxt environment variables
- ✅ **Database Scripts** - Backup, restore, migrations
- ✅ **Docker Services** - All container configurations

### Environment Variables Explained

```bash
# Database - Used by postgres container and backend
DB_HOST=postgres              # Database hostname (use 'postgres' for Docker)
DB_PORT=5432                  # Database port
DB_USER=playability_user      # Database username
DB_PASSWORD=***               # Database password (CHANGE THIS!)
DB_NAME=playability           # Database name
DB_SSLMODE=disable            # SSL mode (use 'require' in production)

# Backend API
PORT=8080                     # API server port
ENVIRONMENT=development       # Environment (development/production)
JWT_SECRET=***                # JWT signing key (GENERATE THIS!)
ALLOWED_ORIGINS=http://localhost:3000  # CORS allowed origins
FRONTEND_URL=http://localhost:3000     # Frontend URL

# Frontend
NUXT_PUBLIC_API_URL=http://localhost:8080  # Backend API URL

# External APIs
IGDB_ACCESS_TOKEN=***         # IGDB API token
IGDB_CLIENT_SECRET=***        # IGDB client secret
CLAUDE_API_KEY=***            # Claude AI API key

# PostgreSQL (Docker specific)
POSTGRES_DB=playability       # Must match DB_NAME
POSTGRES_USER=playability_user # Must match DB_USER
POSTGRES_PASSWORD=***         # Must match DB_PASSWORD
```

### Verification

After setup, verify everything works:

```bash
# Check that .env exists
ls -la .env

# Start services
make dev

# Check logs
make logs

# Verify all services are running
make ps
```

### Troubleshooting

**"Error loading .env file"**
- Make sure `.env` exists in the project root
- Run: `cp .env.example .env`

**"Environment variable not set"**
- Check that all required variables are filled in `.env`
- No quotes needed around values
- No spaces around `=` sign

**"Cannot connect to database"**
- Verify `POSTGRES_PASSWORD` matches `DB_PASSWORD`
- Check `DB_HOST=postgres` (not `localhost` when using Docker)

### Need Help?

See complete documentation:
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Project organization
- [DOCKER.md](DOCKER.md) - Docker usage guide
- [DATABASE.md](DATABASE.md) - Database documentation

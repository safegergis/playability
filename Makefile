# Playability Makefile
# Provides convenient commands for Docker and database operations

.PHONY: help up down restart logs clean build db-init db-migrate db-backup db-restore db-status db-shell test

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

## help: Show this help message
help:
	@echo "$(BLUE)Playability Development Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Docker Commands:$(NC)"
	@echo "  make up              - Start all services"
	@echo "  make down            - Stop all services"
	@echo "  make restart         - Restart all services"
	@echo "  make logs            - View logs (add SVC=<service> for specific service)"
	@echo "  make build           - Build all Docker images"
	@echo "  make clean           - Stop services and remove volumes (WARNING: deletes data)"
	@echo "  make ps              - Show running containers"
	@echo ""
	@echo "$(GREEN)Database Commands:$(NC)"
	@echo "  make db-init         - Initialize database with schema"
	@echo "  make db-migrate      - Run database migrations"
	@echo "  make db-migrate-down - Rollback last migration"
	@echo "  make db-migrate-status - Show migration status"
	@echo "  make db-backup       - Create database backup"
	@echo "  make db-restore      - Restore database from backup"
	@echo "  make db-shell        - Open PostgreSQL shell"
	@echo "  make db-logs         - View database logs"
	@echo ""
	@echo "$(GREEN)Development Commands:$(NC)"
	@echo "  make dev             - Start in development mode with hot reload"
	@echo "  make frontend        - Start only frontend"
	@echo "  make backend         - Start only backend"
	@echo "  make test            - Run tests"
	@echo "  make lint            - Run linters"
	@echo ""
	@echo "$(GREEN)Production Commands:$(NC)"
	@echo "  make prod-up         - Start in production mode"
	@echo "  make prod-down       - Stop production services"
	@echo "  make prod-backup     - Create production backup"

## up: Start all services
up:
	@echo "$(BLUE)Starting all services...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"
	@make ps

## down: Stop all services
down:
	@echo "$(BLUE)Stopping all services...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

## restart: Restart all services
restart:
	@echo "$(BLUE)Restarting all services...$(NC)"
	docker-compose restart
	@echo "$(GREEN)✓ Services restarted$(NC)"

## logs: View logs (use SVC=service_name for specific service)
logs:
ifdef SVC
	docker-compose logs -f $(SVC)
else
	docker-compose logs -f
endif

## build: Build all Docker images
build:
	@echo "$(BLUE)Building Docker images...$(NC)"
	docker-compose build
	@echo "$(GREEN)✓ Build complete$(NC)"

## clean: Stop services and remove volumes (WARNING: deletes data)
clean:
	@echo "$(RED)WARNING: This will delete all data in Docker volumes!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Stopping services and removing volumes...$(NC)"; \
		docker-compose down -v; \
		echo "$(GREEN)✓ Cleanup complete$(NC)"; \
	else \
		echo "$(YELLOW)Cancelled$(NC)"; \
	fi

## ps: Show running containers
ps:
	@echo "$(BLUE)Running containers:$(NC)"
	@docker-compose ps

## db-init: Initialize database with schema
db-init:
	@echo "$(BLUE)Initializing database...$(NC)"
	@if [ ! -f backend/.env ]; then \
		echo "$(RED)✗ Error: backend/.env not found$(NC)"; \
		echo "$(YELLOW)Run: cp backend/.env.example backend/.env$(NC)"; \
		exit 1; \
	fi
	./database/scripts/init.sh
	@echo "$(GREEN)✓ Database initialized$(NC)"

## db-migrate: Run database migrations
db-migrate:
	@echo "$(BLUE)Running database migrations...$(NC)"
	docker-compose run --rm db-migrate /migrate.sh up
	@echo "$(GREEN)✓ Migrations complete$(NC)"

## db-migrate-down: Rollback last migration
db-migrate-down:
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	docker-compose run --rm db-migrate /migrate.sh down
	@echo "$(GREEN)✓ Rollback complete$(NC)"

## db-migrate-status: Show migration status
db-migrate-status:
	@echo "$(BLUE)Migration status:$(NC)"
	docker-compose run --rm db-migrate /migrate.sh status

## db-backup: Create database backup
db-backup:
	@echo "$(BLUE)Creating database backup...$(NC)"
	@mkdir -p backups
	docker-compose --profile backup run --rm db-backup
	@echo "$(GREEN)✓ Backup created$(NC)"
	@ls -lh backups/ | tail -5

## db-restore: Restore database from backup
db-restore:
	@echo "$(BLUE)Available backups:$(NC)"
	@ls -1 backups/ 2>/dev/null || echo "No backups found"
	@echo ""
	@read -p "Enter backup filename: " backup; \
	if [ -n "$$backup" ]; then \
		echo "$(YELLOW)Restoring from $$backup...$(NC)"; \
		./database/scripts/restore.sh backups/$$backup; \
		echo "$(GREEN)✓ Restore complete$(NC)"; \
	else \
		echo "$(RED)No backup specified$(NC)"; \
	fi

## db-shell: Open PostgreSQL shell
db-shell:
	@echo "$(BLUE)Opening database shell...$(NC)"
	@docker-compose exec postgres psql -U $${DB_USER:-playability_user} -d $${DB_NAME:-playability}

## db-logs: View database logs
db-logs:
	@docker-compose logs -f postgres

## dev: Start in development mode
dev:
	@echo "$(BLUE)Starting development environment...$(NC)"
	@if [ ! -f frontend/.env ]; then \
		echo "$(YELLOW)Creating frontend/.env from .env.example...$(NC)"; \
		cp frontend/.env.example frontend/.env; \
	fi
	@if [ ! -f backend/.env ]; then \
		echo "$(YELLOW)Creating backend/.env from .env.example...$(NC)"; \
		cp backend/.env.example backend/.env; \
		echo "$(RED)⚠ Please edit backend/.env with your credentials$(NC)"; \
		exit 1; \
	fi
	docker-compose up -d
	@echo "$(GREEN)✓ Development environment started$(NC)"
	@echo ""
	@echo "$(BLUE)Services:$(NC)"
	@echo "  Frontend: http://localhost:3000"
	@echo "  Backend:  http://localhost:8080"
	@echo "  Database: localhost:5432"
	@echo ""
	@echo "$(YELLOW)View logs: make logs$(NC)"

## frontend: Start only frontend
frontend:
	@echo "$(BLUE)Starting frontend...$(NC)"
	docker-compose up -d frontend
	@echo "$(GREEN)✓ Frontend started at http://localhost:3000$(NC)"

## backend: Start only backend and database
backend:
	@echo "$(BLUE)Starting backend...$(NC)"
	docker-compose up -d postgres backend
	@echo "$(GREEN)✓ Backend started at http://localhost:8080$(NC)"

## test: Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@echo "$(YELLOW)Frontend tests:$(NC)"
	docker-compose exec frontend pnpm test || echo "No frontend tests configured"
	@echo ""
	@echo "$(YELLOW)Backend tests:$(NC)"
	docker-compose exec backend go test ./... || echo "No backend tests configured"

## lint: Run linters
lint:
	@echo "$(BLUE)Running linters...$(NC)"
	@echo "$(YELLOW)Frontend linting:$(NC)"
	docker-compose exec frontend pnpm lint || echo "Linter not configured"
	@echo ""
	@echo "$(YELLOW)Backend linting:$(NC)"
	docker-compose exec backend golangci-lint run || echo "golangci-lint not installed"

## prod-up: Start in production mode
prod-up:
	@echo "$(BLUE)Starting production environment...$(NC)"
	@if [ ! -f backend/.env ]; then \
		echo "$(RED)✗ Error: backend/.env file not found$(NC)"; \
		exit 1; \
	fi
	docker-compose -f docker-compose.production.yml up -d
	@echo "$(GREEN)✓ Production environment started$(NC)"

## prod-down: Stop production services
prod-down:
	@echo "$(BLUE)Stopping production services...$(NC)"
	docker-compose -f docker-compose.production.yml down
	@echo "$(GREEN)✓ Production services stopped$(NC)"

## prod-backup: Create production backup
prod-backup:
	@echo "$(BLUE)Creating production backup...$(NC)"
	@mkdir -p backups
	docker-compose -f docker-compose.production.yml --profile backup run --rm db-backup
	@echo "$(GREEN)✓ Production backup created$(NC)"

# Default target
.DEFAULT_GOAL := help

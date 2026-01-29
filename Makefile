.PHONY: up down restart logs install dev dev-backend dev-frontend help

# Default target
help:
	@echo "Available commands:"
	@echo "  make install       - Install dependencies for all projects"
	@echo "  make up            - Start database (Docker) in background"
	@echo "  make down          - Stop all docker services"
	@echo "  make restart       - Restart all docker services"
	@echo "  make logs          - View database logs"
	@echo "  make dev           - One-shot: start backend + frontend (DB must be running)"
	@echo "  make dev-backend   - Run backend services only"
	@echo "  make dev-frontend  - Run frontend services only"
	@echo "  make clean         - Remove project temporary files and build artifacts"

install:
	@echo "Installing Backend Dependencies..."
	cd src/backend && go mod download
	@echo "Installing Frontend (Teams App) Dependencies..."
	cd src/frontend/teams-app && pnpm install
	@echo "Installing Frontend (Web Dashboard) Dependencies..."
	cd src/frontend/web-dashboard && pnpm install

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

dev:
	@echo "Starting all services (backend + frontend)..."
	@echo "Ensure your database is running and .env is configured."
	(cd src/backend && go run cmd/card-service/main.go) & \
	(cd src/backend && go run cmd/analytics-service/main.go) & \
	(cd src/frontend/teams-app && pnpm dev) & \
	(cd src/frontend/web-dashboard && pnpm dev) & \
	wait

dev-backend:
	@echo "Starting Backend Services..."
	@echo "Ensure your database is running and .env is configured."
	(cd src/backend && go run cmd/card-service/main.go) & \
	(cd src/backend && go run cmd/analytics-service/main.go) & \
	wait

dev-frontend:
	@echo "Starting Frontend Services..."
	(cd src/frontend/teams-app && pnpm dev) & \
	(cd src/frontend/web-dashboard && pnpm dev) & \
	wait

clean:
	rm -rf src/frontend/teams-app/node_modules
	rm -rf src/frontend/teams-app/dist
	rm -rf src/frontend/web-dashboard/node_modules
	rm -rf src/frontend/web-dashboard/dist

# Development Guide

This guide outlines how to set up the local development environment for the Thank You Card application using Docker and Make to streamline the process.

## 📋 Prerequisites

Ensure you have the following installed:

-   **Docker** & **Docker Compose**
-   **Go** (>= 1.21)
-   **Node.js** (>= 22.12.0)
-   **pnpm** (>= 9.x)
-   **Make** (Standard on macOS/Linux)

## 🚀 Quick Start

We provide a `Makefile` to simplify common tasks.

### 1. Install Dependencies
Install dependencies for both backend and frontend:
```bash
make install
```

### 2. Start Infrastructure (Database)
Start PostgreSQL with **automatic schema migration**:
```bash
make up
```
> This starts a Docker container for Postgres exposed on port `5432`.
> **Note**: The database is automatically seeded with tables via the mounted migration scripts.

### 3. Start Application Services

**Option A: One-shot (single terminal)**  
Start backend and frontend together:
```bash
make dev
```
> Runs Card Service (`:8080`), Analytics Service (`:8081`), Teams App (`:5173`), and Web Dashboard (`:3000`). Ensure `.env` is configured and the database is running.

**Option B: Separate terminals**

**Terminal 1: Backend**
```bash
make dev-backend
```
> Runs Card Service (`:8080`) and Analytics Service (`:8081`).

**Terminal 2: Frontend**
```bash
make dev-frontend
```
> Runs Teams App (`:5173`) and Web Dashboard (`:3000`).

---

## 🖥 Local development without Docker

If you use your own local database (no Docker), configure `src/backend/.env` with your connection (e.g. `localhost:5432`), then start the app:

```bash
make dev
```
> Backend and frontend run in one terminal; stop with Ctrl+C.

---

## 🔧 Configuration Details

### Backend Configuration
1.  Navigate to `src/backend`.
2.  Copy the example env file:
    ```bash
    cp .env.example .env
    ```
3.  The defaults in `.env.example` are configured to work **out-of-the-box** with the `make up` database command.
    -   Host: `localhost`
    -   Port: `5432`
    -   User: `thankyoucard`
    -   Password: `password`
    -   DB: `thankyoucard`

### Frontend Configuration
Navigate to `src/frontend/teams-app` and `src/frontend/web-dashboard` and copy their respective `.env.example` to `.env` if needed.
-   **Teams App**: `VITE_API_BASE_URL` should point to `http://localhost:8080`.

## 🗄️ Database Management

-   **View Logs**: `make logs`
-   **Reset Database**: `make down` (Warning: this destroys data if volumes are removed, use `docker-compose down -v` to fully reset).
-   **Manual Connection**:
    ```bash
    psql postgres://thankyoucard:password@localhost:5432/thankyoucard
    ```

## 🛠 Troubleshooting

-   **Port Conflicts**: If port 5432 is in use, modify `docker-compose.yml` to map to a different host port (e.g., `5433:5432`) and update `DB_PORT` in your `.env`.
-   **Migration Errors**: If migrations fail to run, try `docker-compose down -v` to clear the volume and start fresh.

## 📂 Project Structure

-   `src/backend`: Go microservices.
    -   `cmd/`: Entry points.
    -   `internal/`: Business logic.
    -   `scripts/migrations`: SQL migration files.
-   `src/frontend`: React applications.
    -   `teams-app`: MS Teams Tab application.
    -   `web-dashboard`: Admin dashboard.
-   `docker-compose.yml`: Infrastructure definition.
-   `Makefile`: Task automation.

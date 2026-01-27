# Thank You Card Backend

Go + Gin + Gorm + PostgreSQL backend for the internal employee recognition system.

## Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin
- **ORM**: Gorm
- **Database**: PostgreSQL
- **Environment**: godotenv

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entrypoint
├── db/
│   └── migrate.go           # Database migrations and seeding
├── internal/
│   ├── handlers/            # HTTP handlers (controller layer)
│   │   ├── dto/            # Data transfer objects
│   │   └── card_handler.go
│   ├── middleware/          # HTTP middleware
│   │   └── auth.go         # Mock user & HR admin middleware
│   ├── models/              # Gorm models
│   │   └── models.go
│   ├── repositories/        # Data access layer
│   │   └── repositories.go
│   ├── router/              # Route configuration
│   │   └── router.go
│   └── services/            # Business logic layer
│       └── card_service.go
├── script/
│   └── init.sql            # Database initialization script
├── .env.example            # Environment variables template
├── docker-compose.yml      # PostgreSQL Docker setup
├── Makefile                # Common commands
└── README.md              # This file
```

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 12+ (or use Docker Compose)
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env and configure your database connection
# DATABASE_DSN=postgres://postgres:postgres@localhost:5432/thankyou?sslmode=disable
```

### 2. Start Database

**Option A: Using Docker Compose (Recommended)**

```bash
make docker-up
# or
docker-compose up -d
```

**Option B: Using Local PostgreSQL**

```bash
# Create database manually
psql -U postgres -f script/init.sql

# Or use Makefile
make db-init
```

### 3. Install Dependencies

```bash
make deps
# or
go mod download
go mod tidy
```

### 4. Run the Application

```bash
make run
# or
go run ./cmd/server
```

The server will start on `http://localhost:8056` (or the port specified in `.env`).

### 5. Verify Setup

```bash
# Health check
curl http://localhost:8056/api/cards/feed

# Should return empty array if no cards exist yet
```

## Makefile Commands

```bash
make help          # Show all available commands
make build         # Build the application binary
make run           # Run the application
make test          # Run tests
make test-coverage # Run tests with coverage report
make clean         # Clean build artifacts
make deps          # Download and tidy dependencies
make fmt           # Format Go code
make vet           # Run go vet
make lint          # Run fmt + vet

# Database
make db-init       # Initialize database
make db-reset      # Reset database (WARNING: drops database)

# Docker
make docker-up     # Start PostgreSQL container
make docker-down   # Stop PostgreSQL container
make docker-logs   # View container logs

# Development
make dev           # Start database and run application
make setup         # Complete setup for new developers
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_DSN` | PostgreSQL connection string | Required |
| `PORT` | HTTP server port | `8056` |
| `MOCK_USER_ID` | Mock user ID for MVP (no real auth) | `1` |
| `ENV` | Environment (development/production) | `development` |

## API Endpoints

### Cards

- `POST /api/cards` - Create a thank-you card
- `GET /api/cards/feed` - Get company-wide feed
- `GET /api/cards/me/received` - Get personal received cards
- `GET /api/cards/me/sent` - Get personal sent cards
- `GET /api/cards/me/stats` - Get personal statistics (sent/received counts, top values)
- `GET /api/cards/top-recipients` - Get top 10 recognized employees (default: last 30 days)

### Reactions

- `POST /api/cards/:id/reactions` - Add/update emoji reaction
- `DELETE /api/cards/:id/reactions` - Remove emoji reaction

### Teams Bot Endpoints (Simplified)

- `GET /api/teams/feed` - Get simplified feed for Teams Bot
- `GET /api/teams/cards/:id` - Get simplified card detail for Teams Bot

### HR Analytics (Requires HR Admin)

- `GET /api/analytics/dashboard` - Get dashboard overview
- `GET /api/analytics/recognizers` - Get most active recognizers
- `GET /api/analytics/teams` - Get team recognition patterns by department
- `GET /api/analytics/values` - Get company values distribution
- `GET /api/analytics/export` - Export cards data as CSV

### Query Parameters

- `page` - Page number (default: 1)
- `pageSize` - Items per page (default: 20, max: 100)
- `from` - Start date (RFC3339 format)
- `to` - End date (RFC3339 format)
- `q` - Search keyword (searches in reason field)
- `values[]` - Filter by company value IDs
- `senderIds[]` - Filter by sender IDs
- `recipientIds[]` - Filter by recipient IDs

## Example API Usage

### Create a Card

```bash
curl -X POST http://localhost:8056/api/cards \
  -H "Content-Type: application/json" \
  -d '{
    "recipientIds": [2],
    "valueIds": [1, 2],
    "reason": "Thank you for your excellent work on the project!"
  }'
```

### Get Company Feed

```bash
curl http://localhost:8056/api/cards/feed?page=1&pageSize=20
```

### Get Personal Statistics

```bash
curl http://localhost:8056/api/cards/me/stats
```

### Get Top 10 Recognized Employees

```bash
# Last 30 days (default)
curl http://localhost:8056/api/cards/top-recipients

# Custom time range
curl "http://localhost:8056/api/cards/top-recipients?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z&limit=10"
```

### Add Emoji Reaction

```bash
curl -X POST http://localhost:8056/api/cards/1/reactions \
  -H "Content-Type: application/json" \
  -d '{"emojiCode": "👍"}'
```

### Teams Bot - Get Feed

```bash
curl http://localhost:8056/api/teams/feed?page=1&pageSize=20
```

### HR Analytics - Get Dashboard (Requires HR Admin)

```bash
# Set MOCK_USER_ID=2 in .env to use HR admin account
curl "http://localhost:8056/api/analytics/dashboard?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z"
```

### HR Analytics - Export Cards as CSV

```bash
curl "http://localhost:8056/api/analytics/export?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z" -o cards_export.csv
```

## Database Schema

The application uses Gorm AutoMigrate to create tables. Key models:

- **employees** - Employee information
- **company_values** - Company values and credos
- **cards** - Thank-you cards
- **card_recipients** - Many-to-many: cards to recipients
- **card_values** - Many-to-many: cards to company values
- **emoji_reactions** - User reactions on cards

## Seeding Data

On first run, the application automatically:

1. Creates all database tables via AutoMigrate
2. Seeds company values/credos from user stories
3. Seeds mock employees:
   - `employee@example.com` (ID: 1, regular employee)
   - `hradmin@example.com` (ID: 2, HR admin)

## Mock User Authentication

For MVP, the application uses mock user authentication:

- Set `MOCK_USER_ID` in `.env` to switch between users
- The middleware attaches the user to the request context
- No real authentication is performed

To test as HR admin, set `MOCK_USER_ID=2` in `.env`.

## Development Workflow

1. Start database: `make docker-up`
2. Run application: `make run`
3. Make changes and restart
4. Run tests: `make test`
5. Format code: `make fmt`

## Troubleshooting

### Database Connection Error

```bash
# Check if PostgreSQL is running
docker-compose ps

# Check connection string in .env
# Ensure DATABASE_DSN is correct
```

### Port Already in Use

```bash
# Change PORT in .env or kill the process using port 8056
lsof -ti:8056 | xargs kill -9
```

### Migration Errors

```bash
# Reset database (WARNING: deletes all data)
make db-reset

# Or manually drop and recreate
psql -U postgres -c "DROP DATABASE IF EXISTS thankyou;"
make db-init
```

## Implementation Status

### Completed Features

- ✅ Card creation and management
- ✅ Company-wide feed
- ✅ Personal card views (received/sent)
- ✅ Personal statistics
- ✅ Top 10 recognized employees
- ✅ Emoji reactions
- ✅ HR Analytics dashboard
- ✅ Most active recognizers
- ✅ Team recognition patterns
- ✅ Company values distribution
- ✅ CSV export functionality
- ✅ Teams Bot simplified endpoints

### Future Enhancements

- [ ] Real authentication (Azure AD / Teams integration)
- [ ] Teams notification service (HTTP push to Teams channel)
- [ ] Audit logging
- [ ] Rate limiting
- [ ] Request validation middleware
- [ ] Structured logging (zap/logrus)
- [ ] Card filtering and search UI integration
- [ ] Real-time notifications via WebSocket

## License

Internal use only.

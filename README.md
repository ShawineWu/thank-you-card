# Thank You Card - Employee Recognition System

Internal employee recognition system that enables employees to send public digital thank-you cards to colleagues, displays them in a company-wide feed, and provides HR analytics capabilities.

## Tech Stack

### Backend
- Go 1.22+
- Gin (HTTP framework)
- Gorm (ORM)
- PostgreSQL

### Frontend
- React 18 + TypeScript
- Vite
- Tailwind CSS
- shadcn/ui
- React Query
- Recharts

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL (running locally)
- Make

### Database Setup

The project expects PostgreSQL to be running with:
- **Database**: `thankyoucard`
- **User**: `postgres`
- **Host**: `localhost`
- **Port**: `5432`
- **Password**: `123456`

### One-Command Setup

```bash
# Complete setup (installs dependencies, initializes DB)
make setup
```

### Development

**Terminal 1 - Backend:**
```bash
make backend-dev
```

**Terminal 2 - Frontend:**
```bash
make frontend-dev
```

Then open http://localhost:3000 in your browser.

**Note:** Backend runs on port 8056, frontend on port 3000.

## Makefile Commands

### Quick Start
- `make setup` - Complete setup for new developers
- `make start` - Initialize database and show instructions
- `make dev-full` - Start backend (frontend needs separate terminal)

### Backend
- `make backend-dev` - Start backend development server
- `make backend-build` - Build backend binary
- `make backend-run` - Build and run backend

### Frontend
- `make frontend-install` - Install frontend dependencies
- `make frontend-dev` - Start frontend development server
- `make frontend-build` - Build frontend for production

### Database
- `make db-init` - Initialize database (create if not exists)
- `make db-reset` - Reset database (WARNING: drops database)

### Utilities
- `make clean` - Clean build artifacts
- `make help` - Show all available commands

## Project Structure

```
.
├── backend/              # Go backend
│   ├── cmd/server/      # Application entrypoint
│   ├── db/              # Database migrations and seeding
│   ├── internal/        # Internal packages
│   │   ├── handlers/    # HTTP handlers
│   │   ├── models/      # Data models
│   │   ├── repositories/# Data access layer
│   │   ├── services/    # Business logic
│   │   └── router/      # Route configuration
│   └── script/          # Database scripts
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API services
│   │   └── types/       # TypeScript types
│   └── public/
└── Makefile            # Build and run commands
```

## Features

### For All Employees
- Create and send thank-you cards
- View company-wide feed
- View personal received/sent cards
- Personal statistics dashboard
- Top 10 recognized employees

### For HR Administrators
- Analytics dashboard
- Most active recognizers
- Team recognition patterns
- Company values distribution
- CSV export functionality

## API Endpoints

### Cards
- `POST /api/cards` - Create a card
- `GET /api/cards/feed` - Get company feed
- `GET /api/cards/me/received` - Get received cards
- `GET /api/cards/me/sent` - Get sent cards
- `GET /api/cards/me/stats` - Get personal stats
- `GET /api/cards/top-recipients` - Get top 10 employees

### Reactions
- `POST /api/cards/:id/reactions` - Add reaction
- `DELETE /api/cards/:id/reactions` - Remove reaction

### HR Analytics (requires HR admin)
- `GET /api/analytics/dashboard` - Dashboard overview
- `GET /api/analytics/recognizers` - Most active recognizers
- `GET /api/analytics/teams` - Team patterns
- `GET /api/analytics/values` - Values distribution
- `GET /api/analytics/export` - Export CSV

## Environment Variables

### Backend (.env)
```env
DATABASE_DSN=postgres://postgres:123456@localhost:5432/thankyoucard?sslmode=disable
PORT=8056
MOCK_USER_ID=1
```

### Frontend (.env)
```env
VITE_API_BASE_URL=http://localhost:8056/api
VITE_MOCK_USER_ID=1
```

## Database

The backend automatically:
1. Creates all tables via Gorm AutoMigrate
2. Seeds company values/credos
3. Seeds mock employees (ID: 1 = regular, ID: 2 = HR admin)

## Testing HR Features

To test HR admin features, set `MOCK_USER_ID=2` in backend `.env` file.

## License

Internal use only.

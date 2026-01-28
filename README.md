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

### Install Backend Dependencies

Before starting, install required Go packages:

```bash
cd backend
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

### Development

**Option 1: Start both services together (Recommended)**
```bash
# Automatically cleans ports and starts both backend and frontend
make dev-all
```

**Option 2: Start services separately**

**Terminal 1 - Backend:**
```bash
make backend-dev
```

**Terminal 2 - Frontend:**
```bash
make frontend-dev
```

Then open http://localhost:3000 in your browser.

**Note:** 
- Backend runs on port 8056, frontend on port 3000
- All start commands automatically clean ports before starting to prevent conflicts

## Makefile Commands

### Quick Start
- `make setup` - Complete setup for new developers
- `make start` - Initialize database and show instructions
- `make dev-all` - Start both backend and frontend (same terminal, recommended)
- `make dev-full` - Start backend (frontend needs separate terminal)

### Service Management
- `make stop` - Stop all services (backend and frontend)
- `make stop-backend` - Stop backend server only
- `make stop-frontend` - Stop frontend server only

### Backend
- `make backend-dev` - Start backend development server (auto-cleans port)
- `make backend-build` - Build backend binary
- `make backend-run` - Build and run backend

### Frontend
- `make frontend-install` - Install frontend dependencies
- `make frontend-dev` - Start frontend development server (auto-cleans port)
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

## Authentication

The system uses JWT-based authentication. All API endpoints (except login) require authentication.

### Login

**Endpoint:** `POST /api/auth/login`

**Request:**
```json
{
  "username": "hr",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "username": "hr",
    "role": "HR",
    "employee": {
      "id": 2,
      "name": "Mock HR Admin",
      "email": "hradmin@example.com",
      "department": "HR"
    }
  }
}
```

**Using the Token:**
Include the token in the `Authorization` header for all subsequent requests:
```
Authorization: Bearer <token>
```

### Test Accounts

The system automatically creates the following test accounts (all passwords: `password123`):

- **employee** / password123 - Employee role (can access all basic features)
- **hr** / password123 - HR role (can access all features + HR Analytics)
- **admin** / password123 - Admin role (can access all features + HR Analytics)

### Role Permissions

- **EMPLOYEE**: Feed, Create Card, My Cards, Stats, Top 10, Company Values
- **HR**: All employee features + HR Analytics
- **ADMIN**: All employee features + HR Analytics

## API Endpoints

### Authentication
- `POST /api/auth/login` - Login and get JWT token

### Company Values
- `GET /api/company-values` - Get all company values (supports `?type=VALUE` or `?type=CREDO`)
- `GET /api/company-values/:id` - Get company value details

### Cards (requires authentication)
- `POST /api/cards` - Create a card
- `GET /api/cards/feed` - Get company feed
- `GET /api/cards/me/received` - Get received cards
- `GET /api/cards/me/sent` - Get sent cards
- `GET /api/cards/me/stats` - Get personal stats
- `GET /api/cards/top-recipients` - Get top 10 employees

### Reactions (requires authentication)
- `POST /api/cards/:id/reactions` - Add reaction
- `DELETE /api/cards/:id/reactions` - Remove reaction

### HR Analytics (requires HR or ADMIN role)
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
JWT_SECRET=your-secret-key-change-in-production  # Optional, has default for dev
```

**Note:** `MOCK_USER_ID` is no longer used. Authentication is now handled via login.

### Frontend (.env)
```env
VITE_API_BASE_URL=http://localhost:8056/api
```

**Note:** `VITE_MOCK_USER_ID` is no longer used. Users must log in.

## Database

The backend automatically:
1. Creates all tables via Gorm AutoMigrate (including `users` table)
2. Seeds company values/credos (5 values + 10 credos)
3. Seeds mock employees (ID: 1 = regular, ID: 2 = HR admin)
4. Seeds test user accounts (employee, hr, admin)

## Usage Examples

### Starting the Application

```bash
# Start both services (recommended)
make dev-all

# Or start separately
make backend-dev    # Terminal 1
make frontend-dev   # Terminal 2
```

### Stopping Services

```bash
# Stop all services
make stop

# Stop individual services
make stop-backend
make stop-frontend
```

### Testing Different Roles

1. **As Employee:**
   - Login with: `employee` / `password123`
   - Can access: Feed, Create Card, My Cards, Stats, Top 10, Values
   - Cannot access: HR Analytics

2. **As HR:**
   - Login with: `hr` / `password123`
   - Can access: All features including HR Analytics

3. **As Admin:**
   - Login with: `admin` / `password123`
   - Can access: All features including HR Analytics

### Creating a Card

1. Login to the application
2. Navigate to "Create Card"
3. Select recipients (employees)
4. Select company values/credos
5. Write a reason/message
6. Submit

### Viewing Analytics (HR/Admin only)

1. Login with HR or Admin account
2. Click "Analytics" in the navigation bar
3. Explore:
   - Dashboard: Overview statistics
   - Recognizers: Most active recognizers
   - Teams: Team recognition patterns
   - Values: Company values distribution
   - Export: Download data as CSV

## License

Internal use only.

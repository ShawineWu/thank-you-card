# Thank You Card System

An internal employee recognition system that enables sending thank-you cards aligned with company values, with public display and HR analytics capabilities.

## Project Structure

```
├── backend/                 # Go backend service
│   ├── cmd/server/         # Application entry point
│   ├── internal/           # Internal application code
│   │   ├── domain/         # Domain models and business logic
│   │   ├── application/    # Application services
│   │   ├── infrastructure/ # External integrations
│   │   └── presentation/   # HTTP handlers and routes
│   ├── pkg/               # Shared packages
│   ├── migrations/        # Database migration files
│   └── scripts/           # Build and deployment scripts
├── frontend/              # React frontend application
│   ├── src/
│   │   ├── components/    # Reusable UI components
│   │   ├── pages/         # Page components
│   │   ├── stores/        # Zustand state management
│   │   ├── services/      # API and external services
│   │   ├── hooks/         # Custom React hooks
│   │   ├── utils/         # Utility functions
│   │   └── types/         # TypeScript type definitions
│   └── tests/             # Frontend tests
├── construction/          # Design documents
├── inception/             # Requirements and user stories
└── docs/                  # Additional documentation
```

## Technology Stack

### Backend
- **Language**: Go 1.21
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM
- **Architecture**: Domain-Driven Design with Clean Architecture

### Frontend
- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite (Rolldown)
- **State Management**: Zustand
- **HTTP Client**: Axios
- **Testing**: Jest with Testing Library

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 13 or higher

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Copy environment configuration:
   ```bash
   cp .env.example .env
   ```

3. Update database configuration in `.env`

4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

The backend will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Copy environment configuration:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

4. Start the development server:
   ```bash
   npm run dev
   ```

The frontend will be available at `http://localhost:5173`

### Database Setup

1. Create a PostgreSQL database named `thank_you_card`
2. Update the database configuration in `backend/.env`
3. The application will automatically run migrations on startup

## API Endpoints

### Health Check
- `GET /health` - Basic health check

### Cards (Planned)
- `POST /api/v1/cards` - Create a new card
- `GET /api/v1/cards` - Get company-wide card feed
- `GET /api/v1/cards/{id}` - Get specific card
- `GET /api/v1/cards/received` - Get received cards
- `GET /api/v1/cards/sent` - Get sent cards

### Analytics (Planned)
- `GET /api/v1/analytics/dashboard` - Analytics dashboard

## Development Status

### ✅ Completed (Phase 1)
- [x] Project structure and monorepo setup
- [x] Backend development environment (Go, Gin, GORM, PostgreSQL)
- [x] Frontend development environment (React, Vite, TypeScript, Zustand)
- [x] Database schema and migrations
- [x] Basic API structure and health checks
- [x] Company values and credos seeding

### 🚧 In Progress (Phase 2)
- [ ] Card Service domain layer implementation
- [ ] Card Service infrastructure layer
- [ ] Card Service application layer
- [ ] Card Service presentation layer (REST API)

### 📋 Planned
- [ ] Analytics Service implementation
- [ ] Frontend UI components and pages
- [ ] Event-driven integration with Teams
- [ ] Authentication and authorization
- [ ] Testing and quality assurance
- [ ] Deployment and monitoring

## Company Values & Credos

The system includes 15 predefined company values and credos:

**Values:**
- Make an Impact
- Strive for Excellence
- Stand Together
- Be Open-Minded
- Stay Grounded

**Credos:**
- Bias for Action
- Customer Centric
- Think Strategically
- Deep Dive
- Invent and Simplify
- Earn Trust
- Take Ownership
- Challenge, Disagree, and Commit
- Learn and Be Curious
- Do More with Less

## Contributing

1. Follow the established project structure
2. Use conventional commit messages
3. Ensure all tests pass before submitting
4. Update documentation as needed

## License

Internal company project - All rights reserved.
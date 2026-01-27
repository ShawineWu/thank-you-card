# Thank You Card Project Introduction

## Overview
Internal Microsoft Teams application for employee recognition and appreciation.

## Project Structure

```
src/
├── frontend/
│   ├── teams-app/          # Microsoft Teams App
│   └── web-dashboard/      # Web Management Dashboard
└── backend/                # Backend Services (Go)
    ├── cmd/                # Service Entry Points
    ├── internal/           # Business Logic
    └── pkg/                # Reusable Packages
```

## Tech Stack

### Frontend
- **Teams App**: React + TypeScript + Fluent UI + Vite
- **Web Dashboard**: React + TypeScript + Zustand + Vite

### Backend
- **Language**: Go
- **Framework**: Gin
- **ORM**: Gorm
- **Database**: PostgreSQL

### Authentication
- **Teams App**: Azure AD SSO
- **Web Dashboard**: Hydra OAuth2 (reusing existing client ID)
- **User Management**: PLAT UMS

## Getting Started

See documentation in `/docs`:
- `deployment.md` - Deployment guide
- `user-guide.md` - User manual
- Authentication design in brain artifacts

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Microsoft Teams Developer Account

### Local Setup
```bash
# Backend
cd src/backend
go mod init thank-you-card
go mod tidy
go run cmd/card-service/main.go

# Frontend - Teams App
cd src/frontend/teams-app
npm install
npm run dev

# Frontend - Web Dashboard
cd src/frontend/web-dashboard
npm install
npm run dev
```

## Documentation

- Inception documents: `/inception/units/`
- Construction design: `/construction/`
- Implementation plan: See brain artifacts
- API documentation: `/docs/api.md`

## License
Internal use only - Castlery

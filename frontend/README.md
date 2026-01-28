# Thank You Card Frontend

React + TypeScript + Vite frontend for the Thank You Card employee recognition system.

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **shadcn/ui** for UI components
- **React Router** for routing
- **TanStack Query** for data fetching and caching
- **Axios** for HTTP requests
- **Recharts** for data visualization
- **date-fns** for date formatting

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The app will be available at `http://localhost:8099`

## Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_MOCK_USER_ID=1
```

## Features

### For All Employees

- **Feed Page**: View company-wide thank you card feed
- **Create Card**: Send thank you cards to colleagues
- **My Cards**: View received and sent cards
- **My Stats**: Personal statistics dashboard
- **Top 10 Employees**: See most recognized employees

### For HR Administrators

- **Analytics Dashboard**: Overview of recognition activity
- **Most Active Recognizers**: Employees who send the most cards
- **Team Patterns**: Department-level recognition statistics
- **Values Distribution**: Analysis of company values recognition
- **CSV Export**: Export card data for analysis

## Project Structure

```
frontend/
├── src/
│   ├── components/       # Reusable components
│   │   ├── ui/          # shadcn/ui components
│   │   ├── Card/        # Card-related components
│   │   ├── Forms/       # Form components
│   │   └── Filters/     # Filter components
│   ├── pages/           # Page components
│   │   ├── Feed/
│   │   ├── CreateCard/
│   │   ├── MyCards/
│   │   ├── Stats/
│   │   ├── TopEmployees/
│   │   └── Analytics/
│   ├── services/        # API service layer
│   ├── types/           # TypeScript type definitions
│   ├── utils/           # Utility functions
│   ├── lib/             # Library configurations
│   └── App.tsx          # Main app component
├── public/
└── package.json
```

## API Integration

The frontend communicates with the backend API at the base URL specified in `VITE_API_BASE_URL`. All API calls are handled through service functions in `src/services/`.

## Development

### Adding New Components

Use shadcn/ui CLI to add new components:

```bash
npx shadcn-ui@latest add [component-name]
```

### Code Style

- Use TypeScript for type safety
- Follow React best practices
- Use Tailwind CSS for styling
- Keep components small and focused

## License

Internal use only.

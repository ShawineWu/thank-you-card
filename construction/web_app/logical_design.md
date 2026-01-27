# Web App Unit - Logical Design

## Document Information
- **Application**: Web App Unit (Presentation Layer)
- **Architecture**: Single Page Application (SPA), Component-Based
- **Technology Stack**: React, Zustand, Vite, Jest
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## Application Architecture Overview

The Web App is designed as a modern single-page application (SPA) that serves as the primary user interface for the employee recognition system. It follows a component-based architecture with clear separation between presentation logic and business logic, consuming backend APIs for all data operations.

### Architecture Principles
- **Presentation Layer Only**: No business logic, pure UI concerns
- **Component-Based Design**: Reusable, composable UI components
- **State Management**: Centralized state with Zustand
- **API-First**: All data operations via backend APIs
- **Responsive Design**: Mobile-first, responsive user interface
- **Performance Optimized**: Efficient rendering and loading

---

## Frontend Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    Web App Architecture                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    UI Layer                             │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │    Pages    │ │ Components  │ │   Layouts   │       │   │
│  │  │             │ │             │ │             │       │   │
│  │  │ • Dashboard │ │ • CardForm  │ │ • MainLayout│       │   │
│  │  │ • CreateCard│ │ • CardList  │ │ • AuthLayout│       │   │
│  │  │ • MyCards   │ │ • CardItem  │ │ • HRLayout  │       │   │
│  │  │ • Analytics │ │ • Filters   │ │             │       │   │
│  │  │ • Top10     │ │ • Search    │ │             │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                State Management                          │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │   Zustand   │ │    Stores   │ │   Actions   │       │   │
│  │  │   Store     │ │             │ │             │       │   │
│  │  │             │ │ • Auth      │ │ • API Calls │       │   │
│  │  │ • Global    │ │ • Cards     │ │ • State     │       │   │
│  │  │   State     │ │ • User      │ │   Updates   │       │   │
│  │  │ • Actions   │ │ • UI        │ │ • Error     │       │   │
│  │  │ • Selectors │ │ • Analytics │ │   Handling  │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                API Integration                           │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │ API Clients │ │ HTTP Client │ │    Auth     │       │   │
│  │  │             │ │             │ │  Interceptor│       │   │
│  │  │ • Card API  │ │ • Axios     │ │             │       │   │
│  │  │ • Analytics │ │ • Request   │ │ • JWT Token │       │   │
│  │  │   API       │ │   Config    │ │ • Refresh   │       │   │
│  │  │ • Employee  │ │ • Response  │ │ • Logout    │       │   │
│  │  │   API       │ │   Handling  │ │             │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │    Error    │ │   Retry     │ │   Loading   │       │   │
│  │  │  Handling   │ │   Logic     │ │    States   │       │   │
│  │  │             │ │             │ │             │       │   │
│  │  │ • Network   │ │ • Exponential│ │ • Spinners  │       │   │
│  │  │   Errors    │ │   Backoff   │ │ • Skeletons │       │   │
│  │  │ • API       │ │ • Circuit   │ │ • Progress  │       │   │
│  │  │   Errors    │ │   Breaker   │ │   Bars      │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                 Utilities & Helpers                     │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │ Validation  │ │ Formatting  │ │   Routing   │       │   │
│  │  │             │ │             │ │             │       │   │
│  │  │ • Form      │ │ • Dates     │ │ • React     │       │   │
│  │  │   Validation│ │ • Numbers   │ │   Router    │       │   │
│  │  │ • Input     │ │ • Currency  │ │ • Protected │       │   │
│  │  │   Sanitize  │ │ • Text      │ │   Routes    │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │   Hooks     │ │   Context   │ │    Utils    │       │   │
│  │  │             │ │             │ │             │       │   │
│  │  │ • useAPI    │ │ • Theme     │ │ • Constants │       │   │
│  │  │ • useAuth   │ │ • Auth      │ │ • Helpers   │       │   │
│  │  │ • useDebounce│ │ • Toast     │ │ • Types     │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP API Calls
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Backend Services                           │
│  ┌─────────────────┐         ┌─────────────────┐               │
│  │  Card Service   │         │ Analytics       │               │
│  │     APIs        │         │ Service APIs    │               │
│  └─────────────────┘         └─────────────────┘               │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Design and Organization

### Page Components (Route-Level)

#### Dashboard Page
**Purpose**: Main landing page with overview and quick actions

**Components**:
- **DashboardOverview**: Welcome message and quick stats
- **QuickActions**: Create card, view cards buttons
- **RecentActivity**: Recent cards received/sent
- **TopRecognized**: Mini top 10 preview

**Responsibilities**:
- Display user dashboard overview
- Provide navigation to main features
- Show personalized content
- Quick access to common actions

#### Create Card Page
**Purpose**: Multi-step card creation workflow

**Components**:
- **CardCreationWizard**: Step-by-step card creation
- **RecipientSelector**: Employee search and selection
- **ReasonInput**: Recognition reason text input
- **ValuesSelector**: Company values/credos selection
- **CardPreview**: Preview before submission

**Responsibilities**:
- Guide user through card creation process
- Validate input at each step
- Provide clear feedback and guidance
- Handle form submission and errors

#### My Cards Pages
**Purpose**: Personal card views (received/sent)

**Components**:
- **CardTabs**: Switch between received/sent views
- **CardList**: Paginated list of cards
- **CardFilters**: Filter and search options
- **PersonalStats**: Personal statistics display

**Responsibilities**:
- Display user's personal cards
- Provide filtering and search capabilities
- Show personal statistics
- Handle pagination and loading states

#### Analytics Dashboard (HR Only)
**Purpose**: HR analytics and reporting interface

**Components**:
- **AnalyticsDashboard**: Main analytics overview
- **TeamAnalytics**: Team recognition patterns
- **ValuesDistribution**: Values usage charts
- **ExportTools**: Data export functionality

**Responsibilities**:
- Display comprehensive analytics
- Provide interactive charts and graphs
- Enable data export operations
- Restrict access to HR administrators

### Reusable Components

#### Card Components
**CardItem**:
- Display individual card information
- Handle emoji reactions (future)
- Show sender, recipients, reason, values
- Responsive card layout

**CardList**:
- Display paginated list of cards
- Handle loading and empty states
- Support different card layouts
- Infinite scroll or pagination

**CardForm**:
- Reusable card creation form
- Multi-step form handling
- Validation and error display
- Form state management

#### Input Components
**EmployeeSearch**:
- Autocomplete employee search
- Multi-select capability
- Display employee details
- Handle search debouncing

**ValuesSelector**:
- Company values/credos selection
- Visual selection interface
- Enforce selection limits (1-3)
- Display value descriptions

**TextInput**:
- Enhanced text input with validation
- Character count display
- Error state handling
- Accessibility support

#### UI Components
**Button**:
- Consistent button styling
- Loading states
- Different variants (primary, secondary, etc.)
- Icon support

**Modal**:
- Reusable modal component
- Backdrop click handling
- Keyboard navigation
- Focus management

**Toast**:
- Success/error notifications
- Auto-dismiss functionality
- Queue management
- Positioning options

### Layout Components

#### MainLayout
**Purpose**: Primary application layout for authenticated users

**Components**:
- **Header**: Navigation, user menu, notifications
- **Sidebar**: Main navigation menu
- **Content**: Page content area
- **Footer**: Application footer

**Responsibilities**:
- Provide consistent layout structure
- Handle responsive design breakpoints
- Manage navigation state
- Display user information

#### AuthLayout
**Purpose**: Layout for authentication pages

**Components**:
- **LoginForm**: Azure AD/Teams login integration
- **LoadingSpinner**: Authentication loading state
- **ErrorMessage**: Authentication error display

**Responsibilities**:
- Handle authentication flow
- Redirect after successful login
- Display authentication errors
- Manage loading states

---

## State Management Architecture

### Zustand Store Design

#### Global Store Structure
```javascript
// Main store combining all slices
const useStore = create((set, get) => ({
  ...authSlice(set, get),
  ...cardsSlice(set, get),
  ...userSlice(set, get),
  ...uiSlice(set, get),
  ...analyticsSlice(set, get),
}));
```

#### Auth Slice
**State**:
- `user`: Current user information
- `token`: JWT authentication token
- `isAuthenticated`: Authentication status
- `isLoading`: Authentication loading state

**Actions**:
- `login()`: Handle user login
- `logout()`: Handle user logout
- `refreshToken()`: Refresh authentication token
- `checkAuth()`: Verify authentication status

#### Cards Slice
**State**:
- `cards`: Card data arrays (received, sent, feed)
- `currentCard`: Currently selected card
- `filters`: Active filters and search terms
- `pagination`: Pagination state
- `loading`: Loading states for different operations

**Actions**:
- `fetchCards()`: Fetch cards from API
- `createCard()`: Create new card
- `setFilters()`: Update filter state
- `clearCards()`: Clear card data
- `updatePagination()`: Update pagination state

#### UI Slice
**State**:
- `sidebarOpen`: Sidebar visibility state
- `theme`: Current theme (light/dark)
- `notifications`: Toast notifications
- `modals`: Modal visibility state

**Actions**:
- `toggleSidebar()`: Toggle sidebar visibility
- `setTheme()`: Change application theme
- `showNotification()`: Display toast notification
- `openModal()`: Open modal dialog

#### Analytics Slice (HR Only)
**State**:
- `dashboardData`: Analytics dashboard data
- `teamData`: Team analytics data
- `valuesData`: Values distribution data
- `exportStatus`: Export operation status

**Actions**:
- `fetchAnalytics()`: Fetch analytics data
- `exportData()`: Initiate data export
- `setTimePeriod()`: Update analytics time period
- `clearAnalytics()`: Clear analytics data

### State Management Patterns

#### Async State Handling
```javascript
// Pattern for handling async operations
const fetchCards = async (filters) => {
  set({ loading: { ...get().loading, cards: true } });
  try {
    const response = await cardAPI.getCards(filters);
    set({ 
      cards: response.data,
      pagination: response.pagination,
      loading: { ...get().loading, cards: false }
    });
  } catch (error) {
    set({ 
      error: error.message,
      loading: { ...get().loading, cards: false }
    });
  }
};
```

#### Optimistic Updates
```javascript
// Pattern for optimistic UI updates
const createCard = async (cardData) => {
  const tempCard = { ...cardData, id: 'temp-' + Date.now() };
  
  // Optimistically add card to UI
  set({ cards: [tempCard, ...get().cards] });
  
  try {
    const response = await cardAPI.createCard(cardData);
    // Replace temp card with real card
    set({ 
      cards: get().cards.map(card => 
        card.id === tempCard.id ? response.data : card
      )
    });
  } catch (error) {
    // Remove temp card on error
    set({ 
      cards: get().cards.filter(card => card.id !== tempCard.id),
      error: error.message
    });
  }
};
```

---

## API Integration Architecture

### HTTP Client Configuration

#### Axios Setup
```javascript
// Base API client configuration
const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});
```

#### Request Interceptor
```javascript
// Add authentication token to requests
apiClient.interceptors.request.use(
  (config) => {
    const token = useStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);
```

#### Response Interceptor
```javascript
// Handle common response patterns
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Handle token expiration
      useStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### API Client Services

#### Card API Service
```javascript
export const cardAPI = {
  // Get cards with filters and pagination
  getCards: (filters = {}, pagination = {}) => 
    apiClient.get('/api/v1/cards', { params: { ...filters, ...pagination } }),
  
  // Create new card
  createCard: (cardData) => 
    apiClient.post('/api/v1/cards', cardData),
  
  // Get personal received cards
  getReceivedCards: (pagination = {}) => 
    apiClient.get('/api/v1/cards/received', { params: pagination }),
  
  // Get personal sent cards
  getSentCards: (pagination = {}) => 
    apiClient.get('/api/v1/cards/sent', { params: pagination }),
  
  // Search cards
  searchCards: (query, filters = {}) => 
    apiClient.get('/api/v1/cards/search', { params: { query, ...filters } }),
  
  // Filter cards
  filterCards: (filters) => 
    apiClient.get('/api/v1/cards/filter', { params: filters }),
};
```

#### Analytics API Service (HR Only)
```javascript
export const analyticsAPI = {
  // Get dashboard analytics
  getDashboard: (timePeriod = 30) => 
    apiClient.get('/api/v1/analytics/dashboard', { params: { timePeriod } }),
  
  // Get top recognizers
  getTopRecognizers: (timePeriod = 30, limit = 10) => 
    apiClient.get('/api/v1/analytics/recognizers/top', { params: { timePeriod, limit } }),
  
  // Get team analytics
  getTeamAnalytics: (timePeriod = 30) => 
    apiClient.get('/api/v1/analytics/teams', { params: { timePeriod } }),
  
  // Get values distribution
  getValuesDistribution: (timePeriod = 30) => 
    apiClient.get('/api/v1/analytics/values/distribution', { params: { timePeriod } }),
  
  // Export data
  exportData: (exportRequest) => 
    apiClient.post('/api/v1/analytics/export', exportRequest),
};
```

### Error Handling Patterns

#### API Error Handling
```javascript
// Centralized error handling utility
export const handleAPIError = (error) => {
  if (error.response) {
    // Server responded with error status
    const { status, data } = error.response;
    switch (status) {
      case 400:
        return { message: data.message || 'Invalid request', type: 'validation' };
      case 401:
        return { message: 'Authentication required', type: 'auth' };
      case 403:
        return { message: 'Access denied', type: 'permission' };
      case 404:
        return { message: 'Resource not found', type: 'notfound' };
      case 500:
        return { message: 'Server error. Please try again.', type: 'server' };
      default:
        return { message: 'An unexpected error occurred', type: 'unknown' };
    }
  } else if (error.request) {
    // Network error
    return { message: 'Network error. Please check your connection.', type: 'network' };
  } else {
    // Other error
    return { message: error.message || 'An error occurred', type: 'unknown' };
  }
};
```

#### Retry Logic
```javascript
// Retry failed requests with exponential backoff
export const retryRequest = async (requestFn, maxRetries = 3) => {
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await requestFn();
    } catch (error) {
      if (attempt === maxRetries || error.response?.status < 500) {
        throw error;
      }
      
      // Wait before retry (exponential backoff)
      const delay = Math.pow(2, attempt) * 1000;
      await new Promise(resolve => setTimeout(resolve, delay));
    }
  }
};
```

---

## User Experience Flow Diagrams

### Card Creation Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   User      │    │  Web App    │    │ Card Service│    │  Employee   │
│             │    │             │    │     API     │    │ Directory   │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                  │                  │
       │ Click "Create Card"                 │                  │
       ├─────────────────►│                  │                  │
       │                  │                  │                  │
       │                  │ Load Values/Credos                  │
       │                  ├─────────────────►│                  │
       │                  │                  │                  │
       │ Search Recipients│                  │                  │
       ├─────────────────►│                  │                  │
       │                  │                  │                  │
       │                  │ Search Employees │                  │
       │                  ├─────────────────►│                  │
       │                  │                  │                  │
       │                  │                  │ Query Directory  │
       │                  │                  ├─────────────────►│
       │                  │                  │                  │
       │                  │                  │ Employee Results │
       │                  │                  │◄─────────────────┤
       │                  │                  │                  │
       │                  │ Search Results   │                  │
       │                  │◄─────────────────┤                  │
       │                  │                  │                  │
       │ Employee List    │                  │                  │
       │◄─────────────────┤                  │                  │
       │                  │                  │                  │
       │ Select Recipients│                  │                  │
       │ Enter Reason     │                  │                  │
       │ Select Values    │                  │                  │
       ├─────────────────►│                  │                  │
       │                  │                  │                  │
       │                  │ Validate & Preview                  │
       │                  │                  │                  │
       │ Preview Card     │                  │                  │
       │◄─────────────────┤                  │                  │
       │                  │                  │                  │
       │ Submit Card      │                  │                  │
       ├─────────────────►│                  │                  │
       │                  │                  │                  │
       │                  │ Create Card      │                  │
       │                  ├─────────────────►│                  │
       │                  │                  │                  │
       │                  │ Success Response │                  │
       │                  │◄─────────────────┤                  │
       │                  │                  │                  │
       │ Success Message  │                  │                  │
       │◄─────────────────┤                  │                  │
       │                  │                  │                  │
```

### Analytics Dashboard Flow (HR Only)

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  HR Admin   │    │  Web App    │    │ Analytics   │
│             │    │             │    │ Service API │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                  │
       │ Access Analytics │                  │
       ├─────────────────►│                  │
       │                  │                  │
       │                  │ Check HR Role    │
       │                  │ (JWT Token)      │
       │                  │                  │
       │                  │ Load Dashboard   │
       │                  ├─────────────────►│
       │                  │                  │
       │                  │ Analytics Data   │
       │                  │◄─────────────────┤
       │                  │                  │
       │ Dashboard View   │                  │
       │◄─────────────────┤                  │
       │                  │                  │
       │ Change Time Period                  │
       ├─────────────────►│                  │
       │                  │                  │
       │                  │ Refresh Data     │
       │                  ├─────────────────►│
       │                  │                  │
       │                  │ Updated Data     │
       │                  │◄─────────────────┤
       │                  │                  │
       │ Updated Charts   │                  │
       │◄─────────────────┤                  │
       │                  │                  │
       │ Export Data      │                  │
       ├─────────────────►│                  │
       │                  │                  │
       │                  │ Generate Export  │
       │                  ├─────────────────►│
       │                  │                  │
       │                  │ Download URL     │
       │                  │◄─────────────────┤
       │                  │                  │
       │ Download Link    │                  │
       │◄─────────────────┤                  │
       │                  │                  │
```

---

## Responsive Design Architecture

### Breakpoint Strategy
```css
/* Mobile-first responsive breakpoints */
:root {
  --breakpoint-sm: 640px;   /* Small devices */
  --breakpoint-md: 768px;   /* Medium devices */
  --breakpoint-lg: 1024px;  /* Large devices */
  --breakpoint-xl: 1280px;  /* Extra large devices */
}
```

### Component Responsiveness

#### Layout Adaptations
**Mobile (< 768px)**:
- Single column layout
- Collapsible sidebar navigation
- Touch-optimized buttons and inputs
- Simplified card display

**Tablet (768px - 1024px)**:
- Two-column layout where appropriate
- Sidebar overlay on smaller tablets
- Optimized touch targets
- Condensed navigation

**Desktop (> 1024px)**:
- Multi-column layouts
- Persistent sidebar navigation
- Hover states and interactions
- Full feature set display

#### Component-Level Responsiveness
```javascript
// Responsive card list component
const CardList = ({ cards }) => {
  const [isMobile] = useMediaQuery('(max-width: 768px)');
  
  return (
    <div className={`card-grid ${isMobile ? 'mobile' : 'desktop'}`}>
      {cards.map(card => (
        <CardItem 
          key={card.id} 
          card={card} 
          compact={isMobile}
        />
      ))}
    </div>
  );
};
```

### Performance Optimization Strategies

#### Code Splitting
```javascript
// Route-based code splitting
const Dashboard = lazy(() => import('./pages/Dashboard'));
const CreateCard = lazy(() => import('./pages/CreateCard'));
const Analytics = lazy(() => import('./pages/Analytics'));

// Component-based code splitting
const HeavyChart = lazy(() => import('./components/HeavyChart'));
```

#### Image Optimization
```javascript
// Responsive image component
const ResponsiveImage = ({ src, alt, sizes }) => (
  <img
    src={src}
    alt={alt}
    sizes={sizes}
    loading="lazy"
    decoding="async"
  />
);
```

#### Virtual Scrolling
```javascript
// Virtual scrolling for large card lists
const VirtualCardList = ({ cards }) => {
  const [visibleRange, setVisibleRange] = useState({ start: 0, end: 20 });
  
  return (
    <VirtualList
      items={cards}
      itemHeight={120}
      visibleRange={visibleRange}
      onRangeChange={setVisibleRange}
      renderItem={({ item }) => <CardItem card={item} />}
    />
  );
};
```

---

## Security Implementation in Frontend

### Authentication Flow

#### Azure AD/Teams Integration
```javascript
// Authentication service
export const authService = {
  // Initialize authentication
  initialize: async () => {
    const token = localStorage.getItem('auth_token');
    if (token && !isTokenExpired(token)) {
      return { success: true, token };
    }
    return { success: false };
  },
  
  // Handle login redirect from Teams/Azure AD
  handleLoginCallback: async (code) => {
    try {
      const response = await apiClient.post('/auth/callback', { code });
      const { token, user } = response.data;
      
      localStorage.setItem('auth_token', token);
      useStore.getState().setAuth({ token, user });
      
      return { success: true };
    } catch (error) {
      return { success: false, error: error.message };
    }
  },
  
  // Logout user
  logout: () => {
    localStorage.removeItem('auth_token');
    useStore.getState().clearAuth();
    window.location.href = '/login';
  }
};
```

#### Protected Routes
```javascript
// Route protection component
const ProtectedRoute = ({ children, requiredRole }) => {
  const { isAuthenticated, user } = useStore();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  
  if (requiredRole && !user.roles.includes(requiredRole)) {
    return <Navigate to="/unauthorized" replace />;
  }
  
  return children;
};

// Usage in routing
<Route path="/analytics" element={
  <ProtectedRoute requiredRole="hr-admin">
    <Analytics />
  </ProtectedRoute>
} />
```

### Input Validation and Sanitization

#### Form Validation
```javascript
// Card creation form validation
const validateCardForm = (formData) => {
  const errors = {};
  
  // Recipients validation
  if (!formData.recipients || formData.recipients.length === 0) {
    errors.recipients = 'At least one recipient is required';
  }
  
  // Reason validation
  if (!formData.reason || formData.reason.trim().length < 10) {
    errors.reason = 'Recognition reason must be at least 10 characters';
  }
  
  if (formData.reason && formData.reason.length > 1000) {
    errors.reason = 'Recognition reason cannot exceed 1000 characters';
  }
  
  // Values validation
  if (!formData.values || formData.values.length === 0) {
    errors.values = 'At least one company value must be selected';
  }
  
  if (formData.values && formData.values.length > 3) {
    errors.values = 'Maximum 3 company values can be selected';
  }
  
  return { isValid: Object.keys(errors).length === 0, errors };
};
```

#### Input Sanitization
```javascript
// Sanitize user input
export const sanitizeInput = (input) => {
  if (typeof input !== 'string') return input;
  
  return input
    .trim()
    .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '') // Remove scripts
    .replace(/javascript:/gi, '') // Remove javascript: URLs
    .replace(/on\w+\s*=/gi, ''); // Remove event handlers
};
```

### Data Protection

#### Sensitive Data Handling
```javascript
// Secure storage utility
export const secureStorage = {
  // Store sensitive data with encryption (if needed)
  setItem: (key, value) => {
    try {
      const serialized = JSON.stringify(value);
      localStorage.setItem(key, serialized);
    } catch (error) {
      console.error('Failed to store data:', error);
    }
  },
  
  // Retrieve sensitive data
  getItem: (key) => {
    try {
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : null;
    } catch (error) {
      console.error('Failed to retrieve data:', error);
      return null;
    }
  },
  
  // Remove sensitive data
  removeItem: (key) => {
    localStorage.removeItem(key);
  }
};
```

#### Content Security Policy
```html
<!-- CSP headers for security -->
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'self'; 
               script-src 'self' 'unsafe-inline'; 
               style-src 'self' 'unsafe-inline'; 
               img-src 'self' data: https:; 
               connect-src 'self' https://api.company.com;">
```

---

## Performance Optimization Strategies

### Bundle Optimization

#### Webpack/Vite Configuration
```javascript
// Vite configuration for optimization
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          router: ['react-router-dom'],
          state: ['zustand'],
          ui: ['@headlessui/react', '@heroicons/react']
        }
      }
    },
    chunkSizeWarningLimit: 1000
  },
  optimizeDeps: {
    include: ['react', 'react-dom', 'zustand']
  }
});
```

#### Tree Shaking
```javascript
// Import only what's needed
import { debounce } from 'lodash/debounce'; // Instead of entire lodash
import { format } from 'date-fns/format'; // Instead of entire date-fns
```

### Runtime Performance

#### Memoization
```javascript
// Memoize expensive calculations
const ExpensiveComponent = memo(({ data }) => {
  const processedData = useMemo(() => {
    return data.map(item => ({
      ...item,
      calculated: expensiveCalculation(item)
    }));
  }, [data]);
  
  return <div>{/* Render processed data */}</div>;
});
```

#### Debouncing
```javascript
// Debounce search input
const useDebounce = (value, delay) => {
  const [debouncedValue, setDebouncedValue] = useState(value);
  
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);
    
    return () => clearTimeout(handler);
  }, [value, delay]);
  
  return debouncedValue;
};
```

#### Lazy Loading
```javascript
// Lazy load components and images
const LazyImage = ({ src, alt, ...props }) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [isInView, setIsInView] = useState(false);
  const imgRef = useRef();
  
  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true);
          observer.disconnect();
        }
      },
      { threshold: 0.1 }
    );
    
    if (imgRef.current) {
      observer.observe(imgRef.current);
    }
    
    return () => observer.disconnect();
  }, []);
  
  return (
    <div ref={imgRef} {...props}>
      {isInView && (
        <img
          src={src}
          alt={alt}
          onLoad={() => setIsLoaded(true)}
          style={{ opacity: isLoaded ? 1 : 0 }}
        />
      )}
    </div>
  );
};
```

---

## Testing Strategy

### Unit Testing

#### Component Testing
```javascript
// Example component test
describe('CardItem Component', () => {
  const mockCard = {
    id: '1',
    senderName: 'John Doe',
    recipients: [{ name: 'Jane Smith' }],
    reason: 'Great teamwork',
    values: [{ name: 'Make an Impact' }],
    createdAt: '2026-01-27T10:30:00Z'
  };
  
  test('renders card information correctly', () => {
    render(<CardItem card={mockCard} />);
    
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    expect(screen.getByText('Great teamwork')).toBeInTheDocument();
    expect(screen.getByText('Make an Impact')).toBeInTheDocument();
  });
  
  test('handles missing data gracefully', () => {
    const incompleteCard = { ...mockCard, recipients: [] };
    render(<CardItem card={incompleteCard} />);
    
    expect(screen.getByText('No recipients')).toBeInTheDocument();
  });
});
```

#### Hook Testing
```javascript
// Example hook test
describe('useAPI Hook', () => {
  test('handles successful API call', async () => {
    const mockData = { data: [{ id: 1, name: 'Test' }] };
    jest.spyOn(cardAPI, 'getCards').mockResolvedValue(mockData);
    
    const { result } = renderHook(() => useAPI(cardAPI.getCards));
    
    act(() => {
      result.current.execute();
    });
    
    await waitFor(() => {
      expect(result.current.data).toEqual(mockData.data);
      expect(result.current.loading).toBe(false);
    });
  });
  
  test('handles API error', async () => {
    const mockError = new Error('API Error');
    jest.spyOn(cardAPI, 'getCards').mockRejectedValue(mockError);
    
    const { result } = renderHook(() => useAPI(cardAPI.getCards));
    
    act(() => {
      result.current.execute();
    });
    
    await waitFor(() => {
      expect(result.current.error).toBe('API Error');
      expect(result.current.loading).toBe(false);
    });
  });
});
```

### Integration Testing

#### API Integration Tests
```javascript
// Example API integration test
describe('Card Creation Integration', () => {
  test('creates card successfully', async () => {
    const cardData = {
      recipients: ['emp123'],
      reason: 'Excellent work on the project',
      values: ['val123', 'val456']
    };
    
    render(<CreateCard />);
    
    // Fill form
    fireEvent.change(screen.getByLabelText('Recognition Reason'), {
      target: { value: cardData.reason }
    });
    
    // Submit form
    fireEvent.click(screen.getByText('Create Card'));
    
    // Wait for success message
    await waitFor(() => {
      expect(screen.getByText('Card created successfully!')).toBeInTheDocument();
    });
  });
});
```

### End-to-End Testing

#### User Workflow Tests
```javascript
// Example E2E test with Cypress
describe('Card Creation Workflow', () => {
  beforeEach(() => {
    cy.login('employee@company.com');
    cy.visit('/create-card');
  });
  
  it('creates a card successfully', () => {
    // Search and select recipient
    cy.get('[data-testid="recipient-search"]').type('John Doe');
    cy.get('[data-testid="recipient-option"]').first().click();
    
    // Enter recognition reason
    cy.get('[data-testid="reason-input"]').type('Great teamwork on the project');
    
    // Select company values
    cy.get('[data-testid="value-Make an Impact"]').click();
    cy.get('[data-testid="value-Bias for Action"]').click();
    
    // Preview and submit
    cy.get('[data-testid="preview-button"]').click();
    cy.get('[data-testid="submit-button"]').click();
    
    // Verify success
    cy.get('[data-testid="success-message"]').should('be.visible');
    cy.url().should('include', '/dashboard');
  });
});
```

---

## Implementation Roadmap

### Phase 1: Core Infrastructure (Weeks 1-2)
- **Project Setup**: Vite, React, TypeScript configuration
- **Authentication**: Azure AD/Teams integration
- **Routing**: React Router setup with protected routes
- **State Management**: Zustand store configuration
- **API Client**: Axios setup with interceptors

### Phase 2: Core Features (Weeks 3-4)
- **Dashboard**: Main dashboard with overview
- **Card Creation**: Complete card creation workflow
- **Card Display**: Card list and individual card components
- **Personal Views**: My received/sent cards pages

### Phase 3: Advanced Features (Weeks 5-6)
- **Search and Filtering**: Advanced card search and filtering
- **Personal Statistics**: Personal stats dashboard
- **Top 10 List**: Top recognized employees display
- **Responsive Design**: Mobile and tablet optimization

### Phase 4: HR Analytics (Weeks 7-8)
- **Analytics Dashboard**: HR analytics interface
- **Charts and Visualizations**: Interactive charts and graphs
- **Data Export**: CSV export functionality
- **Role-Based Access**: HR-only access controls

### Phase 5: Polish and Testing (Weeks 9-10)
- **Performance Optimization**: Bundle optimization, lazy loading
- **Testing**: Comprehensive unit and integration tests
- **Accessibility**: WCAG compliance and accessibility improvements
- **Documentation**: User guides and technical documentation

---

## Success Criteria and Metrics

### Functional Requirements
- ✅ All user stories implemented with intuitive UI
- ✅ Responsive design working on all device sizes
- ✅ Complete integration with backend APIs
- ✅ Role-based access control for HR features

### Performance Requirements
- **Load Time**: Initial page load < 3 seconds
- **Interaction Response**: UI interactions < 100ms
- **Bundle Size**: Main bundle < 500KB gzipped
- **Lighthouse Score**: > 90 for Performance, Accessibility, Best Practices

### User Experience Requirements
- **Usability**: Intuitive navigation and workflows
- **Accessibility**: WCAG 2.1 AA compliance
- **Mobile Experience**: Fully functional on mobile devices
- **Error Handling**: Clear error messages and recovery options

### Technical Requirements
- **Browser Support**: Modern browsers (Chrome, Firefox, Safari, Edge)
- **Security**: Secure authentication and data handling
- **Maintainability**: Clean, well-documented code
- **Testability**: Comprehensive test coverage (>80%)

This logical design provides a comprehensive blueprint for implementing the Web App Unit as a modern, performant, and user-friendly single-page application that serves as the primary interface for the employee recognition system.
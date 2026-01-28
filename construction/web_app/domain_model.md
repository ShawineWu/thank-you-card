# Web App Unit - Domain Model

## Document Information
- **Unit Type**: Presentation Layer
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## Domain Model Analysis

### Conclusion: No Domain Model Required

The **Web App Unit** is purely a **presentation layer** that consumes backend APIs and does not contain domain logic. In Domain Driven Design, frontend applications typically do not have their own domain models when they serve as thin clients.

---

## Architectural Role

### Presentation Layer Responsibilities
- **User Interface**: Provides web-based UI for all user interactions
- **API Integration**: Consumes REST APIs from Card Service and Analytics Service
- **Client-Side Validation**: Basic input validation for user experience
- **State Management**: Manages UI state and user session
- **Routing**: Handles navigation between different views
- **Authentication Integration**: Integrates with Azure AD/Teams SSO

### What Web App Does NOT Do
- **Business Logic**: All business rules handled by backend services
- **Data Persistence**: No direct database access
- **Domain Events**: Does not publish or handle domain events
- **Business Calculations**: All calculations performed by backend services
- **Data Validation**: Business validation handled by backend APIs

---

## Integration Pattern

### API Consumption Model
The Web App follows a **thin client** pattern, consuming APIs from:

#### Card Service Unit APIs
- Card creation and management
- Personal card views (sent/received)
- Personal statistics
- Employee search
- Card filtering and search
- Emoji reactions (note: actual reactions handled by Teams)

#### Analytics Service Unit APIs (HR Admin Only)
- HR analytics dashboard
- Data export functionality
- Team recognition patterns
- Values distribution analysis
- Most active recognizers

### Data Flow Pattern
```
User Input → Web App → Backend API → Domain Logic → Response → Web App → UI Update
```

---

## Frontend Architecture Components

### UI Components (Not Domain Objects)
- **Card Display Components**: Visual representation of cards
- **Form Components**: Input forms for card creation
- **Navigation Components**: Menu and routing components
- **Chart Components**: Data visualization for analytics
- **Filter Components**: UI for filtering and search

### State Management (Not Domain State)
- **UI State**: Form states, loading states, error states
- **Session State**: User authentication, preferences
- **Cache State**: Temporary storage of API responses
- **Navigation State**: Current route, breadcrumbs

### Client-Side Models (DTOs, Not Domain Models)
These are data transfer objects that mirror API responses, not domain models:

```
// These are DTOs, not domain entities
interface CardDTO {
  id: string;
  senderName: string;
  recipients: RecipientDTO[];
  reason: string;
  values: ValueDTO[];
  createdAt: string;
}

interface EmployeeDTO {
  id: string;
  name: string;
  email: string;
  department: string;
}
```

---

## Validation Strategy

### Client-Side Validation (UX Only)
- **Input Format Validation**: Email format, required fields
- **Length Validation**: Character limits for text inputs
- **Selection Validation**: Minimum/maximum selections
- **Real-time Feedback**: Immediate user feedback

**Important**: All client-side validation is for user experience only. Business validation occurs in backend services.

### Server-Side Validation (Authoritative)
- **Business Rules**: Enforced by Card Service domain model
- **Data Integrity**: Ensured by backend services
- **Security Validation**: Authentication and authorization
- **Cross-System Validation**: Employee existence, value validity

---

## Security Considerations

### Authentication
- **Azure AD/Teams SSO**: Integrated authentication
- **Token Management**: JWT token handling
- **Session Management**: Secure session handling

### Authorization
- **Role-Based Access**: HR admin vs regular employee
- **API Security**: All API calls include authentication tokens
- **Route Protection**: Protected routes for HR analytics

### Data Security
- **No Sensitive Data Storage**: No persistent storage of sensitive data
- **HTTPS Only**: All communication over secure channels
- **XSS Protection**: Input sanitization and output encoding

---

## Why No Domain Model?

### DDD Principles Applied
1. **Domain Logic Location**: All domain logic resides in backend services
2. **Bounded Context Separation**: Web App is not a bounded context
3. **Thin Client Pattern**: Frontend focuses on presentation only
4. **API-First Design**: Backend services expose complete APIs

### Benefits of This Approach
- **Separation of Concerns**: Clear separation between presentation and domain
- **Scalability**: Backend services can be scaled independently
- **Technology Flexibility**: Frontend technology can change without affecting domain
- **Testability**: Domain logic tested in backend, UI logic tested separately

### Alternative Patterns Not Used
- **Rich Client**: Would require domain logic in frontend
- **CQRS with Frontend**: Would require separate read models
- **Micro-frontends**: Single frontend application sufficient for this scope

---

## Integration Contracts

### API Dependencies
The Web App depends on well-defined API contracts from:

#### Card Service APIs
- All card management operations
- Employee search functionality
- Personal statistics calculations
- Filtering and search capabilities

#### Analytics Service APIs
- HR dashboard data
- Export functionality
- Team analytics
- Values distribution

### Error Handling
- **API Error Mapping**: Maps backend errors to user-friendly messages
- **Retry Logic**: Handles temporary network failures
- **Fallback UI**: Graceful degradation when services unavailable

---

## Development Approach

### Frontend-Specific Concerns
- **Responsive Design**: Mobile and desktop compatibility
- **Accessibility**: WCAG compliance for inclusive design
- **Performance**: Optimized loading and rendering
- **User Experience**: Intuitive navigation and interactions

### Testing Strategy
- **Unit Tests**: Component logic and utility functions
- **Integration Tests**: API integration and data flow
- **E2E Tests**: Complete user workflows
- **Visual Tests**: UI consistency and responsive design

### Technology Considerations
- **Framework Choice**: React/Angular/Vue (implementation detail)
- **State Management**: Redux/MobX/Context (implementation detail)
- **Build Tools**: Webpack/Vite (implementation detail)
- **Testing Tools**: Jest/Cypress (implementation detail)

---

## Summary

The Web App Unit serves as a **pure presentation layer** without domain logic, following DDD principles by:

1. **Delegating Domain Logic**: All business logic handled by backend services
2. **Consuming APIs**: Acts as a client to well-defined backend APIs
3. **Focusing on UX**: Concentrates on user experience and interface design
4. **Maintaining Separation**: Clear separation between presentation and domain concerns

This approach ensures:
- **Maintainability**: Changes to domain logic don't affect frontend
- **Scalability**: Backend and frontend can scale independently
- **Testability**: Domain and presentation logic tested separately
- **Flexibility**: Frontend technology choices independent of domain design

**No domain model is required or recommended for this unit.**
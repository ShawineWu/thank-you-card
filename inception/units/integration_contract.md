# Integration Contract

## Document Information
- **Feature**: Internal Employee Thank You Card System
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## Overview

This document defines the integration contracts between the three main units of the Thank You Card system:
1. **Card Service Unit** - Backend service for card management
2. **Web App Unit** - Frontend web application
3. **Analytics Service Unit** - Backend service for HR analytics

Additionally, this document covers integration with external systems:
- **Teams Channel** - Microsoft Teams integration
- **Employee Directory Service** - External employee data source

---

## Architecture Overview

```
┌─────────────────┐
│  Teams Channel  │
│   (External)    │
└────────┬────────┘
         │ HTTP API
         │
┌────────▼────────────────────────────────────────┐
│                                                  │
│  ┌──────────────┐         ┌─────────────────┐  │
│  │   Web App    │◄────────┤  Card Service   │  │
│  │     Unit     │  HTTP   │      Unit       │  │
│  └──────┬───────┘  API    └────────┬────────┘  │
│         │                           │           │
│         │                           │           │
│         │  HTTP API                 │ Read      │
│         │                           │ Access    │
│         │                           │           │
│  ┌──────▼───────────────────────────▼────────┐  │
│  │        Analytics Service Unit             │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
└──────────────────────────────────────────────────┘
         │
         │ HTTP API
         │
┌────────▼────────┐
│   Employee      │
│   Directory     │
│   (External)    │
└─────────────────┘
```

---

## Unit 1: Card Service Unit

### Exposed API Endpoints

#### 1.1 Card Management

##### POST /api/cards
**Description**: Create a new thank-you card

**Request Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "recipientIds": ["string"],
  "reason": "string",
  "valueIds": ["string"]
}
```

**Validation Rules**:
- `recipientIds`: Required, array of 1+ valid employee IDs
- `reason`: Required, non-empty string, max 1000 characters
- `valueIds`: Required, array of 1-3 valid value/credo IDs

**Response** (201 Created):
```json
{
  "id": "uuid",
  "senderId": "string",
  "senderName": "string",
  "senderEmail": "string",
  "senderDepartment": "string",
  "recipients": [
    {
      "id": "string",
      "name": "string",
      "email": "string",
      "department": "string"
    }
  ],
  "reason": "string",
  "values": [
    {
      "id": "string",
      "name": "string",
      "type": "value | credo"
    }
  ],
  "createdAt": "2026-01-27T10:30:00Z",
  "reactions": []
}
```

**Error Responses**:
- 400 Bad Request: Invalid input data
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: Invalid recipient or value ID
- 500 Internal Server Error

---

##### GET /api/cards
**Description**: Get company-wide card feed with pagination

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page`: integer, default 1
- `pageSize`: integer, default 20, max 100
- `sortOrder`: "desc" | "asc", default "desc"

**Response** (200 OK):
```json
{
  "cards": [
    {
      "id": "uuid",
      "senderId": "string",
      "senderName": "string",
      "senderEmail": "string",
      "senderDepartment": "string",
      "recipients": [...],
      "reason": "string",
      "values": [...],
      "createdAt": "2026-01-27T10:30:00Z",
      "reactions": [
        {
          "emoji": "👍",
          "userId": "string",
          "userName": "string",
          "timestamp": "2026-01-27T11:00:00Z"
        }
      ]
    }
  ],
  "pagination": {
    "currentPage": 1,
    "pageSize": 20,
    "totalPages": 15,
    "totalCards": 300
  }
}
```

---

##### GET /api/cards/{id}
**Description**: Get specific card details by ID

**Request Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id`: string (UUID)

**Response** (200 OK):
```json
{
  "id": "uuid",
  "senderId": "string",
  "senderName": "string",
  "senderEmail": "string",
  "senderDepartment": "string",
  "recipients": [...],
  "reason": "string",
  "values": [...],
  "createdAt": "2026-01-27T10:30:00Z",
  "reactions": [...]
}
```

**Error Responses**:
- 404 Not Found: Card ID does not exist

---

##### GET /api/cards/received
**Description**: Get cards received by authenticated user

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page`: integer, default 1
- `pageSize`: integer, default 20, max 100

**Response** (200 OK):
```json
{
  "cards": [...],
  "pagination": {
    "currentPage": 1,
    "pageSize": 20,
    "totalPages": 3,
    "totalCards": 45
  },
  "totalReceived": 45
}
```

---

##### GET /api/cards/sent
**Description**: Get cards sent by authenticated user

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page`: integer, default 1
- `pageSize`: integer, default 20, max 100

**Response** (200 OK):
```json
{
  "cards": [...],
  "pagination": {
    "currentPage": 1,
    "pageSize": 20,
    "totalPages": 2,
    "totalCards": 32
  },
  "totalSent": 32
}
```

---

#### 1.2 Employee Search

##### GET /api/employees/search
**Description**: Search employees by name, email, or department

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `query`: string, required, min 2 characters
- `limit`: integer, default 10, max 50

**Response** (200 OK):
```json
{
  "employees": [
    {
      "id": "string",
      "name": "string",
      "email": "string",
      "department": "string",
      "profilePicture": "string (URL)",
      "isActive": true
    }
  ],
  "count": 5
}
```

---

#### 1.3 Values & Credos

##### GET /api/values
**Description**: Get all company values and credos

**Request Headers**:
```
Authorization: Bearer {token}
```

**Response** (200 OK):
```json
{
  "values": [
    {
      "id": "string",
      "name": "Make an Impact",
      "description": "Be driven by the desire to build something that can touch millions of lives.",
      "type": "value"
    },
    {
      "id": "string",
      "name": "Strive for Excellence",
      "description": "Today's great is not good enough for tomorrow.",
      "type": "value"
    }
  ],
  "credos": [
    {
      "id": "string",
      "name": "Bias for Action",
      "description": "Speed matters, take action to deliver a high quality result...",
      "type": "credo"
    }
  ]
}
```

---

#### 1.4 Emoji Reactions

##### POST /api/cards/{id}/reactions
**Description**: Add emoji reaction to a card

**Request Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id`: string (UUID) - Card ID

**Request Body**:
```json
{
  "emoji": "👍"
}
```

**Validation Rules**:
- User can only have one emoji reaction per card
- If user already has a reaction, it will be replaced

**Response** (201 Created):
```json
{
  "cardId": "uuid",
  "reaction": {
    "emoji": "👍",
    "userId": "string",
    "userName": "string",
    "timestamp": "2026-01-27T11:00:00Z"
  }
}
```

**Error Responses**:
- 400 Bad Request: Invalid emoji
- 404 Not Found: Card not found

---

##### DELETE /api/cards/{id}/reactions
**Description**: Remove user's emoji reaction from a card

**Request Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id`: string (UUID) - Card ID

**Response** (204 No Content)

**Error Responses**:
- 404 Not Found: Card not found or no reaction exists

---

#### 1.5 Statistics

##### GET /api/statistics/personal
**Description**: Get personal card statistics for authenticated user

**Request Headers**:
```
Authorization: Bearer {token}
```

**Response** (200 OK):
```json
{
  "userId": "string",
  "cardsSent": 32,
  "cardsReceived": 45,
  "mostUsedValues": [
    {
      "id": "string",
      "name": "Make an Impact",
      "count": 12,
      "percentage": 37.5
    }
  ],
  "mostReceivedValues": [
    {
      "id": "string",
      "name": "Bias for Action",
      "count": 18,
      "percentage": 40.0
    }
  ]
}
```

---

##### GET /api/statistics/top10
**Description**: Get top 10 recognized employees

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `timePeriod`: integer, days (default 30)

**Response** (200 OK):
```json
{
  "employees": [
    {
      "rank": 1,
      "employeeId": "string",
      "employeeName": "string",
      "department": "string",
      "cardCount": 28
    }
  ],
  "period": {
    "startDate": "2025-12-28",
    "endDate": "2026-01-27",
    "days": 30
  }
}
```

---

#### 1.6 Filtering & Search

##### GET /api/cards/filter
**Description**: Filter cards by multiple criteria

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `valueIds`: string[], comma-separated value/credo IDs
- `senderIds`: string[], comma-separated sender IDs
- `receiverIds`: string[], comma-separated receiver IDs
- `startDate`: string (ISO 8601 date)
- `endDate`: string (ISO 8601 date)
- `page`: integer, default 1
- `pageSize`: integer, default 20

**Response** (200 OK):
```json
{
  "cards": [...],
  "pagination": {...},
  "appliedFilters": {
    "valueIds": ["id1", "id2"],
    "senderIds": ["id3"],
    "receiverIds": [],
    "dateRange": {
      "startDate": "2026-01-01",
      "endDate": "2026-01-27"
    }
  }
}
```

---

##### GET /api/cards/search
**Description**: Search cards by keywords

**Request Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `query`: string, required, min 2 characters
- `page`: integer, default 1
- `pageSize`: integer, default 20
- Can be combined with filter parameters

**Response** (200 OK):
```json
{
  "cards": [...],
  "pagination": {...},
  "searchQuery": "teamwork",
  "matchCount": 15
}
```

---

#### 1.7 Milestones

##### GET /api/milestones/{userId}
**Description**: Get user milestones

**Request Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `userId`: string

**Response** (200 OK):
```json
{
  "userId": "string",
  "milestones": [
    {
      "id": "string",
      "type": "cards_sent | cards_received",
      "threshold": 10,
      "achievedAt": "2026-01-15T10:30:00Z",
      "title": "10 Cards Sent",
      "description": "You've sent 10 thank-you cards!"
    }
  ]
}
```

---

## Unit 2: Analytics Service Unit

### Exposed API Endpoints

#### 2.1 Dashboard Overview

##### GET /api/analytics/dashboard
**Description**: Get HR analytics dashboard overview

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
```

**Query Parameters**:
- `startDate`: string (ISO 8601 date), optional
- `endDate`: string (ISO 8601 date), optional

**Response** (200 OK):
```json
{
  "totalCards": {
    "allTime": 1250,
    "selectedPeriod": 180
  },
  "activeUsers": {
    "senders": 85,
    "receivers": 120,
    "total": 150
  },
  "averageCardsPerEmployee": 8.3,
  "activityTrend": {
    "current": 180,
    "previous": 150,
    "percentageChange": 20.0,
    "direction": "increasing"
  },
  "period": {
    "startDate": "2026-01-01",
    "endDate": "2026-01-27"
  }
}
```

**Error Responses**:
- 403 Forbidden: User is not HR admin

---

#### 2.2 Data Export

##### POST /api/analytics/export
**Description**: Export card data to CSV

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
Content-Type: application/json
```

**Request Body**:
```json
{
  "startDate": "2026-01-01",
  "endDate": "2026-01-27",
  "filters": {
    "valueIds": ["string"],
    "senderIds": ["string"],
    "receiverIds": ["string"]
  }
}
```

**Response** (200 OK):
```json
{
  "downloadUrl": "https://storage.example.com/exports/cards-2026-01-27.csv",
  "expiresAt": "2026-01-27T23:59:59Z",
  "recordCount": 180,
  "fileSize": "45KB"
}
```

**CSV Format**:
```csv
Card ID,Sender Name,Sender Email,Recipient Name,Recipient Email,Recognition Reason,Values/Credos,Created At,Emoji Reactions Count
uuid-1,John Doe,john@company.com,Jane Smith,jane@company.com,Great teamwork,"Make an Impact, Bias for Action",2026-01-15 10:30:00,5
```

---

#### 2.3 Most Active Recognizers

##### GET /api/analytics/recognizers/top
**Description**: Get most active recognizers

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
```

**Query Parameters**:
- `timePeriod`: integer, days (default 30)
- `limit`: integer, default 10, max 100

**Response** (200 OK):
```json
{
  "recognizers": [
    {
      "rank": 1,
      "employeeId": "string",
      "employeeName": "string",
      "department": "string",
      "cardsSent": 45,
      "trend": {
        "previous": 38,
        "percentageChange": 18.4,
        "direction": "increasing"
      }
    }
  ],
  "period": {
    "startDate": "2025-12-28",
    "endDate": "2026-01-27",
    "days": 30
  }
}
```

---

#### 2.4 Team Analytics

##### GET /api/analytics/teams
**Description**: Get team recognition patterns

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
```

**Query Parameters**:
- `timePeriod`: integer, days (default 30)

**Response** (200 OK):
```json
{
  "teams": [
    {
      "teamId": "string",
      "teamName": "string",
      "cardsSent": 120,
      "cardsReceived": 95,
      "employeeCount": 15,
      "averageCardsPerEmployee": 8.0,
      "activityLevel": "high"
    }
  ],
  "period": {
    "startDate": "2025-12-28",
    "endDate": "2026-01-27"
  }
}
```

**Activity Level Calculation**:
- High: Above 75th percentile
- Medium: Between 25th and 75th percentile
- Low: Below 25th percentile

---

##### GET /api/analytics/teams/{teamId}
**Description**: Get specific team details

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
```

**Path Parameters**:
- `teamId`: string

**Query Parameters**:
- `timePeriod`: integer, days (default 30)

**Response** (200 OK):
```json
{
  "teamId": "string",
  "teamName": "string",
  "cardsSent": 120,
  "cardsReceived": 95,
  "employeeCount": 15,
  "averageCardsPerEmployee": 8.0,
  "activityLevel": "high",
  "topRecognizers": [
    {
      "employeeId": "string",
      "employeeName": "string",
      "cardsSent": 25
    }
  ],
  "topRecognized": [
    {
      "employeeId": "string",
      "employeeName": "string",
      "cardsReceived": 18
    }
  ],
  "period": {
    "startDate": "2025-12-28",
    "endDate": "2026-01-27"
  }
}
```

---

#### 2.5 Values Distribution

##### GET /api/analytics/values/distribution
**Description**: Get values/credos distribution

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
```

**Query Parameters**:
- `timePeriod`: integer, days (default 30)

**Response** (200 OK):
```json
{
  "values": [
    {
      "id": "string",
      "name": "Make an Impact",
      "type": "value",
      "count": 145,
      "percentage": 18.5,
      "trend": {
        "previous": 120,
        "percentageChange": 20.8,
        "direction": "increasing"
      }
    }
  ],
  "credos": [
    {
      "id": "string",
      "name": "Bias for Action",
      "type": "credo",
      "count": 98,
      "percentage": 12.5,
      "trend": {
        "previous": 85,
        "percentageChange": 15.3,
        "direction": "increasing"
      }
    }
  ],
  "totalCards": 785,
  "underrepresentedValues": [
    {
      "id": "string",
      "name": "Stay Grounded",
      "count": 25,
      "percentage": 3.2
    }
  ],
  "period": {
    "startDate": "2025-12-28",
    "endDate": "2026-01-27"
  }
}
```

---

##### GET /api/analytics/values/trends
**Description**: Get values trends over time

**Request Headers**:
```
Authorization: Bearer {token}
X-Role: hr-admin
```

**Query Parameters**:
- `timePeriod`: integer, days (default 90)
- `granularity`: "daily" | "weekly" | "monthly", default "weekly"

**Response** (200 OK):
```json
{
  "trends": [
    {
      "valueId": "string",
      "valueName": "Make an Impact",
      "type": "value",
      "dataPoints": [
        {
          "date": "2026-01-01",
          "count": 12
        },
        {
          "date": "2026-01-08",
          "count": 15
        }
      ]
    }
  ],
  "period": {
    "startDate": "2025-10-28",
    "endDate": "2026-01-27"
  },
  "granularity": "weekly"
}
```

---

## Unit 3: Web App Unit

### Consumed APIs

The Web App Unit is a frontend application that consumes APIs from:
1. **Card Service Unit**: All endpoints listed in Section 1
2. **Analytics Service Unit**: All endpoints listed in Section 2 (HR admin only)

### No Exposed APIs
The Web App Unit does not expose any backend APIs. It is a client-side application.

---

## External System Integration

### Teams Channel Integration

#### Inbound: Teams to Card Service

##### Webhook: POST /api/webhooks/teams/card-created
**Description**: Receive notification when card is created (for Teams feed update)

**Request Headers**:
```
X-Teams-Signature: {signature}
Content-Type: application/json
```

**Request Body**:
```json
{
  "eventType": "card.created",
  "cardId": "uuid",
  "timestamp": "2026-01-27T10:30:00Z"
}
```

**Response** (200 OK):
```json
{
  "status": "received"
}
```

---

#### Outbound: Card Service to Teams

##### POST https://teams.microsoft.com/api/channels/{channelId}/messages
**Description**: Post card to Teams channel feed

**Request Headers**:
```
Authorization: Bearer {teams-token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "body": {
    "contentType": "html",
    "content": "<card-html-content>"
  },
  "attachments": [
    {
      "contentType": "application/vnd.microsoft.card.adaptive",
      "content": {
        "type": "AdaptiveCard",
        "body": [...]
      }
    }
  ]
}
```

---

### Employee Directory Service Integration

##### GET /api/directory/employees/search
**Description**: Search employees (external service)

**Request Headers**:
```
Authorization: Bearer {service-token}
```

**Query Parameters**:
- `query`: string
- `limit`: integer

**Response** (200 OK):
```json
{
  "employees": [
    {
      "id": "string",
      "name": "string",
      "email": "string",
      "department": "string",
      "profilePicture": "string",
      "isActive": true
    }
  ]
}
```

---

##### GET /api/directory/employees/{id}
**Description**: Get employee details by ID

**Request Headers**:
```
Authorization: Bearer {service-token}
```

**Response** (200 OK):
```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "department": "string",
  "profilePicture": "string",
  "isActive": true,
  "teamId": "string",
  "teamName": "string"
}
```

---

## Authentication & Authorization

### Authentication Flow
1. User authenticates via Teams/Azure AD
2. Teams/Azure AD issues JWT token
3. Token is included in all API requests
4. Services validate token on each request

### Token Format
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Claims
```json
{
  "sub": "user-id",
  "name": "John Doe",
  "email": "john@company.com",
  "roles": ["employee", "hr-admin"],
  "iat": 1706356800,
  "exp": 1706360400
}
```

### Role-Based Access Control
- **employee**: Access to Card Service endpoints (personal views)
- **hr-admin**: Access to Analytics Service endpoints + all employee permissions

---

## Error Handling

### Standard Error Response Format
```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "Validation failed",
    "details": [
      {
        "field": "recipientIds",
        "message": "At least one recipient is required"
      }
    ],
    "timestamp": "2026-01-27T10:30:00Z",
    "requestId": "uuid"
  }
}
```

### Error Codes
- `INVALID_INPUT`: 400 - Validation error
- `UNAUTHORIZED`: 401 - Missing or invalid authentication
- `FORBIDDEN`: 403 - Insufficient permissions
- `NOT_FOUND`: 404 - Resource not found
- `CONFLICT`: 409 - Resource conflict
- `INTERNAL_ERROR`: 500 - Server error
- `SERVICE_UNAVAILABLE`: 503 - Service temporarily unavailable

---

## Data Synchronization

### Real-time Updates
- **Card Creation**: Immediately available in all feeds
- **Emoji Reactions**: Sync within 1 second across platforms
- **Statistics**: Update in real-time for personal views

### Eventual Consistency
- **Analytics Data**: May have up to 5-minute delay for aggregated metrics
- **Top 10 List**: Refreshed every 5 minutes
- **Team Analytics**: Refreshed hourly

---

## Rate Limiting

### Card Service Unit
- Card creation: 10 requests per minute per user
- Card retrieval: 100 requests per minute per user
- Search: 30 requests per minute per user

### Analytics Service Unit
- Dashboard: 20 requests per minute per user
- Export: 5 requests per hour per user
- Reports: 50 requests per minute per user

### Rate Limit Headers
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1706356800
```

---

## Versioning

### API Versioning Strategy
- Version included in URL path: `/api/v1/cards`
- Current version: v1
- Backward compatibility maintained for at least 6 months after new version release

### Breaking Changes
- New major version required for breaking changes
- Deprecation notice: 3 months minimum
- Migration guide provided for version upgrades

---

## Monitoring & Logging

### Request Logging
All API requests should log:
- Request ID (UUID)
- Timestamp
- User ID
- Endpoint
- HTTP method
- Response status
- Response time

### Metrics
- Request count per endpoint
- Response time (p50, p95, p99)
- Error rate
- Active users
- Card creation rate

---

## Security Considerations

### Data Protection
- All API communication over HTTPS
- Sensitive data encrypted at rest
- PII (Personally Identifiable Information) handled per company policy

### Input Validation
- All inputs validated and sanitized
- SQL injection prevention
- XSS prevention
- CSRF protection

### Audit Trail
- All card creations logged
- All data exports logged
- All admin actions logged

---

## Summary

This integration contract defines:
- **15 Card Service API endpoints** for card management, search, statistics
- **7 Analytics Service API endpoints** for HR analytics and reporting
- **External integrations** with Teams Channel and Employee Directory
- **Authentication/Authorization** via Azure AD JWT tokens
- **Error handling** with standard error response format
- **Rate limiting** to prevent abuse
- **Data synchronization** strategies for real-time and eventual consistency

All units communicate via RESTful HTTP APIs with JSON payloads, following standard HTTP status codes and authentication patterns.

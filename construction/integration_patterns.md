# Integration Patterns - Cross-Context Communication

## Document Information
- **Scope**: Cross-Bounded Context Integration
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## Overview

This document defines how the different bounded contexts and units communicate with each other, including integration patterns, shared concepts, and anti-corruption layers where needed.

---

## Bounded Context Map

```
┌─────────────────────────────────────────────────────────────────┐
│                    System Context Map                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐                                           │
│  │   Web App Unit  │                                           │
│  │ (Presentation)  │                                           │
│  └────────┬────────┘                                           │
│           │ HTTP API                                            │
│           │                                                     │
│  ┌────────▼────────────────────────────────────────┐           │
│  │                                                  │           │
│  │  ┌──────────────┐         ┌─────────────────┐  │           │
│  │  │   Card       │         │   Analytics     │  │           │
│  │  │ Recognition  │◄────────┤   Service       │  │           │
│  │  │   Context    │ Read    │   Context       │  │           │
│  │  │              │ Only    │                 │  │           │
│  │  └──────┬───────┘         └─────────────────┘  │           │
│  │         │                                       │           │
│  │         │ HTTP API                              │           │
│  │         │                                       │           │
│  └─────────┼───────────────────────────────────────┘           │
│            │                                                   │
│            │                                                   │
│  ┌─────────▼────────┐         ┌─────────────────┐             │
│  │   Teams Channel  │         │   Employee      │             │
│  │   (External)     │         │   Directory     │             │
│  │                  │         │   (External)    │             │
│  └──────────────────┘         └─────────────────┘             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Context Relationships

### 1. Card Recognition Context ↔ Analytics Service Context

**Relationship Type**: Customer-Supplier (Analytics as Customer)

**Integration Pattern**: Direct Repository Access (Read-Only)

**Description**: Analytics Service has read-only access to Card Recognition Context data for analytics calculations.

#### Data Flow
```
Card Recognition Context → Analytics Service Context
- Card data (read-only)
- Company values data (read-only)
- No reverse data flow
```

#### Implementation Details
- **No Anti-Corruption Layer**: Not needed due to read-only access
- **Shared Repository Interfaces**: Analytics uses read-only views of Card repositories
- **No Domain Events**: Analytics doesn't subscribe to events, queries data on-demand
- **Data Consistency**: Analytics always gets current data from Card context

#### Shared Concepts
- **Card Data Structure**: Analytics reads Card aggregate data directly
- **Company Values**: Shared reference to CompanyValue entities
- **Employee References**: Both contexts reference employees by ID only

---

### 2. Web App Unit ↔ Card Recognition Context

**Relationship Type**: Customer-Supplier (Web App as Customer)

**Integration Pattern**: REST API Consumption

**Description**: Web App consumes Card Service APIs for all card-related operations.

#### Data Flow
```
Web App → Card Recognition Context
- Card creation requests
- Card retrieval requests
- Search and filter requests
- Personal statistics requests

Card Recognition Context → Web App
- Card data responses
- Validation errors
- Success confirmations
```

#### Implementation Details
- **Anti-Corruption Layer**: Web App maps API responses to UI models
- **Error Handling**: Web App translates domain errors to user-friendly messages
- **Data Transformation**: DTOs used for API communication
- **State Management**: Web App manages UI state separately from domain state

#### API Contract
- **RESTful APIs**: Standard HTTP methods and status codes
- **JSON Payloads**: Structured data exchange
- **Authentication**: JWT tokens for all requests
- **Versioning**: API versioning for backward compatibility

---

### 3. Web App Unit ↔ Analytics Service Context

**Relationship Type**: Customer-Supplier (Web App as Customer)

**Integration Pattern**: REST API Consumption (HR Admin Only)

**Description**: Web App consumes Analytics Service APIs for HR dashboard functionality.

#### Data Flow
```
Web App (HR Admin) → Analytics Service Context
- Dashboard data requests
- Export requests
- Team analytics requests
- Values distribution requests

Analytics Service Context → Web App
- Analytics data responses
- Export file URLs
- Calculated metrics
```

#### Implementation Details
- **Role-Based Access**: Only HR admin users can access analytics APIs
- **Anti-Corruption Layer**: Web App maps analytics responses to chart/dashboard models
- **Large Data Handling**: Streaming for large exports, pagination for lists
- **Caching**: Web App may cache analytics data with appropriate TTL

---

### 4. Card Recognition Context ↔ External Systems

#### 4.1 Employee Directory Service Integration

**Relationship Type**: Customer-Supplier (Card Service as Customer)

**Integration Pattern**: REST API Consumption with Anti-Corruption Layer

**Description**: Card Service validates employees and retrieves employee details from external Employee Directory.

##### Data Flow
```
Card Recognition Context → Employee Directory Service
- Employee search requests
- Employee validation requests
- Team/department queries

Employee Directory Service → Card Recognition Context
- Employee details
- Validation results
- Team/department information
```

##### Implementation Details
- **Anti-Corruption Layer**: Required to protect Card domain from external changes
- **Employee Value Object**: Card context uses its own Employee value object
- **Caching Strategy**: No caching - always fetch fresh data as per requirements
- **Failure Handling**: Graceful degradation when Employee Directory unavailable

##### Anti-Corruption Layer Design
```
External Employee API → Employee ACL → Card Domain Employee Value Object

External Format:
{
  "employeeId": "123",
  "fullName": "John Doe",
  "emailAddress": "john@company.com",
  "departmentName": "Engineering"
}

Domain Format:
Sender/Recipient Value Object:
- EmployeeId: "123"
(Name, email, department fetched on-demand)
```

#### 4.2 Teams Channel Integration

**Relationship Type**: Customer-Supplier (Teams as Customer)

**Integration Pattern**: Event-Driven Integration (Synchronous)

**Description**: Card Service publishes events to Teams Channel for feed updates and notifications.

##### Data Flow
```
Card Recognition Context → Teams Channel
- CardCreated events
- CardShared events
- MilestoneAchieved events

Teams Channel → Card Recognition Context
- Card display requests (via API)
- Share action requests
```

##### Implementation Details
- **Synchronous Events**: Events published immediately via HTTP API
- **Event Payload**: Complete card data included in events
- **No Anti-Corruption Layer**: Teams consumes Card Service API directly
- **Retry Logic**: Implement retry for failed event publishing

---

## Shared Kernel

### Shared Value Objects

#### EmployeeId
**Shared Across**: All contexts
**Definition**: String identifier referencing external Employee Directory
**Validation**: Non-empty string, valid format

#### CompanyValueId  
**Shared Across**: Card Recognition Context, Analytics Service Context
**Definition**: UUID identifier for company values and credos
**Validation**: Valid UUID format, must exist in CompanyValue aggregate

### Shared Enums

#### ValueType
```
enum ValueType {
  Value,
  Credo
}
```

#### MilestoneType
```
enum MilestoneType {
  CardsSent,
  CardsReceived
}
```

### Shared Constants

#### Business Rules Constants
- **MAX_RECOGNITION_REASON_LENGTH**: 1000 characters
- **MIN_RECOGNITION_REASON_LENGTH**: 10 characters
- **MAX_VALUES_PER_CARD**: 3
- **MIN_VALUES_PER_CARD**: 1
- **DEFAULT_ANALYTICS_PERIOD_DAYS**: 30

---

## Event-Driven Integration

### Domain Events Published

#### From Card Recognition Context

##### CardCreated
```
{
  "eventType": "CardCreated",
  "eventId": "uuid",
  "timestamp": "2026-01-27T10:30:00Z",
  "cardId": "uuid",
  "senderId": "string",
  "recipientIds": ["string"],
  "recognitionReason": "string",
  "selectedValueIds": ["uuid"],
  "createdAt": "2026-01-27T10:30:00Z"
}
```

**Consumers**:
- Teams Channel (for feed updates)
- Milestone Tracking Service (asynchronous)

##### CardShared
```
{
  "eventType": "CardShared",
  "eventId": "uuid",
  "timestamp": "2026-01-27T10:30:00Z",
  "cardId": "uuid",
  "sharedBy": "string",
  "targetChannel": "string",
  "sharedAt": "2026-01-27T10:30:00Z"
}
```

**Consumers**:
- Analytics Service (for tracking sharing metrics)

##### MilestoneAchieved
```
{
  "eventType": "MilestoneAchieved",
  "eventId": "uuid",
  "timestamp": "2026-01-27T10:30:00Z",
  "employeeId": "string",
  "milestoneType": "CardsSent | CardsReceived",
  "threshold": 10,
  "achievedAt": "2026-01-27T10:30:00Z",
  "title": "10 Cards Sent",
  "description": "You've sent 10 thank-you cards!"
}
```

**Consumers**:
- Teams Channel (for milestone announcements)
- Notification Service (external)

### Event Processing Patterns

#### Synchronous Processing
- **Teams Channel Integration**: Immediate feed updates
- **Card Creation Response**: Immediate confirmation to user

#### Asynchronous Processing
- **Milestone Detection**: Processed after card creation completes
- **Analytics Updates**: No real-time updates needed

---

## Data Consistency Patterns

### Eventual Consistency
- **Analytics Data**: May have slight delay (acceptable for reporting)
- **Milestone Achievements**: Processed asynchronously
- **Teams Channel Feed**: Immediate consistency via synchronous events

### Strong Consistency
- **Card Creation**: ACID transactions within Card aggregate
- **Employee Validation**: Synchronous validation during card creation
- **Value Selection**: Immediate validation against CompanyValue aggregate

---

## Error Handling Patterns

### Cross-Context Error Handling

#### API Error Mapping
```
Domain Error → HTTP Status → UI Message

ValidationException → 400 Bad Request → "Please check your input"
EmployeeNotFoundException → 404 Not Found → "Employee not found"
UnauthorizedException → 401 Unauthorized → "Please log in"
```

#### Retry Patterns
- **External Service Calls**: Exponential backoff with circuit breaker
- **Event Publishing**: Retry with dead letter queue
- **Database Operations**: Transient error retry

#### Fallback Strategies
- **Employee Directory Unavailable**: Show cached employee names, disable new card creation
- **Teams Channel Unavailable**: Store events for later replay
- **Analytics Service Unavailable**: Show cached dashboard data

---

## Security Integration

### Authentication Flow
```
User → Azure AD/Teams → JWT Token → Web App → Backend Services
```

### Authorization Patterns
- **Role-Based Access Control**: HR admin vs regular employee
- **Context-Level Security**: Each context validates permissions
- **API Gateway Pattern**: Centralized authentication/authorization (future)

### Data Protection
- **Encryption in Transit**: HTTPS for all communications
- **Token Validation**: JWT signature validation in each service
- **Audit Logging**: Cross-context audit trail

---

## Performance Considerations

### Caching Strategies
- **Web App**: Cache API responses with appropriate TTL
- **Analytics Service**: Cache calculation results (infrastructure level)
- **Employee Data**: No caching (always fresh as per requirements)

### Database Patterns
- **Read Replicas**: Analytics queries use read-only replicas
- **Connection Pooling**: Efficient database connection management
- **Query Optimization**: Indexes on common query patterns

### API Performance
- **Pagination**: All list APIs support pagination
- **Compression**: Response compression for large payloads
- **Rate Limiting**: Prevent abuse of API endpoints

---

## Monitoring and Observability

### Cross-Context Tracing
- **Correlation IDs**: Track requests across contexts
- **Distributed Tracing**: End-to-end request tracing
- **Event Tracking**: Monitor event publishing and consumption

### Health Checks
- **Service Health**: Each context exposes health endpoints
- **Dependency Health**: Monitor external service availability
- **Circuit Breaker Status**: Monitor circuit breaker states

### Metrics Collection
- **API Metrics**: Response times, error rates, throughput
- **Business Metrics**: Cards created, milestones achieved
- **Integration Metrics**: External service call success rates

---

## Future Evolution Patterns

### Bounded Context Evolution
- **Context Splitting**: Guidelines for splitting contexts as they grow
- **Context Merging**: Criteria for merging related contexts
- **API Versioning**: Backward compatibility strategies

### Integration Pattern Evolution
- **Event Sourcing**: Potential future pattern for audit trail
- **CQRS**: Separate read/write models if needed
- **Microservices**: Service decomposition strategies

### Technology Evolution
- **Message Queues**: Async event processing infrastructure
- **API Gateway**: Centralized API management
- **Service Mesh**: Advanced service-to-service communication

---

## Summary

The integration patterns defined ensure:

1. **Loose Coupling**: Contexts communicate via well-defined interfaces
2. **High Cohesion**: Each context maintains its domain integrity
3. **Scalability**: Integration patterns support independent scaling
4. **Reliability**: Error handling and retry patterns ensure robustness
5. **Security**: Authentication and authorization enforced across contexts
6. **Observability**: Comprehensive monitoring and tracing capabilities

These patterns support the current requirements while providing flexibility for future evolution and scaling needs.
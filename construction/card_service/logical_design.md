# Card Service Unit - Logical Design

## Document Information
- **Service**: Card Service Unit (Card Recognition Context)
- **Architecture**: Event-Driven, Domain-Driven Design
- **Technology Stack**: Go, PostgreSQL, Gin Framework, GORM
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## System Architecture Overview

The Card Service is designed as a domain-driven, event-driven backend service that manages the complete lifecycle of recognition cards. It follows a layered architecture with clear separation of concerns and implements event-driven patterns for integration with external systems.

### Architecture Principles
- **Domain-Driven Design**: Business logic encapsulated in domain layer
- **Event-Driven Architecture**: HTTP webhooks for external integration
- **Clean Architecture**: Clear separation between layers
- **SOLID Principles**: Maintainable and extensible design
- **Fail-Fast**: Early validation and error detection

---

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    Card Service Architecture                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                 Presentation Layer                       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │   REST API  │ │  Webhooks   │ │   Health    │       │   │
│  │  │ Controllers │ │  Handlers   │ │   Checks    │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                Application Layer                         │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │   Command   │ │    Query    │ │   Event     │       │   │
│  │  │  Handlers   │ │  Handlers   │ │  Handlers   │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │    DTOs     │ │ Validation  │ │ Mapping     │       │   │
│  │  │             │ │   Logic     │ │   Logic     │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  Domain Layer                           │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │ Aggregates  │ │   Domain    │ │   Domain    │       │   │
│  │  │             │ │  Services   │ │   Events    │       │   │
│  │  │ • Card      │ │             │ │             │       │   │
│  │  │ • Value     │ │ • Creation  │ │ • Created   │       │   │
│  │  │ • Milestone │ │ • Milestone │ │ • Shared    │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │   Value     │ │ Factories   │ │ Policies    │       │   │
│  │  │  Objects    │ │             │ │             │       │   │
│  │  │             │ │ • Card      │ │ • Creation  │       │   │
│  │  │ • Sender    │ │ • Milestone │ │ • Milestone │       │   │
│  │  │ • Recipient │ │             │ │ • Values    │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │               Infrastructure Layer                       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │ Repositories│ │   External  │ │    Event    │       │   │
│  │  │             │ │  Services   │ │  Publishers │       │   │
│  │  │ • Card      │ │             │ │             │       │   │
│  │  │ • Value     │ │ • Employee  │ │ • HTTP      │       │   │
│  │  │ • Milestone │ │   Directory │ │   Webhooks  │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │   │
│  │  │  Database   │ │   Logging   │ │ Monitoring  │       │   │
│  │  │ (PostgreSQL)│ │             │ │             │       │   │
│  │  │   + GORM    │ │ • Structured│ │ • Metrics   │       │   │
│  │  │             │ │ • Levels    │ │ • Health    │       │   │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Design and Responsibilities

### Presentation Layer

#### REST API Controllers
**Purpose**: Handle HTTP requests and responses, coordinate with application layer

**Components**:
- **CardController**: Card CRUD operations, filtering, search
- **EmployeeController**: Employee search and validation
- **ValueController**: Company values and credos management
- **StatisticsController**: Personal statistics and top 10 lists
- **MilestoneController**: Milestone retrieval and management

**Responsibilities**:
- HTTP request/response handling
- Input validation (format, required fields)
- Authentication token validation
- Response formatting and status codes
- Error handling and user-friendly messages

#### Webhook Handlers
**Purpose**: Handle incoming webhooks and external system notifications

**Components**:
- **TeamsWebhookHandler**: Handle Teams Channel integration events
- **HealthCheckHandler**: System health and readiness checks

**Responsibilities**:
- Webhook signature validation
- Event payload processing
- Asynchronous event handling
- Integration with external systems

### Application Layer

#### Command Handlers
**Purpose**: Process commands that modify system state

**Components**:
- **CreateCardCommandHandler**: Handle card creation requests
- **AddReactionCommandHandler**: Handle emoji reaction additions (future)
- **RemoveReactionCommandHandler**: Handle emoji reaction removal (future)
- **ShareCardCommandHandler**: Handle card sharing requests

**Responsibilities**:
- Command validation and processing
- Coordinate with domain services
- Transaction management
- Event publishing coordination

#### Query Handlers
**Purpose**: Process queries that retrieve system data

**Components**:
- **GetCardQueryHandler**: Retrieve individual cards
- **GetCardFeedQueryHandler**: Retrieve card feeds with pagination
- **GetPersonalCardsQueryHandler**: Retrieve user's sent/received cards
- **SearchCardsQueryHandler**: Handle card search and filtering
- **GetStatisticsQueryHandler**: Calculate personal statistics
- **GetTop10QueryHandler**: Calculate top recognized employees

**Responsibilities**:
- Query processing and optimization
- Data aggregation and calculation
- Pagination and sorting
- Result formatting

#### Event Handlers
**Purpose**: Handle domain events and coordinate side effects

**Components**:
- **CardCreatedEventHandler**: Handle card creation events
- **CardSharedEventHandler**: Handle card sharing events
- **MilestoneAchievedEventHandler**: Handle milestone achievements

**Responsibilities**:
- Asynchronous event processing
- External system integration
- Milestone detection and recording
- Notification coordination

### Domain Layer

#### Aggregates
**Card Aggregate**:
- **Card** (Root): Core card entity with business logic
- **Sender**: Value object representing card sender
- **Recipients**: Collection of recipient value objects
- **RecognitionReason**: Value object with validation
- **SelectedValues**: References to company values

**CompanyValue Aggregate**:
- **CompanyValue** (Root): Company values and credos entity

**EmployeeMilestone Aggregate**:
- **EmployeeMilestone** (Root): Milestone tracking entity

#### Domain Services
**CardCreationService**:
- Orchestrate card creation process
- Validate recipients and values
- Enforce business rules
- Coordinate with external services

**MilestoneTrackingService**:
- Calculate milestone achievements
- Track user activity across cards
- Generate milestone events

**PersonalStatisticsService**:
- Calculate personal recognition statistics
- Aggregate value usage patterns
- Generate statistical insights

#### Factories
**CardFactory**:
- Create Card aggregates with validation
- Generate unique identifiers
- Enforce creation business rules

**MilestoneFactory**:
- Create milestone instances
- Generate milestone titles and descriptions

### Infrastructure Layer

#### Repositories
**CardRepository**:
- CRUD operations for Card aggregate
- Complex queries for filtering and search
- Pagination and sorting support
- Performance optimization

**CompanyValueRepository**:
- Manage company values and credos
- Reference data operations
- Validation support

**EmployeeMilestoneRepository**:
- Milestone persistence and retrieval
- Achievement tracking
- Historical milestone data

#### External Services
**EmployeeDirectoryService**:
- Employee search and validation
- Fresh employee data retrieval
- Department and team information
- Anti-corruption layer implementation

#### Event Publishers
**HTTPWebhookPublisher**:
- Synchronous HTTP webhook publishing
- Teams Channel integration
- Retry logic and error handling
- Event payload formatting

---

## Database Design and Data Access Patterns

### Database Schema Design

#### Tables Structure

**cards**
```sql
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id VARCHAR(255) NOT NULL,
    recognition_reason TEXT NOT NULL CHECK (length(recognition_reason) >= 10 AND length(recognition_reason) <= 1000),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cards_sender_id ON cards(sender_id);
CREATE INDEX idx_cards_created_at ON cards(created_at DESC);
```

**card_recipients**
```sql
CREATE TABLE card_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    recipient_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_card_recipients_card_id ON card_recipients(card_id);
CREATE INDEX idx_card_recipients_recipient_id ON card_recipients(recipient_id);
CREATE UNIQUE INDEX idx_card_recipients_unique ON card_recipients(card_id, recipient_id);
```

**card_values**
```sql
CREATE TABLE card_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    value_id UUID NOT NULL REFERENCES company_values(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_card_values_card_id ON card_values(card_id);
CREATE INDEX idx_card_values_value_id ON card_values(value_id);
CREATE UNIQUE INDEX idx_card_values_unique ON card_values(card_id, value_id);
```

**company_values**
```sql
CREATE TABLE company_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('Value', 'Credo')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_company_values_type ON company_values(type);
CREATE INDEX idx_company_values_name ON company_values(name);
```

**employee_milestones**
```sql
CREATE TABLE employee_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id VARCHAR(255) NOT NULL,
    milestone_type VARCHAR(50) NOT NULL CHECK (milestone_type IN ('CardsSent', 'CardsReceived')),
    threshold INTEGER NOT NULL,
    achieved_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

CREATE INDEX idx_employee_milestones_employee_id ON employee_milestones(employee_id);
CREATE INDEX idx_employee_milestones_type ON employee_milestones(milestone_type);
CREATE UNIQUE INDEX idx_employee_milestones_unique ON employee_milestones(employee_id, milestone_type, threshold);
```

### Data Access Patterns

#### Repository Implementation with GORM

**Performance Optimizations**:
- **Eager Loading**: Load related data in single queries
- **Batch Operations**: Bulk inserts for recipients and values
- **Query Optimization**: Use appropriate indexes and query patterns
- **Connection Pooling**: Efficient database connection management

**Query Patterns**:
- **Card Feed**: Paginated queries with joins for complete card data
- **Personal Views**: Filtered queries by sender/recipient with pagination
- **Search and Filter**: Complex WHERE clauses with full-text search
- **Statistics**: Aggregation queries with GROUP BY and COUNT operations

#### Transaction Management
- **Card Creation**: Single transaction for card, recipients, and values
- **Milestone Recording**: Separate transaction for milestone achievements
- **Atomic Operations**: Ensure data consistency across related tables

---

## API Design and Endpoint Specifications

### RESTful API Design Principles
- **Resource-Based URLs**: Clear resource identification
- **HTTP Methods**: Proper use of GET, POST, PUT, DELETE
- **Status Codes**: Meaningful HTTP status codes
- **JSON Payloads**: Consistent request/response formats
- **Pagination**: Cursor-based pagination for large datasets
- **Versioning**: API versioning for backward compatibility

### Core API Endpoints

#### Card Management
```
POST   /api/v1/cards                    # Create new card
GET    /api/v1/cards                    # Get company-wide feed
GET    /api/v1/cards/{id}               # Get specific card
GET    /api/v1/cards/received           # Get received cards
GET    /api/v1/cards/sent               # Get sent cards
GET    /api/v1/cards/filter             # Filter cards
GET    /api/v1/cards/search             # Search cards
```

#### Employee Operations
```
GET    /api/v1/employees/search         # Search employees
```

#### Company Values
```
GET    /api/v1/values                   # Get all values/credos
GET    /api/v1/values/{id}              # Get specific value
```

#### Statistics
```
GET    /api/v1/statistics/personal      # Personal statistics
GET    /api/v1/statistics/top10         # Top 10 recognized
```

#### Milestones
```
GET    /api/v1/milestones/{userId}      # User milestones
```

### Request/Response Patterns

#### Standard Response Format
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "timestamp": "2026-01-27T10:30:00Z",
  "requestId": "uuid"
}
```

#### Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "recipients",
        "message": "At least one recipient is required"
      }
    ]
  },
  "timestamp": "2026-01-27T10:30:00Z",
  "requestId": "uuid"
}
```

#### Pagination Format
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "totalPages": 15,
    "totalItems": 300,
    "hasNext": true,
    "hasPrevious": false
  }
}
```

---

## Event Flow Diagrams

### Card Creation Event Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Web App   │    │Card Service │    │  Employee   │    │   Teams     │
│             │    │             │    │  Directory  │    │  Channel    │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                  │                  │
       │ POST /cards      │                  │                  │
       ├─────────────────►│                  │                  │
       │                  │                  │                  │
       │                  │ Validate Recipients                 │
       │                  ├─────────────────►│                  │
       │                  │                  │                  │
       │                  │ Employee Details │                  │
       │                  │◄─────────────────┤                  │
       │                  │                  │                  │
       │                  │ Create Card      │                  │
       │                  │ (Transaction)    │                  │
       │                  │                  │                  │
       │                  │ Publish CardCreated Event           │
       │                  ├─────────────────────────────────────►│
       │                  │                  │                  │
       │ Card Created     │                  │                  │
       │◄─────────────────┤                  │                  │
       │                  │                  │                  │
       │                  │ Async: Check Milestones             │
       │                  │ (Background)     │                  │
       │                  │                  │                  │
```

### Milestone Achievement Event Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│Card Service │    │  Milestone  │    │   Teams     │
│             │    │   Service   │    │  Channel    │
└──────┬──────┘    └──────┬──────┘    └──────┬──────┘
       │                  │                  │
       │ CardCreated      │                  │
       │ Event            │                  │
       ├─────────────────►│                  │
       │                  │                  │
       │                  │ Count User Cards │
       │                  │ Check Thresholds │
       │                  │                  │
       │                  │ Record Milestone │
       │                  │ (if achieved)    │
       │                  │                  │
       │                  │ Publish MilestoneAchieved Event
       │                  ├─────────────────►│
       │                  │                  │
```

---

## Event-Driven Components

### Event Publishing Architecture

#### HTTP Webhook Publisher
**Purpose**: Synchronous event publishing to external systems

**Implementation Pattern**:
- **Immediate Publishing**: Events published synchronously after aggregate changes
- **Retry Logic**: Exponential backoff for failed webhook calls
- **Circuit Breaker**: Prevent cascade failures
- **Event Serialization**: JSON payload formatting
- **Signature Validation**: Secure webhook authentication

#### Event Schema Design

**CardCreated Event**:
```json
{
  "eventType": "CardCreated",
  "eventId": "uuid",
  "timestamp": "2026-01-27T10:30:00Z",
  "version": "1.0",
  "data": {
    "cardId": "uuid",
    "senderId": "string",
    "recipientIds": ["string"],
    "recognitionReason": "string",
    "selectedValueIds": ["uuid"],
    "createdAt": "2026-01-27T10:30:00Z"
  }
}
```

**MilestoneAchieved Event**:
```json
{
  "eventType": "MilestoneAchieved",
  "eventId": "uuid",
  "timestamp": "2026-01-27T10:30:00Z",
  "version": "1.0",
  "data": {
    "employeeId": "string",
    "milestoneType": "CardsSent",
    "threshold": 10,
    "achievedAt": "2026-01-27T10:30:00Z",
    "title": "10 Cards Sent",
    "description": "You've sent 10 thank-you cards!"
  }
}
```

### Asynchronous Processing Patterns

#### Background Job Processing
**Milestone Detection**:
- **Trigger**: CardCreated event
- **Processing**: Asynchronous background job
- **Pattern**: Event-driven background processing
- **Reliability**: Job retry and dead letter handling

**Implementation Approach**:
- **Go Routines**: Lightweight concurrent processing
- **Channel-Based**: Event passing between goroutines
- **Worker Pool**: Controlled concurrency for milestone processing
- **Error Handling**: Graceful error recovery and logging

---

## Scalability and Performance Considerations

### Database Optimization

#### Indexing Strategy
- **Primary Indexes**: All primary keys (UUID)
- **Foreign Key Indexes**: All foreign key relationships
- **Query-Specific Indexes**: Based on common query patterns
- **Composite Indexes**: Multi-column indexes for complex queries
- **Partial Indexes**: For filtered queries (e.g., active records)

#### Query Optimization
- **Eager Loading**: Reduce N+1 query problems
- **Query Planning**: Analyze and optimize slow queries
- **Connection Pooling**: Efficient database connection management
- **Read Replicas**: Separate read and write operations (future)

### Application Performance

#### Memory Management
- **Go Garbage Collection**: Efficient memory management
- **Object Pooling**: Reuse expensive objects
- **Streaming**: Handle large datasets without loading into memory
- **Pagination**: Limit result set sizes

#### Concurrency Patterns
- **Goroutines**: Lightweight concurrent processing
- **Channel Communication**: Safe data passing between goroutines
- **Context Cancellation**: Proper request cancellation handling
- **Rate Limiting**: Prevent resource exhaustion

### Horizontal Scaling Patterns

#### Stateless Design
- **No Session State**: All state in database or external systems
- **Idempotent Operations**: Safe to retry operations
- **Load Balancer Ready**: Can run multiple instances
- **Configuration Externalization**: Environment-based configuration

#### Database Scaling (Future)
- **Read Replicas**: Scale read operations
- **Sharding**: Partition data across multiple databases
- **Connection Pooling**: Efficient connection management
- **Caching Layer**: Reduce database load (when needed)

---

## Security Implementation Details

### Authentication and Authorization

#### JWT Token Validation
- **Token Verification**: Validate JWT signatures
- **Claims Extraction**: Extract user identity and roles
- **Token Expiration**: Handle expired tokens gracefully
- **Refresh Token Support**: Token renewal mechanism

#### Role-Based Access Control
- **Employee Role**: Standard card operations
- **HR Admin Role**: Analytics and export operations
- **Permission Checking**: Endpoint-level authorization
- **Resource-Level Security**: User can only access own data

### Data Protection

#### Input Validation
- **Schema Validation**: Validate request payloads
- **SQL Injection Prevention**: Parameterized queries with GORM
- **XSS Prevention**: Input sanitization
- **Length Limits**: Prevent buffer overflow attacks

#### Data Encryption
- **HTTPS Only**: All communication encrypted in transit
- **Database Encryption**: Sensitive data encrypted at rest
- **Password Hashing**: Secure password storage (if applicable)
- **API Key Management**: Secure external service credentials

### Security Monitoring

#### Audit Logging
- **Request Logging**: All API requests logged
- **Authentication Events**: Login/logout events
- **Data Access**: Sensitive data access logging
- **Error Logging**: Security-related errors

#### Threat Detection
- **Rate Limiting**: Prevent brute force attacks
- **Anomaly Detection**: Unusual access patterns
- **IP Whitelisting**: Restrict access by IP (if needed)
- **Security Headers**: Proper HTTP security headers

---

## Monitoring and Observability Design

### Logging Architecture

#### Structured Logging
- **JSON Format**: Machine-readable log format
- **Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Contextual Information**: Request ID, user ID, correlation ID
- **Performance Metrics**: Response times, query durations

#### Log Categories
- **Application Logs**: Business logic and flow
- **Access Logs**: HTTP request/response logging
- **Error Logs**: Exception and error tracking
- **Audit Logs**: Security and compliance events
- **Performance Logs**: Slow queries and operations

### Metrics Collection

#### Application Metrics
- **Request Metrics**: Count, duration, status codes
- **Business Metrics**: Cards created, users active, milestones achieved
- **Error Metrics**: Error rates, exception counts
- **Performance Metrics**: Response times, throughput

#### System Metrics
- **Resource Usage**: CPU, memory, disk usage
- **Database Metrics**: Connection pool, query performance
- **Network Metrics**: Bandwidth, latency
- **Goroutine Metrics**: Concurrent processing monitoring

### Health Checks

#### Health Check Endpoints
```
GET /health          # Basic health check
GET /health/ready    # Readiness check (dependencies)
GET /health/live     # Liveness check (application)
```

#### Dependency Checks
- **Database Connectivity**: PostgreSQL connection test
- **External Services**: Employee Directory availability
- **Disk Space**: Available storage check
- **Memory Usage**: Memory consumption check

### Alerting Strategy

#### Critical Alerts
- **Service Down**: Application unavailable
- **Database Issues**: Connection failures, slow queries
- **High Error Rates**: Unusual error patterns
- **Security Events**: Authentication failures, suspicious activity

#### Warning Alerts
- **Performance Degradation**: Slow response times
- **Resource Usage**: High CPU/memory usage
- **External Service Issues**: Employee Directory problems
- **Business Metrics**: Unusual activity patterns

---

## Deployment and Operational Considerations

### On-Premises Deployment

#### Application Deployment
- **Binary Deployment**: Single Go binary deployment
- **Configuration Management**: Environment variables and config files
- **Service Management**: Systemd service configuration
- **Log Rotation**: Automated log file management

#### Database Setup
- **PostgreSQL Installation**: Database server setup
- **Schema Migration**: Database schema versioning
- **Backup Strategy**: Regular database backups
- **Performance Tuning**: Database optimization

### Configuration Management

#### Environment Configuration
- **Development**: Local development settings
- **Staging**: Pre-production testing environment
- **Production**: Production environment settings
- **Configuration Validation**: Startup configuration checks

#### Secret Management
- **Environment Variables**: Sensitive configuration
- **File-Based Secrets**: Secure file storage
- **Rotation Strategy**: Regular secret rotation
- **Access Control**: Limited secret access

### Maintenance and Operations

#### Backup and Recovery
- **Database Backups**: Regular automated backups
- **Point-in-Time Recovery**: Transaction log backups
- **Disaster Recovery**: Recovery procedures and testing
- **Data Retention**: Backup retention policies

#### Performance Monitoring
- **Query Performance**: Slow query identification
- **Resource Monitoring**: System resource usage
- **Capacity Planning**: Growth trend analysis
- **Optimization**: Performance tuning recommendations

---

## Implementation Roadmap

### Phase 1: Core Foundation (Weeks 1-2)
- **Domain Model Implementation**: Aggregates, entities, value objects
- **Repository Layer**: Database access with GORM
- **Basic API Endpoints**: Card CRUD operations
- **Authentication**: JWT token validation

### Phase 2: Business Logic (Weeks 3-4)
- **Card Creation Flow**: Complete card creation with validation
- **Employee Integration**: Employee Directory Service integration
- **Company Values**: Values and credos management
- **Basic Statistics**: Personal statistics calculation

### Phase 3: Advanced Features (Weeks 5-6)
- **Search and Filtering**: Advanced card search capabilities
- **Event Publishing**: HTTP webhook integration
- **Milestone Tracking**: Asynchronous milestone detection
- **Top 10 Lists**: Recognition leaderboards

### Phase 4: Production Readiness (Weeks 7-8)
- **Monitoring and Logging**: Comprehensive observability
- **Security Hardening**: Security review and improvements
- **Performance Optimization**: Query and application optimization
- **Documentation**: API documentation and deployment guides

---

## Success Criteria and Metrics

### Functional Requirements
- ✅ All user stories implemented and tested
- ✅ API endpoints meet integration contract specifications
- ✅ Domain model supports all business rules
- ✅ Event publishing works reliably

### Non-Functional Requirements
- **Performance**: API response times < 200ms (95th percentile)
- **Reliability**: 99.9% uptime availability
- **Scalability**: Support 1000+ concurrent users
- **Security**: Pass security audit and penetration testing

### Business Metrics
- **Card Creation Rate**: Track daily card creation volume
- **User Engagement**: Active users sending/receiving cards
- **Milestone Achievements**: Milestone completion rates
- **System Usage**: API endpoint usage patterns

This logical design provides a comprehensive blueprint for implementing the Card Service Unit with proper architecture, scalability, and operational considerations while maintaining simplicity and avoiding over-engineering.
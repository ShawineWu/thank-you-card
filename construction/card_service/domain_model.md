# Card Service Unit - Domain Model

## Document Information
- **Bounded Context**: Card Recognition Context
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## Bounded Context Description

The **Card Recognition Context** is responsible for managing the core business logic of employee recognition through digital thank-you cards. This bounded context encapsulates all operations related to creating, storing, retrieving, and managing recognition cards, including their associated values, recipients, and milestone tracking.

**Context Boundaries:**
- Card lifecycle management (creation, retrieval, filtering, searching)
- Company values and credos management
- Employee milestone tracking
- Recognition statistics calculation
- Integration with external Employee Directory Service

**What's NOT in this context:**
- Employee management (handled by Employee Directory Service)
- HR analytics and reporting (handled by Analytics Service)
- User interface concerns (handled by Web App Unit)
- Teams Channel integration (external system)

---

## Ubiquitous Language Glossary

### Core Domain Terms

**Recognition Card (Card)**: A digital thank-you message sent from one employee to one or more recipients, acknowledging specific contributions and aligning them with company values.

**Sender**: The employee who creates and sends a recognition card.

**Recipient**: An employee who receives a recognition card.

**Recognition Reason**: The textual explanation of why the recognition is being given, describing the specific action or behavior being appreciated.

**Company Value**: One of the five core organizational values that guide company culture (Make an Impact, Strive for Excellence, Stand Together, Be Open-Minded, Stay Grounded).

**Company Credo**: One of the ten operational principles that define how work should be done (Bias for Action, Customer Centric, Think Strategically, etc.).

**Milestone**: A significant achievement in an employee's recognition activity (e.g., first card sent, 10th card received).

**Card Feed**: A chronologically ordered collection of recognition cards visible to all employees.

**Personal Statistics**: Aggregated data about an individual employee's recognition activity (cards sent/received, frequently used values).

### Technical Terms

**Card ID**: Unique identifier for each recognition card (UUID format).

**Employee ID**: External reference to an employee in the Employee Directory Service.

**Value ID**: Unique identifier for company values and credos.

**Timestamp**: The exact date and time when a card was created.

---

## Aggregates

### 1. Card Aggregate

**Aggregate Root**: `Card`

**Description**: The Card aggregate represents a complete recognition transaction, including all information about the sender, recipients, reason, and associated company values. This aggregate maintains consistency around the card creation and modification operations.

**Aggregate Boundary**: 
- Card entity (root)
- Sender value object
- Recipients collection (value objects)
- Recognition reason value object
- Selected values collection (entity references)
- Creation timestamp

**Invariants**:
- A card must have exactly one sender
- A card must have at least one recipient
- A card must have a non-empty recognition reason (max 1000 characters)
- A card must have 1-3 selected company values/credos
- A card's creation timestamp cannot be modified after creation
- Recipients must be valid, active employees (validated during creation)

**Business Rules**:
- Cards are immutable after creation (no updates allowed)
- All recipients must be different from the sender
- Selected values must be valid company values or credos
- Recognition reason must not be empty or only whitespace

---

### 2. CompanyValue Aggregate

**Aggregate Root**: `CompanyValue`

**Description**: Represents the predefined company values and credos that can be associated with recognition cards. These are reference data that rarely change.

**Aggregate Boundary**:
- CompanyValue entity (root)
- Value name
- Value description
- Value type (Value or Credo)

**Invariants**:
- Each company value has a unique ID
- Value names are unique within their type
- Values are immutable once created
- Value type must be either "Value" or "Credo"

**Business Rules**:
- Only predefined values can be selected for cards
- Values cannot be deleted if referenced by existing cards
- Maximum 3 values can be selected per card

---

### 3. EmployeeMilestone Aggregate

**Aggregate Root**: `EmployeeMilestone`

**Description**: Tracks milestone achievements for individual employees based on their recognition activity. Managed by domain service for cross-card calculations.

**Aggregate Boundary**:
- EmployeeMilestone entity (root)
- Employee ID reference
- Milestone type (cards sent/received)
- Achievement thresholds
- Achievement timestamps

**Invariants**:
- Each employee can have multiple milestones
- Milestone achievements are immutable once recorded
- Achievement timestamps must be chronologically consistent

**Business Rules**:
- Milestones are calculated based on card counts
- Standard thresholds: 1, 5, 10, 25, 50, 100 cards
- Both sending and receiving milestones are tracked

---

## Entities

### Card (Aggregate Root)
**Identity**: CardId (UUID)
**Attributes**:
- CardId: Unique identifier
- Sender: Sender value object
- Recipients: Collection of Recipient value objects
- RecognitionReason: RecognitionReason value object
- SelectedValues: Collection of CompanyValue entity references
- CreatedAt: Timestamp
- Status: CardStatus (Active)

**Behaviors**:
- Create new card (via factory)
- Retrieve card details
- Validate card completeness

### CompanyValue (Aggregate Root)
**Identity**: ValueId (UUID)
**Attributes**:
- ValueId: Unique identifier
- Name: String (e.g., "Make an Impact")
- Description: String (detailed explanation)
- Type: ValueType enum (Value, Credo)

**Behaviors**:
- Retrieve value details
- Validate value selection

### EmployeeMilestone (Aggregate Root)
**Identity**: MilestoneId (UUID)
**Attributes**:
- MilestoneId: Unique identifier
- EmployeeId: External employee reference
- MilestoneType: MilestoneType enum (CardsSent, CardsReceived)
- Threshold: Integer (milestone level)
- AchievedAt: Timestamp
- Title: String (display title)
- Description: String (achievement description)

**Behaviors**:
- Record milestone achievement
- Check if milestone threshold reached

---

## Value Objects

### Sender
**Attributes**:
- EmployeeId: String (external reference)

**Validation Rules**:
- Must reference a valid, active employee
- Cannot be null or empty

### Recipient
**Attributes**:
- EmployeeId: String (external reference)

**Validation Rules**:
- Must reference a valid, active employee
- Cannot be the same as sender
- Cannot be null or empty

### RecognitionReason
**Attributes**:
- Text: String

**Validation Rules**:
- Cannot be null, empty, or only whitespace
- Maximum 1000 characters
- Minimum 10 characters

### CardId
**Attributes**:
- Value: UUID

**Validation Rules**:
- Must be a valid UUID format
- Immutable once assigned

### ValueId
**Attributes**:
- Value: UUID

**Validation Rules**:
- Must be a valid UUID format
- Must reference existing company value

---

## Domain Events

### CardCreated
**Description**: Published when a new recognition card is successfully created
**Attributes**:
- CardId: UUID
- SenderId: String
- RecipientIds: List<String>
- RecognitionReason: String
- SelectedValueIds: List<UUID>
- CreatedAt: Timestamp

**Consumers**:
- Teams Channel integration (for feed updates)
- Milestone tracking service
- Analytics service (for statistics)

### CardShared
**Description**: Published when a card is shared to other Teams channels
**Attributes**:
- CardId: UUID
- SharedBy: String (employee ID)
- TargetChannel: String
- SharedAt: Timestamp

**Consumers**:
- Analytics service (for tracking sharing activity)

### MilestoneAchieved
**Description**: Published when an employee reaches a recognition milestone
**Attributes**:
- EmployeeId: String
- MilestoneType: MilestoneType enum
- Threshold: Integer
- AchievedAt: Timestamp
- Title: String
- Description: String

**Consumers**:
- Notification service (for milestone notifications)
- Teams Channel integration (for milestone announcements)

---

## Domain Services

### CardCreationService
**Purpose**: Orchestrates the complex card creation process including validation, value verification, and employee validation

**Operations**:
- `CreateCard(sender, recipients, reason, valueIds)`: Creates a new card with full validation
- `ValidateRecipients(recipientIds)`: Validates recipients against Employee Directory
- `ValidateValues(valueIds)`: Validates selected values against available options

**Dependencies**:
- EmployeeDirectoryService (external)
- CompanyValueRepository
- CardRepository

### MilestoneTrackingService
**Purpose**: Manages milestone detection and achievement recording across multiple cards

**Operations**:
- `CheckMilestones(employeeId)`: Evaluates if employee has reached new milestones
- `RecordMilestone(employeeId, milestoneType, threshold)`: Records milestone achievement
- `GetEmployeeMilestones(employeeId)`: Retrieves all milestones for an employee

**Dependencies**:
- CardRepository (for counting cards)
- EmployeeMilestoneRepository

### PersonalStatisticsService
**Purpose**: Calculates personal recognition statistics for individual employees

**Operations**:
- `CalculatePersonalStats(employeeId)`: Computes cards sent/received and value frequencies
- `GetMostUsedValues(employeeId)`: Returns frequently used values when sending
- `GetMostReceivedValues(employeeId)`: Returns frequently received values

**Dependencies**:
- CardRepository
- CompanyValueRepository

---

## Repositories

### CardRepository
**Purpose**: Manages persistence and retrieval of Card aggregates

**Interface**:
- `Save(card)`: Persists a new card
- `FindById(cardId)`: Retrieves card by ID
- `FindBySender(senderId, pagination)`: Gets cards sent by employee
- `FindByRecipient(recipientId, pagination)`: Gets cards received by employee
- `FindAll(pagination)`: Gets company-wide card feed
- `FindByFilters(filters, pagination)`: Gets filtered cards
- `SearchByKeywords(keywords, pagination)`: Gets cards matching search terms
- `CountBySender(senderId)`: Counts cards sent by employee
- `CountByRecipient(recipientId)`: Counts cards received by employee

### CompanyValueRepository
**Purpose**: Manages persistence and retrieval of CompanyValue aggregates

**Interface**:
- `FindAll()`: Gets all company values and credos
- `FindById(valueId)`: Gets specific value by ID
- `FindByIds(valueIds)`: Gets multiple values by IDs
- `FindByType(valueType)`: Gets values by type (Value or Credo)

### EmployeeMilestoneRepository
**Purpose**: Manages persistence and retrieval of EmployeeMilestone aggregates

**Interface**:
- `Save(milestone)`: Persists milestone achievement
- `FindByEmployee(employeeId)`: Gets all milestones for employee
- `FindByEmployeeAndType(employeeId, milestoneType)`: Gets specific milestone type
- `ExistsByEmployeeAndThreshold(employeeId, milestoneType, threshold)`: Checks if milestone exists

---

## Factories

### CardFactory
**Purpose**: Creates Card aggregates with proper validation and business rule enforcement

**Operations**:
- `CreateCard(senderId, recipientIds, reason, valueIds)`: Creates new card instance
- `ValidateCardData(cardData)`: Validates all card creation data
- `GenerateCardId()`: Generates unique card identifier

**Validation Logic**:
- Sender validation (active employee)
- Recipients validation (active employees, different from sender)
- Recognition reason validation (length, content)
- Values validation (1-3 valid values)
- Business rule enforcement

### MilestoneFactory
**Purpose**: Creates EmployeeMilestone aggregates for achievement tracking

**Operations**:
- `CreateMilestone(employeeId, milestoneType, threshold)`: Creates milestone instance
- `GenerateMilestoneTitle(milestoneType, threshold)`: Creates display title
- `GenerateMilestoneDescription(milestoneType, threshold)`: Creates description

---

## Policies

### CardCreationPolicy
**Rules**:
- Cards can only be created by active employees
- Recipients must be active employees different from sender
- Recognition reason must be meaningful (10-1000 characters)
- Must select 1-3 company values/credos
- No duplicate recipients allowed
- Cards are immutable after creation

### MilestonePolicy
**Rules**:
- Milestones are calculated asynchronously after card creation
- Standard thresholds: 1, 5, 10, 25, 50, 100 cards
- Both sending and receiving milestones tracked
- Milestone achievements are permanent
- Duplicate milestones not allowed

### ValueSelectionPolicy
**Rules**:
- Only predefined company values/credos can be selected
- Minimum 1, maximum 3 values per card
- Values must be currently active
- Both Values and Credos can be mixed in selection

---

## Business Rules and Invariants

### Card Aggregate Invariants
1. **Card Completeness**: Every card must have sender, recipients, reason, and values
2. **Recipient Validity**: All recipients must be valid, active employees
3. **Sender-Recipient Separation**: Sender cannot be a recipient of the same card
4. **Value Count Constraint**: Must have 1-3 selected values
5. **Immutability**: Cards cannot be modified after creation
6. **Reason Length**: Recognition reason must be 10-1000 characters

### Cross-Aggregate Business Rules
1. **Employee Validation**: All employee references validated against Employee Directory
2. **Value Reference Integrity**: Selected values must exist in CompanyValue aggregate
3. **Milestone Uniqueness**: Each milestone threshold achieved only once per employee
4. **Chronological Consistency**: Card timestamps must be chronologically consistent

### System-Wide Constraints
1. **No Card Limits**: No restrictions on card sending frequency or volume
2. **Public Visibility**: All cards are publicly visible company-wide
3. **Immediate Availability**: Cards available in feed immediately after creation
4. **External Dependencies**: Employee data always fetched fresh from external service

---

## Integration Points

### External Services
- **Employee Directory Service**: Validates employee IDs, retrieves employee details
- **Teams Channel Service**: Receives card creation events for feed updates

### Internal Services
- **Analytics Service**: Reads card data for analytics and reporting
- **Web App**: Consumes all Card Service APIs for user interface

### Event Publishing
- Domain events published synchronously to Teams HTTP interface
- Events contain all necessary data for consumers
- Event publishing is part of the aggregate transaction

---

## Domain Model Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    Card Recognition Context                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐         ┌─────────────────┐               │
│  │   Card          │         │  CompanyValue   │               │
│  │   (Aggregate)   │◄────────┤  (Aggregate)    │               │
│  │                 │         │                 │               │
│  │ - CardId        │         │ - ValueId       │               │
│  │ - Sender        │         │ - Name          │               │
│  │ - Recipients[]  │         │ - Description   │               │
│  │ - Reason        │         │ - Type          │               │
│  │ - Values[]      │         │                 │               │
│  │ - CreatedAt     │         └─────────────────┘               │
│  └─────────────────┘                                           │
│           │                                                     │
│           │                                                     │
│  ┌─────────────────┐                                           │
│  │EmployeeMilestone│                                           │
│  │   (Aggregate)   │                                           │
│  │                 │                                           │
│  │ - MilestoneId   │                                           │
│  │ - EmployeeId    │                                           │
│  │ - Type          │                                           │
│  │ - Threshold     │                                           │
│  │ - AchievedAt    │                                           │
│  └─────────────────┘                                           │
│                                                                 │
│  Value Objects:                                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐              │
│  │   Sender    │ │  Recipient  │ │Recognition  │              │
│  │             │ │             │ │   Reason    │              │
│  │-EmployeeId  │ │-EmployeeId  │ │   - Text    │              │
│  └─────────────┘ └─────────────┘ └─────────────┘              │
│                                                                 │
│  Domain Services:                                               │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ CardCreationService                                     │   │
│  │ MilestoneTrackingService                               │   │
│  │ PersonalStatisticsService                              │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Repositories:                                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ CardRepository                                          │   │
│  │ CompanyValueRepository                                  │   │
│  │ EmployeeMilestoneRepository                            │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Domain Events:                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ CardCreated, CardShared, MilestoneAchieved             │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation Notes

### Aggregate Design Decisions
- **Card as Aggregate Root**: Card contains all related data for consistency
- **Small Aggregates**: Each aggregate focuses on single responsibility
- **External References**: Employee data referenced by ID only, not embedded
- **Immutable Cards**: Simplifies concurrency and maintains audit trail

### Event-Driven Architecture
- **Synchronous Events**: Immediate Teams Channel integration
- **Asynchronous Milestones**: Milestone calculation doesn't block card creation
- **Event Sourcing**: Not implemented but events provide audit trail

### Performance Considerations
- **Fresh Employee Data**: Always fetch from external service for accuracy
- **Pagination**: All list operations support pagination
- **Indexing**: Repository implementations should index on common query fields
- **Caching**: No domain-level caching, handled at infrastructure level

### Future Extensibility
- **Card Approval**: Can be added as Card state transitions
- **Custom Templates**: Can be added as Card value objects
- **Advanced Reactions**: Can be modeled as separate aggregate if needed
- **Audit Trail**: Domain events provide foundation for audit logging

---

## Validation Summary

This domain model supports all user stories from the Card Service Unit:
- ✅ Card creation with validation (US-2.1 to US-2.5)
- ✅ Personal card views and statistics (US-3.1 to US-3.4)
- ✅ Card filtering and search (US-4.2, US-4.3)
- ✅ Teams Channel integration support (US-1.1, US-1.4)
- ✅ Milestone tracking (US-6.1)

All business rules and invariants are properly enforced through aggregate boundaries, domain services, and policies.
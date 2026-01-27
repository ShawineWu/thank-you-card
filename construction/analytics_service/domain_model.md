# Analytics Service Unit - Domain Model

## Document Information
- **Bounded Context**: Recognition Analytics Context
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

---

## Bounded Context Description

The **Recognition Analytics Context** is responsible for providing comprehensive insights and reporting capabilities for the employee recognition system. This bounded context focuses on aggregating, analyzing, and presenting recognition data to HR administrators for strategic decision-making.

**Context Boundaries:**
- HR analytics dashboard data aggregation
- Recognition pattern analysis and reporting
- Team and department performance analytics
- Company values distribution analysis
- Data export functionality for external analysis
- Time-based trend calculations

**What's NOT in this context:**
- Card creation or modification (handled by Card Recognition Context)
- Employee management (handled by Employee Directory Service)
- User interface concerns (handled by Web App Unit)
- Real-time card operations (handled by Card Service)

**Integration Pattern:**
- Read-only access to Card Recognition Context data
- On-demand computation of analytics
- Direct querying of Card Service domain model

---

## Ubiquitous Language Glossary

### Core Analytics Terms

**Recognition Analytics**: The comprehensive analysis of employee recognition patterns, trends, and behaviors within the organization.

**Active Recognizer**: An employee who has sent one or more recognition cards within a specified time period.

**Recognition Activity**: The overall level of card sending and receiving within teams, departments, or the entire organization.

**Value Distribution**: The frequency and percentage breakdown of how often each company value or credo is recognized.

**Team Recognition Pattern**: The analysis of recognition behaviors within specific teams or departments, including sending/receiving ratios and activity levels.

**Recognition Trend**: The change in recognition activity over time, showing increases, decreases, or stability in recognition behaviors.

**Activity Level**: A categorization of team or individual recognition activity (High, Medium, Low) based on statistical percentiles.

**Export Dataset**: A comprehensive data extract containing all recognition card information for external analysis.

### Measurement Terms

**Cards Sent Count**: The total number of recognition cards sent by an employee or team within a time period.

**Cards Received Count**: The total number of recognition cards received by an employee or team within a time period.

**Average Cards Per Employee**: The mean number of cards (sent or received) per employee within a team or organization.

**Recognition Frequency**: How often recognition occurs within a specific time frame.

**Value Usage Frequency**: How often specific company values or credos are selected in recognition cards.

**Percentile Ranking**: Statistical ranking of teams or individuals based on their recognition activity.

### Time-Based Terms

**Analysis Period**: The specific time range for which analytics are calculated (default: 30 days).

**Trend Period**: The comparison period used to calculate trends (typically the previous period of equal duration).

**Time Granularity**: The level of time detail for trend analysis (daily, weekly, monthly).

---

## Domain Services

Since Analytics Service uses on-demand computation and directly queries the Card Service domain model, it primarily consists of domain services rather than aggregates and entities.

### RecognitionAnalyticsService
**Purpose**: Provides comprehensive analytics calculations and insights for HR administrators

**Operations**:
- `CalculateDashboardOverview(timePeriod)`: Computes key metrics for dashboard display
- `CalculateActivityTrend(currentPeriod, previousPeriod)`: Determines recognition activity trends
- `IdentifyActiveUsers(timePeriod)`: Counts distinct senders and receivers
- `CalculateAverageCardsPerEmployee(timePeriod)`: Computes organization-wide averages

**Dependencies**:
- CardRepository (from Card Service)
- EmployeeDirectoryService (external)

**Business Logic**:
- Aggregates card data across specified time periods
- Calculates percentage changes between periods
- Determines trend directions (increasing/decreasing/stable)

### MostActiveRecognizersService
**Purpose**: Identifies and ranks employees who are most active in sending recognition

**Operations**:
- `GetTopRecognizers(timePeriod, limit)`: Returns ranked list of most active senders
- `CalculateRecognizerTrends(employeeId, currentPeriod, previousPeriod)`: Computes individual trends
- `RankRecognizersByActivity(recognizers)`: Applies ranking logic with tie handling

**Dependencies**:
- CardRepository (from Card Service)
- EmployeeDirectoryService (external)

**Business Logic**:
- Counts cards sent per employee within time period
- Calculates trend comparisons with previous periods
- Handles ranking ties appropriately
- Includes employee details (name, department)

### TeamRecognitionAnalyticsService
**Purpose**: Analyzes recognition patterns at team and department levels

**Operations**:
- `AnalyzeTeamPatterns(timePeriod)`: Computes team-level recognition metrics
- `GetTeamDetails(teamId, timePeriod)`: Provides detailed team analysis
- `CategorizeTeamActivityLevels(teams)`: Assigns activity level categories
- `CompareTeamPerformance(teams)`: Enables side-by-side team comparisons

**Dependencies**:
- CardRepository (from Card Service)
- EmployeeDirectoryService (external)

**Business Logic**:
- Aggregates cards sent/received per team
- Calculates average cards per employee per team
- Determines activity levels using percentile thresholds:
  - High: Above 75th percentile
  - Medium: Between 25th and 75th percentile
  - Low: Below 25th percentile
- Identifies teams needing recognition encouragement

### ValuesDistributionAnalyticsService
**Purpose**: Analyzes the distribution and trends of company values recognition

**Operations**:
- `CalculateValuesDistribution(timePeriod)`: Computes count and percentage for each value
- `AnalyzeValuesTrends(timePeriod, granularity)`: Tracks value usage over time
- `IdentifyUnderrepresentedValues(distribution)`: Finds values below average usage
- `CompareValuesPeriods(currentPeriod, previousPeriod)`: Calculates period-over-period changes

**Dependencies**:
- CardRepository (from Card Service)
- CompanyValueRepository (from Card Service)

**Business Logic**:
- Counts value selections across all cards
- Calculates percentage distribution
- Identifies values needing more emphasis
- Tracks trends for strategic insights

### DataExportService
**Purpose**: Generates comprehensive data exports for external analysis

**Operations**:
- `ExportCardData(timePeriod, filters)`: Creates CSV export with all card details
- `ApplyExportFilters(cards, filters)`: Filters data based on criteria
- `FormatExportData(cards)`: Formats data for CSV output
- `GenerateDownloadUrl(exportFile)`: Creates secure download link

**Dependencies**:
- CardRepository (from Card Service)
- EmployeeDirectoryService (external)
- FileStorageService (infrastructure)

**Business Logic**:
- Includes all required fields: sender, recipient, reason, values, timestamp, reactions
- Supports date range filtering
- Supports value, sender, receiver filtering
- Handles large datasets efficiently
- Provides secure, time-limited download URLs

---

## Value Objects

### AnalyticsPeriod
**Purpose**: Represents a time period for analytics calculations

**Attributes**:
- StartDate: DateTime
- EndDate: DateTime
- Days: Integer (calculated)

**Validation Rules**:
- StartDate must be before EndDate
- Period cannot be longer than 2 years
- Dates must be valid

**Behaviors**:
- Calculate period duration
- Validate period boundaries
- Generate previous period of same duration

### ActivityTrend
**Purpose**: Represents trend analysis between two periods

**Attributes**:
- Current: Integer (current period value)
- Previous: Integer (previous period value)
- PercentageChange: Decimal
- Direction: TrendDirection enum (Increasing, Decreasing, Stable)

**Validation Rules**:
- Current and Previous must be non-negative
- PercentageChange calculated as ((Current - Previous) / Previous) * 100
- Direction based on percentage change thresholds

**Behaviors**:
- Calculate percentage change
- Determine trend direction
- Format trend display

### TeamActivityLevel
**Purpose**: Categorizes team recognition activity

**Attributes**:
- Level: ActivityLevel enum (High, Medium, Low)
- Threshold: Decimal (percentile threshold)
- Description: String

**Validation Rules**:
- Level must be valid enum value
- Threshold must be between 0 and 100

**Behaviors**:
- Determine activity level from percentile
- Generate level description

### ValueDistribution
**Purpose**: Represents the distribution of a company value usage

**Attributes**:
- ValueId: UUID
- ValueName: String
- ValueType: ValueType enum (Value, Credo)
- Count: Integer
- Percentage: Decimal
- Trend: ActivityTrend

**Validation Rules**:
- Count must be non-negative
- Percentage must be between 0 and 100
- ValueId must reference valid company value

**Behaviors**:
- Calculate percentage from total
- Compare with previous period
- Identify if underrepresented

---

## Repositories (Read-Only Interfaces)

### AnalyticsCardRepository
**Purpose**: Read-only interface to Card data for analytics calculations

**Interface**:
- `CountCardsByPeriod(startDate, endDate)`: Total cards in period
- `CountCardsBySender(senderId, startDate, endDate)`: Cards sent by employee
- `CountCardsByRecipient(recipientId, startDate, endDate)`: Cards received by employee
- `GetCardsByPeriod(startDate, endDate)`: All cards in period
- `GetCardsByTeam(teamId, startDate, endDate)`: Team-specific cards
- `GetValueUsageByPeriod(startDate, endDate)`: Value selection counts
- `GetDistinctSenders(startDate, endDate)`: Unique senders in period
- `GetDistinctRecipients(startDate, endDate)`: Unique recipients in period

**Implementation Note**: This is a read-only view of the CardRepository from Card Service

### AnalyticsEmployeeRepository
**Purpose**: Read-only interface to Employee data for analytics

**Interface**:
- `GetEmployeeDetails(employeeId)`: Employee information
- `GetEmployeesByTeam(teamId)`: Team member list
- `GetAllTeams()`: Organization structure
- `GetTeamDetails(teamId)`: Team information
- `CountActiveEmployees()`: Total employee count

**Implementation Note**: This interfaces with the external Employee Directory Service

---

## Policies

### AnalyticsAccessPolicy
**Rules**:
- Only HR administrators can access analytics endpoints
- All analytics operations are read-only
- No modification of card data allowed
- Access logging required for audit purposes

### DataExportPolicy
**Rules**:
- Exports limited to HR administrators
- Maximum export period: 2 years
- Export files have 24-hour expiration
- Export actions logged for compliance
- Personal employee data included (authorized for HR)

### AnalyticsCalculationPolicy
**Rules**:
- All calculations performed on-demand
- No pre-aggregated data storage
- Fresh data fetched for each request
- Time periods default to 30 days
- Trend calculations use equal-duration previous periods

### ActivityLevelPolicy
**Rules**:
- High activity: Above 75th percentile
- Medium activity: Between 25th and 75th percentile
- Low activity: Below 25th percentile
- Percentiles calculated across all teams
- Minimum 3 teams required for meaningful percentiles

---

## Business Rules and Invariants

### Analytics Calculation Rules
1. **Time Period Boundaries**: All analytics respect specified time boundaries
2. **Data Freshness**: Always use current data from Card Service
3. **Trend Consistency**: Trend periods must be equal duration for valid comparison
4. **Percentile Accuracy**: Activity levels calculated using proper statistical methods
5. **Zero Handling**: Handle zero values gracefully in percentage calculations

### Data Export Rules
1. **Complete Data**: Exports include all required fields per specification
2. **Filter Consistency**: Applied filters must be logically consistent
3. **Data Privacy**: Employee data included only for authorized HR users
4. **File Security**: Export files secured with time-limited access
5. **Audit Trail**: All export operations logged with user and timestamp

### Access Control Rules
1. **HR Only**: Analytics endpoints restricted to HR administrator role
2. **Read Only**: No write operations allowed in analytics context
3. **Audit Logging**: All access attempts logged for security
4. **Data Scope**: Access to all employee recognition data within role permissions

---

## Integration Points

### Card Recognition Context Integration
- **Read-Only Access**: Analytics Service reads but never modifies card data
- **Direct Repository Access**: Uses Card Service repositories for data access
- **No Domain Events**: Analytics doesn't publish events, only consumes data
- **Anti-Corruption Layer**: Not needed due to read-only access pattern

### External Service Integration
- **Employee Directory Service**: For employee and team information
- **File Storage Service**: For export file management
- **Authentication Service**: For HR role validation

### Web App Integration
- **API Consumption**: Web App consumes all Analytics Service endpoints
- **Role-Based UI**: Web App enforces HR-only access in user interface
- **Data Visualization**: Web App handles chart and graph rendering

---

## Domain Model Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                Recognition Analytics Context                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Domain Services (Core Business Logic):                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ RecognitionAnalyticsService                             │   │
│  │ - CalculateDashboardOverview()                          │   │
│  │ - CalculateActivityTrend()                              │   │
│  │ - IdentifyActiveUsers()                                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ MostActiveRecognizersService                            │   │
│  │ - GetTopRecognizers()                                   │   │
│  │ - CalculateRecognizerTrends()                           │   │
│  │ - RankRecognizersByActivity()                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ TeamRecognitionAnalyticsService                         │   │
│  │ - AnalyzeTeamPatterns()                                 │   │
│  │ - GetTeamDetails()                                      │   │
│  │ - CategorizeTeamActivityLevels()                        │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ ValuesDistributionAnalyticsService                      │   │
│  │ - CalculateValuesDistribution()                         │   │
│  │ - AnalyzeValuesTrends()                                 │   │
│  │ - IdentifyUnderrepresentedValues()                      │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ DataExportService                                       │   │
│  │ - ExportCardData()                                      │   │
│  │ - ApplyExportFilters()                                  │   │
│  │ - GenerateDownloadUrl()                                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Value Objects:                                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐              │
│  │ Analytics   │ │  Activity   │ │    Value    │              │
│  │   Period    │ │    Trend    │ │Distribution │              │
│  │             │ │             │ │             │              │
│  │-StartDate   │ │-Current     │ │-ValueId     │              │
│  │-EndDate     │ │-Previous    │ │-Count       │              │
│  │-Days        │ │-Percentage  │ │-Percentage  │              │
│  └─────────────┘ └─────────────┘ └─────────────┘              │
│                                                                 │
│  Read-Only Repositories:                                        │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ AnalyticsCardRepository (→ Card Service)               │   │
│  │ AnalyticsEmployeeRepository (→ Employee Directory)     │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Policies:                                                      │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ AnalyticsAccessPolicy, DataExportPolicy                │   │
│  │ AnalyticsCalculationPolicy, ActivityLevelPolicy        │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ Read-Only Access
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                Card Recognition Context                          │
│                     (Card Service)                              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation Notes

### Architecture Decisions
- **Service-Oriented Design**: Analytics primarily uses domain services rather than aggregates
- **Read-Only Pattern**: No data modification, only computation and analysis
- **On-Demand Calculation**: No pre-aggregated data storage for real-time accuracy
- **Direct Repository Access**: Efficient data access without unnecessary abstraction layers

### Performance Considerations
- **Query Optimization**: Repository implementations should optimize for analytics queries
- **Caching Strategy**: Consider caching at infrastructure level for frequently accessed data
- **Large Dataset Handling**: Export service handles large datasets with streaming
- **Pagination**: All list operations support pagination for performance

### Scalability Patterns
- **Read Replicas**: Analytics queries can use read-only database replicas
- **Async Processing**: Large exports processed asynchronously
- **Result Caching**: Cache analytics results at application level with appropriate TTL
- **Query Batching**: Batch multiple analytics calculations when possible

### Future Extensibility
- **Advanced Analytics**: Foundation for machine learning and predictive analytics
- **Custom Reports**: Framework supports additional report types
- **Real-Time Dashboards**: Event-driven updates can be added
- **Data Warehouse Integration**: Export functionality supports data warehouse loading

---

## Validation Summary

This domain model supports all user stories from the Analytics Service Unit:
- ✅ HR Analytics Dashboard (US-5.1)
- ✅ Data Export functionality (US-5.2)
- ✅ Most Active Recognizers analysis (US-5.3)
- ✅ Team Recognition Patterns analysis (US-5.4)
- ✅ Company Values Distribution analysis (US-5.5)

All business rules for HR-only access, data privacy, and analytics calculations are properly enforced through domain services and policies.

The read-only integration pattern with Card Service ensures data consistency while maintaining bounded context independence.
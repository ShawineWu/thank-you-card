# Analytics Service Unit

## Unit Overview
**Description**: Backend service dedicated to HR analytics, reporting, and data export functionality. This unit provides comprehensive insights into recognition patterns, team behaviors, and company values distribution. It serves HR administrators with tools to analyze and export recognition data.

**Team**: Backend/Data Analytics Team
**Technology Stack**: RESTful API service with analytics/reporting capabilities
**Dependencies**: 
- Card Service Unit (reads card data for analysis)
- Employee Directory Service (for team/department information)

---

## Key Responsibilities
- HR analytics dashboard data aggregation
- Recognition pattern analysis
- Team and department analytics
- Company values distribution analysis
- Data export functionality (CSV)
- Most active recognizers identification
- Time-based trend analysis
- Report generation

---

## User Stories

### Epic 5: HR Analytics and Reporting

#### US-5.1: Access HR Analytics Dashboard
**As an** HR administrator  
**I want to** access a dedicated analytics dashboard in the web application  
**So that** I can view comprehensive recognition data and insights

**Acceptance Criteria:**
- HR admin can access analytics dashboard in web application
- Dashboard is accessible only to HR administrators
- Dashboard provides overview of key metrics
- Dashboard includes navigation to detailed reports
- Dashboard loads efficiently with large datasets
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for dashboard overview metrics
- Aggregate key metrics (total cards, active users, etc.)
- Implement role-based access control for HR admins
- Optimize queries for large datasets
- Return dashboard data in structured format

**Frontend Responsibilities (Web App Unit):**
- Implement HR analytics dashboard UI
- Display key metrics and visualizations
- Provide navigation to detailed reports
- Restrict access to HR admin role

---

#### US-5.2: Export Thank You Card Data
**As an** HR administrator  
**I want to** export thank-you card data to a file  
**So that** I can analyze recognition patterns and trends outside the system

**Acceptance Criteria:**
- HR admin can access export functionality from analytics dashboard in web application
- Export includes: sender name, sender email, recipient name, recipient email, recognition reason, selected values/credos, timestamp, and emoji reaction counts
- HR admin can select date range for export
- Export is available in CSV format
- Export file is downloadable
- Export result based on the filter conditions applied
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for data export
- Support date range parameter
- Support filter parameters (values, senders, receivers)
- Generate CSV file with all required fields
- Include emoji reaction counts in export
- Stream large exports efficiently
- Return downloadable file or download URL

---

#### US-5.3: View Most Active Recognizers
**As an** HR administrator  
**I want to** view a list of employees who send the most thank-you cards  
**So that** I can identify culture champions and encourage recognition behaviors

**Acceptance Criteria:**
- HR admin can view list of most active recognizers in analytics dashboard in web application
- List shows employee name, department, and total cards sent
- HR admin can select time period for analysis
- List displays at least top 10 active recognizers
- List updates based on selected time period
- Dashboard shows trends over time (increasing/decreasing activity)
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for most active recognizers
- Support time period parameter
- Calculate cards sent per employee
- Return top N recognizers (configurable, default 10)
- Include employee name, department, and card count
- Calculate trend data (comparison with previous period)
- Sort by card count in descending order

---

#### US-5.4: Analyze Team Recognition Patterns
**As an** HR administrator  
**I want to** view recognition patterns by team or department  
**So that** I can identify which teams have strong recognition cultures and which may need encouragement

**Acceptance Criteria:**
- HR admin can view recognition data grouped by team/department in analytics dashboard in web application
- Dashboard shows: total cards sent per team, total cards received per team, and average cards per employee
- HR admin can compare teams side-by-side
- Data can be visualized in charts or graphs
- HR admin can drill down into specific team details
- Time period selection applies to team analysis
- Dashboard identifies teams with low recognition activity
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for team recognition patterns
- Support time period parameter
- Aggregate cards sent per team/department
- Aggregate cards received per team/department
- Calculate average cards per employee per team
- Support team comparison queries
- Identify teams with low activity (below threshold)
- Return data suitable for visualization
- Support drill-down queries for specific teams

---

#### US-5.5: View Company Values Distribution
**As an** HR administrator  
**I want to** see which company values and credos are most frequently recognized  
**So that** I can understand which values are being actively demonstrated and which may need more emphasis

**Acceptance Criteria:**
- HR admin can view distribution of values/credos across all cards in analytics dashboard in web application
- Dashboard shows count and percentage for each value/credo
- Data is visualized in charts (bar chart, pie chart, etc.)
- HR admin can select time period for analysis
- Dashboard shows trends over time for each value/credo
- HR admin can identify underrepresented values
- This feature is available only in web application

**Priority:** Should Have

**Backend Responsibilities:**
- Provide API endpoint for values distribution
- Support time period parameter
- Aggregate count for each value/credo
- Calculate percentage distribution
- Calculate trend data over time
- Identify underrepresented values (below average)
- Return data suitable for chart visualization
- Support comparison between time periods

---

## API Endpoints

### Dashboard Overview
- `GET /api/analytics/dashboard` - Get dashboard overview metrics
  - Total cards sent (all time and selected period)
  - Total active users (senders and receivers)
  - Average cards per employee
  - Recognition activity trend

### Data Export
- `POST /api/analytics/export` - Export card data to CSV
  - Parameters: dateRange, filters (values, senders, receivers)
  - Returns: CSV file or download URL

### Most Active Recognizers
- `GET /api/analytics/recognizers/top` - Get most active recognizers
  - Parameters: timePeriod, limit (default 10)
  - Returns: List of employees with card counts and trends

### Team Analytics
- `GET /api/analytics/teams` - Get team recognition patterns
  - Parameters: timePeriod
  - Returns: Team-level aggregated data
- `GET /api/analytics/teams/{teamId}` - Get specific team details
  - Parameters: timePeriod
  - Returns: Detailed team recognition data

### Values Distribution
- `GET /api/analytics/values/distribution` - Get values/credos distribution
  - Parameters: timePeriod
  - Returns: Count and percentage for each value/credo
- `GET /api/analytics/values/trends` - Get values trends over time
  - Parameters: timePeriod, granularity (daily/weekly/monthly)
  - Returns: Time-series data for each value/credo

---

## Data Models

### Dashboard Overview Response
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
    "percentageChange": 20.0
  },
  "period": {
    "startDate": "2026-01-01",
    "endDate": "2026-01-27"
  }
}
```

### Most Active Recognizers Response
```json
{
  "recognizers": [
    {
      "employeeId": "string",
      "employeeName": "string",
      "department": "string",
      "cardsSent": 45,
      "trend": {
        "previous": 38,
        "percentageChange": 18.4
      }
    }
  ],
  "period": {
    "startDate": "2026-01-01",
    "endDate": "2026-01-27"
  }
}
```

### Team Recognition Patterns Response
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
      "activityLevel": "high | medium | low"
    }
  ],
  "period": {
    "startDate": "2026-01-01",
    "endDate": "2026-01-27"
  }
}
```

### Values Distribution Response
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
        "percentageChange": 20.8
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
        "percentageChange": 15.3
      }
    }
  ],
  "totalCards": 785,
  "period": {
    "startDate": "2026-01-01",
    "endDate": "2026-01-27"
  }
}
```

### Export Data CSV Format
```csv
Card ID,Sender Name,Sender Email,Recipient Name,Recipient Email,Recognition Reason,Values/Credos,Created At,Emoji Reactions Count
uuid-1,John Doe,john@company.com,Jane Smith,jane@company.com,Great teamwork on project X,"Make an Impact, Bias for Action",2026-01-15 10:30:00,5
```

---

## Analytics Calculations

### Key Metrics
1. **Total Cards**: Count of all cards in selected period
2. **Active Users**: Distinct count of senders and receivers
3. **Average Cards Per Employee**: Total cards / Total employees
4. **Activity Trend**: Percentage change compared to previous period
5. **Team Activity Level**: 
   - High: Above 75th percentile
   - Medium: Between 25th and 75th percentile
   - Low: Below 25th percentile

### Time Period Options
- Last 7 days
- Last 30 days (default)
- Last 90 days
- Last 6 months
- Last 12 months
- Custom date range
- All time

### Trend Calculations
- Compare current period with previous period of same duration
- Calculate percentage change: ((current - previous) / previous) * 100
- Indicate increasing/decreasing/stable trends

---

## Performance Considerations

### Optimization Strategies
- Implement caching for frequently accessed analytics
- Use database indexes on timestamp, sender, receiver, values fields
- Pre-aggregate common metrics (daily/weekly batch jobs)
- Implement pagination for large result sets
- Use database views for complex aggregations
- Consider read replicas for analytics queries

### Data Refresh
- Real-time data for current period metrics
- Cached data for historical trends (refresh hourly/daily)
- Export generation: Async processing for large datasets

---

## Access Control

### HR Administrator Role
- Only HR administrators can access analytics endpoints
- Implement role-based access control (RBAC)
- Validate HR admin role on every API request
- Log all analytics access for audit purposes

### Data Privacy
- Aggregate data only (no individual employee tracking beyond top lists)
- Export data includes employee information (authorized for HR)
- Comply with company data privacy policies

---

## Dependencies on Other Units

### Card Service Unit
- Read access to card data
- Read access to employee information
- Read access to reaction data
- No write operations to card data

### Employee Directory Service
- Read access to team/department structure
- Read access to employee count per team
- Read access to employee status (active/inactive)

### Web App Unit
- Provides UI for HR analytics dashboard
- Consumes all Analytics Service APIs
- Implements role-based access control in UI

---

## Integration Points

### Data Flow
1. Card Service Unit stores card data
2. Analytics Service Unit reads card data
3. Analytics Service Unit aggregates and calculates metrics
4. Web App Unit displays analytics to HR admins

### Real-time vs Batch Processing
- **Real-time**: Dashboard overview, current period metrics
- **Batch**: Historical trends, complex aggregations (run daily)
- **On-demand**: Data exports, drill-down queries

---

## Out of Scope (Future Enhancements)
- US-7.2: Audit Logging (basic logging included)
- Advanced predictive analytics
- Machine learning-based insights
- Automated report scheduling
- Email report delivery
- Custom report builder
- Advanced data visualization (beyond basic charts)
- Sentiment analysis on recognition reasons
- Correlation analysis between values and business metrics

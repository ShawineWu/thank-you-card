# Card Service Unit

## Unit Overview
**Description**: Backend service responsible for core card management operations including creation, storage, retrieval, filtering, searching, and emoji reactions. This unit serves as the central data layer for all card-related operations and is consumed by both the Web App and Teams Channel integration.

**Team**: Backend Team
**Technology Stack**: RESTful API service
**Dependencies**: 
- Employee Directory Service (for recipient validation and search)
- Database (for card storage)

---

## Key Responsibilities
- Card CRUD operations (Create, Read, Update, Delete)
- Card data persistence and retrieval
- Employee search and validation
- Card filtering and search functionality
- Emoji reaction management
- Data validation and business logic enforcement
- Serving card data to Web App and Teams Channel

---

## User Stories

### Epic 2: Card Creation

#### US-2.1: Access Card Creation Interface
**As an** employee  
**I want to** access the card creation interface from the web application  
**So that** I can create and send digital thank-you cards to colleagues

**Acceptance Criteria:**
- User can access the web application via link from Teams channel
- Card creation interface is easily accessible from web app home page
- User can initiate card creation with a clear call-to-action button
- Card creation form is intuitive and easy to use
- User receives confirmation when card creation is initiated

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for card creation initialization
- Return company values/credos list for selection

---

#### US-2.2: Select Card Recipients
**As an** employee sending a thank-you card  
**I want to** search and select a recipient from the company directory  
**So that** I can recognize the specific person who deserves appreciation

**Acceptance Criteria:**
- User can search for recipients by name, email and department
- Search results display relevant employee information (name, department, profile picture)
- User can select multi recipients per card
- Selected recipients are clearly displayed before submission
- System validates that the recipients are active employee

**Priority:** Must Have

**Backend Responsibilities:**
- Provide employee search API endpoint
- Support search by name, email, and department
- Validate recipient employee status (active)
- Return employee details (name, department, profile picture)

---

#### US-2.3: Specify Recognition Reason
**As an** employee sending a thank-you card  
**I want to** specify the reason for recognition  
**So that** the recipients and others understand what specific action or behavior is being appreciated

**Acceptance Criteria:**
- User can enter a recognition reason in a text field
- Text field has appropriate character limits (maximum)
- User receives guidance on what to include in the reason
- Reason field is mandatory before card submission
- System validates that reason is not empty or only whitespace

**Priority:** Must Have

**Backend Responsibilities:**
- Validate reason field is not empty or whitespace
- Enforce character limits on reason text
- Store reason text with card data

---

#### US-2.4: Align Recognition with Company Values
**As an** employee sending a thank-you card  
**I want to** select which company value(s) the recognition aligns with  
**So that** we can reinforce our organizational culture and values

**Acceptance Criteria:**
- User can view all company values in the selection interface:
  - **Values**: Make an Impact, Strive for Excellence, Stand Together, Be Open-Minded, Stay Grounded
  - **Credos**: Bias for Action, Customer Centric, Think Strategically, Deep Dive, Invent and Simplify, Earn Trust, Take Ownership, Challenge Disagree and Commit, Learn and Be Curious, Do More with Less
- User can select one or multiple values/credos (No more than 3)
- Selected values are clearly highlighted
- At least one value/credo must be selected before submission
- Values and credos are displayed with their descriptions for reference

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint to retrieve all company values and credos
- Validate that 1-3 values/credos are selected
- Store selected values/credos with card data

---

#### US-2.5: Preview and Submit Thank You Card
**As an** employee sending a thank-you card  
**I want to** preview the complete card before submitting  
**So that** I can ensure all information is correct and the card looks as intended

**Acceptance Criteria:**
- User can preview the card with all entered information
- Preview shows: recipient names, recognition reason, and selected values/credos
- User can edit any field from the preview screen
- User can submit the card from the preview screen
- User can cancel card creation at any time
- Upon successful submission, user receives confirmation
- Submitted card is immediately published to the company-wide feed

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint to create and store card
- Validate all required fields (recipients, reason, values)
- Generate unique card ID and timestamp
- Store sender information
- Return success confirmation with card ID
- Make card immediately available in feed queries

---

### Epic 3: Personal Card Views

#### US-3.1: View Personal Received Cards
**As an** employee who has received thank-you cards  
**I want to** view all cards that have been sent to me  
**So that** I can appreciate the recognition and keep track of positive feedback

**Acceptance Criteria:**
- User can access personal view of all cards received in web application
- Personal view is separate from the company-wide feed (Teams channel)
- Cards are displayed in reverse chronological order
- Each card shows complete information (sender, reason, values, timestamp)
- User can easily navigate between personal received/sent views
- Personal view shows a count of total cards received
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint to retrieve cards by recipient
- Return cards in reverse chronological order
- Include complete card information
- Return total count of cards received

---

#### US-3.2: View Personal Sent Cards
**As an** employee who has sent thank-you cards  
**I want to** view all cards that I have sent to others  
**So that** I can acknowledge who assisted me in which projects

**Acceptance Criteria:**
- User can access personal view of all cards sent in web application
- Personal view is separate from the company-wide feed (Teams channel)
- Cards are displayed in reverse chronological order
- Each card shows complete information (sender, receiver, reason, values, timestamp)
- User can easily navigate between personal received/sent views
- Personal view shows a count of total cards sent
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint to retrieve cards by sender
- Return cards in reverse chronological order
- Include complete card information
- Return total count of cards sent

---

#### US-3.3: View Personal Card Statistics
**As an** employee  
**I want to** see statistics about my thank-you card activity in the web application  
**So that** I can track my recognition giving and receiving

**Acceptance Criteria:**
- User can view personal statistics dashboard in web application
- Dashboard shows total cards sent count
- Dashboard shows total cards received count
- Dashboard shows most frequently used values when sending
- Dashboard shows most frequently received values
- Statistics update in real-time
- Statistics are viewable by the user only (private)

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for personal statistics
- Calculate total cards sent and received
- Aggregate most frequently used values (sending)
- Aggregate most frequently received values
- Return statistics in structured format

---

### Epic 4: Card Interaction

#### US-4.1: React to Thank You Cards with Emojis
**As an** employee viewing a thank-you card  
**I want to** react to cards using emoji reactions in both Teams channel and web application  
**So that** I can show support and appreciation for the recognition

**Acceptance Criteria:**
- User can add emoji reactions to any thank-you card in Teams channel or web app
- System provides a predefined set of emoji options
- User can only select one emoji on the same card
- User can remove their own emoji reactions
- Card displays count of each emoji reaction type
- User can see who reacted with each emoji
- Reactions are visible to all employees viewing the card
- Reactions sync between Teams channel and web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint to add emoji reaction
- Provide API endpoint to remove emoji reaction
- Enforce one emoji per user per card rule
- Store emoji reactions with user ID and timestamp
- Return emoji reaction counts and user lists
- Ensure reactions are immediately available across all platforms

---

#### US-4.2: Filter Thank You Cards
**As an** employee or HR admin  
**I want to** filter thank-you cards by specific criteria in the web application  
**So that** I can find cards that match specific values, people, or time periods

**Acceptance Criteria:**
- User can access filter options in the web application
- User can select one or multiple values/credos to filter
- User can select one or multiple senders to filter
- User can select one or multiple receivers to filter
- User can select period of create time to filter
- List updates to show only cards matching selected conditions
- User can clear filters to return to full view
- Active filters are clearly displayed
- Filter state persists during the session
- This feature is available only in web application

**Priority:** Should Have

**Backend Responsibilities:**
- Provide API endpoint with filter parameters
- Support filtering by values/credos (multiple)
- Support filtering by sender (multiple)
- Support filtering by receiver (multiple)
- Support filtering by date range
- Return filtered results efficiently
- Support combination of multiple filters

---

#### US-4.3: Search Thank You Cards
**As an** employee  
**I want to** search for specific thank-you cards in the web application  
**So that** I can find cards related to specific people, values, or keywords

**Acceptance Criteria:**
- User can access search functionality in the web application
- User can search by recipient name, sender name, or keywords in reason
- Search results display matching cards
- Search is case-insensitive
- User can clear search to return to full view
- Search works in combination with filters
- This feature is available only in web application

**Priority:** Should Have

**Backend Responsibilities:**
- Provide API endpoint with search parameter
- Support search by recipient name
- Support search by sender name
- Support search by keywords in reason text
- Implement case-insensitive search
- Support search combined with filters
- Return matching cards efficiently

---

### Epic 1: Teams Channel Integration (Backend Support)

#### US-1.1: Access Thank You Card Feed from Teams Channel
**As an** employee  
**I want to** view the company-wide thank-you card feed directly in a Teams channel  
**So that** I can see recognition happening across the company without leaving Teams

**Acceptance Criteria:**
- Dedicated Teams channel is accessible to all employees
- Channel displays company-wide thank-you card feed
- Feed shows cards in reverse chronological order (newest first)
- Each card displays: sender name, recipient name(s), recognition reason, values/credos, and timestamp
- Feed updates automatically when new cards are posted
- Feed has infinite scroll or pagination for browsing older cards
- Channel interface is consistent with Teams design patterns

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint to retrieve company-wide card feed
- Support pagination for feed retrieval
- Return cards in reverse chronological order
- Include all card details (sender, recipients, reason, values, timestamp, reactions)
- Support real-time updates (webhook or polling)

---

#### US-1.2: Access Web Application Link from Teams Channel
**As an** employee  
**I want to** access a link to the web application from the Teams channel  
**So that** I can create cards, view my personal statistics, and access additional features

**Acceptance Criteria:**
- Web application link is prominently displayed in Teams channel
- Link is easily accessible (pinned message or channel tab)
- Clicking link opens web application in browser
- Link includes clear description of web app features
- All employees can access the link

**Priority:** Must Have

**Backend Responsibilities:**
- No backend responsibilities (Teams channel configuration)

---

#### US-1.4: Share Thank You Cards to Other Teams Channels
**As an** employee  
**I want to** share a thank-you card to other Teams channels  
**So that** I can amplify recognition within specific team contexts

**Acceptance Criteria:**
- User can share any thank-you card to other Teams channels they have access to
- Share action is available from card view in Teams channel
- Shared card displays in target channel with link to original
- Shared card maintains all original information
- User can add context when sharing
- Share action is tracked for analytics

**Priority:** Could Have

**Backend Responsibilities:**
- Provide API endpoint to retrieve card details by ID
- Track share actions for analytics
- Return complete card information for sharing

---

### Epic 3: Personal Card Views (Additional)

#### US-3.4: View Top 10 Recognized Employees
**As an** employee  
**I want to** view a list of the top 10 employees by thank-you cards received count in the web application  
**So that** I can identify employees who are consistently demonstrating company values

**Acceptance Criteria:**
- Employee can access web application to view the top 10 list
- List displays top 10 employees by number of cards received in descending order
- List shows employee name, department, and total card count
- User can select time period for analysis (default: recent 30 calendar days)
- List updates based on selected time period
- Ties in ranking are handled appropriately (employees with same count ranked equally)
- Multiple employees can hold the same rank
- This feature is available only in web application

**Priority:** Must Have

**Backend Responsibilities:**
- Provide API endpoint for top 10 recognized employees
- Support time period parameter (default: 30 days)
- Calculate card counts per employee
- Return top 10 in descending order
- Include employee name, department, and card count
- Handle ranking ties appropriately

---

### Epic 6: User Experience Enhancements

#### US-6.1: Celebrate Recognition Milestones
**As an** employee  
**I want to** receive recognition when I reach card milestones  
**So that** I feel appreciated for my participation in the recognition culture

**Acceptance Criteria:**
- System tracks milestones (e.g., 1st card sent, 10th card received, 50 cards sent)
- User receives notification when reaching a milestone
- Milestone achievements are displayed on user profile
- Milestones are celebratory and positive in tone
- System supports both sending and receiving milestones

**Priority:** Could Have

**Backend Responsibilities:**
- Track card counts for each user (sent and received)
- Detect milestone achievements
- Provide API endpoint to retrieve user milestones
- Trigger milestone events for notification service

---

## API Endpoints (Summary)

### Card Management
- `POST /api/cards` - Create new card
- `GET /api/cards` - Get company-wide feed (with pagination)
- `GET /api/cards/{id}` - Get specific card details
- `GET /api/cards/received` - Get cards received by user
- `GET /api/cards/sent` - Get cards sent by user

### Employee Search
- `GET /api/employees/search` - Search employees by name, email, department

### Values & Credos
- `GET /api/values` - Get all company values and credos

### Reactions
- `POST /api/cards/{id}/reactions` - Add emoji reaction
- `DELETE /api/cards/{id}/reactions` - Remove emoji reaction

### Statistics
- `GET /api/statistics/personal` - Get personal card statistics
- `GET /api/statistics/top10` - Get top 10 recognized employees

### Filtering & Search
- `GET /api/cards/filter` - Filter cards by criteria
- `GET /api/cards/search` - Search cards by keywords

### Milestones
- `GET /api/milestones/{userId}` - Get user milestones

---

## Data Models

### Card
```json
{
  "id": "string (UUID)",
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
  "values": ["string"],
  "createdAt": "timestamp",
  "reactions": [
    {
      "emoji": "string",
      "userId": "string",
      "userName": "string",
      "timestamp": "timestamp"
    }
  ]
}
```

### Employee
```json
{
  "id": "string",
  "name": "string",
  "email": "string",
  "department": "string",
  "profilePicture": "string (URL)",
  "isActive": "boolean"
}
```

### Value/Credo
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "type": "value | credo"
}
```

---

## Dependencies on Other Units
- **Web App Unit**: Consumes all Card Service APIs
- **Analytics Service Unit**: Reads card data for analytics and reporting
- **Employee Directory Service**: External dependency for employee search and validation
- **Teams Channel**: Consumes Card Service APIs for feed display

---

## Out of Scope (Future Enhancements)
- US-7.1: User Authentication (handled by Teams/Azure AD)
- US-7.2: Audit Logging
- US-7.3: Error Handling (basic error handling included)
- US-7.4: Card Approval Workflow
- US-7.5: Custom Card Templates
- US-7.6: Seasonal Card Templates
- US-7.7: Recipient Suggestions
- US-1.3: Notifications (handled by Teams Channel existing functionality)

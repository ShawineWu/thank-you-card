# Web App Unit

## Unit Overview
**Description**: Frontend web application that provides the user interface for creating thank-you cards, viewing personal card history, accessing statistics, and interacting with cards. This unit serves as the primary interface for employees to engage with the thank-you card system beyond the Teams channel feed.

**Team**: Frontend Team
**Technology Stack**: Web application (React/Angular/Vue or similar)
**Dependencies**: 
- Card Service Unit (for all card operations)
- Analytics Service Unit (for statistics and top 10 list)

---

## Key Responsibilities
- User interface for card creation workflow
- Personal card views (received and sent)
- Personal statistics dashboard
- Top 10 recognized employees display
- Card filtering and search interface
- Emoji reaction interface
- Responsive design and user experience
- Integration with Card Service and Analytics Service APIs

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

**Frontend Responsibilities:**
- Design and implement card creation interface
- Provide prominent "Create Card" call-to-action
- Display intuitive multi-step or single-page form
- Show loading states during API calls
- Display confirmation messages

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

**Frontend Responsibilities:**
- Implement employee search interface with autocomplete
- Display search results with employee details
- Support multi-recipient selection
- Show selected recipients with ability to remove
- Call Card Service API for employee search
- Display validation errors for inactive employees

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

**Frontend Responsibilities:**
- Implement text area for reason input
- Display character count and limit
- Show helpful placeholder text or guidance
- Validate reason field before submission
- Display validation errors

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

**Frontend Responsibilities:**
- Fetch and display all values/credos from Card Service API
- Implement selection interface (checkboxes, chips, or cards)
- Show value/credo descriptions on hover or expand
- Enforce 1-3 selection limit
- Highlight selected values/credos
- Validate at least one is selected before submission

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

**Frontend Responsibilities:**
- Implement preview screen with card visualization
- Display all entered information in card format
- Provide "Edit" and "Submit" buttons
- Handle submission via Card Service API
- Show success confirmation message
- Redirect to appropriate view after submission
- Handle submission errors gracefully

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

**Frontend Responsibilities:**
- Implement "Received Cards" view
- Fetch cards from Card Service API
- Display cards in reverse chronological order
- Show complete card information
- Display total count of received cards
- Provide navigation between received/sent views
- Implement pagination or infinite scroll

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

**Frontend Responsibilities:**
- Implement "Sent Cards" view
- Fetch cards from Card Service API
- Display cards in reverse chronological order
- Show complete card information
- Display total count of sent cards
- Provide navigation between received/sent views
- Implement pagination or infinite scroll

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

**Frontend Responsibilities:**
- Implement personal statistics dashboard
- Fetch statistics from Card Service API
- Display cards sent and received counts
- Visualize most frequently used values (charts/graphs)
- Visualize most frequently received values
- Ensure statistics are private to the user
- Refresh statistics on page load

---

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

**Frontend Responsibilities:**
- Implement Top 10 recognized employees view
- Fetch data from Card Service API
- Display employees in ranked list format
- Show employee name, department, and card count
- Provide time period selector (default: 30 days)
- Update list when time period changes
- Handle ranking ties in display
- Make view accessible from main navigation

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

**Frontend Responsibilities:**
- Implement emoji reaction picker interface
- Display available emoji options
- Call Card Service API to add/remove reactions
- Show reaction counts on cards
- Display list of users who reacted (on hover/click)
- Enforce one emoji per user per card
- Update UI immediately after reaction
- Sync reactions across views

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

**Frontend Responsibilities:**
- Implement filter panel/sidebar
- Provide multi-select for values/credos
- Provide multi-select for senders
- Provide multi-select for receivers
- Provide date range picker
- Call Card Service API with filter parameters
- Display active filters with clear indicators
- Provide "Clear All Filters" button
- Persist filter state in session
- Update card list based on filters

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

**Frontend Responsibilities:**
- Implement search bar interface
- Call Card Service API with search query
- Display search results
- Show "No results" message when appropriate
- Provide clear search button
- Support search combined with filters
- Implement debouncing for search input
- Show search query in UI

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

**Frontend Responsibilities:**
- Display milestone achievements on user profile
- Show milestone notifications/badges
- Fetch milestones from Card Service API
- Design celebratory UI for milestones
- Display milestone progress indicators

---

## Navigation Structure

### Main Navigation
- Home / Dashboard
- Create Card
- My Cards
  - Received
  - Sent
- My Statistics
- Top 10 Recognized
- (HR Admin only) Analytics Dashboard → Links to Analytics Service Unit

### Card Views
- Company-wide feed (accessible via Teams Channel)
- Personal received cards
- Personal sent cards
- Filtered/searched card results

---

## UI Components

### Core Components
- Card Display Component (shows card with all details)
- Card Creation Form (multi-step or single page)
- Employee Search Component (autocomplete)
- Value/Credo Selector Component
- Emoji Reaction Picker Component
- Filter Panel Component
- Search Bar Component
- Statistics Dashboard Component
- Top 10 List Component
- Navigation Component

### Shared Components
- Loading Spinner
- Error Message Display
- Success Notification
- Confirmation Dialog
- Pagination Component

---

## API Integration Points

### Card Service Unit APIs
- `POST /api/cards` - Create card
- `GET /api/cards/received` - Get received cards
- `GET /api/cards/sent` - Get sent cards
- `GET /api/employees/search` - Search employees
- `GET /api/values` - Get values/credos
- `POST /api/cards/{id}/reactions` - Add reaction
- `DELETE /api/cards/{id}/reactions` - Remove reaction
- `GET /api/cards/filter` - Filter cards
- `GET /api/cards/search` - Search cards
- `GET /api/statistics/personal` - Get personal statistics
- `GET /api/statistics/top10` - Get top 10 employees
- `GET /api/milestones/{userId}` - Get milestones

### Analytics Service Unit APIs
- (HR Admin only) Analytics dashboard endpoints

---

## Dependencies on Other Units
- **Card Service Unit**: Primary dependency for all card operations
- **Analytics Service Unit**: For HR admin analytics features (separate view)
- **Authentication Service**: For user authentication (Azure AD/Teams SSO)

---

## Out of Scope (Future Enhancements)
- US-7.1: User Authentication UI (handled by Teams/Azure AD SSO)
- US-7.3: Advanced Error Handling (basic error handling included)
- US-7.4: Card Approval Workflow UI
- US-7.5: Custom Card Templates
- US-7.6: Seasonal Card Templates
- US-7.7: Recipient Suggestions UI
- US-1.3: Notification UI (handled by Teams)
- US-1.4: Share to Teams Channels (Teams-specific feature)

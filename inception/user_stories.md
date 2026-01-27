# Thank You Card Feature - User Stories

## Document Information
- **Feature**: Internal Employee Thank You Card System
- **Version**: 1.0
- **Date**: January 27, 2026
- **Status**: Draft

## Overview
This document contains user stories for an internal employee recognition system that enables employees to send public digital thank-you cards to colleagues, displays them in a company-wide feed, and provides HR analytics capabilities.

---

## User Personas

1. **Employee**: Any employee who wants to send or receive thank-you cards
2. **HR Administrator**: HR team member responsible for analytics, reporting, and dashboard viewing

---

## Platform Overview

**Teams Channel:**
- View company-wide thank-you card feed
- Access link to web application
- Receive real-time card updates
- Available to all employees

**Web Application:**
- Create and send thank-you cards
- View personal received and sent cards
- View Top 10 recognized employees list
- View personal card statistics (cards received/sent count)
- HR Admin: Access analytics dashboard and reporting tools
- HR Admin: Export data and view team recognition patterns

---

## Epic 1: Teams Channel Integration

### US-1.1: Access Thank You Card Feed from Teams Channel
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

---

### US-1.2: Access Web Application Link from Teams Channel
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

---

### US-1.3: Receive Notifications for Received Cards
**As an** employee  
**I want to** receive a Teams notification when someone sends me a thank-you card  
**So that** I am immediately aware of the recognition

**Acceptance Criteria:**
- User receives Teams notification when a card is sent to them
- Notification includes sender name and brief preview
- Clicking notification navigates to the card in Teams channel or web app
- User can configure notification preferences
- Notifications respect user's Teams notification settings
- Notification is sent in real-time

**Priority:** Should Have

---

### US-1.4: Share Thank You Cards to Other Teams Channels
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

---

## Epic 2: Web App - Card Creation

### US-2.1: Access Card Creation Interface
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

---

### US-2.2: Select Card Recipients
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

---

### US-2.3: Specify Recognition Reason
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

---

### US-2.4: Align Recognition with Company Values
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

---

### US-2.5: Preview and Submit Thank You Card
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

---

## Epic 3: Web App - Personal Card Views

### US-3.1: View Personal Received Cards
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

---

### US-3.2: View Personal Sent Cards
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

---

### US-3.3: View Personal Card Statistics
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

---

### US-3.4: View Top 10 Recognized Employees
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

---

## Epic 4: Web App - Card Interaction

### US-4.1: React to Thank You Cards with Emojis
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

---

### US-4.2: Filter Thank You Cards
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

---

### US-4.3: Search Thank You Cards
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

---

## Epic 5: Web App - HR Analytics and Reporting

### US-5.1: Access HR Analytics Dashboard
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

---

### US-5.2: Export Thank You Card Data
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

---

### US-5.3: View Most Active Recognizers
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

---

### US-5.4: Analyze Team Recognition Patterns
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

---

### US-5.5: View Company Values Distribution
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

---

## Epic 6: User Experience Enhancements

### US-6.1: Celebrate Recognition Milestones
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

---

## Epic 7: Future Enhancements (Out of MVP Scope)

### US-7.1: Authenticate Users
**As a** system  
**I want to** authenticate users through Teams/Azure AD  
**So that** only authorized employees can access the thank-you card feature

**Acceptance Criteria:**
- System integrates with Teams/Azure AD authentication
- Users are automatically authenticated when accessing from Teams
- User identity is verified before any action
- Session management follows security best practices
- Unauthorized users cannot access the system

**Priority:** Future Enhancement

---

### US-7.2: Audit Card Activities
**As a** system administrator  
**I want to** maintain audit logs of all card-related activities  
**So that** we can track system usage and investigate issues if needed

**Acceptance Criteria:**
- System logs all card creation events with timestamp and user ID
- System logs all data export events by HR admins
- System logs all emoji reactions
- Logs include sufficient detail for troubleshooting
- Logs are stored securely and retained per company policy
- Logs are accessible only to authorized administrators

**Priority:** Future Enhancement

---

### US-7.3: Handle System Errors Gracefully
**As an** employee  
**I want to** receive clear error messages when something goes wrong  
**So that** I understand what happened and what to do next

**Acceptance Criteria:**
- System displays user-friendly error messages
- Error messages provide actionable guidance
- Technical errors are logged for administrator review
- System prevents data loss during errors
- Users can retry failed actions
- Critical errors are escalated to support team

**Priority:** Future Enhancement

---

### US-7.4: Implement Card Approval Workflow
**As an** HR administrator  
**I want to** review and approve thank-you cards before they are published  
**So that** we can ensure content meets company standards

**Acceptance Criteria:**
- Cards enter pending state after submission
- HR admin receives notification of pending cards
- HR admin can approve or reject cards with reason
- Sender is notified of approval/rejection status
- Approved cards are published to feed
- Rejected cards are returned to sender for revision

**Priority:** Future Enhancement

---

### US-7.5: Create Custom Card Templates
**As an** employee  
**I want to** choose from different card templates or designs  
**So that** I can personalize the visual appearance of my recognition

**Acceptance Criteria:**
- System provides multiple card template options
- User can preview templates before selection
- Templates are visually distinct and professional
- Selected template is applied to the card
- All templates display required information clearly

**Priority:** Future Enhancement

---

### US-7.6: Add Seasonal Card Templates
**As an** employee  
**I want to** choose from seasonal or themed card templates  
**So that** I can make recognition more festive and engaging

**Acceptance Criteria:**
- System provides seasonal templates (holidays, celebrations)
- Templates include different color schemes and designs
- User can preview templates before selection
- Templates maintain professional appearance
- All templates display required information clearly

**Priority:** Future Enhancement

---

### US-7.7: Suggest Recipients Based on Collaboration
**As an** employee creating a thank-you card  
**I want to** see suggested recipients based on my recent collaborations  
**So that** I can easily recognize people I work with frequently

**Acceptance Criteria:**
- System suggests recipients based on Teams chat/meeting history
- Suggestions appear when user starts creating a card
- User can choose from suggestions or search manually
- Suggestions are relevant and recent
- User can dismiss suggestions
- Suggestions respect privacy settings

**Priority:** Future Enhancement

---

## Summary

**Total User Stories:** 27
- **Must Have:** 15 stories
- **Should Have:** 5 stories
- **Could Have:** 2 stories
- **Future Enhancements:** 7 stories

**Epics:**
1. Teams Channel Integration (4 stories - 2 Must Have, 1 Should Have, 1 Could Have)
2. Web App - Card Creation (5 stories - all Must Have)
3. Web App - Personal Card Views (4 stories - all Must Have)
4. Web App - Card Interaction (3 stories - 1 Must Have, 2 Should Have)
5. Web App - HR Analytics and Reporting (5 stories - 4 Must Have, 1 Should Have)
6. User Experience Enhancements (1 story - Could Have)
7. Future Enhancements (7 stories - System Administration, Security, and Advanced Features)

---

## Platform Feature Distribution

**Teams Channel Features:**
- View company-wide thank-you card feed (real-time updates)
- Access link to web application
- Receive notifications for received cards
- React to cards with emojis (synced with web app)
- Share cards to other Teams channels

**Web Application Features:**

*For All Employees:*
- Create and send thank-you cards
- View personal received cards with count
- View personal sent cards with count
- View personal card statistics (sent/received counts, frequently used/received values)
- View Top 10 recognized employees list
- React to cards with emojis (synced with Teams)
- Filter cards by values/credos, sender, receiver, time period
- Search cards by name or keywords

*For HR Administrators:*
- Access dedicated analytics dashboard
- Export card data to CSV
- View most active recognizers
- Analyze team recognition patterns
- View company values distribution
- All filtering and search capabilities

---

## Notes

**MVP Scope:**
- All cards are publicly visible to the entire company
- Multiple recipients supported per card (max not specified)
- Maximum 3 values/credos can be selected per card
- Users can only select one emoji reaction per card
- No personalized message feature in MVP
- No limitations on card sending frequency or volume in MVP
- Emoji reactions sync between Teams channel and web application

**Platform Architecture:**
- **Teams Channel**: Primary interface for viewing company-wide feed and accessing web app
- **Web Application**: Full-featured platform for card creation, personal views, statistics, and HR analytics
- Seamless integration between both platforms

**Default Settings:**
- Top 10 list default time period: Recent 30 calendar days
- Ranking ties: Multiple employees can hold the same rank

**Future Enhancements:**
- Authentication and security features (Azure AD integration)
- Audit logging capabilities
- Card approval workflow
- Custom and seasonal card templates
- Recipient suggestions based on collaboration
- Advanced error handling

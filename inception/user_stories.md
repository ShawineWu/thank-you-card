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

1. **Employee (Sender)**: Any employee who wants to recognize a colleague
2. **Employee (Recipient)**: Any employee who receives a thank-you card
3. **Employee (Viewer)**: Any employee browsing the company-wide feed
4. **HR Administrator**: HR team member responsible for analytics and reporting

---

## Epic 1: Send Thank You Cards

### US-1.1: Create a Thank You Card
**As an** employee  
**I want to** create a digital thank-you card for a colleague  
**So that** I can publicly recognize their contributions and efforts

**Acceptance Criteria:**
- User can access the card creation interface from web page which the link will be displayed in teams channel
- User can initiate card creation with a clear call-to-action
- Card creation form is intuitive and easy to use
- User receives confirmation when card creation is initiated

**Priority:** Must Have

---

### US-1.2: Select Card Recipient
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

### US-1.3: Specify Recognition Reason
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

### US-1.4: Align Recognition with Company Values
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

### US-1.5: Preview and Submit Thank You Card
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

## Epic 2: View and Interact with Thank You Cards

### US-2.1: View Company-Wide Thank You Card Feed
**As an** employee  
**I want to** view a feed of all thank-you cards sent across the company  
**So that** I can see how colleagues are being recognized and feel inspired by the positive culture

**Acceptance Criteria:**
- User can access the thank-you card feed from Teams channel
- Feed displays all thank-you cards in reverse chronological order (newest first)
- Each card displays: sender name, recipient name, recognition reason, values/credos, and timestamp
- Feed is accessible to all employees in the company
- Feed updates automatically when new cards are posted
- Feed has infinite scroll or pagination for browsing older cards

**Priority:** Must Have

---

### US-2.2: View Personal Received Cards
**As an** employee who has received thank-you cards  
**I want to** view all cards that have been sent to me  
**So that** I can appreciate the recognition and keep track of positive feedback

**Acceptance Criteria:**
- User can access a personal view of all cards received
- Personal view is separate from the company-wide feed
- Cards are displayed in reverse chronological order
- Each card shows complete information (sender, reason, values, timestamp)
- User can easily navigate between personal view and company-wide feed
- Personal view shows a count of total cards received
- This feature is managed in web page but not the teams channel

**Priority:** Must Have

---

### US-2.3: View Personal Sent Cards
**As an** employee who has sent thank-you cards  
**I want to** view all cards that I have sent to others  
**So that** I can acknowledge who assisted me in which projects

**Acceptance Criteria:**
- User can access a personal view of all cards sent
- Personal view is separate from the company-wide feed
- Cards are displayed in reverse chronological order
- Each card shows complete information (sender, receiver, reason, values, timestamp)
- User can easily navigate between personal view and company-wide feed
- Personal view shows a count of total cards sent
- This feature is managed in web page but not the teams channel

**Priority:** Must Have

---

### US-2.4: React to Thank You Cards with Emojis
**As an** employee viewing a thank-you card  
**I want to** react to cards using emoji reactions  
**So that** I can show support and appreciation for the recognition

**Acceptance Criteria:**
- User can add emoji reactions to any thank-you card
- System provides a predefined set of emoji options
- User can only select one emojis on the same card
- User can remove their own emoji reactions
- Card displays count of each emoji reaction type
- User can see who reacted with each emoji
- Reactions are visible to all employees viewing the card

**Priority:** Must Have

---

### US-2.5: Filter Thank You Cards by Values/Credos
**As an** HR Admin  
**I want to** filter thank-you cards by specific company values or credos / Sender / Receiver / Create time
**So that** I can see examples of how specific values are being demonstrated

**Acceptance Criteria:**
- User can access filter options in the web page
- User can select one or multiple values/credos to filter
- User can select one or multiple sender to filter
- User can select one or multiple receiver to filter
- User can select period of create time to filter
- List updates to show only cards tagged with selected conditions
- User can clear filters to return
- Active filters are clearly displayed
- Filter state persists during the session
- This feature is managed in web page but not the teams channel

**Priority:** Should Have

---

### US-2.6: Search Thank You Cards
**As an** employee  
**I want to** search for specific thank-you cards  
**So that** I can find cards related to specific people, values, or keywords

**Acceptance Criteria:**
- User can access search functionality in the feed
- User can search by recipient name, sender name, or keywords in reason
- Search results display matching cards
- Search is case-insensitive
- User can clear search to return to full feed
- Search works in combination with filters

**Priority:** Should Have

---

## Epic 3: HR Analytics and Reporting

### US-3.1: Export Thank You Card Data
**As an** HR administrator  
**I want to** export thank-you card data to a file  
**So that** I can analyze recognition patterns and trends outside the system

**Acceptance Criteria:**
- HR admin can access export functionality from admin panel
- Export includes: sender name, sender email, recipient name, recipient email, recognition reason, selected values/credos, timestamp, and emoji reaction counts
- HR admin can select date range for export
- Export is available in CSV format
- Export file is downloadable
- Export result based on the filter condition in US-2.4: Filter Thank You Cards by Values/Credos

**Priority:** Must Have

---

### US-3.2: View Top 10 Thank-you Cards Received Count
**As an** employee  
**I want to** view a list of the top 10 employees by thank-you cards received count  
**So that** I can identify employees who are consistently demonstrating company values

**Acceptance Criteria:**
- Employee can access web page to view the top 10 list
- List displays top 10 employees by number of cards received in descending order
- List shows employee name, department, and total card count
- Can select time period for analysis (default will be recent 30 calendar days)
- List updates based on selected time period
- Ties in ranking are handled appropriately

**Priority:** Must Have

---

### US-3.3: View Most Active Recognizers
**As an** HR administrator  
**I want to** view a list of employees who send the most thank-you cards  
**So that** I can identify culture champions and encourage recognition behaviors

**Acceptance Criteria:**
- HR admin can view list of most active recognizers in analytics dashboard
- List shows employee name, department, and total cards sent
- HR admin can select time period for analysis
- List displays at least top 10 active recognizers
- List updates based on selected time period
- Dashboard shows trends over time (increasing/decreasing activity)

**Priority:** Must Have

---

### US-3.4: Analyze Team Recognition Patterns
**As an** HR administrator  
**I want to** view recognition patterns by team or department  
**So that** I can identify which teams have strong recognition cultures and which may need encouragement

**Acceptance Criteria:**
- HR admin can view recognition data grouped by team/department
- Dashboard shows: total cards sent per team, total cards received per team, and average cards per employee
- HR admin can compare teams side-by-side
- Data can be visualized in charts or graphs
- HR admin can drill down into specific team details
- Time period selection applies to team analysis
- Dashboard identifies teams with low recognition activity

**Priority:** Must Have

---

### US-3.5: View Company Values Distribution
**As an** HR administrator  
**I want to** see which company values and credos are most frequently recognized  
**So that** I can understand which values are being actively demonstrated and which may need more emphasis

**Acceptance Criteria:**
- HR admin can view distribution of values/credos across all cards
- Dashboard shows count and percentage for each value/credo
- Data is visualized in charts (bar chart, pie chart, etc.)
- HR admin can select time period for analysis
- Dashboard shows trends over time for each value/credo
- HR admin can identify underrepresented values

**Priority:** Should Have

---

## Epic 4: Teams Integration

### US-4.1: Access Thank You Cards from Teams Channel
**As an** employee  
**I want to** access the thank-you card feature directly from a Teams channel  
**So that** I can easily view cards without leaving my workflow

**Acceptance Criteria:**
- Thank-you card feature is integrated into a dedicated Teams channel
- Channel is accessible to all employees
- Users can view the company-wide feed within the Teams channel
- Channel displays real-time updates when new cards are posted
- Channel interface is consistent with Teams design patterns

**Priority:** Must Have

---

### US-4.2: Receive Notifications for Received Cards
**As an** employee  
**I want to** receive a notification when someone sends me a thank-you card  
**So that** I am immediately aware of the recognition

**Acceptance Criteria:**
- User receives Teams notification when a card is sent to them
- Notification includes sender name and brief preview
- Clicking notification navigates to the specific card
- User can configure notification preferences
- Notifications respect user's Teams notification settings
- Notification is sent in real-time

**Priority:** Should Have

---

### US-4.3: Share Thank You Cards to Other Teams Channels
**As an** employee  
**I want to** share a thank-you card to other Teams channels  
**So that** I can amplify recognition within specific team contexts

**Acceptance Criteria:**
- User can share any thank-you card to other Teams channels they have access to
- Share action is available from card view
- Shared card displays in target channel with link to original
- Shared card maintains all original information
- User can add context when sharing
- Share action is tracked for analytics

**Priority:** Could Have

---

## Epic 5: User Experience Enhancements

### US-5.1: View Card Statistics on Profile
**As an** employee  
**I want to** see statistics about my thank-you card activity on my profile  
**So that** I can track my recognition giving and receiving

**Acceptance Criteria:**
- User profile shows total cards sent
- User profile shows total cards received
- User profile shows most frequently used values when sending
- User profile shows most frequently received values
- Statistics update in real-time
- Profile is viewable by the user only (private statistics)

**Priority:** Should Have

---

### US-5.2: Celebrate Recognition Milestones
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

## Epic 6: Future Enhancements (Out of MVP Scope)

### US-6.1: Authenticate Users
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

### US-6.2: Audit Card Activities
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

### US-6.3: Handle System Errors Gracefully
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

### US-6.4: Implement Card Approval Workflow
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

### US-6.5: Create Custom Card Templates
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

### US-6.6: Add Seasonal Card Templates
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

### US-6.7: Suggest Recipients Based on Collaboration
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

**Total User Stories:** 26
- **Must Have:** 14 stories
- **Should Have:** 5 stories
- **Could Have:** 2 stories
- **Future Enhancements:** 7 stories

**Epics:**
1. Send Thank You Cards (5 stories - all Must Have)
2. View and Interact with Thank You Cards (6 stories - 4 Must Have, 2 Should Have)
3. HR Analytics and Reporting (5 stories - 4 Must Have, 1 Should Have)
4. Teams Integration (3 stories - 1 Must Have, 1 Should Have, 1 Could Have)
5. User Experience Enhancements (2 stories - 1 Should Have, 1 Could Have)
6. Future Enhancements (7 stories - System Administration, Security, and Advanced Features)

---

## Notes

**MVP Scope:**
- All cards are publicly visible to the entire company (no privacy options in MVP)
- Multiple recipients supported per card (US-1.2)
- Maximum 3 values/credos can be selected per card (US-1.4)
- Users can only select one emoji reaction per card (US-2.4)
- No personalized message feature in MVP
- No limitations on card sending frequency or volume in MVP

**Platform Architecture:**
- Web application for card creation, personal views, filtering, and analytics
- Teams channel integration for viewing company-wide feed and real-time updates
- Link to web application displayed in Teams channel for full functionality

**Future Enhancements:**
- Authentication and security features (Azure AD integration)
- Audit logging capabilities
- Card approval workflow
- Custom and seasonal card templates
- Recipient suggestions based on collaboration
- Advanced error handling

**Analytics Focus:**
- Top 10 recognized employees (default: last 30 days)
- Most active recognizers
- Team recognition patterns
- Company values distribution

# Plan: Thank You Card Feature - User Stories Development

## Overview
Creating user stories for an internal employee Thank You card feature in Teams application that enables recognition, public display, and HR analytics.

## Steps

- [x] 1. Analyze the feature requirements and identify key user personas
  - Employees (senders)
  - Employees (recipients)
  - HR Administrators
  - General viewers (company-wide feed)

- [x] 2. Define core user stories for sending thank-you cards
  - Card creation workflow
  - Recipient selection
  - Recognition reason specification
  - Company values alignment
  - Optional personalized messages

- [x] 3. Define user stories for viewing and interacting with thank-you cards
  - Company-wide feed display
  - Recipient's personal view of received cards
  - Response/interaction capabilities

- [x] 4. Define user stories for HR administrator analytics and reporting
  - Data export functionality
  - Recognition behavior analysis
  - Core value implementation tracking

- [x] 5. Create the /inception/ directory structure

- [x] 6. Write comprehensive user stories to user_stories.md
  - Follow standard user story format: "As a [persona], I want to [action], so that [benefit]"
  - Include acceptance criteria for each story
  - Prioritize stories (Must Have, Should Have, Could Have)

- [x] 7. Review and validate completeness of user stories

## Questions for Clarification

### [Question] 1. Company Values
What are the specific company values that should be available for selection when sending a thank-you card? (e.g., Innovation, Teamwork, Customer Focus, Integrity, etc.)
[Answer] 
Our Values
- Make an Impact: Be driven by the desire to build something that can touch millions of lives.
- Strive for Excellence: Today's great is not good enough for tomorrow.
- Stand Together: Embrace the spirit of solidarity & commitment through thick and thin.
- Be Open-Minded: Diversity is strength. Magic often happens at the intersection of two different worlds.
- Stay Grounded: Always remember our roots and higher purpose. Plus, life is too short to be around jerks.

Our Credos
- Bias for Action: Speed matters, take action to deliver a high quality result with calculated risk taking. Be time conscious and set deadlines on initiatives. Promptly ask for help when faced with roadblocks, be the roadblock-remover where we can. It is better to be moving than not, most consequences are manageable.
- Customer Centric: We start with the customer experience, centring our processes and decisions around it. Never stick to what has been done for the sake of tradition. Advocate for our customers, add value wherever possible and build trust with sincere interactions.
- Think Strategically: We know our business inside out, acutely aware of the market dynamics and crystal clear of our positioning and what we offer. We stay ahead of the curve by being flexible enough to pivot quickly, being nimble in problem solving and being obsessed with value creation for our customers, employees and partners. We connect the dots where others don't.
- Deep Dive: No detail is too small, no task too unimportant. Always ask why, validate with data, slicing for insights and looking for opportunities. Only when we know the nuts and bolts of everything we work on, can we troubleshoot effectively and innovate creatively. Identify root causes. When in doubt, keep questioning constructively.
- Invent and Simplify: There is no need to over-complicate. Get creative with solutions, streamline where we can and tap on the expertise of different teams. Innovation and invention are expected and achieved by leveraging all available resources.
- Earn Trust: Always deliver on promises. We listen attentively, speak candidly and maintain the highest standards of ethics - towards our customers and our team.
- Take Ownership: Each of us represents the company, beyond just ourselves or our team. When times are good, we celebrate together. When times are bad, we stay and face the adversity as one. Act with a view of the long term, today's great piece of work will have a lasting impact on the bigger picture. The work we produce speaks for the company, be proud of it, own it.
- Challenge, Disagree, and Commit: When in doubt, challenge respectfully and objectively, even if it is uncomfortable. Unemotionally review the objective of the things we are doing and whether how we are doing it is the best way to do it. Present facts instead of succumbing to feelings. Ideas evolve and improve when scrutinised. Once a decision is determined, commit wholeheartedly and own the consequences together.
- Learn and Be Curious: Learn from team mates, customers and competitors. Always be curious about the whys and how things are done, be eager to go deeper and bring our learnings back to our work.
- Do More with Less: Know where we should be investing, and invest it better. Everyone is empowered to implement efficient processes and reduce spending. Do so wisely because we cut costs but we don't cut corners.


### [Question] 2. Card Visibility and Privacy
Should all thank-you cards be publicly visible to everyone in the company, or should there be privacy options (e.g., team-only, department-only)?
[Answer] Yes, all thank-you cards be publicly visible to everyone in the company.

### [Question] 3. Response Capabilities
What types of responses should recipients be able to make to thank-you cards? (e.g., simple "like", comment, emoji reactions, thank-you reply)
[Answer] Emoji reactions.

### [Question] 4. Card Approval Workflow
Should thank-you cards be published immediately, or should there be an approval/moderation process before they appear in the company-wide feed?
[Answer] In MVP version, there is no approval flow, but for the future we might add this flow.

### [Question] 5. Analytics Scope
For HR analytics, what specific metrics or insights are most important? (e.g., most recognized employees, most active recognizers, value distribution, team recognition patterns, time-based trends)
[Answer] Most recognized employees will be group in a top 10 list, most active recognizers, team recognition patterns

### [Question] 6. Integration with Teams
Should this feature integrate with existing Teams features like notifications, @mentions, or Teams channels? What level of integration is expected?
[Answer] Yes， this feature should be able to integrate with existing teams channels, everyone can go to the channel to check the public thank-you cards.

### [Question] 7. Card Limitations
Should there be any limitations on sending cards? (e.g., maximum cards per day/week, cannot send to same person repeatedly within X days, character limits on messages)
[Answer] No limitation.

---

**Status**: ✅ Complete - All user stories have been created and documented in /inception/user_stories.md


---

# Plan: Grouping User Stories into Independent Units

## Overview
Analyze user stories from `/inception/user_stories.md` and group them into loosely coupled, highly cohesive units that can be built independently by separate teams.

## Steps

### Phase 1: Analysis and Planning
- [x] **Step 1.1**: Analyze all user stories and identify natural boundaries based on:
  - Functional domains (card creation, analytics, notifications, etc.)
  - Platform boundaries (Teams Channel vs Web App)
  - User roles (Employee vs HR Admin)
  - Data dependencies and coupling

- [x] **Step 1.2**: Propose unit groupings with rationale
  - [Question] Should the Teams Channel integration be a separate unit from the Web App, or should they be combined based on functional domains?
  - [Answer] No, the web app only uses HTTP API to interact with Teams channel. Teams channel is not a separate functional domain. 
  
  - [Question] Should HR Analytics be a completely separate unit/service, or part of the main Web App unit?
  - [Answer] A separate unit 
  
  - [Question] For the MVP scope, should we include "Could Have" priority stories (US-1.4: Share cards, US-6.1: Milestones) in the units, or focus only on "Must Have" and "Should Have"?
  - [Answer] Need to focus on all of them 
  
  - [Question] Should notification functionality (US-1.3) be part of the Teams Channel unit or a separate Notification Service unit?
  - [Answer] Notification functionality does not need to be included in any of the above units, it belongs to Teams channel's existing functional unit which exlcluded in this project. 

### Phase 2: Unit Definition
- [x] **Step 2.1**: Create `/inception/units/` folder structure

- [x] **Step 2.2**: For each identified unit, create individual `.md` files containing:
  - Unit name and description
  - Relevant user stories with full acceptance criteria
  - Dependencies on other units
  - Key responsibilities

### Phase 3: Integration Contract
- [x] **Step 3.1**: Identify all inter-unit communication points

- [x] **Step 3.2**: Define API endpoints for each unit:
  - HTTP methods (GET, POST, PUT, DELETE)
  - Endpoint paths
  - Request/response schemas
  - Authentication requirements

- [x] **Step 3.3**: Create `/inception/units/integration_contract.md` with:
  - API specifications for each unit
  - Data models shared across units
  - Event/notification contracts (if applicable)
  - Synchronization requirements between Teams and Web App

### Phase 4: Review and Finalization
- [x] **Step 4.1**: Review all unit files for completeness and consistency

- [x] **Step 4.2**: Validate that units are:
  - Loosely coupled (minimal dependencies)
  - Highly cohesive (related functionality grouped together)
  - Independently buildable and deployable

- [x] **Step 4.3**: Present final unit structure for approval

## Initial Unit Proposal (Pending Clarification)

Based on initial analysis, here are potential unit groupings:

### Option A: Platform-Based Separation
1. **Teams Channel Integration Unit** - All Teams-specific features
2. **Web App Core Unit** - Card creation, personal views, interactions
3. **HR Analytics Unit** - All HR dashboard and reporting features
4. **Notification Service Unit** (Optional) - Cross-platform notifications

### Option B: Domain-Based Separation
1. **Card Management Unit** - Card creation, viewing, filtering (both platforms)
2. **Recognition Analytics Unit** - Personal stats, Top 10, HR analytics
3. **Teams Integration Unit** - Teams channel feed, notifications, sharing
4. **Interaction Unit** - Emoji reactions, search, filters

### Option C: Hybrid Approach
1. **Card Service Unit** - Core card CRUD operations, storage
2. **Teams Channel Unit** - Teams-specific UI and integration
3. **Web App Unit** - Web application UI and features
4. **Analytics Service Unit** - All analytics, reporting, exports for both platforms

[Question] Which unit grouping approach (A, B, C, or a different approach) would you prefer? Please consider your team structure, deployment preferences, and technical architecture.
[Answer] C, but we won't have a separate Teams Channel unit 

---

## Notes
- All "Future Enhancement" stories (Epic 7) will be documented but marked as out of scope for current unit planning
- Integration contracts will focus on REST APIs unless event-driven architecture is preferred
- Emoji reaction synchronization between Teams and Web App will require careful contract definition


---

## Final Unit Structure Summary

### ✅ Completed Units

#### 1. Card Service Unit (`card_service_unit.md`)
- **18 User Stories** covering card management, search, reactions, and statistics
- **15 API Endpoints** for complete card operations
- **Backend service** handling all card-related business logic
- Dependencies: Employee Directory Service, Database

#### 2. Web App Unit (`web_app_unit.md`)
- **13 User Stories** covering UI for card creation, personal views, and interactions
- **Frontend application** consuming Card Service and Analytics Service APIs
- **No exposed APIs** (client-side only)
- Dependencies: Card Service Unit, Analytics Service Unit

#### 3. Analytics Service Unit (`analytics_service_unit.md`)
- **5 User Stories** covering HR analytics, reporting, and data export
- **7 API Endpoints** for analytics dashboard and reports
- **Backend service** with read-only access to card data
- Dependencies: Card Service Unit (read-only), Employee Directory Service

### ✅ Integration Contract (`integration_contract.md`)
- Complete API specifications for all 22 endpoints
- Request/response schemas with validation rules
- Authentication and authorization specifications
- Error handling standards
- Rate limiting policies
- Data synchronization strategies
- External system integration (Teams Channel, Employee Directory)

### ✅ Documentation (`README.md`)
- Unit overview and architecture diagram
- User story distribution by unit and priority
- Development approach and recommended build order
- Testing strategy
- Next steps for implementation

---

## Key Decisions Made

Based on your answers:

1. **No separate Teams Channel Unit**: Teams Channel uses HTTP API to interact with Card Service, not a separate functional domain
2. **HR Analytics as separate unit**: Analytics Service Unit is independent from Web App
3. **All priorities included**: Must Have, Should Have, and Could Have stories all included in units
4. **Notifications excluded**: Handled by Teams Channel's existing functionality, not part of these units

---

## Unit Characteristics Validation

### ✅ Loosely Coupled
- Units communicate only via REST APIs
- No direct database sharing
- Clear API contracts defined
- Independent deployment possible

### ✅ Highly Cohesive
- Card Service: All card operations together
- Web App: All UI components together
- Analytics Service: All analytics together

### ✅ Independently Buildable
- Each unit can be developed by separate team
- Clear dependencies documented
- API versioning supports parallel development

---

**Status**: ✅ Complete - All units defined, integration contracts created, ready for team assignment and development


---

# Plan: Domain Driven Design - Domain Model Design

## Overview
Design comprehensive Domain Driven Design (DDD) domain models for each software unit, including all tactical components: Aggregates, Entities, Value Objects, Domain Events, Policies, Repositories, and Domain Services.

## Steps

### Phase 1: Analysis and Domain Understanding
- [x] **Step 1.1**: Analyze Card Service Unit user stories and identify:
  - Core domain concepts and business rules
  - Bounded context boundaries
  - Ubiquitous language terms
  - Invariants and business constraints

- [x] **Step 1.2**: Analyze Analytics Service Unit user stories and identify:
  - Core domain concepts and business rules
  - Bounded context boundaries
  - Ubiquitous language terms
  - Invariants and business constraints

- [x] **Step 1.3**: Identify cross-cutting concerns and shared concepts between units

### Phase 2: Domain Model Design Questions

#### Card Service Unit Questions

- [ ] **Step 2.1**: Card Aggregate Design
  - [Question] Should Card be the aggregate root with Recipients and Reactions as entities within the aggregate, or should Recipients be value objects since they reference external Employee entities?
  - [Answer] Card is the aggregate root, Sender and Recipient are value objects
  
  - [Question] Should emoji Reactions be part of the Card aggregate or a separate aggregate? Consider: reactions can be added/removed independently and may have high concurrency.
  - [Answer] Not applicable, this is an external Teams system feature, not an object within our system
  
  - [Question] For card creation, should we use a Factory pattern or a domain service? The card creation involves validation of recipients, values, and business rules.
  - [Answer] Factory pattern 

- [ ] **Step 2.2**: Value Objects vs Entities
  - [Question] Should RecognitionReason be a simple value object (string) or a rich value object with validation rules (character limits, content validation)?
  - [Answer] Simple value object (string)
  
  - [Question] Should CompanyValue/Credo be an entity (with ID) or a value object? They are predefined and immutable, but need to be referenced by ID in the API.
  - [Answer] Entity (with ID) 

- [ ] **Step 2.3**: Domain Events
  - [Question] Which domain events should be published? Suggestions: CardCreated, CardShared, ReactionAdded, ReactionRemoved, MilestoneAchieved. Should we include all of these or only critical ones?
  - [Answer] Include all
  
  - [Question] Should domain events be published synchronously or asynchronously? This affects whether Teams Channel notifications happen in real-time.
  - [Answer] Synchronous, published synchronously through Teams HTTP interface 

- [ ] **Step 2.4**: Milestone Tracking
  - [Question] Should Milestone be part of the Card aggregate, a separate aggregate, or managed by a domain service? Milestones track user activity across multiple cards.
  - [Answer] Managed by domain service, milestones track the count of received cards per user dimension
  
  - [Question] Should milestone detection happen synchronously when a card is created (domain event handler) or asynchronously (batch process)?
  - [Answer] Asynchronously 

#### Analytics Service Unit Questions

- [ ] **Step 2.5**: Analytics Bounded Context
  - [Question] Should Analytics Service have its own domain model (separate bounded context) or share the Card domain model? Analytics reads card data but doesn't modify it.
  - [Answer] Analytics Service has its own separate domain model
  
  - [Question] Should analytics aggregations (team stats, value distributions) be computed on-demand or pre-aggregated and stored? This affects whether we need aggregate entities for analytics.
  - [Answer] On-demand computation 

- [ ] **Step 2.6**: Read Model vs Domain Model
  - [Question] Should Analytics Service use a CQRS pattern with a separate read model, or query the Card Service domain model directly?
  - [Answer] Directly query the Card Service domain model 

#### Shared Concepts

- [ ] **Step 2.7**: Employee Reference
  - [Question] How should we model Employee in the domain? As an entity, value object, or just an ID reference to external Employee Directory Service?
  - [Answer] Only as an ID reference to external Employee Directory Service
  
  - [Question] Should we cache employee information (name, department) in the Card aggregate to avoid external service calls, or always fetch fresh data?
  - [Answer] Always fetch fresh data 

### Phase 3: Domain Model Design - Card Service Unit

- [x] **Step 3.1**: Create `/construction/` folder structure

- [x] **Step 3.2**: Design Card Service domain model components:
  - Identify and document all Aggregates with their boundaries
  - Identify and document all Entities within aggregates
  - Identify and document all Value Objects
  - Define all Domain Events
  - Define Domain Services (if needed)
  - Define Repositories
  - Define Policies/Business Rules
  - Document invariants and constraints

- [x] **Step 3.3**: Create `/construction/card_service/domain_model.md` with:
  - Bounded context description
  - Ubiquitous language glossary
  - Aggregate designs with relationships
  - Entity and Value Object specifications
  - Domain Events catalog
  - Domain Services descriptions
  - Repository interfaces
  - Business rules and invariants
  - Domain model diagram (text-based)

### Phase 4: Domain Model Design - Analytics Service Unit

- [x] **Step 4.1**: Design Analytics Service domain model components:
  - Identify and document all Aggregates (if any)
  - Identify and document all Entities
  - Identify and document all Value Objects
  - Define all Domain Events (if any)
  - Define Domain Services
  - Define Repositories
  - Define calculation/aggregation logic

- [x] **Step 4.2**: Create `/construction/analytics_service/domain_model.md` with:
  - Bounded context description
  - Ubiquitous language glossary
  - Aggregate designs (if applicable)
  - Entity and Value Object specifications
  - Domain Services descriptions
  - Repository interfaces
  - Analytics calculation specifications
  - Domain model diagram (text-based)

### Phase 5: Domain Model Design - Web App Unit

- [x] **Step 5.1**: Analyze Web App Unit requirements
  - [Question] Should Web App Unit have its own domain model, or is it purely a presentation layer consuming backend APIs? In DDD, frontend typically doesn't have a domain model.
  - [Answer] Purely a presentation layer consuming backend APIs 

- [x] **Step 5.2**: If Web App needs domain model, create `/construction/web_app/domain_model.md`
  - Otherwise, document that Web App is a presentation layer without domain logic

### Phase 6: Cross-Cutting Concerns

- [x] **Step 6.1**: Document shared kernel (if any):
  - Shared value objects across bounded contexts
  - Shared domain events
  - Integration patterns between contexts

- [x] **Step 6.2**: Create `/construction/integration_patterns.md`:
  - How bounded contexts communicate
  - Anti-corruption layers (if needed)
  - Event-driven integration patterns
  - Shared data models vs separate models

### Phase 7: Review and Validation

- [x] **Step 7.1**: Review all domain models for:
  - Proper aggregate boundaries (consistency boundaries)
  - Correct entity vs value object classifications
  - Complete domain event coverage
  - Repository interface completeness
  - Business rule enforcement locations

- [x] **Step 7.2**: Validate domain models against user stories:
  - All user stories can be implemented with the domain model
  - All business rules are captured
  - All invariants are enforced

- [x] **Step 7.3**: Present final domain models for approval

---

## Key DDD Principles to Follow

### Aggregates
- Each aggregate has one root entity
- Aggregates are consistency boundaries
- External references only by ID
- Small aggregates preferred

### Entities
- Have unique identity
- Mutable over time
- Identity remains constant

### Value Objects
- Immutable
- No identity (equality by value)
- Can be shared
- Prefer value objects over entities when possible

### Domain Events
- Past tense naming (CardCreated, not CreateCard)
- Immutable
- Contain all relevant data
- Published after aggregate state change

### Repositories
- One repository per aggregate root
- Collection-like interface
- Hide persistence details

### Domain Services
- Stateless operations
- Operations that don't belong to any entity
- Coordinate multiple aggregates

---

## Notes
- NO code snippets will be generated (as per requirements)
- Focus on conceptual design and relationships
- Use text-based diagrams (ASCII art or Mermaid syntax)
- Document business rules and invariants clearly
- Use ubiquitous language from the domain

---

## Deliverables
1. `/construction/card_service/domain_model.md` - Complete DDD domain model
2. `/construction/analytics_service/domain_model.md` - Complete DDD domain model
3. `/construction/web_app/domain_model.md` - Domain model or documentation of presentation layer
4. `/construction/integration_patterns.md` - Cross-context integration patterns

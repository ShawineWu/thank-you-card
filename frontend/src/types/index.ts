// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message: string;
  timestamp: string;
  requestId: string;
}

export interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: Array<{
      field: string;
      message: string;
    }>;
  };
  timestamp: string;
  requestId: string;
}

// Card Types
export interface Card {
  id: string;
  sender_id: string;
  reason: string;
  created_at: string;
  updated_at: string;
  recipients: CardRecipient[];
  values: CardValue[];
}

export interface CardRecipient {
  id: string;
  card_id: string;
  recipient_id: string;
  created_at: string;
}

export interface CardValue {
  id: string;
  card_id: string;
  value_id: string;
  created_at: string;
}

export interface CompanyValue {
  id: string;
  name: string;
  description: string;
  type: 'Value' | 'Credo';
  created_at: string;
  updated_at: string;
}

// Employee Types
export interface Employee {
  id: string;
  name: string;
  email: string;
  department?: string;
  team?: string;
}

// Statistics Types
export interface PersonalStatistics {
  cardsSent: number;
  cardsReceived: number;
  valueDistribution: ValueDistribution[];
  recentActivity: RecentActivity[];
}

export interface ValueDistribution {
  valueId: string;
  valueName: string;
  count: number;
  percentage: number;
}

export interface RecentActivity {
  type: 'sent' | 'received';
  cardId: string;
  date: string;
  otherParty: string;
}

export interface Top10Employee {
  employeeId: string;
  employeeName: string;
  cardCount: number;
  rank: number;
}

// Milestone Types
export interface EmployeeMilestone {
  id: string;
  employee_id: string;
  milestone_type: 'CardsSent' | 'CardsReceived';
  threshold: number;
  achieved_at: string;
  title: string;
  description: string;
}

// Form Types
export interface CreateCardForm {
  recipients: string[];
  reason: string;
  selectedValues: string[];
}

export interface SearchFilters {
  sender?: string;
  recipient?: string;
  values?: string[];
  dateFrom?: string;
  dateTo?: string;
  searchText?: string;
}

// Pagination Types
export interface PaginationInfo {
  page: number;
  pageSize: number;
  totalPages: number;
  totalItems: number;
  hasNext: boolean;
  hasPrevious: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationInfo;
}

// Auth Types
export interface User {
  id: string;
  name: string;
  email: string;
  role: 'Employee' | 'HR';
  department?: string;
  team?: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

// Store Types
export interface CardStore {
  cards: Card[];
  currentCard: Card | null;
  isLoading: boolean;
  error: string | null;
  filters: SearchFilters;
  pagination: PaginationInfo | null;
}

export interface AnalyticsStore {
  personalStats: PersonalStatistics | null;
  top10: Top10Employee[];
  isLoading: boolean;
  error: string | null;
}

export interface UIStore {
  sidebarOpen: boolean;
  theme: 'light' | 'dark';
  notifications: Notification[];
}

export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
}
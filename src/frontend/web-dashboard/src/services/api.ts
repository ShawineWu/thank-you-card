import { authorizedAxios } from "../utils/axios";

export interface UserStats {
  employeeId: string;
  cardsSent: number;
  cardsReceived: number;
  valueStats: Array<{
    valueId: string;
    valueName: string;
    count: number;
  }>;
  milestones: Array<{
    type: string;
    threshold: number;
    title: string;
    description: string;
    achievedAt: string;
  }>;
}

export interface TopEmployee {
  employeeId: string;
  cardsReceived: number;
  rank: number;
}

export interface ValueResponse {
  id: string;
  name: string;
  description: string;
  type: string;
}

export interface Employee {
  id: string;
  aadId?: string;
  name: string;
  email: string;
  department: string;
}

export interface CreateCardRequest {
  recipients: Array<{ id: string; name: string }>;
  recognitionReason: string;
  valueIds: string[];
}

export interface Recipient {
  id: string;
  name: string;
}

export interface Card {
  id: string;
  senderId: string;
  senderName: string;
  recipients: Recipient[];
  recognitionReason: string;
  selectedValues: ValueResponse[];
  createdAt: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    pageSize: number;
    totalPages: number;
    totalItems: number;
    hasNext: boolean;
    hasPrevious: boolean;
  };
}

export interface CardFilter {
  page?: number;
  pageSize?: number;
  senderId?: string;
  recipientId?: string;
  valueIds?: string[];
  startDate?: string;
  endDate?: string;
  search?: string;
}

// Analytics types
export interface DashboardAnalytics {
  totalCards: number;
  activeUsers: number;
  topValues: Array<{
    valueId: string;
    valueName: string;
    count: number;
  }>;
  cardsTrend: {
    thisMonth: number;
    lastMonth: number;
  };
  engagementRate: number;
}

export interface TopRecognizer {
  employeeId: string;
  cardsSent: number;
  rank: number;
}

export interface ValueDistribution {
  distribution: Array<{
    valueId: string;
    valueName: string;
    count: number;
    percentage: number;
  }>;
  totalCards: number;
}

export interface ExportRequest {
  startDate: string;
  endDate: string;
  format: "csv";
}

export interface StandardResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  timestamp: string;
  requestId: string;
}

export const api = {
  getCards: (params: CardFilter = {}) => {
    const { page = 1, pageSize = 10, ...rest } = params;
    const query = new URLSearchParams({
      page: page.toString(),
      pageSize: pageSize.toString(),
    });

    Object.entries(rest).forEach(([key, value]) => {
      if (value) {
        if (Array.isArray(value)) {
          value.forEach((v) => query.append(`${key}[]`, v));
        } else {
          query.append(key, value.toString());
        }
      }
    });

    return authorizedAxios.get<StandardResponse<PaginatedResponse<Card>>>(
      `/cards?${query.toString()}`,
    );
  },

  getValues: () =>
    authorizedAxios.get<StandardResponse<ValueResponse[]>>("/values"),

  getUserStats: (employeeId: string) =>
    authorizedAxios.get<StandardResponse<UserStats>>(
      `/statistics/user/${employeeId}`,
    ),

  getTopEmployees: (limit = 10) =>
    authorizedAxios.get<StandardResponse<TopEmployee[]>>(
      `/statistics/top10?limit=${limit}`,
    ),

  // Analytics APIs
  getAnalyticsDashboard: () =>
    authorizedAxios.get<StandardResponse<DashboardAnalytics>>(
      "/analytics/dashboard",
    ),

  getTopRecognizers: (limit = 10) =>
    authorizedAxios.get<StandardResponse<{ recognizers: TopRecognizer[] }>>(
      `/analytics/recognizers/top?limit=${limit}`,
    ),

  getValueDistribution: () =>
    authorizedAxios.get<StandardResponse<ValueDistribution>>(
      "/analytics/values/distribution",
    ),

  exportAnalytics: (request: ExportRequest) =>
    authorizedAxios.post("/analytics/export", request, {
      responseType: "blob",
    }),

  // Card creation and employee search
  createCard: (request: CreateCardRequest) =>
    authorizedAxios.post<StandardResponse<Card>>("/cards", request),

  generateRecognitionReason: (request: {
    keyInfo: string;
    valueNames: string[];
    senderName: string;
    recipientNames: string[];
  }) =>
    authorizedAxios.post<StandardResponse<{ recognitionReason: string }>>(
      "/cards/generate-reason",
      request,
    ),

  searchEmployees: (query: string) =>
    authorizedAxios.get<StandardResponse<Employee[]>>(
      `/employees/search?q=${encodeURIComponent(query)}`,
    ),
};

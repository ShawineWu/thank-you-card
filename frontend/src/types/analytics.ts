export interface DashboardOverview {
  totalCards: number;
  totalEmployees: number;
  totalDepartments: number;
  activeRecognizers: number;
  period: string;
}

export interface TeamPattern {
  department: string;
  totalSent: number;
  totalReceived: number;
  averageCardsPerEmployee: number;
}

export interface ValueDistribution {
  valueId: number;
  code: string;
  name: string;
  type: 'VALUE' | 'CREDO';
  count: number;
  percentage: number;
}

export interface AnalyticsListResponse<T> {
  items: T[];
  total?: number;
}

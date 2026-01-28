export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface TimeRange {
  from?: string;
  to?: string;
}

export interface ApiError {
  error: string;
}

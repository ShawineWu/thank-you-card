import api from './api';
import { DashboardOverview, TeamPattern, ValueDistribution, AnalyticsListResponse } from '@/types/analytics';
import { TopEmployee } from '@/types/stats';
import { TimeRange } from '@/types/common';

export const analyticsApi = {
  // Get dashboard overview
  getDashboard: async (timeRange?: TimeRange): Promise<DashboardOverview> => {
    const params = new URLSearchParams();
    if (timeRange?.from) params.append('from', timeRange.from);
    if (timeRange?.to) params.append('to', timeRange.to);

    const response = await api.get<DashboardOverview>(`/analytics/dashboard?${params.toString()}`);
    return response.data;
  },

  // Get most active recognizers
  getMostActiveRecognizers: async (timeRange?: TimeRange, limit?: number): Promise<TopEmployee[]> => {
    const params = new URLSearchParams();
    if (timeRange?.from) params.append('from', timeRange.from);
    if (timeRange?.to) params.append('to', timeRange.to);
    if (limit) params.append('limit', limit.toString());

    const response = await api.get<AnalyticsListResponse<TopEmployee>>(`/analytics/recognizers?${params.toString()}`);
    return response.data.items;
  },

  // Get team recognition patterns
  getTeamPatterns: async (timeRange?: TimeRange): Promise<TeamPattern[]> => {
    const params = new URLSearchParams();
    if (timeRange?.from) params.append('from', timeRange.from);
    if (timeRange?.to) params.append('to', timeRange.to);

    const response = await api.get<AnalyticsListResponse<TeamPattern>>(`/analytics/teams?${params.toString()}`);
    return response.data.items;
  },

  // Get values distribution
  getValuesDistribution: async (timeRange?: TimeRange): Promise<ValueDistribution[]> => {
    const params = new URLSearchParams();
    if (timeRange?.from) params.append('from', timeRange.from);
    if (timeRange?.to) params.append('to', timeRange.to);

    const response = await api.get<AnalyticsListResponse<ValueDistribution>>(`/analytics/values?${params.toString()}`);
    return response.data.items;
  },

  // Export cards as CSV
  exportCards: async (timeRange?: TimeRange, filters?: Record<string, any>): Promise<Blob> => {
    const params = new URLSearchParams();
    if (timeRange?.from) params.append('from', timeRange.from);
    if (timeRange?.to) params.append('to', timeRange.to);
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value) params.append(key, value.toString());
      });
    }

    const response = await api.get(`/analytics/export?${params.toString()}`, {
      responseType: 'blob',
    });
    return response.data;
  },
};

import api from './api';
import { PersonalStats, TopEmployee } from '@/types/stats';
import { TimeRange } from '@/types/common';

export const statsApi = {
  // Get personal statistics
  getPersonalStats: async (): Promise<PersonalStats> => {
    const response = await api.get<PersonalStats>('/cards/me/stats');
    return response.data;
  },

  // Get top recognized employees
  getTopRecipients: async (timeRange?: TimeRange, limit?: number): Promise<TopEmployee[]> => {
    const params = new URLSearchParams();
    if (timeRange?.from) params.append('from', timeRange.from);
    if (timeRange?.to) params.append('to', timeRange.to);
    if (limit) params.append('limit', limit.toString());

    const response = await api.get<{ items: TopEmployee[] }>(`/cards/top-recipients?${params.toString()}`);
    return response.data.items;
  },
};

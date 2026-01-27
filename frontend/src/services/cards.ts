import api from './api';
import { CardResponse, CreateCardRequest, CardListResponse, CardFilters } from '@/types/card';

export const cardsApi = {
  // Create a new card
  createCard: async (data: CreateCardRequest): Promise<CardResponse> => {
    const response = await api.post<CardResponse>('/cards', data);
    return response.data;
  },

  // Get company-wide feed
  getFeed: async (filters?: CardFilters): Promise<CardListResponse> => {
    const params = new URLSearchParams();
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.pageSize) params.append('pageSize', filters.pageSize.toString());
    if (filters?.from) params.append('from', filters.from);
    if (filters?.to) params.append('to', filters.to);
    if (filters?.q) params.append('q', filters.q);
    if (filters?.values?.length) {
      filters.values.forEach(v => params.append('values[]', v.toString()));
    }
    if (filters?.senderIds?.length) {
      filters.senderIds.forEach(id => params.append('senderIds[]', id.toString()));
    }
    if (filters?.recipientIds?.length) {
      filters.recipientIds.forEach(id => params.append('recipientIds[]', id.toString()));
    }

    const response = await api.get<CardListResponse>(`/cards/feed?${params.toString()}`);
    return response.data;
  },

  // Get personal received cards
  getMyReceived: async (filters?: CardFilters): Promise<CardListResponse> => {
    const params = new URLSearchParams();
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.pageSize) params.append('pageSize', filters.pageSize.toString());
    if (filters?.from) params.append('from', filters.from);
    if (filters?.to) params.append('to', filters.to);
    if (filters?.q) params.append('q', filters.q);

    const response = await api.get<CardListResponse>(`/cards/me/received?${params.toString()}`);
    return response.data;
  },

  // Get personal sent cards
  getMySent: async (filters?: CardFilters): Promise<CardListResponse> => {
    const params = new URLSearchParams();
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.pageSize) params.append('pageSize', filters.pageSize.toString());
    if (filters?.from) params.append('from', filters.from);
    if (filters?.to) params.append('to', filters.to);
    if (filters?.q) params.append('q', filters.q);

    const response = await api.get<CardListResponse>(`/cards/me/sent?${params.toString()}`);
    return response.data;
  },
};

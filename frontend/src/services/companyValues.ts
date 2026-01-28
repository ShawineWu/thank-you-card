import api from './api';
import { CompanyValueSummary, CompanyValueDetail } from '@/types/card';

export interface CompanyValueListResponse {
  items: CompanyValueDetail[];
  total: number;
}

export const companyValuesApi = {
  // Get all company values and credos
  getAll: async (): Promise<CompanyValueListResponse> => {
    const response = await api.get<CompanyValueListResponse>('/company-values');
    return response.data;
  },

  // Get company values by type (VALUE or CREDO)
  getByType: async (type?: 'VALUE' | 'CREDO'): Promise<CompanyValueListResponse> => {
    const params = type ? `?type=${type}` : '';
    const response = await api.get<CompanyValueListResponse>(`/company-values${params}`);
    return response.data;
  },

  // Get a single company value by ID
  getById: async (id: number): Promise<CompanyValueDetail> => {
    const response = await api.get<CompanyValueDetail>(`/company-values/${id}`);
    return response.data;
  },
};

// Legacy mock functions for backward compatibility
// These will be replaced by API calls in components
export const getCompanyValues = (): CompanyValueSummary[] => [];
export const getValues = (): CompanyValueSummary[] => [];
export const getCredos = (): CompanyValueSummary[] => [];

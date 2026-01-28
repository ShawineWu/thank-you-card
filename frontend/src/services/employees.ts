import api from './api';
import { EmployeeSummary } from '@/types/card';

export interface EmployeeListResponse {
  items: EmployeeSummary[];
  total: number;
}

export const employeesApi = {
  // Get all employees
  getAll: async (): Promise<EmployeeListResponse> => {
    const response = await api.get<EmployeeListResponse>('/employees');
    return response.data;
  },

  // Search employees by name, department, or email
  searchEmployees: async (query?: string): Promise<EmployeeSummary[]> => {
    const params = new URLSearchParams();
    if (query) {
      params.append('q', query);
    }
    const response = await api.get<EmployeeListResponse>(`/employees/search?${params.toString()}`);
    return response.data.items;
  },

  // Get employee by ID (fallback to search if needed)
  getEmployeeById: async (id: number): Promise<EmployeeSummary | null> => {
    try {
      const allEmployees = await employeesApi.getAll();
      return allEmployees.items.find(emp => emp.id === id) || null;
    } catch {
      return null;
    }
  },
};

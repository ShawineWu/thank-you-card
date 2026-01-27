import { EmployeeSummary } from '@/types/card';

// Mock employee data - in production, this would come from an API
// For MVP, we'll use mock data for recipient selection
export const MOCK_EMPLOYEES: EmployeeSummary[] = [
  { id: 1, name: 'Mock Employee', department: 'Engineering' },
  { id: 2, name: 'Mock HR Admin', department: 'HR' },
  { id: 3, name: 'Alice Johnson', department: 'Engineering' },
  { id: 4, name: 'Bob Smith', department: 'Product' },
  { id: 5, name: 'Carol White', department: 'Design' },
  { id: 6, name: 'David Brown', department: 'Marketing' },
  { id: 7, name: 'Eva Green', department: 'Sales' },
  { id: 8, name: 'Frank Miller', department: 'Engineering' },
  { id: 9, name: 'Grace Lee', department: 'Product' },
  { id: 10, name: 'Henry Wilson', department: 'Operations' },
];

export const employeesApi = {
  // Search employees (mock implementation)
  searchEmployees: async (query?: string): Promise<EmployeeSummary[]> => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 200));
    
    if (!query) {
      return MOCK_EMPLOYEES;
    }
    
    const lowerQuery = query.toLowerCase();
    return MOCK_EMPLOYEES.filter(
      emp => 
        emp.name.toLowerCase().includes(lowerQuery) ||
        emp.department.toLowerCase().includes(lowerQuery)
    );
  },

  // Get employee by ID
  getEmployeeById: async (id: number): Promise<EmployeeSummary | null> => {
    await new Promise(resolve => setTimeout(resolve, 100));
    return MOCK_EMPLOYEES.find(emp => emp.id === id) || null;
  },
};

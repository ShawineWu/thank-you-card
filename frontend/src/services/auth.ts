import api from './api';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    role: 'EMPLOYEE' | 'HR' | 'ADMIN';
    employee?: {
      id: number;
      name: string;
      email: string;
      department: string;
    };
  };
}

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data);
    return response.data;
  },
};

// Token management
export const tokenStorage = {
  get: (): string | null => {
    return localStorage.getItem('auth_token');
  },
  set: (token: string): void => {
    localStorage.setItem('auth_token', token);
  },
  remove: (): void => {
    localStorage.removeItem('auth_token');
  },
};

// User info management
export const userStorage = {
  get: (): LoginResponse['user'] | null => {
    const userStr = localStorage.getItem('auth_user');
    if (!userStr) return null;
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  },
  set: (user: LoginResponse['user']): void => {
    localStorage.setItem('auth_user', JSON.stringify(user));
  },
  remove: (): void => {
    localStorage.removeItem('auth_user');
  },
};

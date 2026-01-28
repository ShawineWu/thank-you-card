import axios, { AxiosInstance, AxiosError } from 'axios';
import { ApiError } from '@/types/common';

// Use relative URL to leverage Vite proxy, or absolute URL if specified
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

export const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token to requests
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiError>) => {
    // Handle common errors
    if (error.response) {
      const status = error.response.status;
      // Handle 401 Unauthorized - clear auth and redirect to login
      if (status === 401) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
        // Redirect to login if not already there
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
      }
      const message = error.response.data?.error || error.message;
      console.error('API Error:', message);
    } else if (error.request) {
      console.error('Network Error:', error.request);
    } else {
      console.error('Error:', error.message);
    }
    return Promise.reject(error);
  }
);

export default api;

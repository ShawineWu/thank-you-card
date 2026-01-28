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
      const errorMessage = error.response.data?.error || error.message;
      
      // Log detailed error for debugging
      console.error('API Error:', {
        status,
        message: errorMessage,
        url: error.config?.url,
        method: error.config?.method,
      });
      
      // Handle 401 Unauthorized - clear auth and redirect to login
      if (status === 401) {
        // Only redirect if it's a real authentication error, not a missing employee error
        // Some endpoints might return 401 for missing employee, which is a different issue
        const isAuthError = errorMessage === 'unauthorized' || 
                           errorMessage === 'invalid token' || 
                           errorMessage === 'invalid authorization header';
        
        if (isAuthError) {
          localStorage.removeItem('auth_token');
          localStorage.removeItem('auth_user');
          // Redirect to login if not already there
          if (window.location.pathname !== '/login') {
            window.location.href = '/login';
          }
        } else {
          // For other 401 errors (like missing employee), show error but don't redirect
          console.error('Authentication issue (possibly missing employee):', errorMessage);
        }
      }
    } else if (error.request) {
      console.error('Network Error:', error.request);
    } else {
      console.error('Error:', error.message);
    }
    return Promise.reject(error);
  }
);

export default api;

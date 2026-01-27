import axios, {
  type AxiosInstance,
  type InternalAxiosRequestConfig,
  type AxiosError,
} from "axios";
import { initConfigProcess } from "../app";

/**
 * Authorized Axios instance with automatic token injection and refresh
 */
class AuthorizedAxios {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create({
      timeout: 30000,
      baseURL: typeof API_BASE_URL !== "undefined" ? API_BASE_URL : undefined,
    });

    // Request interceptor - add Authorization header
    this.axiosInstance.interceptors.request.use(
      async (config: InternalAxiosRequestConfig) => {
        const auth = await initConfigProcess;
        const user = await auth.getUser();

        if (user?.access_token) {
          config.headers.Authorization = `Bearer ${user.access_token}`;
        }

        return config;
      },
      (error) => {
        return Promise.reject(error);
      },
    );

    // Response interceptor - handle 401 errors
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as InternalAxiosRequestConfig & {
          _retry?: boolean;
        };

        // Handle 401 errors and retry with refreshed token
        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            const auth = await initConfigProcess;
            await auth.renewToken();

            // Retry original request with new token
            const user = await auth.getUser();
            if (user?.access_token && originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${user.access_token}`;
            }

            return this.axiosInstance(originalRequest);
          } catch (refreshError) {
            // Token refresh failed, redirect to login
            const auth = await initConfigProcess;
            await auth.login();
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      },
    );
  }

  getInstance(): AxiosInstance {
    return this.axiosInstance;
  }
}

// Export singleton instance
export const authorizedAxios = new AuthorizedAxios().getInstance();

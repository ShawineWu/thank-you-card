// Global configuration variables injected at runtime
declare const REACT_APP_ENV: "local" | "dev" | "test" | "prod";
declare const AUTH_SERVER: string;
declare const CLIENT_ID: string;
declare const LOGIN_CALLBACK_URL: string;
declare const LOGOUT_CALLBACK_URL: string;
declare const PLAT_UMS_SERVER: string;
declare const API_BASE_URL: string;

// Window object extensions
interface Window {
  REACT_APP_ENV?: "local" | "dev" | "test" | "prod";
  AUTH_SERVER?: string;
  CLIENT_ID?: string;
  LOGIN_CALLBACK_URL?: string;
  LOGOUT_CALLBACK_URL?: string;
  PLAT_UMS_SERVER?: string;
  API_BASE_URL?: string;
}

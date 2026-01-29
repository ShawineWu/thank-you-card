import type {
  CompanyValue,
  CreateCardRequest,
  ApiResponse,
  Employee,
} from "../types";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1";
const AUTH_TOKEN = import.meta.env.VITE_AUTH_TOKEN || "";

const headers: HeadersInit = {
  "Content-Type": "application/json",
};

if (AUTH_TOKEN) {
  headers["Authorization"] = `Bearer ${AUTH_TOKEN}`;
}

export const api = {
  async getCompanyValues(): Promise<CompanyValue[]> {
    const response = await fetch(`${API_BASE_URL}/values`, { headers });
    const data: ApiResponse<CompanyValue[]> = await response.json();

    if (!data.success) {
      throw new Error(data.error.message);
    }

    return data.data;
  },

  async searchEmployees(query: string): Promise<Employee[]> {
    const response = await fetch(
      `${API_BASE_URL}/employees/search?q=${encodeURIComponent(query)}`,
      { headers },
    );
    const data: ApiResponse<Employee[]> = await response.json();

    if (!data.success) {
      throw new Error(data.error.message);
    }

    return data.data;
  },

  async createCard(cardData: CreateCardRequest): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/cards`, {
      method: "POST",
      headers,
      body: JSON.stringify(cardData),
    });

    const data: ApiResponse<any> = await response.json();

    if (!data.success) {
      throw new Error(data.error.message);
    }
  },
};

export interface CompanyValue {
  id: string;
  name: string;
  description: string;
  type: string;
}

export interface Employee {
  id: string;
  aadId: string;
  name: string;
  email: string;
  department?: string;
}

export interface Recipient {
  id: string;
  name: string;
}

export interface CardFormData {
  recipients: Recipient[];
  recognitionReason: string;
  valueIds: string[];
}

export interface CreateCardRequest {
  recipients: Recipient[];
  recognitionReason: string;
  valueIds: string[];
}

export interface ApiError {
  success: false;
  error: {
    code: string;
    message: string;
    details?: any;
  };
  timestamp: string;
  requestId: string;
}

export interface ApiSuccess<T> {
  success: true;
  data: T;
  message?: string;
  timestamp: string;
  requestId: string;
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError;

import api from './api';

export interface GenerateTextRequest {
  recipientName: string;
  userInput: string;
  valueNames: string[];
}

export interface GenerateTextResponse {
  text: string;
}

export const aiApi = {
  // Generate recognition text using AI
  generateText: async (data: GenerateTextRequest): Promise<GenerateTextResponse> => {
    const response = await api.post<GenerateTextResponse>('/ai/generate-text', data);
    return response.data;
  },
};

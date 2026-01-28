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
    console.log('AI API: Sending request to /ai/generate-text', data);
    try {
      const response = await api.post<GenerateTextResponse>('/ai/generate-text', data);
      console.log('AI API: Response received', response.data);
      return response.data;
    } catch (error: any) {
      console.error('AI API: Error occurred', error);
      console.error('AI API: Error response', error.response?.data);
      throw error;
    }
  },
};

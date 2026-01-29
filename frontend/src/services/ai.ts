import api from './api';

export interface GenerateTextRequest {
  recipientName: string;
  senderName?: string;
  senderDepartment?: string;
  userInput: string;
  valueNames: string[];
  language?: string; // Language code: "en" (English, default), "zh" (中文), "ja" (日本語), "ko" (한국어), "fr" (Français), "de" (Deutsch), "es" (Español), "pt" (Português), "it" (Italiano), "ru" (Русский)
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

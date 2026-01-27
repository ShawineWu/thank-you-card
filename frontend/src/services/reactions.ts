import api from './api';

export const reactionsApi = {
  // Add or update emoji reaction
  addReaction: async (cardId: number, emojiCode: string): Promise<void> => {
    await api.post(`/cards/${cardId}/reactions`, { emojiCode });
  },

  // Remove emoji reaction
  removeReaction: async (cardId: number): Promise<void> => {
    await api.delete(`/cards/${cardId}/reactions`);
  },
};

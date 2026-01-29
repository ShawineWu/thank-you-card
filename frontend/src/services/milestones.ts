import api from './api';
import { Milestone, UserAchievement } from '@/types/milestone';

export const milestonesApi = {
  // Get all milestone definitions
  getAllMilestones: async (): Promise<Milestone[]> => {
    const response = await api.get<{ items: Milestone[] }>('/milestones');
    return response.data.items;
  },

  // Get current user's achievements
  getMyAchievements: async (): Promise<UserAchievement[]> => {
    const response = await api.get<{ items: UserAchievement[] }>('/milestones/me');
    return response.data.items;
  },
};

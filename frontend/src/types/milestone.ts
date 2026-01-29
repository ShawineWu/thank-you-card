export interface Milestone {
  id: number;
  code: string;
  name: string;
  description: string;
  type: 'SENT' | 'RECEIVED' | 'TOTAL';
  threshold: number;
  iconUrl?: string;
}

export interface UserAchievement {
  milestone: Milestone;
  achievedAt: string; // ISO date string
}

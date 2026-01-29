import React from "react";
import { cn } from "@/lib/utils";

export type AchievementType = "received" | "sent" | "total";

export interface Achievement {
  id: string;
  type: AchievementType;
  threshold: number;
  title: string;
  description: string;
  achievedAt?: string;
  unlocked: boolean;
}

interface AchievementBadgeProps {
  achievement: Achievement;
  onClick?: () => void;
}

// Tier based on threshold - like game ranks
type Tier = "bronze" | "silver" | "gold" | "diamond";

const getTier = (threshold: number): Tier => {
  if (threshold <= 5) return "bronze";
  if (threshold <= 10) return "silver";
  if (threshold <= 50) return "gold";
  return "diamond";
};

// Tier colors and gradients
const tierStyles: Record<Tier, { gradient: string; border: string; glow: string; accent: string }> = {
  bronze: {
    gradient: "from-amber-600 via-amber-500 to-yellow-600",
    border: "border-amber-400",
    glow: "shadow-amber-500/50",
    accent: "#CD7F32",
  },
  silver: {
    gradient: "from-slate-300 via-gray-200 to-slate-400",
    border: "border-slate-300",
    glow: "shadow-slate-400/50",
    accent: "#C0C0C0",
  },
  gold: {
    gradient: "from-yellow-400 via-amber-300 to-yellow-500",
    border: "border-yellow-400",
    glow: "shadow-yellow-500/50",
    accent: "#FFD700",
  },
  diamond: {
    gradient: "from-cyan-300 via-blue-200 to-purple-300",
    border: "border-cyan-300",
    glow: "shadow-cyan-400/50",
    accent: "#B9F2FF",
  },
};

const tierLabels: Record<Tier, string> = {
  bronze: "Bronze",
  silver: "Silver",
  gold: "Gold",
  diamond: "Diamond",
};

const typeIcons: Record<AchievementType, React.ReactNode> = {
  received: (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
      <path d="M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 14H4V8l8 5 8-5v10zm-8-7L4 6h16l-8 5z"/>
    </svg>
  ),
  sent: (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
      <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
    </svg>
  ),
  total: (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
      <path d="M19 5h-2V3H7v2H5c-1.1 0-2 .9-2 2v1c0 2.55 1.92 4.63 4.39 4.94.63 1.5 1.98 2.63 3.61 2.96V19H7v2h10v-2h-4v-3.1c1.63-.33 2.98-1.46 3.61-2.96C19.08 12.63 21 10.55 21 8V7c0-1.1-.9-2-2-2zm-7 9c-1.86 0-3.41-1.28-3.86-3h7.72c-.45 1.72-2 3-3.86 3z"/>
    </svg>
  ),
};

export const AchievementBadge: React.FC<AchievementBadgeProps> = ({
  achievement,
  onClick,
}) => {
  const { type, threshold, unlocked } = achievement;
  const tier = getTier(threshold);
  const styles = tierStyles[tier];

  return (
    <button
      onClick={onClick}
      className={cn(
        "group flex flex-col items-center gap-2 p-2 rounded-xl transition-all duration-300",
        "hover:scale-105 focus:outline-none focus:ring-2 focus:ring-primary/50",
        unlocked ? "cursor-pointer" : "cursor-default"
      )}
    >
      {/* Badge Container */}
      <div className="relative">
        {/* Outer glow effect for unlocked */}
        {unlocked && (
          <div
            className={cn(
              "absolute inset-0 rounded-2xl blur-md opacity-60 transition-opacity",
              `bg-gradient-to-br ${styles.gradient}`,
              "group-hover:opacity-80"
            )}
          />
        )}
        
        {/* Main badge shape - shield/medal style */}
        <div
          className={cn(
            "relative w-20 h-24 flex flex-col items-center justify-center",
            "rounded-t-2xl rounded-b-[40%]",
            "border-2 transition-all duration-300",
            unlocked
              ? `bg-gradient-to-br ${styles.gradient} ${styles.border} shadow-lg ${styles.glow}`
              : "bg-gray-200 dark:bg-gray-700 border-gray-300 dark:border-gray-600"
          )}
        >
          {/* Shine effect */}
          {unlocked && (
            <div className="absolute inset-0 rounded-t-2xl rounded-b-[40%] overflow-hidden">
              <div className="absolute top-0 left-0 w-full h-1/2 bg-gradient-to-b from-white/30 to-transparent" />
            </div>
          )}
          
          {/* Icon */}
          <div
            className={cn(
              "relative z-10 mb-1",
              unlocked ? "text-white drop-shadow-md" : "text-gray-400 dark:text-gray-500"
            )}
          >
            {typeIcons[type]}
          </div>
          
          {/* Threshold number */}
          <div
            className={cn(
              "relative z-10 text-lg font-bold",
              unlocked ? "text-white drop-shadow-md" : "text-gray-400 dark:text-gray-500"
            )}
          >
            {threshold}
          </div>

          {/* Tier ribbon */}
          {unlocked && (
            <div
              className={cn(
                "absolute -bottom-1 left-1/2 -translate-x-1/2",
                "px-2 py-0.5 rounded-full text-[10px] font-semibold",
                "bg-white/90 shadow-sm",
                tier === "bronze" && "text-amber-700",
                tier === "silver" && "text-slate-600",
                tier === "gold" && "text-yellow-700",
                tier === "diamond" && "text-cyan-600"
              )}
            >
              {tierLabels[tier]}
            </div>
          )}
        </div>

        {/* Lock icon for locked achievements */}
        {!unlocked && (
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="absolute inset-0 bg-gray-900/20 rounded-t-2xl rounded-b-[40%]" />
            <svg
              className="relative z-10 w-8 h-8 text-gray-400 dark:text-gray-500"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"/>
            </svg>
          </div>
        )}
      </div>

      {/* Label */}
      <span
        className={cn(
          "text-xs font-medium text-center max-w-[90px] leading-tight",
          unlocked ? "text-foreground" : "text-muted-foreground"
        )}
      >
        {achievement.title}
      </span>
    </button>
  );
};

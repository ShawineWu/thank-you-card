import React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { Achievement, AchievementType } from "./AchievementBadge";
import { CheckCircle, Lock, Calendar, Target } from "lucide-react";

interface AchievementDialogProps {
  achievement: Achievement | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

type Tier = "bronze" | "silver" | "gold" | "diamond";

const getTier = (threshold: number): Tier => {
  if (threshold <= 5) return "bronze";
  if (threshold <= 10) return "silver";
  if (threshold <= 50) return "gold";
  return "diamond";
};

const tierStyles: Record<Tier, { gradient: string; text: string }> = {
  bronze: { gradient: "from-amber-600 via-amber-500 to-yellow-600", text: "text-amber-700" },
  silver: { gradient: "from-slate-300 via-gray-200 to-slate-400", text: "text-slate-600" },
  gold: { gradient: "from-yellow-400 via-amber-300 to-yellow-500", text: "text-yellow-700" },
  diamond: { gradient: "from-cyan-300 via-blue-200 to-purple-300", text: "text-cyan-600" },
};

const tierLabels: Record<Tier, string> = {
  bronze: "Bronze",
  silver: "Silver", 
  gold: "Gold",
  diamond: "Diamond",
};

const typeLabels: Record<AchievementType, string> = {
  received: "Cards Received",
  sent: "Cards Sent",
  total: "Total Cards",
};

export const AchievementDialog: React.FC<AchievementDialogProps> = ({
  achievement,
  open,
  onOpenChange,
}) => {
  if (!achievement) return null;

  const { type, threshold, title, description, achievedAt, unlocked } = achievement;
  const tier = getTier(threshold);
  const styles = tierStyles[tier];

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <div className="flex items-center gap-4">
            {/* Badge preview - shield style */}
            <div
              className={cn(
                "w-16 h-20 flex flex-col items-center justify-center",
                "rounded-t-xl rounded-b-[40%] border-2",
                unlocked
                  ? `bg-gradient-to-br ${styles.gradient} border-white/50 shadow-lg`
                  : "bg-gray-200 dark:bg-gray-700 border-gray-300"
              )}
            >
              <span className={cn(
                "text-xl font-bold",
                unlocked ? "text-white drop-shadow" : "text-gray-400"
              )}>
                {threshold}
              </span>
              {unlocked && (
                <span className={cn("text-[10px] font-semibold text-white/90")}>
                  {tierLabels[tier]}
                </span>
              )}
            </div>
            <div>
              <DialogTitle className="text-xl">{title}</DialogTitle>
              <div className="flex items-center gap-2 mt-1">
                <Badge
                  variant="secondary"
                  className={cn(
                    unlocked && "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100"
                  )}
                >
                  {unlocked ? "Unlocked" : "Locked"}
                </Badge>
                {unlocked && (
                  <Badge variant="outline" className={styles.text}>
                    {tierLabels[tier]} Tier
                  </Badge>
                )}
              </div>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-4 mt-4">
          {/* Description */}
          <DialogDescription className="text-base">
            {description}
          </DialogDescription>

          {/* Achievement details */}
          <div className="space-y-3 pt-2 border-t">
            <div className="flex items-center gap-2 text-sm">
              <Target className="w-4 h-4 text-muted-foreground" />
              <span className="text-muted-foreground">Goal:</span>
              <span className="font-medium">
                {typeLabels[type]} - {threshold}
              </span>
            </div>

            {unlocked && achievedAt ? (
              <div className="flex items-center gap-2 text-sm">
                <Calendar className="w-4 h-4 text-muted-foreground" />
                <span className="text-muted-foreground">Achieved:</span>
                <span className="font-medium">{formatDate(achievedAt)}</span>
              </div>
            ) : (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Lock className="w-4 h-4" />
                <span>Keep going to unlock this achievement!</span>
              </div>
            )}

            {unlocked && (
              <div className="flex items-center gap-2 text-sm text-green-600 dark:text-green-400">
                <CheckCircle className="w-4 h-4" />
                <span>Congratulations on this achievement!</span>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};

import React, { useState, useMemo } from "react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Trophy } from "lucide-react";
import { AchievementBadge, type Achievement, type AchievementType } from "./AchievementBadge";
import { AchievementDialog } from "./AchievementDialog";

// Achievement definitions - thresholds for each type
const ACHIEVEMENT_THRESHOLDS = [5, 10, 50, 100];

interface UserMilestone {
  type: string;
  threshold: number;
  title: string;
  description: string;
  achievedAt: string;
}

interface AchievementsSectionProps {
  cardsSent: number;
  cardsReceived: number;
  milestones: UserMilestone[];
  loading?: boolean;
}

// Generate all possible achievements
const generateAllAchievements = (
  cardsSent: number,
  cardsReceived: number,
  milestones: UserMilestone[]
): Achievement[] => {
  const achievements: Achievement[] = [];
  const totalCards = cardsSent + cardsReceived;

  // Helper to find milestone data
  const findMilestone = (type: string, threshold: number) =>
    milestones.find((m) => m.type === type && m.threshold === threshold);

  // Received achievements
  ACHIEVEMENT_THRESHOLDS.forEach((threshold) => {
    const milestone = findMilestone("CardsReceived", threshold);
    const unlocked = cardsReceived >= threshold;
    achievements.push({
      id: `received-${threshold}`,
      type: "received" as AchievementType,
      threshold,
      title: milestone?.title || `Received ${threshold} Cards`,
      description:
        milestone?.description ||
        `Receive ${threshold} recognition cards from your colleagues.`,
      achievedAt: milestone?.achievedAt,
      unlocked,
    });
  });

  // Sent achievements
  ACHIEVEMENT_THRESHOLDS.forEach((threshold) => {
    const milestone = findMilestone("CardsSent", threshold);
    const unlocked = cardsSent >= threshold;
    achievements.push({
      id: `sent-${threshold}`,
      type: "sent" as AchievementType,
      threshold,
      title: milestone?.title || `Sent ${threshold} Cards`,
      description:
        milestone?.description ||
        `Send ${threshold} recognition cards to appreciate your colleagues.`,
      achievedAt: milestone?.achievedAt,
      unlocked,
    });
  });

  // Total achievements
  ACHIEVEMENT_THRESHOLDS.forEach((threshold) => {
    const unlocked = totalCards >= threshold;
    achievements.push({
      id: `total-${threshold}`,
      type: "total" as AchievementType,
      threshold,
      title: `${threshold} Cards Total`,
      description: `Reach a total of ${threshold} cards (sent + received combined).`,
      achievedAt: unlocked ? new Date().toISOString() : undefined,
      unlocked,
    });
  });

  return achievements;
};

export const AchievementsSection: React.FC<AchievementsSectionProps> = ({
  cardsSent,
  cardsReceived,
  milestones,
  loading = false,
}) => {
  const [selectedAchievement, setSelectedAchievement] = useState<Achievement | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const achievements = useMemo(
    () => generateAllAchievements(cardsSent, cardsReceived, milestones),
    [cardsSent, cardsReceived, milestones]
  );

  const unlockedCount = achievements.filter((a) => a.unlocked).length;

  const handleBadgeClick = (achievement: Achievement) => {
    setSelectedAchievement(achievement);
    setDialogOpen(true);
  };

  if (loading) {
    return (
      <Card className="bg-glass-bg backdrop-blur-md border-glass-border">
        <CardHeader>
          <Skeleton className="h-6 w-40" />
          <Skeleton className="h-4 w-60" />
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-12 gap-2">
            {Array.from({ length: 12 }).map((_, i) => (
              <Skeleton key={i} className="h-24 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  // Group achievements by type for better display
  const receivedAchievements = achievements.filter((a) => a.type === "received");
  const sentAchievements = achievements.filter((a) => a.type === "sent");
  const totalAchievements = achievements.filter((a) => a.type === "total");

  return (
    <>
      <Card className="bg-glass-bg backdrop-blur-md border-glass-border">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Trophy className="h-5 w-5 text-yellow-500" />
            My Achievements
          </CardTitle>
          <CardDescription>
            {unlockedCount} of {achievements.length} achievements unlocked
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            {/* Received Row */}
            <div>
              <h4 className="text-sm font-medium text-muted-foreground mb-3">
                Cards Received
              </h4>
              <div className="flex flex-wrap gap-2">
                {receivedAchievements.map((achievement) => (
                  <AchievementBadge
                    key={achievement.id}
                    achievement={achievement}
                    onClick={() => handleBadgeClick(achievement)}
                  />
                ))}
              </div>
            </div>

            {/* Sent Row */}
            <div>
              <h4 className="text-sm font-medium text-muted-foreground mb-3">
                Cards Sent
              </h4>
              <div className="flex flex-wrap gap-2">
                {sentAchievements.map((achievement) => (
                  <AchievementBadge
                    key={achievement.id}
                    achievement={achievement}
                    onClick={() => handleBadgeClick(achievement)}
                  />
                ))}
              </div>
            </div>

            {/* Total Row */}
            <div>
              <h4 className="text-sm font-medium text-muted-foreground mb-3">
                Total Activity
              </h4>
              <div className="flex flex-wrap gap-2">
                {totalAchievements.map((achievement) => (
                  <AchievementBadge
                    key={achievement.id}
                    achievement={achievement}
                    onClick={() => handleBadgeClick(achievement)}
                  />
                ))}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <AchievementDialog
        achievement={selectedAchievement}
        open={dialogOpen}
        onOpenChange={setDialogOpen}
      />
    </>
  );
};

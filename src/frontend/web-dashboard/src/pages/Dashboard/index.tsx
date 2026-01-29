import React, { useEffect, useState } from "react";
import { useAuthStore } from "@/store/authStore";
import { api } from "@/services/api";
import type { UserStats, TopEmployee } from "@/services/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Send, Inbox, Trophy, TrendingUp, Award } from "lucide-react";

export const Dashboard: React.FC = () => {
  const { user } = useAuthStore();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [topEmployees, setTopEmployees] = useState<TopEmployee[]>([]);
  const [loading, setLoading] = useState(true);

  const employeeId = user?.profile?.sub || "";

  useEffect(() => {
    const fetchData = async () => {
      if (!employeeId) return;

      try {
        setLoading(true);
        const [statsRes, topRes] = await Promise.all([
          api.getUserStats(employeeId),
          api.getTopEmployees(5),
        ]);

        setStats(statsRes.data.data);
        setTopEmployees(topRes.data.data);
      } catch (error) {
        console.error("Failed to fetch dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [employeeId]);

  const getInitials = (id: string) => id.substring(0, 2).toUpperCase();

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          Welcome back, {user?.profile.name || "User"}!
        </h1>
        <p className="text-muted-foreground">
          Here's what's happening with your recognitions.
        </p>
      </div>

      {/* Stats Overview */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-glass-bg backdrop-blur-md border-glass-border">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Cards Received
            </CardTitle>
            <Inbox className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold">
                {stats?.cardsReceived || 0}
              </div>
            )}
            <p className="text-xs text-muted-foreground">Lifetime total</p>
          </CardContent>
        </Card>
        <Card className="bg-glass-bg backdrop-blur-md border-glass-border">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cards Sent</CardTitle>
            <Send className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold">{stats?.cardsSent || 0}</div>
            )}
            <p className="text-xs text-muted-foreground">Recognition given</p>
          </CardContent>
        </Card>
        <Card className="bg-glass-bg backdrop-blur-md border-glass-border">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Milestones</CardTitle>
            <Trophy className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-8 w-16" />
            ) : (
              <div className="text-2xl font-bold">
                {stats?.milestones.length || 0}
              </div>
            )}
            <p className="text-xs text-muted-foreground">Achievements earned</p>
          </CardContent>
        </Card>
        <Card className="bg-glass-bg backdrop-blur-md border-glass-border">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Top Value</CardTitle>
            <TrendingUp className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-8 w-24" />
            ) : (
              <div className="text-xl font-bold truncate">
                {stats?.valueStats[0]?.valueName || "None yet"}
              </div>
            )}
            <p className="text-xs text-muted-foreground">Most recognized for</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        {/* Leaderboard */}
        <Card className="col-span-full bg-glass-bg backdrop-blur-md border-glass-border">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="h-5 w-5 text-yellow-500" />
              Leaderboard
            </CardTitle>
            <CardDescription>
              Most recognized employees this period.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-4">
                {[1, 2, 3, 4].map((i) => (
                  <Skeleton key={i} className="h-10 w-full" />
                ))}
              </div>
            ) : (
              <div className="space-y-6">
                {topEmployees.map((emp, index) => (
                  <div key={emp.employeeId} className="flex items-center">
                    <div className="mr-4 flex h-8 w-8 items-center justify-center font-bold text-muted-foreground">
                      #{index + 1}
                    </div>
                    <Avatar className="h-9 w-9">
                      <AvatarFallback>
                        {getInitials(emp.employeeId)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="ml-4 space-y-1">
                      <p className="text-sm font-medium leading-none">
                        {emp.employeeId}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Recognition received
                      </p>
                    </div>
                    <div className="ml-auto font-medium">
                      <Badge
                        variant="secondary"
                        className="bg-primary/10 text-primary hover:bg-primary/20 border-none"
                      >
                        {emp.cardsReceived}
                      </Badge>
                    </div>
                  </div>
                ))}
                {topEmployees.length === 0 && (
                  <div className="text-center py-8 text-muted-foreground">
                    No data available yet.
                  </div>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

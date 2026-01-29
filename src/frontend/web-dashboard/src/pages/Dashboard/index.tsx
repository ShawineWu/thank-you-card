import React, { useEffect, useState } from "react";
import { useAuthStore } from "@/store/authStore";
import { api } from "@/services/api";
import type {
  UserStats,
  TopEmployee,
  Card as RecognitionCard,
} from "@/services/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Send, Inbox, Trophy, TrendingUp, Clock, Award } from "lucide-react";
import { format } from "date-fns";

export const Dashboard: React.FC = () => {
  const { user } = useAuthStore();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [topEmployees, setTopEmployees] = useState<TopEmployee[]>([]);
  const [recentCards, setRecentCards] = useState<RecognitionCard[]>([]);
  const [loading, setLoading] = useState(true);

  const employeeId = user?.profile?.sub || "";

  useEffect(() => {
    const fetchData = async () => {
      if (!employeeId) return;

      try {
        setLoading(true);
        const [statsRes, topRes, cardsRes] = await Promise.all([
          api.getUserStats(employeeId),
          api.getTopEmployees(5),
          api.getCards({ page: 1, pageSize: 5 }),
        ]);

        setStats(statsRes.data.data);
        setTopEmployees(topRes.data.data);
        setRecentCards(cardsRes.data.data.data);
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
        {/* Recent Recognitions */}
        <Card className="col-span-4 bg-glass-bg backdrop-blur-md border-glass-border">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5" />
              Recent Feed
            </CardTitle>
            <CardDescription>
              Latest thank you cards from across the company.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Sender</TableHead>
                    <TableHead>Reason</TableHead>
                    <TableHead className="text-right">Date</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recentCards.map((card) => (
                    <TableRow key={card.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <Avatar className="h-8 w-8">
                            <AvatarFallback className="bg-primary/20 text-xs">
                              {getInitials(card.senderId)}
                            </AvatarFallback>
                          </Avatar>
                          <span>{card.senderId}</span>
                        </div>
                      </TableCell>
                      <TableCell className="max-w-[200px] truncate">
                        {card.recognitionReason}
                      </TableCell>
                      <TableCell className="text-right text-muted-foreground text-xs">
                        {format(new Date(card.createdAt), "MMM d, yyyy")}
                      </TableCell>
                    </TableRow>
                  ))}
                  {recentCards.length === 0 && (
                    <TableRow>
                      <TableCell
                        colSpan={3}
                        className="text-center py-8 text-muted-foreground"
                      >
                        No recognitions yet.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>

        {/* Top 10 Recognized */}
        <Card className="col-span-3 bg-glass-bg backdrop-blur-md border-glass-border">
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

import React, { useEffect, useState, useMemo, useCallback } from "react";
import { useAuthStore } from "@/store/authStore";
import { api } from "@/services/api";
import type {
  UserStats,
  Card as CardType,
  PaginatedResponse,
  CardFilter,
} from "@/services/api";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { CardList } from "@/components/Cards/CardList";
import { CardFilters } from "@/components/Cards/CardFilters";
import {
  User,
  Send,
  Inbox,
  Trophy,
  TrendingUp,
} from "lucide-react";
import { AchievementsSection } from "@/components/Achievements";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  PieChart,
  Pie,
} from "recharts";

export const Profile: React.FC = () => {
  const { user } = useAuthStore();
  const employeeId = user?.profile?.sub || "";

  // Overview data
  const [stats, setStats] = useState<UserStats | null>(null);
  const [loadingStats, setLoadingStats] = useState(true);

  // Received cards data
  const [receivedData, setReceivedData] = useState<PaginatedResponse<CardType> | null>(null);
  const [loadingReceived, setLoadingReceived] = useState(true);
  const [receivedPage, setReceivedPage] = useState(1);
  const [receivedFilters, setReceivedFilters] = useState<Partial<CardFilter>>({});

  // Sent cards data
  const [sentData, setSentData] = useState<PaginatedResponse<CardType> | null>(null);
  const [loadingSent, setLoadingSent] = useState(true);
  const [sentPage, setSentPage] = useState(1);
  const [sentFilters, setSentFilters] = useState<Partial<CardFilter>>({});

  // Fetch stats for Overview tab
  useEffect(() => {
    const fetchStats = async () => {
      if (!employeeId) return;
      try {
        setLoadingStats(true);
        const res = await api.getUserStats(employeeId);
        setStats(res.data.data);
      } catch (error) {
        console.error("Failed to fetch user stats:", error);
      } finally {
        setLoadingStats(false);
      }
    };

    fetchStats();
  }, [employeeId]);

  // Fetch received cards
  const fetchReceivedCards = useCallback(
    async (pageNum: number, currentFilters: Partial<CardFilter>) => {
      if (!employeeId) return;

      try {
        setLoadingReceived(true);
        const res = await api.getCards({
          ...currentFilters,
          page: pageNum,
          pageSize: 10,
          recipientId: employeeId,
        });
        setReceivedData(res.data.data);
      } catch (error) {
        console.error("Failed to fetch received cards:", error);
      } finally {
        setLoadingReceived(false);
      }
    },
    [employeeId]
  );

  // Fetch sent cards
  const fetchSentCards = useCallback(
    async (pageNum: number, currentFilters: Partial<CardFilter>) => {
      if (!employeeId) return;

      try {
        setLoadingSent(true);
        const res = await api.getCards({
          ...currentFilters,
          page: pageNum,
          pageSize: 10,
          senderId: employeeId,
        });
        setSentData(res.data.data);
      } catch (error) {
        console.error("Failed to fetch sent cards:", error);
      } finally {
        setLoadingSent(false);
      }
    },
    [employeeId]
  );

  // Debounce received filters
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchReceivedCards(receivedPage, receivedFilters);
    }, 400);

    return () => clearTimeout(timer);
  }, [receivedPage, receivedFilters, fetchReceivedCards]);

  // Debounce sent filters
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchSentCards(sentPage, sentFilters);
    }, 400);

    return () => clearTimeout(timer);
  }, [sentPage, sentFilters, fetchSentCards]);

  const handleReceivedFilterChange = (newFilters: Partial<CardFilter>) => {
    setReceivedFilters(newFilters);
    setReceivedPage(1);
  };

  const handleSentFilterChange = (newFilters: Partial<CardFilter>) => {
    setSentFilters(newFilters);
    setSentPage(1);
  };

  const chartData = useMemo(() => {
    if (!stats) return [];
    return stats.valueStats
      .map((v) => ({
        name: v.valueName,
        count: v.count,
      }))
      .sort((a, b) => b.count - a.count);
  }, [stats]);

  const pieData = useMemo(() => {
    if (!stats) return [];
    return [
      { name: "Sent", value: stats.cardsSent, color: "#6366f1" },
      { name: "Received", value: stats.cardsReceived, color: "#ec4899" },
    ];
  }, [stats]);

  const COLORS = [
    "#6366f1",
    "#8b5cf6",
    "#ec4899",
    "#f43f5e",
    "#f59e0b",
    "#10b981",
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <div className="p-3 rounded-2xl bg-primary/10 text-primary">
          <User size={24} />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            {user?.profile.name || "User Profile"}
          </h1>
          <p className="text-muted-foreground">
            Your recognition journey and impact
          </p>
        </div>
      </div>

      <Tabs defaultValue="overview" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="received">Received Cards</TabsTrigger>
          <TabsTrigger value="sent">Sent Cards</TabsTrigger>
          <TabsTrigger value="account">Account</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card className="bg-white border border-gray-100 shadow-sm hover:shadow-md transition-shadow rounded-xl">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-gray-500">
                  Cards Received
                </CardTitle>
                <div className="p-2 rounded-lg bg-indigo-50">
                  <Inbox className="h-4 w-4 text-indigo-500" />
                </div>
              </CardHeader>
              <CardContent>
                {loadingStats ? (
                  <Skeleton className="h-9 w-20" />
                ) : (
                  <div className="text-3xl font-bold text-gray-900">{stats?.cardsReceived || 0}</div>
                )}
                <p className="text-xs text-gray-400 mt-1">
                  Lifetime total
                </p>
              </CardContent>
            </Card>
            <Card className="bg-white border border-gray-100 shadow-sm hover:shadow-md transition-shadow rounded-xl">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-gray-500">
                  Cards Sent
                </CardTitle>
                <div className="p-2 rounded-lg bg-blue-50">
                  <Send className="h-4 w-4 text-blue-500" />
                </div>
              </CardHeader>
              <CardContent>
                {loadingStats ? (
                  <Skeleton className="h-9 w-20" />
                ) : (
                  <div className="text-3xl font-bold text-gray-900">
                    {stats?.cardsSent || 0}
                  </div>
                )}
                <p className="text-xs text-gray-400 mt-1">
                  Recognition given
                </p>
              </CardContent>
            </Card>
            <Card className="bg-white border border-gray-100 shadow-sm hover:shadow-md transition-shadow rounded-xl">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-gray-500">
                  Milestones
                </CardTitle>
                <div className="p-2 rounded-lg bg-amber-50">
                  <Trophy className="h-4 w-4 text-amber-500" />
                </div>
              </CardHeader>
              <CardContent>
                {loadingStats ? (
                  <Skeleton className="h-9 w-20" />
                ) : (
                  <div className="text-3xl font-bold text-gray-900">
                    {stats?.milestones.length || 0}
                  </div>
                )}
                <p className="text-xs text-gray-400 mt-1">Achievements earned</p>
              </CardContent>
            </Card>
            <Card className="bg-white border border-gray-100 shadow-sm hover:shadow-md transition-shadow rounded-xl">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-gray-500">
                  Top Value
                </CardTitle>
                <div className="p-2 rounded-lg bg-emerald-50">
                  <TrendingUp className="h-4 w-4 text-emerald-500" />
                </div>
              </CardHeader>
              <CardContent>
                {loadingStats ? (
                  <Skeleton className="h-9 w-24" />
                ) : (
                  <div className="text-xl font-bold text-gray-900 truncate">
                    {chartData[0]?.name || "None yet"}
                  </div>
                )}
                <p className="text-xs text-gray-400 mt-1">
                  Most recognized for
                </p>
              </CardContent>
            </Card>
          </div>

          <div className="grid gap-6 md:grid-cols-2">
            <Card className="bg-white border border-gray-100 shadow-sm rounded-xl">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-gray-800">
                  <div className="p-1.5 rounded-lg bg-indigo-50">
                    <TrendingUp size={16} className="text-indigo-500" />
                  </div>
                  Recent Feed
                </CardTitle>
                <CardDescription className="text-gray-400">
                  Latest thank you cards from across the company.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="min-h-[200px]">
                  {loadingStats ? (
                    <Skeleton className="h-full w-full rounded-xl" />
                  ) : chartData.length > 0 ? (
                    <div className="space-y-3">
                      <div className="grid grid-cols-3 text-xs font-medium text-gray-500 pb-2 border-b border-gray-100">
                        <span>Sender</span>
                        <span>Reason</span>
                        <span className="text-right">Date</span>
                      </div>
                      <ResponsiveContainer width="100%" height={200}>
                        <BarChart
                          data={chartData}
                          layout="vertical"
                          margin={{ left: 20, right: 20 }}
                        >
                          <CartesianGrid
                            strokeDasharray="3 3"
                            horizontal={true}
                            vertical={false}
                            stroke="#f1f5f9"
                          />
                          <XAxis type="number" hide />
                          <YAxis
                            dataKey="name"
                            type="category"
                            width={80}
                            axisLine={false}
                            tickLine={false}
                            tick={{ fill: "#64748b", fontSize: 12 }}
                          />
                          <Tooltip
                            contentStyle={{
                              backgroundColor: "#ffffff",
                              border: "1px solid #e2e8f0",
                              borderRadius: "8px",
                              boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                            }}
                            itemStyle={{ color: "#1e293b" }}
                            cursor={{ fill: "rgba(0,0,0,0.02)" }}
                          />
                          <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={20}>
                            {chartData.map((_, index) => (
                              <Cell
                                key={`cell-${index}`}
                                fill={COLORS[index % COLORS.length]}
                              />
                            ))}
                          </Bar>
                        </BarChart>
                      </ResponsiveContainer>
                    </div>
                  ) : (
                    <div className="h-full flex items-center justify-center text-gray-400 italic py-12">
                      No recognitions yet.
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            <Card className="bg-white border border-gray-100 shadow-sm rounded-xl">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-gray-800">
                  <div className="p-1.5 rounded-lg bg-amber-50">
                    <Trophy size={16} className="text-amber-500" />
                  </div>
                  Leaderboard
                </CardTitle>
                <CardDescription className="text-gray-400">
                  Most recognized employees this period.
                </CardDescription>
              </CardHeader>
              <CardContent className="flex flex-col items-center">
                <div className="min-h-[200px] w-full">
                  {loadingStats ? (
                    <Skeleton className="h-full w-2/3 mx-auto rounded-full" />
                  ) : (stats?.cardsSent || 0) + (stats?.cardsReceived || 0) >
                    0 ? (
                    <ResponsiveContainer width="100%" height={200}>
                      <PieChart>
                        <Pie
                          data={pieData}
                          cx="50%"
                          cy="50%"
                          innerRadius={50}
                          outerRadius={70}
                          paddingAngle={5}
                          dataKey="value"
                        >
                          {pieData.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={entry.color} />
                          ))}
                        </Pie>
                        <Tooltip
                          contentStyle={{
                            backgroundColor: "#ffffff",
                            border: "1px solid #e2e8f0",
                            borderRadius: "8px",
                            boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                          }}
                          itemStyle={{ color: "#1e293b" }}
                        />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <div className="h-full flex items-center justify-center text-gray-400 italic py-12">
                      No data available yet.
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Achievements Section */}
          <AchievementsSection
            cardsSent={stats?.cardsSent || 0}
            cardsReceived={stats?.cardsReceived || 0}
            milestones={stats?.milestones || []}
            loading={loadingStats}
          />
        </TabsContent>

        {/* Received Cards Tab */}
        <TabsContent value="received" className="space-y-6">
          <CardFilters
            onFilterChange={handleReceivedFilterChange}
            showRecipientFilter={false}
          />
          <CardList
            data={receivedData}
            loading={loadingReceived}
            onPageChange={setReceivedPage}
            showRecipient={false}
          />
        </TabsContent>

        {/* Sent Cards Tab */}
        <TabsContent value="sent" className="space-y-6">
          <CardFilters
            onFilterChange={handleSentFilterChange}
            showSenderFilter={false}
          />
          <CardList
            data={sentData}
            loading={loadingSent}
            onPageChange={setSentPage}
            showSender={false}
          />
        </TabsContent>

        {/* Account Tab */}
        <TabsContent value="account" className="space-y-6">
          <Card className="bg-white border border-gray-100 shadow-sm rounded-xl">
            <CardHeader>
              <CardTitle className="text-gray-800">Account Information</CardTitle>
              <CardDescription className="text-gray-400">
                Your profile and authentication details
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-4">
                <Avatar className="h-16 w-16 bg-gradient-to-br from-indigo-400 to-purple-500">
                  <AvatarFallback className="text-xl text-white bg-transparent">
                    {user?.profile?.name?.substring(0, 2).toUpperCase() || "U"}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <h3 className="text-lg font-semibold text-gray-800">
                    {user?.profile?.name || "User"}
                  </h3>
                  <p className="text-sm text-gray-400">
                    {user?.profile?.email || "No email"}
                  </p>
                </div>
              </div>

              <div className="grid gap-3 pt-4 border-t border-gray-100">
                <div className="grid grid-cols-3 gap-4 py-2">
                  <div className="text-sm font-medium text-gray-400">
                    User ID
                  </div>
                  <div className="col-span-2 text-sm text-gray-600">
                    {user?.profile?.sub || "N/A"}
                  </div>
                </div>
                <div className="grid grid-cols-3 gap-4 py-2">
                  <div className="text-sm font-medium text-gray-400">
                    Email
                  </div>
                  <div className="col-span-2 text-sm text-gray-600">
                    {user?.profile?.email || "N/A"}
                  </div>
                </div>
                <div className="grid grid-cols-3 gap-4 py-2">
                  <div className="text-sm font-medium text-gray-400">
                    Username
                  </div>
                  <div className="col-span-2 text-sm text-gray-600">
                    {user?.profile?.preferred_username || "N/A"}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

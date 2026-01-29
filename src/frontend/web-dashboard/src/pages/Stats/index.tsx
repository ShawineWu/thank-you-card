import React, { useEffect, useState, useMemo } from "react";
import { useAuthStore } from "@/store/authStore";
import { api } from "@/services/api";
import type { UserStats } from "@/services/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
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
import {
  BarChart3,
  Trophy,
  Award,
  TrendingUp,
  Send,
  Inbox,
  Star,
} from "lucide-react";

export const Stats: React.FC = () => {
  const { user } = useAuthStore();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [loading, setLoading] = useState(true);

  const employeeId = user?.profile?.sub || "";

  useEffect(() => {
    const fetchStats = async () => {
      if (!employeeId) return;
      try {
        setLoading(true);
        const res = await api.getUserStats(employeeId);
        setStats(res.data.data);
      } catch (error) {
        console.error("Failed to fetch user stats:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [employeeId]);

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
    <div className="space-y-8 pb-10">
      <div className="flex items-center gap-3">
        <div className="p-3 rounded-2xl bg-primary/10 text-primary">
          <BarChart3 size={24} />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Your Statistics</h1>
          <p className="text-muted-foreground">
            Insights into your impact and recognitions within the company.
          </p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Total Sent
            </CardTitle>
            <Send className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">{stats?.cardsSent}</div>
            )}
            <p className="text-xs text-slate-500 mt-1">
              Recognitions given to others
            </p>
          </CardContent>
        </Card>
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Total Received
            </CardTitle>
            <Inbox className="h-4 w-4 text-pink-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">{stats?.cardsReceived}</div>
            )}
            <p className="text-xs text-slate-500 mt-1">
              Recognitions received from peers
            </p>
          </CardContent>
        </Card>
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Achievements
            </CardTitle>
            <Trophy className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">
                {stats?.milestones.length}
              </div>
            )}
            <p className="text-xs text-slate-500 mt-1">Milestones reached</p>
          </CardContent>
        </Card>
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Lead Value
            </CardTitle>
            <Award className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-24" />
            ) : (
              <div className="text-xl font-bold truncate">
                {chartData[0]?.name || "N/A"}
              </div>
            )}
            <p className="text-xs text-slate-500 mt-1">
              Your strongest attribute
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4 bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp size={20} className="text-primary" />
              Value Recognition Distribution
            </CardTitle>
            <CardDescription>
              Breakdown of which company values you are recognized for.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px] w-full mt-4">
              {loading ? (
                <Skeleton className="h-full w-full rounded-xl" />
              ) : chartData.length > 0 ? (
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={chartData}
                    layout="vertical"
                    margin={{ left: 40, right: 30 }}
                  >
                    <CartesianGrid
                      strokeDasharray="3 3"
                      horizontal={true}
                      vertical={false}
                      stroke="#1e293b"
                    />
                    <XAxis type="number" hide />
                    <YAxis
                      dataKey="name"
                      type="category"
                      width={100}
                      axisLine={false}
                      tickLine={false}
                      tick={{ fill: "#94a3b8", fontSize: 12 }}
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #1e293b",
                        borderRadius: "8px",
                      }}
                      itemStyle={{ color: "#f8fafc" }}
                      cursor={{ fill: "rgba(255,255,255,0.05)" }}
                    />
                    <Bar dataKey="count" radius={[0, 4, 4, 0]} barSize={24}>
                      {chartData.map((_, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={COLORS[index % COLORS.length]}
                        />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              ) : (
                <div className="h-full flex items-center justify-center text-slate-500 italic">
                  No data points yet.
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-3 bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Star size={20} className="text-secondary" />
              Impact Ratio
            </CardTitle>
            <CardDescription>Sent vs Received recognitions.</CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col items-center">
            <div className="h-[200px] w-full mt-2">
              {loading ? (
                <Skeleton className="h-full w-2/3 mx-auto rounded-full" />
              ) : (stats?.cardsSent || 0) + (stats?.cardsReceived || 0) > 0 ? (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={pieData}
                      cx="50%"
                      cy="50%"
                      innerRadius={60}
                      outerRadius={80}
                      paddingAngle={5}
                      dataKey="value"
                    >
                      {pieData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #1e293b",
                        borderRadius: "8px",
                      }}
                      itemStyle={{ color: "#f8fafc" }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              ) : (
                <div className="h-full flex items-center justify-center text-slate-500 italic">
                  No activity yet.
                </div>
              )}
            </div>
            <div className="grid grid-cols-2 gap-8 mt-4 w-full px-6">
              <div className="text-center">
                <div className="flex items-center justify-center gap-2 mb-1">
                  <div className="w-3 h-3 rounded-full bg-[#6366f1]" />
                  <span className="text-sm font-medium text-slate-300">
                    Sent
                  </span>
                </div>
                <div className="text-2xl font-bold">
                  {stats?.cardsSent || 0}
                </div>
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center gap-2 mb-1">
                  <div className="w-3 h-3 rounded-full bg-[#ec4899]" />
                  <span className="text-sm font-medium text-slate-300">
                    Received
                  </span>
                </div>
                <div className="text-2xl font-bold">
                  {stats?.cardsReceived || 0}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="space-y-4">
        <h2 className="text-xl font-bold flex items-center gap-2">
          <Trophy size={20} className="text-yellow-500" />
          Milestones & Achievements
        </h2>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {loading ? (
            [1, 2, 3].map((i) => (
              <Skeleton key={i} className="h-40 w-full rounded-2xl" />
            ))
          ) : stats?.milestones && stats.milestones.length > 0 ? (
            stats.milestones.map((m, idx) => (
              <Card
                key={idx}
                className="bg-slate-950 border-slate-800 overflow-hidden relative group transition-all hover:border-primary/50"
              >
                <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
                  <Trophy size={64} className="text-primary" />
                </div>
                <CardHeader>
                  <Badge className="w-fit bg-primary/20 text-primary border-primary/30 mb-2">
                    {m.type}
                  </Badge>
                  <CardTitle className="text-lg">{m.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-slate-400 line-clamp-2">
                    {m.description}
                  </p>
                </CardContent>
              </Card>
            ))
          ) : (
            <div className="col-span-full py-12 text-center text-slate-500 border border-dashed border-slate-800 rounded-3xl">
              Keep engaging to earn your first milestone!
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

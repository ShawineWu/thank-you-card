import React, { useEffect, useState } from "react";
import { api } from "@/services/api";
import type {
  DashboardAnalytics,
  TopRecognizer,
  ValueDistribution,
} from "@/services/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Input } from "@/components/ui/input";
import {
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import {
  BarChart3,
  Download,
  TrendingUp,
  Users,
  Award,
  FileText,
  AlertCircle,
} from "lucide-react";
import { Badge } from "@/components/ui/badge";

export const HRAnalytics: React.FC = () => {
  const [dashboard, setDashboard] = useState<DashboardAnalytics | null>(null);
  const [topRecognizers, setTopRecognizers] = useState<TopRecognizer[]>([]);
  const [valueDistribution, setValueDistribution] =
    useState<ValueDistribution | null>(null);
  const [loading, setLoading] = useState(true);
  const [exportLoading, setExportLoading] = useState(false);
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [dashboardRes, recognizersRes, distributionRes] =
          await Promise.all([
            api.getAnalyticsDashboard(),
            api.getTopRecognizers(10),
            api.getValueDistribution(),
          ]);

        setDashboard(dashboardRes.data.data);
        setTopRecognizers(recognizersRes.data.data.recognizers);
        setValueDistribution(distributionRes.data.data);
      } catch (error) {
        console.error("Failed to fetch analytics:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleExport = async () => {
    if (!startDate || !endDate) {
      alert("Please select both start and end dates");
      return;
    }

    try {
      setExportLoading(true);
      const response = await api.exportAnalytics({
        startDate,
        endDate,
        format: "csv",
      });

      // Create download link
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute(
        "download",
        `recognition_cards_${startDate}_${endDate}.csv`,
      );
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error("Export failed:", error);
      alert("Failed to export data");
    } finally {
      setExportLoading(false);
    }
  };

  const COLORS = ["#6366f1", "#8b5cf6", "#ec4899", "#f43f5e", "#f59e0b"];

  const trendData = dashboard
    ? [
        { name: "Last Month", cards: dashboard.cardsTrend.lastMonth },
        { name: "This Month", cards: dashboard.cardsTrend.thisMonth },
      ]
    : [];

  return (
    <div className="space-y-8 pb-10">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-3 rounded-2xl bg-primary/10 text-primary">
          <BarChart3 size={24} />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            HR Analytics Dashboard
          </h1>
          <p className="text-muted-foreground">
            Comprehensive insights into recognition patterns and team
            engagement.
          </p>
        </div>
      </div>

      {/* Overview Stats */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Total Cards
            </CardTitle>
            <FileText className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">
                {dashboard?.totalCards.toLocaleString()}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Active Users
            </CardTitle>
            <Users className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">
                {dashboard?.activeUsers.toLocaleString()}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              Engagement Rate
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-pink-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">
                {dashboard?.engagementRate.toFixed(1)}%
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">
              This Month
            </CardTitle>
            <Award className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-20" />
            ) : (
              <div className="text-3xl font-bold">
                {dashboard?.cardsTrend.thisMonth.toLocaleString()}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Charts Row */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Monthly Trend */}
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Monthly Trend</CardTitle>
            <CardDescription>Recognition cards sent over time</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px] w-full">
              {loading ? (
                <Skeleton className="h-full w-full rounded-xl" />
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={trendData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                    <XAxis dataKey="name" stroke="#94a3b8" />
                    <YAxis stroke="#94a3b8" />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #1e293b",
                        borderRadius: "8px",
                      }}
                      itemStyle={{ color: "#f8fafc" }}
                    />
                    <Bar dataKey="cards" fill="#6366f1" radius={[8, 8, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Value Distribution */}
        <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
          <CardHeader>
            <CardTitle>Value Distribution</CardTitle>
            <CardDescription>Most recognized company values</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px] w-full">
              {loading ? (
                <Skeleton className="h-full w-full rounded-xl" />
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={valueDistribution?.distribution.slice(0, 5)}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={(props) => {
                        const data = valueDistribution?.distribution.slice(
                          0,
                          5,
                        )[props.index];
                        return data
                          ? `${data.valueName} (${data.percentage.toFixed(0)}%)`
                          : "";
                      }}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="count"
                    >
                      {valueDistribution?.distribution
                        .slice(0, 5)
                        .map((_entry, index) => (
                          <Cell
                            key={`cell-${index}`}
                            fill={COLORS[index % COLORS.length]}
                          />
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
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Recognizers */}
      <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Award size={20} className="text-yellow-500" />
            Top Recognizers
          </CardTitle>
          <CardDescription>
            Most active employees sending recognition cards
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-12 w-full rounded-xl" />
              ))}
            </div>
          ) : (
            <div className="space-y-3">
              {topRecognizers.map((recognizer) => (
                <div
                  key={recognizer.employeeId}
                  className="flex items-center justify-between p-3 rounded-xl bg-slate-900/50 border border-slate-800 hover:border-primary/50 transition-colors"
                >
                  <div className="flex items-center gap-3">
                    <Badge
                      variant={recognizer.rank <= 3 ? "default" : "outline"}
                      className={
                        recognizer.rank <= 3
                          ? "bg-yellow-500/20 text-yellow-500 border-yellow-500/30"
                          : "border-slate-700"
                      }
                    >
                      #{recognizer.rank}
                    </Badge>
                    <span className="font-medium">{recognizer.employeeId}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-2xl font-bold text-primary">
                      {recognizer.cardsSent}
                    </span>
                    <span className="text-sm text-slate-500">cards</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Export Section */}
      <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Download size={20} className="text-primary" />
            Export Data
          </CardTitle>
          <CardDescription>
            Download recognition data as CSV for further analysis
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap items-end gap-4">
            <div className="flex-1 min-w-[200px]">
              <label className="text-sm font-medium text-slate-400 mb-2 block">
                Start Date
              </label>
              <Input
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="bg-slate-950 border-slate-800"
              />
            </div>
            <div className="flex-1 min-w-[200px]">
              <label className="text-sm font-medium text-slate-400 mb-2 block">
                End Date
              </label>
              <Input
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="bg-slate-950 border-slate-800"
              />
            </div>
            <Button
              onClick={handleExport}
              disabled={exportLoading || !startDate || !endDate}
              className="gap-2"
            >
              <Download size={16} />
              {exportLoading ? "Exporting..." : "Export CSV"}
            </Button>
          </div>

          {(!startDate || !endDate) && (
            <div className="flex items-center gap-2 mt-4 text-sm text-slate-500">
              <AlertCircle size={14} />
              Please select both start and end dates to export
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

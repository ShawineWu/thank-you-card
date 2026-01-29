import React, { useEffect, useState, useMemo, useCallback } from "react";
import { api } from "@/services/api";
import type {
  DashboardAnalytics,
  TopRecognizer,
  ValueDistribution,
  Card as RecognitionCard,
} from "@/services/api";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
} from "recharts";
import {
  BarChart3,
  Download,
  Users,
  FileText,
  Building2,
  Calendar,
  TrendingUp,
  Tag,
} from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { getValueBadgeClasses, VALUE_COLOR_MAP, DEFAULT_VALUE_COLOR } from "@/constants/valueColors";

type TabType = "dashboard" | "recognizers" | "teams" | "values";

export const HRAnalytics: React.FC = () => {
  const [activeTab, setActiveTab] = useState<TabType>("dashboard");
  const [dashboard, setDashboard] = useState<DashboardAnalytics | null>(null);
  const [topRecognizers, setTopRecognizers] = useState<TopRecognizer[]>([]);
  const [valueDistribution, setValueDistribution] = useState<ValueDistribution | null>(null);
  const [loading, setLoading] = useState(true);
  const [exportLoading, setExportLoading] = useState(false);
  
  // Date range state - default to last 30 days
  const [startDate, setStartDate] = useState(() => {
    const date = new Date();
    date.setDate(date.getDate() - 30);
    return date.toISOString().split("T")[0];
  });
  const [endDate, setEndDate] = useState(() => {
    return new Date().toISOString().split("T")[0];
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [dashboardRes, recognizersRes, distributionRes] = await Promise.all([
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
    if (!startDate || !endDate) return;

    try {
      setExportLoading(true);
      const response = await api.exportAnalytics({
        startDate,
        endDate,
        format: "csv",
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", `recognition_cards_${startDate}_${endDate}.csv`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error("Export failed:", error);
    } finally {
      setExportLoading(false);
    }
  };

  const handleLast30Days = () => {
    const end = new Date();
    const start = new Date();
    start.setDate(start.getDate() - 30);
    setStartDate(start.toISOString().split("T")[0]);
    setEndDate(end.toISOString().split("T")[0]);
  };

  const tabs: { id: TabType; label: string; icon: React.ReactNode }[] = [
    { id: "dashboard", label: "Dashboard", icon: <BarChart3 size={16} /> },
    { id: "recognizers", label: "Recognizers", icon: <Users size={16} /> },
    { id: "teams", label: "Teams", icon: <Building2 size={16} /> },
    { id: "values", label: "Values", icon: <Tag size={16} /> },
  ];

  // Mock data for departments (in real app, this would come from API)
  const departmentStats = useMemo(() => {
    const departments = new Set<string>();
    if (topRecognizers.length > 0) {
      // Extract department from employeeId or use mock
      departments.add("Engineering");
      departments.add("HR");
      departments.add("Sales");
      departments.add("Marketing");
      departments.add("Product");
      departments.add("Design");
    }
    return departments.size || 6;
  }, [topRecognizers]);

  const formattedDateRange = useMemo(() => {
    if (!startDate || !endDate) return "";
    const start = new Date(startDate);
    const end = new Date(endDate);
    return `${start.toISOString().split("T")[0]} to ${end.toISOString().split("T")[0]}`;
  }, [startDate, endDate]);

  return (
    <div className="space-y-6 pb-10">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">HR Analytics</h1>
        <p className="text-gray-500 text-sm">
          Comprehensive insights into recognition patterns across the company
        </p>
      </div>

      {/* Tab Navigation */}
      <div className="bg-gray-100 p-1 rounded-lg inline-flex w-full">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-md text-sm font-medium transition-all ${
              activeTab === tab.id
                ? "bg-white text-gray-900 shadow-sm"
                : "text-gray-600 hover:text-gray-900"
            }`}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </div>

      {/* Time Range Section */}
      <Card className="border border-gray-200 shadow-sm">
        <CardContent className="p-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-2 text-gray-700">
              <Calendar size={16} />
              <span className="font-medium">Time Range</span>
            </div>
            {activeTab === "dashboard" && (
              <Button
                variant="outline"
                size="sm"
                onClick={handleExport}
                disabled={exportLoading}
                className="gap-2"
              >
                <Download size={14} />
                Export CSV
              </Button>
            )}
          </div>
          <div className="flex items-center gap-4">
            <div className="flex-1">
              <label className="text-sm text-gray-500 mb-1 block">From</label>
              <Input
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="bg-white"
              />
            </div>
            <div className="flex-1">
              <label className="text-sm text-gray-500 mb-1 block">To</label>
              <Input
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="bg-white"
              />
            </div>
            <div className="pt-6">
              <Button variant="outline" onClick={handleLast30Days}>
                Last 30 Days
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tab Content */}
      {activeTab === "dashboard" && (
        <DashboardTab
          loading={loading}
          dashboard={dashboard}
          departmentStats={departmentStats}
          formattedDateRange={formattedDateRange}
        />
      )}

      {activeTab === "recognizers" && (
        <RecognizersTab loading={loading} topRecognizers={topRecognizers} />
      )}

      {activeTab === "teams" && (
        <TeamsTab loading={loading} />
      )}

      {activeTab === "values" && (
        <ValuesTab loading={loading} valueDistribution={valueDistribution} />
      )}
    </div>
  );
};

// Dashboard Tab Component
const DashboardTab: React.FC<{
  loading: boolean;
  dashboard: DashboardAnalytics | null;
  departmentStats: number;
  formattedDateRange: string;
}> = ({ loading, dashboard, departmentStats, formattedDateRange }) => {
  const stats = [
    {
      title: "Total Cards",
      value: dashboard?.totalCards || 0,
      subtitle: "Cards in period",
      icon: <FileText size={16} className="text-gray-400" />,
    },
    {
      title: "Active Recognizers",
      value: dashboard?.activeUsers || 0,
      subtitle: "Employees sending cards",
      icon: <Users size={16} className="text-gray-400" />,
    },
    {
      title: "Departments",
      value: departmentStats,
      subtitle: "Active departments",
      icon: <Building2 size={16} className="text-gray-400" />,
    },
    {
      title: "Period",
      value: formattedDateRange,
      subtitle: "Selected range",
      icon: <TrendingUp size={16} className="text-gray-400" />,
      isText: true,
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, index) => (
        <Card key={index} className="border border-gray-200 shadow-sm">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">
              {stat.title}
            </CardTitle>
            {stat.icon}
          </CardHeader>
          <CardContent>
            {loading ? (
              <Skeleton className="h-9 w-24" />
            ) : (
              <>
                <div className={`font-bold ${stat.isText ? "text-sm" : "text-3xl"} text-indigo-600`}>
                  {stat.isText ? stat.value : stat.value.toLocaleString()}
                </div>
                <p className="text-xs text-gray-400 mt-1">{stat.subtitle}</p>
              </>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
};

// Recognizers Tab Component
const RecognizersTab: React.FC<{
  loading: boolean;
  topRecognizers: TopRecognizer[];
}> = ({ loading, topRecognizers }) => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [selectedEmployee, setSelectedEmployee] = useState<{
    employeeId: string;
    name: string;
  } | null>(null);
  const [employeeCards, setEmployeeCards] = useState<{
    sent: RecognitionCard[];
    received: RecognitionCard[];
  }>({ sent: [], received: [] });
  const [cardsLoading, setCardsLoading] = useState(false);

  const handleCardClick = useCallback(async (employee: {
    employeeId: string;
    name: string;
  }) => {
    setSelectedEmployee(employee);
    setDrawerOpen(true);
    setCardsLoading(true);

    try {
      const [sentRes, receivedRes] = await Promise.all([
        api.getCards({ senderId: employee.employeeId, pageSize: 50 }),
        api.getCards({ recipientId: employee.employeeId, pageSize: 50 }),
      ]);
      setEmployeeCards({
        sent: sentRes.data.data.data,
        received: receivedRes.data.data.data,
      });
    } catch (error) {
      console.error("Failed to fetch employee cards:", error);
    } finally {
      setCardsLoading(false);
    }
  }, []);

  const getRecipientNames = (card: RecognitionCard) => {
    return card.recipients.map((r) => r.name).join(", ");
  };

  const renderCardItem = (card: RecognitionCard) => (
    <div
      key={card.id}
      className="p-4 border border-gray-100 rounded-lg hover:shadow-sm transition-shadow"
    >
      <div className="flex items-center gap-2 mb-1">
        <span className="font-semibold text-indigo-600">
          {card.senderName || card.senderId}
        </span>
        <span className="text-gray-400">→</span>
        <span className="font-semibold text-gray-900">{getRecipientNames(card)}</span>
      </div>
      <p className="text-xs text-gray-400 mb-2">
        {formatDistanceToNow(new Date(card.createdAt), { addSuffix: true })}
      </p>
      <div className="text-sm text-gray-700 whitespace-pre-wrap mb-3">
        {card.recognitionReason}
      </div>
      {card.selectedValues.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {card.selectedValues.map((value) => (
            <Badge
              key={value.id}
              className={`text-xs font-medium px-2 py-0.5 border ${getValueBadgeClasses(value.name)}`}
            >
              {value.name}
            </Badge>
          ))}
        </div>
      )}
    </div>
  );

  return (
    <>
      <Card className="border border-gray-200 shadow-sm">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-gray-800">
            <Users size={18} />
            Most Active Recognizers
          </CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <Skeleton key={i} className="h-16 w-full" />
              ))}
            </div>
          ) : topRecognizers.length > 0 ? (
            <div className="space-y-3">
              {topRecognizers.map((recognizer, index) => (
                <div
                  key={recognizer.employeeId}
                  className="flex items-center justify-between p-4 rounded-lg border border-gray-100 hover:border-gray-200 transition-colors"
                >
                  <div className="flex items-center gap-4">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
                      index < 3 ? "bg-indigo-100 text-indigo-600" : "bg-gray-100 text-gray-600"
                    }`}>
                      {index + 1}
                    </div>
                    <div>
                      <div className="font-medium text-gray-900">
                        {recognizer.employeeName || recognizer.employeeId}
                      </div>
                    </div>
                  </div>
                  <button
                    onClick={() => handleCardClick({
                      employeeId: recognizer.employeeId,
                      name: recognizer.employeeName || recognizer.employeeId,
                    })}
                    className="text-indigo-600 font-medium hover:text-indigo-800 hover:underline cursor-pointer transition-colors"
                  >
                    {recognizer.cardsSent} cards
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <div className="py-8 text-center text-gray-400">No data available</div>
          )}
        </CardContent>
      </Card>

      <Drawer open={drawerOpen} onOpenChange={setDrawerOpen}>
        <DrawerContent>
          <DrawerHeader>
            <DrawerTitle>{selectedEmployee?.name}</DrawerTitle>
          </DrawerHeader>
          <div className="p-6 space-y-6">
            {cardsLoading ? (
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="p-4 border border-gray-100 rounded-lg">
                    <Skeleton className="h-5 w-48 mb-2" />
                    <Skeleton className="h-4 w-32 mb-3" />
                    <Skeleton className="h-16 w-full" />
                  </div>
                ))}
              </div>
            ) : (
              <>
                {/* Sent Cards */}
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-3 flex items-center gap-2">
                    <span className="text-indigo-600">Sent</span>
                    <span className="text-sm font-normal text-gray-500">
                      ({employeeCards.sent.length} cards)
                    </span>
                  </h3>
                  {employeeCards.sent.length > 0 ? (
                    <div className="space-y-3">
                      {employeeCards.sent.map(renderCardItem)}
                    </div>
                  ) : (
                    <p className="text-gray-400 text-sm py-4 text-center">No sent cards</p>
                  )}
                </div>

                {/* Received Cards */}
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-3 flex items-center gap-2">
                    <span className="text-emerald-600">Received</span>
                    <span className="text-sm font-normal text-gray-500">
                      ({employeeCards.received.length} cards)
                    </span>
                  </h3>
                  {employeeCards.received.length > 0 ? (
                    <div className="space-y-3">
                      {employeeCards.received.map(renderCardItem)}
                    </div>
                  ) : (
                    <p className="text-gray-400 text-sm py-4 text-center">No received cards</p>
                  )}
                </div>
              </>
            )}
          </div>
        </DrawerContent>
      </Drawer>
    </>
  );
};

// Teams Tab Component
const TeamsTab: React.FC<{ loading: boolean }> = ({ loading }) => {
  // Mock team data
  const teamStats = [
    { name: "Engineering", cardsSent: 45, cardsReceived: 38, members: 25 },
    { name: "HR", cardsSent: 32, cardsReceived: 28, members: 8 },
    { name: "Sales", cardsSent: 28, cardsReceived: 35, members: 15 },
    { name: "Product", cardsSent: 22, cardsReceived: 24, members: 12 },
    { name: "Marketing", cardsSent: 18, cardsReceived: 20, members: 10 },
    { name: "Design", cardsSent: 15, cardsReceived: 18, members: 6 },
  ];

  return (
    <Card className="border border-gray-200 shadow-sm">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-gray-800">
          <Building2 size={18} />
          Team Recognition Activity
        </CardTitle>
      </CardHeader>
      <CardContent>
        {loading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <Skeleton key={i} className="h-16 w-full" />
            ))}
          </div>
        ) : (
          <div className="space-y-3">
            {teamStats.map((team, index) => (
              <div
                key={team.name}
                className="flex items-center justify-between p-4 rounded-lg border border-gray-100 hover:border-gray-200 transition-colors"
              >
                <div className="flex items-center gap-4">
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
                    index < 3 ? "bg-indigo-100 text-indigo-600" : "bg-gray-100 text-gray-600"
                  }`}>
                    {index + 1}
                  </div>
                  <div>
                    <div className="font-medium text-gray-900">{team.name}</div>
                    <div className="text-sm text-gray-500">{team.members} members</div>
                  </div>
                </div>
                <div className="flex gap-6 text-sm">
                  <div className="text-center">
                    <div className="font-medium text-indigo-600">{team.cardsSent}</div>
                    <div className="text-gray-400">sent</div>
                  </div>
                  <div className="text-center">
                    <div className="font-medium text-emerald-600">{team.cardsReceived}</div>
                    <div className="text-gray-400">received</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};

// Values Tab Component
const ValuesTab: React.FC<{
  loading: boolean;
  valueDistribution: ValueDistribution | null;
}> = ({ loading, valueDistribution }) => {
  const wordCloudData = useMemo(() => {
    if (!valueDistribution) return [];
    return valueDistribution.distribution.map((item) => {
      const colorConfig = VALUE_COLOR_MAP[item.valueName] || DEFAULT_VALUE_COLOR;
      return {
        ...item,
        color: colorConfig.hex,
        fontSize: Math.max(14, Math.min(48, item.percentage * 2)),
      };
    });
  }, [valueDistribution]);

  const getChartColor = (valueName: string) => {
    const colorConfig = VALUE_COLOR_MAP[valueName] || DEFAULT_VALUE_COLOR;
    return colorConfig.hex;
  };

  return (
    <Card className="border border-gray-200 shadow-sm">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-gray-800">
          <Tag size={18} />
          Values Distribution
        </CardTitle>
      </CardHeader>
      <CardContent>
        {loading ? (
          <div className="grid md:grid-cols-2 gap-8">
            <Skeleton className="h-64 w-full" />
            <Skeleton className="h-64 w-full" />
          </div>
        ) : (
          <div className="grid md:grid-cols-2 gap-8">
            {/* Word Cloud */}
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-4">Word Cloud</h3>
              <div className="flex flex-wrap items-center justify-center gap-3 p-4 min-h-[250px]">
                {wordCloudData.map((item) => (
                  <span
                    key={item.valueId}
                    style={{
                      color: item.color,
                      fontSize: `${item.fontSize}px`,
                    }}
                    className="font-medium cursor-default hover:opacity-80 transition-opacity"
                  >
                    {item.valueName}
                  </span>
                ))}
              </div>
            </div>

            {/* Pie Chart */}
            <div>
              <h3 className="text-sm font-medium text-gray-500 mb-4">Distribution by Percentage</h3>
              <div className="h-[250px]">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={valueDistribution?.distribution}
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      dataKey="count"
                    >
                      {valueDistribution?.distribution.map((entry) => (
                        <Cell key={`cell-${entry.valueId}`} fill={getChartColor(entry.valueName)} />
                      ))}
                    </Pie>
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "#fff",
                        border: "1px solid #e5e7eb",
                        borderRadius: "8px",
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

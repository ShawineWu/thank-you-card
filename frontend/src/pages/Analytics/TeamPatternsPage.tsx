import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { analyticsApi } from '@/services/analytics'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2, Building2, Calendar } from 'lucide-react'
import { format, subDays } from 'date-fns'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'

const TeamPatternsPage = () => {
  const [from, setFrom] = useState<string>(
    format(subDays(new Date(), 30), 'yyyy-MM-dd')
  )
  const [to, setTo] = useState<string>(format(new Date(), 'yyyy-MM-dd'))

  const { data, isLoading, error } = useQuery({
    queryKey: ['analytics', 'teams', from, to],
    queryFn: () =>
      analyticsApi.getTeamPatterns({
        from: from ? `${from}T00:00:00Z` : undefined,
        to: to ? `${to}T23:59:59Z` : undefined,
      }),
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-destructive">Failed to load data. Please try again.</p>
      </div>
    )
  }

  const chartData = data?.map(team => ({
    department: team.department,
    sent: team.totalSent,
    received: team.totalReceived,
    average: Number(team.averageCardsPerEmployee.toFixed(1)),
  })) || []

  // Prepare pie chart data (showing distribution of cards sent by department)
  const pieChartData = data?.map(team => ({
    name: team.department,
    value: team.totalSent,
  })) || []

  // Color palette for pie chart
  const COLORS = [
    '#3b82f6', // blue
    '#10b981', // green
    '#f59e0b', // amber
    '#ef4444', // red
    '#8b5cf6', // purple
    '#ec4899', // pink
    '#06b6d4', // cyan
    '#f97316', // orange
    '#6366f1', // indigo
    '#14b8a6', // teal
  ]

  return (
    <div>
      {/* Time Range Selector */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2 text-base">
            <Calendar className="h-4 w-4" />
            <span>Time Range</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <label className="text-sm font-medium mb-2 block">From</label>
              <Input
                type="date"
                value={from}
                onChange={(e) => setFrom(e.target.value)}
              />
            </div>
            <div className="flex-1">
              <label className="text-sm font-medium mb-2 block">To</label>
              <Input
                type="date"
                value={to}
                onChange={(e) => setTo(e.target.value)}
              />
            </div>
            <div className="flex items-end">
              <Button
                variant="outline"
                onClick={() => {
                  setFrom(format(subDays(new Date(), 30), 'yyyy-MM-dd'))
                  setTo(format(new Date(), 'yyyy-MM-dd'))
                }}
              >
                Last 30 Days
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Charts */}
      {chartData.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Building2 className="h-5 w-5" />
              <span>Team Recognition Patterns</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Left: Bar Chart */}
              <div>
                <h3 className="text-sm font-semibold text-muted-foreground mb-4">Cards Sent & Received by Department</h3>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart 
                    data={chartData}
                    margin={{ top: 10, right: 20, left: 0, bottom: 5 }}
                    barCategoryGap="20%"
                  >
                    <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" opacity={0.3} />
                    <XAxis 
                      dataKey="department" 
                      tick={{ fontSize: 12, fill: '#6b7280' }}
                      axisLine={{ stroke: '#e5e7eb' }}
                      tickLine={{ stroke: '#e5e7eb' }}
                      angle={-45}
                      textAnchor="end"
                      height={80}
                    />
                    <YAxis 
                      tick={{ fontSize: 12, fill: '#6b7280' }}
                      axisLine={{ stroke: '#e5e7eb' }}
                      tickLine={{ stroke: '#e5e7eb' }}
                      width={60}
                    />
                    <Tooltip 
                      contentStyle={{
                        backgroundColor: 'rgba(255, 255, 255, 0.95)',
                        border: '1px solid #e5e7eb',
                        borderRadius: '8px',
                        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                        padding: '8px 12px',
                      }}
                      labelStyle={{ color: '#374151', fontWeight: 600, marginBottom: '4px' }}
                      itemStyle={{ color: '#6b7280', padding: '2px 0' }}
                    />
                    <Legend 
                      wrapperStyle={{ paddingTop: '20px' }}
                      iconType="rect"
                      iconSize={12}
                      formatter={(value) => <span style={{ color: '#374151', fontSize: '13px' }}>{value}</span>}
                    />
                    <Bar 
                      dataKey="sent" 
                      fill="#3b82f6" 
                      name="Cards Sent"
                      radius={[4, 4, 0, 0]}
                    />
                    <Bar 
                      dataKey="received" 
                      fill="#10b981" 
                      name="Cards Received"
                      radius={[4, 4, 0, 0]}
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>

              {/* Right: Pie Chart */}
              <div>
                <h3 className="text-sm font-semibold text-muted-foreground mb-4">Cards Sent Distribution</h3>
                <ResponsiveContainer width="100%" height={250}>
                  <PieChart>
                    <Pie
                      data={pieChartData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      outerRadius={90}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {pieChartData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{
                        backgroundColor: 'rgba(255, 255, 255, 0.95)',
                        border: '1px solid #e5e7eb',
                        borderRadius: '8px',
                        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                        padding: '8px 12px',
                      }}
                      formatter={(value: number, name: string, props: any) => {
                        const total = pieChartData.reduce((sum, item) => sum + item.value, 0)
                        const percent = ((value / total) * 100).toFixed(1)
                        return [`${value} cards (${percent}%)`, props.payload.name]
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>
                {/* Custom Legend */}
                <div className="mt-4 space-y-2">
                  {pieChartData.map((entry, index) => {
                    const total = pieChartData.reduce((sum, item) => sum + item.value, 0)
                    const percent = ((entry.value / total) * 100).toFixed(1)
                    return (
                      <div key={index} className="flex items-center justify-between text-sm">
                        <div className="flex items-center space-x-2">
                          <div
                            className="w-4 h-4 rounded"
                            style={{ backgroundColor: COLORS[index % COLORS.length] }}
                          />
                          <span className="text-foreground">{entry.name}</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <span className="text-muted-foreground">{entry.value} cards</span>
                          <span className="text-muted-foreground font-medium">({percent}%)</span>
                        </div>
                      </div>
                    )
                  })}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Table */}
      <Card>
        <CardHeader>
          <CardTitle>Team Statistics</CardTitle>
        </CardHeader>
        <CardContent>
          {data && data.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b">
                    <th className="text-left p-2">Department</th>
                    <th className="text-right p-2">Total Sent</th>
                    <th className="text-right p-2">Total Received</th>
                    <th className="text-right p-2">Avg per Employee</th>
                  </tr>
                </thead>
                <tbody>
                  {data.map((team) => (
                    <tr key={team.department} className="border-b hover:bg-accent/50">
                      <td className="p-2 font-medium">{team.department}</td>
                      <td className="p-2 text-right">{team.totalSent}</td>
                      <td className="p-2 text-right">{team.totalReceived}</td>
                      <td className="p-2 text-right">
                        {team.averageCardsPerEmployee.toFixed(1)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-8">
              No data available for the selected time range
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default TeamPatternsPage

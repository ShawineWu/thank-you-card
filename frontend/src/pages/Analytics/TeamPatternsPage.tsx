import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { analyticsApi } from '@/services/analytics'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2, Building2, Calendar } from 'lucide-react'
import { format, subDays } from 'date-fns'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'

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

      {/* Chart */}
      {chartData.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Building2 className="h-5 w-5" />
              <span>Team Recognition Patterns</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={400}>
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="department" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="sent" fill="#8884d8" name="Cards Sent" />
                <Bar dataKey="received" fill="#82ca9d" name="Cards Received" />
              </BarChart>
            </ResponsiveContainer>
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

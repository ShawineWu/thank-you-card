import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { statsApi } from '@/services/stats'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Loader2, Trophy, Calendar } from 'lucide-react'
import { format, subDays } from 'date-fns'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const TopEmployeesPage = () => {
  const [from, setFrom] = useState<string>(
    format(subDays(new Date(), 30), 'yyyy-MM-dd')
  )
  const [to, setTo] = useState<string>(format(new Date(), 'yyyy-MM-dd'))

  const { data, isLoading, error } = useQuery({
    queryKey: ['top-recipients', from, to],
    queryFn: () =>
      statsApi.getTopRecipients(
        {
          from: from ? `${from}T00:00:00Z` : undefined,
          to: to ? `${to}T23:59:59Z` : undefined,
        },
        10
      ),
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
        <p className="text-destructive">Failed to load top employees. Please try again.</p>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Top 10 Recognized Employees</h1>
        <p className="text-muted-foreground">
          Employees who have received the most thank you cards
        </p>
      </div>

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

      {/* Ranking List */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Trophy className="h-5 w-5 text-yellow-500" />
            <span>Ranking</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {data && data.length > 0 ? (
            <div className="space-y-4">
              {data.map((employee, index) => {
                const rank = index + 1
                const isTopThree = rank <= 3
                return (
                  <div
                    key={employee.employeeId}
                    className="flex items-center space-x-4 p-4 rounded-lg border hover:bg-accent/50 transition-colors"
                  >
                    <div
                      className={`flex items-center justify-center w-10 h-10 rounded-full font-bold ${
                        isTopThree
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      {rank}
                    </div>
                    <div className="flex-1">
                      <div className="font-semibold text-lg">{employee.name}</div>
                      <div className="text-sm text-muted-foreground">
                        {employee.department}
                      </div>
                    </div>
                    <Badge variant="secondary" className="text-base px-3 py-1">
                      {employee.count} {employee.count === 1 ? 'card' : 'cards'}
                    </Badge>
                  </div>
                )
              })}
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

export default TopEmployeesPage

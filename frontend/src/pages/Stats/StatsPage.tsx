import { useQuery } from '@tanstack/react-query'
import { statsApi } from '@/services/stats'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Loader2, Send, Inbox, TrendingUp } from 'lucide-react'
import ValueBadge from '@/components/Card/ValueBadge'
import { CompanyValueSummary } from '@/types/card'

const StatsPage = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['stats', 'personal'],
    queryFn: () => statsApi.getPersonalStats(),
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
        <p className="text-destructive">Failed to load statistics. Please try again.</p>
      </div>
    )
  }

  if (!data) {
    return null
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">My Statistics</h1>
        <p className="text-muted-foreground">Your recognition activity overview</p>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cards Sent</CardTitle>
            <Send className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.totalSent}</div>
            <p className="text-xs text-muted-foreground">
              Total thank you cards you've sent
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cards Received</CardTitle>
            <Inbox className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.totalReceived}</div>
            <p className="text-xs text-muted-foreground">
              Total thank you cards you've received
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Top Values */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Top Sent Values */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <TrendingUp className="h-5 w-5" />
              <span>Most Used Values (When Sending)</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            {data.topSentValues.length > 0 ? (
              <div className="space-y-3">
                {data.topSentValues.map((value) => {
                  const valueSummary: CompanyValueSummary = {
                    id: value.valueId,
                    code: value.code,
                    name: value.name,
                    type: value.type as 'VALUE' | 'CREDO',
                  }
                  return (
                    <div key={value.valueId} className="flex items-center justify-between">
                      <ValueBadge value={valueSummary} />
                      <span className="text-sm font-medium text-muted-foreground">
                        {value.count} {value.count === 1 ? 'time' : 'times'}
                      </span>
                    </div>
                  )
                })}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No data available</p>
            )}
          </CardContent>
        </Card>

        {/* Top Received Values */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <TrendingUp className="h-5 w-5" />
              <span>Most Received Values</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            {data.topReceivedValues.length > 0 ? (
              <div className="space-y-3">
                {data.topReceivedValues.map((value) => {
                  const valueSummary: CompanyValueSummary = {
                    id: value.valueId,
                    code: value.code,
                    name: value.name,
                    type: value.type as 'VALUE' | 'CREDO',
                  }
                  return (
                    <div key={value.valueId} className="flex items-center justify-between">
                      <ValueBadge value={valueSummary} />
                      <span className="text-sm font-medium text-muted-foreground">
                        {value.count} {value.count === 1 ? 'time' : 'times'}
                      </span>
                    </div>
                  )
                })}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No data available</p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default StatsPage

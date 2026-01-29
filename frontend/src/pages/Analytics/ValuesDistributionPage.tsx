import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { analyticsApi } from '@/services/analytics'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Loader2, Calendar, Tag } from 'lucide-react'
import { format, subDays } from 'date-fns'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import ValueCardsDialog from '@/components/Analytics/ValueCardsDialog'
import ValueDetailDialog from '@/components/CompanyValues/ValueDetailDialog'
import { CompanyValueSummary, CompanyValueDetail } from '@/types/card'
import { companyValuesApi } from '@/services/companyValues'
import WordCloud from '@/components/Analytics/WordCloud'
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts'

const ValuesDistributionPage = () => {
  const [from, setFrom] = useState<string>(
    format(subDays(new Date(), 30), 'yyyy-MM-dd')
  )
  const [to, setTo] = useState<string>(format(new Date(), 'yyyy-MM-dd'))
  const [selectedValue, setSelectedValue] = useState<{
    id: number
    name: string
  } | null>(null)
  const [cardsDialogOpen, setCardsDialogOpen] = useState(false)
  const [selectedValueDetail, setSelectedValueDetail] = useState<CompanyValueDetail | null>(null)
  const [valueDetailDialogOpen, setValueDetailDialogOpen] = useState(false)

  const { data, isLoading, error } = useQuery({
    queryKey: ['analytics', 'values', from, to],
    queryFn: () =>
      analyticsApi.getValuesDistribution({
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

  const wordCloudData = data?.map(item => ({
    text: item.name,
    value: item.count,
    percentage: item.percentage,
  })) || []

  // Prepare pie chart data
  const pieChartData = data?.map(item => ({
    name: item.name,
    value: item.count,
    percentage: item.percentage,
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

  const handleCardsClick = (valueId: number, valueName: string) => {
    setSelectedValue({ id: valueId, name: valueName })
    setCardsDialogOpen(true)
  }

  const handleValueClick = async (value: CompanyValueSummary) => {
    try {
      const valueDetail = await companyValuesApi.getById(value.id)
      setSelectedValueDetail(valueDetail)
      setValueDetailDialogOpen(true)
    } catch (error) {
      console.error('Failed to load value details:', error)
    }
  }

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

      {/* Word Cloud and Pie Chart */}
      {wordCloudData.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Tag className="h-5 w-5" />
              <span>Values Distribution</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Left: Word Cloud */}
              <div>
                <h3 className="text-sm font-semibold text-muted-foreground mb-4">Word Cloud</h3>
                <WordCloud
                  data={wordCloudData}
                  onClick={(item) => {
                    const valueItem = data?.find(d => d.name === item.text)
                    if (valueItem) {
                      handleCardsClick(valueItem.valueId, valueItem.name)
                    }
                  }}
                />
              </div>

              {/* Right: Pie Chart */}
              <div>
                <h3 className="text-sm font-semibold text-muted-foreground mb-4">Distribution by Percentage</h3>
                <ResponsiveContainer width="100%" height={400}>
                  <PieChart>
                    <Pie
                      data={pieChartData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      outerRadius={120}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {pieChartData.map((entry, index) => (
                        <Cell 
                          key={`cell-${index}`} 
                          fill={COLORS[index % COLORS.length]}
                          onClick={() => {
                            const valueItem = data?.find(d => d.name === entry.name)
                            if (valueItem) {
                              handleCardsClick(valueItem.valueId, valueItem.name)
                            }
                          }}
                          style={{ cursor: 'pointer' }}
                        />
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
                <div className="mt-4 space-y-2 max-h-[200px] overflow-y-auto">
                  {pieChartData.map((entry, index) => {
                    const total = pieChartData.reduce((sum, item) => sum + item.value, 0)
                    const percent = ((entry.value / total) * 100).toFixed(1)
                    return (
                      <div 
                        key={index} 
                        className="flex items-center justify-between text-sm cursor-pointer hover:bg-accent/50 p-2 rounded transition-colors"
                        onClick={() => {
                          const valueItem = data?.find(d => d.name === entry.name)
                          if (valueItem) {
                            handleCardsClick(valueItem.valueId, valueItem.name)
                          }
                        }}
                      >
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

      {/* List */}
      <Card>
        <CardHeader>
          <CardTitle>Distribution Details</CardTitle>
        </CardHeader>
        <CardContent>
          {data && data.length > 0 ? (
            <div className="space-y-3">
              {data.map((item) => (
                <div
                  key={item.valueId}
                  className="flex items-center justify-between p-3 rounded-lg border hover:bg-accent/50 transition-colors"
                >
                  <div className="flex items-center space-x-3">
                    <Badge variant={item.type === 'VALUE' ? 'default' : 'success'}>
                      {item.type}
                    </Badge>
                    <div>
                      <div className="font-medium">{item.name}</div>
                    </div>
                  </div>
                  <div className="text-right">
                    <Badge 
                      variant="secondary" 
                      className="text-base px-3 py-1 cursor-pointer hover:bg-secondary/80 transition-colors"
                      onClick={() => handleCardsClick(item.valueId, item.name)}
                    >
                      {item.count} {item.count === 1 ? 'card' : 'cards'}
                    </Badge>
                    <div className="text-sm text-muted-foreground mt-1">
                      {item.percentage.toFixed(1)}%
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-8">
              No data available for the selected time range
            </p>
          )}
        </CardContent>
      </Card>

      {/* Value Cards Dialog */}
      {selectedValue && (
        <ValueCardsDialog
          valueId={selectedValue.id}
          valueName={selectedValue.name}
          open={cardsDialogOpen}
          onOpenChange={setCardsDialogOpen}
          from={from}
          to={to}
          onValueClick={handleValueClick}
        />
      )}

      {/* Value Detail Dialog */}
      <ValueDetailDialog
        value={selectedValueDetail}
        open={valueDetailDialogOpen}
        onOpenChange={setValueDetailDialogOpen}
      />
    </div>
  )
}

export default ValuesDistributionPage

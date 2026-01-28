import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { analyticsApi } from '@/services/analytics'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Loader2, PieChart, Calendar } from 'lucide-react'
import { format, subDays } from 'date-fns'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { PieChart as RechartsPieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts'
import ValueCardsDialog from '@/components/Analytics/ValueCardsDialog'
import ValueDetailDialog from '@/components/CompanyValues/ValueDetailDialog'
import { CompanyValueSummary, CompanyValueDetail } from '@/types/card'
import { companyValuesApi } from '@/services/companyValues'

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d', '#ffc658', '#ff7300', '#8dd1e1', '#d084d0']

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

  const chartData = data?.map(item => ({
    name: item.name,
    value: item.count,
    percentage: item.percentage,
  })) || []

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

      {/* Chart */}
      {chartData.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <PieChart className="h-5 w-5" />
              <span>Values Distribution</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={400}>
              <RechartsPieChart>
                <Pie
                  data={chartData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percentage }) => `${name}: ${percentage.toFixed(1)}%`}
                  outerRadius={120}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {chartData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </RechartsPieChart>
            </ResponsiveContainer>
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

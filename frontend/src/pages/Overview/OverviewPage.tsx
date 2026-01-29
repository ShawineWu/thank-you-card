import { useState, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { statsApi } from '@/services/stats'
import { milestonesApi } from '@/services/milestones'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Loader2, Send, Inbox, TrendingUp, Trophy, Calendar, Award, User } from 'lucide-react'
import { format, subDays } from 'date-fns'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import ValueBadge from '@/components/Card/ValueBadge'
import { CompanyValueSummary, CompanyValueDetail } from '@/types/card'
import { companyValuesApi } from '@/services/companyValues'
import Podium from '@/components/Overview/Podium'
import EmployeeCardsDialog from '@/components/TopEmployees/EmployeeCardsDialog'
import ValueDetailDialog from '@/components/CompanyValues/ValueDetailDialog'
import MilestoneBadge from '@/components/Overview/MilestoneBadge'
import MilestoneDetailDialog from '@/components/Overview/MilestoneDetailDialog'
import { Milestone, UserAchievement } from '@/types/milestone'

const OverviewPage = () => {
  const [from, setFrom] = useState<string>(
    format(subDays(new Date(), 30), 'yyyy-MM-dd')
  )
  const [to, setTo] = useState<string>(format(new Date(), 'yyyy-MM-dd'))
  const [selectedEmployee, setSelectedEmployee] = useState<{
    id: number
    name: string
  } | null>(null)
  const [cardsDialogOpen, setCardsDialogOpen] = useState(false)
  const [selectedValue, setSelectedValue] = useState<CompanyValueDetail | null>(null)
  const [valueDialogOpen, setValueDialogOpen] = useState(false)
  const [selectedMilestone, setSelectedMilestone] = useState<Milestone | null>(null)
  const [milestoneDialogOpen, setMilestoneDialogOpen] = useState(false)

  // Fetch Top 10 employees
  const { data: topEmployeesData, isLoading: isLoadingTop, error: topError } = useQuery({
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

  // Fetch personal stats
  const { data: personalStats, isLoading: isLoadingStats, error: statsError } = useQuery({
    queryKey: ['stats', 'personal'],
    queryFn: () => statsApi.getPersonalStats(),
  })

  // Fetch all milestones
  const { data: allMilestones, isLoading: isLoadingMilestones } = useQuery({
    queryKey: ['milestones', 'all'],
    queryFn: () => milestonesApi.getAllMilestones(),
  })

  // Fetch user achievements
  const { data: userAchievements, isLoading: isLoadingAchievements } = useQuery({
    queryKey: ['milestones', 'me'],
    queryFn: () => milestonesApi.getMyAchievements(),
  })

  // Combine milestones with achievements and progress
  const milestonesWithStatus = useMemo(() => {
    if (!allMilestones || !personalStats) return []

    const achievementMap = new Map<number, UserAchievement>()
    if (userAchievements) {
      userAchievements.forEach((achievement) => {
        achievementMap.set(achievement.milestone.id, achievement)
      })
    }

    return allMilestones.map((milestone) => {
      const achievement = achievementMap.get(milestone.id)
      let currentProgress: number | undefined

      if (achievement) {
        currentProgress = milestone.threshold // Already achieved
      } else {
        // Calculate current progress
        switch (milestone.type) {
          case 'SENT':
            currentProgress = personalStats.totalSent
            break
          case 'RECEIVED':
            currentProgress = personalStats.totalReceived
            break
          case 'TOTAL':
            currentProgress = personalStats.totalSent + personalStats.totalReceived
            break
        }
      }

      return {
        milestone,
        achieved: !!achievement,
        achievedAt: achievement?.achievedAt,
        currentProgress: Math.max(0, currentProgress || 0),
      }
    })
  }, [allMilestones, userAchievements, personalStats])

  const handleCardsClick = (employeeId: number, employeeName: string) => {
    setSelectedEmployee({ id: employeeId, name: employeeName })
    setCardsDialogOpen(true)
  }

  const handleValueClick = async (value: CompanyValueSummary) => {
    try {
      const valueDetail = await companyValuesApi.getById(value.id)
      setSelectedValue(valueDetail)
      setValueDialogOpen(true)
    } catch (error) {
      console.error('Failed to load value details:', error)
    }
  }

  const handleMilestoneClick = (milestone: Milestone) => {
    setSelectedMilestone(milestone)
    setMilestoneDialogOpen(true)
  }

  const isLoading = isLoadingTop || isLoadingStats || isLoadingMilestones || isLoadingAchievements
  const error = topError || statsError

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

  const topThree = topEmployeesData?.slice(0, 3) || []
  const remainingEmployees = topEmployeesData?.slice(3) || []

  return (
    <div>
      {/* Page Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Overview</h1>
        <p className="text-muted-foreground">Your recognition activity and company highlights</p>
      </div>

      {/* Time Range Selector for Top 10 */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2 text-base">
            <Calendar className="h-4 w-4" />
            <span>Time Range for Top 10</span>
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

      {/* Top 10 Section */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Trophy className="h-5 w-5 text-yellow-500" />
            <span>Top 10 Recognized Employees</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {topEmployeesData && topEmployeesData.length > 0 ? (
            <div className="space-y-6">
              {/* Top 3 Podium */}
              {topThree.length === 3 && (
                <div className="mb-8">
                  <Podium topThree={topThree} onEmployeeClick={handleCardsClick} />
                </div>
              )}

              {/* Remaining Employees (4-10) */}
              {remainingEmployees.length > 0 && (
                <div className="space-y-3">
                  {remainingEmployees.map((employee, index) => {
                    const rank = index + 4
                    return (
                      <div
                        key={employee.employeeId}
                        className="flex items-center space-x-4 p-3 rounded-lg border hover:bg-accent/50 transition-colors"
                      >
                        <div className="flex items-center justify-center w-10 h-10 rounded-full font-bold bg-muted text-muted-foreground">
                          {rank}
                        </div>
                        <div className="flex-1">
                          <div className="font-semibold">{employee.name}</div>
                          <div className="text-sm text-muted-foreground">
                            {employee.department}
                          </div>
                        </div>
                        <Badge
                          variant="secondary"
                          className="text-base px-3 py-1 cursor-pointer hover:bg-secondary/80 transition-colors"
                          onClick={() => handleCardsClick(employee.employeeId, employee.name)}
                        >
                          {employee.count} {employee.count === 1 ? 'card' : 'cards'}
                        </Badge>
                      </div>
                    )
                  })}
                </div>
              )}
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-8">
              No data available for the selected time range
            </p>
          )}
        </CardContent>
      </Card>

      {/* Divider Section - Separating Company Data from Personal Data */}
      <div className="relative my-10">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t-2 border-dashed border-border/60"></div>
        </div>
        <div className="relative flex justify-center">
          <div className="bg-background px-6 py-2">
            <div className="flex items-center space-x-2 text-base font-semibold text-foreground">
              <div className="p-1.5 rounded-full bg-primary/10">
                <User className="h-4 w-4 text-primary" />
              </div>
              <span>Personal Dashboard</span>
            </div>
          </div>
        </div>
      </div>

      {/* Most Received Values Section */}
      {personalStats && personalStats.topReceivedValues.length > 0 && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <TrendingUp className="h-5 w-5" />
              <span>Most Received Values</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {personalStats.topReceivedValues.map((value) => {
                const valueSummary: CompanyValueSummary = {
                  id: value.valueId,
                  code: value.code,
                  name: value.name,
                  type: value.type as 'VALUE' | 'CREDO',
                }
                return (
                  <div key={value.valueId} className="flex items-center justify-between">
                    <ValueBadge value={valueSummary} onClick={handleValueClick} />
                    <span className="text-sm font-medium text-muted-foreground">
                      {value.count} {value.count === 1 ? 'time' : 'times'}
                    </span>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Personal Statistics Cards */}
      {personalStats && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Cards Sent</CardTitle>
              <Send className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{personalStats.totalSent}</div>
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
              <div className="text-2xl font-bold">{personalStats.totalReceived}</div>
              <p className="text-xs text-muted-foreground">
                Total thank you cards you've received
              </p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Milestones Section */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Award className="h-5 w-5 text-yellow-500" />
            <span>My Achievements</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {milestonesWithStatus.length > 0 ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-8">
              {milestonesWithStatus.map(({ milestone, achieved, achievedAt, currentProgress }) => (
                <MilestoneBadge
                  key={milestone.id}
                  milestone={milestone}
                  achieved={achieved}
                  achievedAt={achievedAt}
                  currentProgress={currentProgress}
                  onClick={() => handleMilestoneClick(milestone)}
                />
              ))}
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-8">
              {isLoadingMilestones || isLoadingAchievements
                ? 'Loading achievements...'
                : 'No milestones available'}
            </p>
          )}
        </CardContent>
      </Card>

      {/* Employee Cards Dialog */}
      {selectedEmployee && (
        <EmployeeCardsDialog
          employeeId={selectedEmployee.id}
          employeeName={selectedEmployee.name}
          open={cardsDialogOpen}
          onOpenChange={setCardsDialogOpen}
          from={from}
          to={to}
          onValueClick={handleValueClick}
          showFilters={false}
        />
      )}

      {/* Value Detail Dialog */}
      <ValueDetailDialog
        value={selectedValue}
        open={valueDialogOpen}
        onOpenChange={setValueDialogOpen}
      />

      {/* Milestone Detail Dialog */}
      {selectedMilestone && (
        <MilestoneDetailDialog
          milestone={selectedMilestone}
          achieved={milestonesWithStatus.find((m) => m.milestone.id === selectedMilestone.id)?.achieved || false}
          achievedAt={milestonesWithStatus.find((m) => m.milestone.id === selectedMilestone.id)?.achievedAt}
          currentProgress={milestonesWithStatus.find((m) => m.milestone.id === selectedMilestone.id)?.currentProgress}
          open={milestoneDialogOpen}
          onOpenChange={setMilestoneDialogOpen}
        />
      )}
    </div>
  )
}

export default OverviewPage

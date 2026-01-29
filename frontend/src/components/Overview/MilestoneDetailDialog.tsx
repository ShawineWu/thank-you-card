import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Milestone } from '@/types/milestone'
import { format } from 'date-fns'
import { Trophy, Medal, Award, Send, Inbox, Target } from 'lucide-react'

interface MilestoneDetailDialogProps {
  milestone: Milestone
  achieved?: boolean
  achievedAt?: string
  currentProgress?: number
  open: boolean
  onOpenChange: (open: boolean) => void
}

const MilestoneDetailDialog = ({
  milestone,
  achieved = false,
  achievedAt,
  currentProgress,
  open,
  onOpenChange,
}: MilestoneDetailDialogProps) => {
  const getTypeIcon = () => {
    switch (milestone.type) {
      case 'SENT':
        return <Send className="h-5 w-5" />
      case 'RECEIVED':
        return <Inbox className="h-5 w-5" />
      case 'TOTAL':
        return <Target className="h-5 w-5" />
      default:
        return <Award className="h-5 w-5" />
    }
  }

  const getTypeLabel = () => {
    switch (milestone.type) {
      case 'SENT':
        return 'Cards Sent'
      case 'RECEIVED':
        return 'Cards Received'
      case 'TOTAL':
        return 'Total Cards'
      default:
        return milestone.type
    }
  }

  const progress = currentProgress !== undefined ? Math.min((currentProgress / milestone.threshold) * 100, 100) : undefined

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            {achieved ? (
              <Trophy className="h-6 w-6 text-yellow-500" />
            ) : (
              <Award className="h-6 w-6 text-gray-400" />
            )}
            <span>{milestone.name}</span>
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 mt-4">
          {/* Description */}
          <div>
            <p className="text-sm text-muted-foreground">{milestone.description}</p>
          </div>

          {/* Achievement Type */}
          <div className="flex items-center space-x-2 p-3 bg-muted rounded-lg">
            {getTypeIcon()}
            <div>
              <div className="text-sm font-medium">{getTypeLabel()}</div>
              <div className="text-xs text-muted-foreground">Achievement Type</div>
            </div>
          </div>

          {/* Threshold */}
          <div className="p-3 bg-muted rounded-lg">
            <div className="text-sm font-medium">Threshold</div>
            <div className="text-2xl font-bold text-primary mt-1">{milestone.threshold} cards</div>
          </div>

          {/* Progress or Achievement Date */}
          {achieved && achievedAt ? (
            <div className="p-3 bg-green-50 dark:bg-green-950 rounded-lg border border-green-200 dark:border-green-800">
              <div className="text-sm font-medium text-green-900 dark:text-green-100">Achievement Unlocked!</div>
              <div className="text-xs text-green-700 dark:text-green-300 mt-1">
                Achieved on {format(new Date(achievedAt), 'MMMM d, yyyy')}
              </div>
            </div>
          ) : (
            currentProgress !== undefined && (
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Progress</span>
                  <span className="font-medium">
                    {currentProgress} / {milestone.threshold}
                  </span>
                </div>
                <div className="w-full bg-muted rounded-full h-2">
                  <div
                    className="bg-primary h-2 rounded-full transition-all duration-300"
                    style={{ width: `${progress}%` }}
                  />
                </div>
                <div className="text-xs text-muted-foreground text-center">
                  {milestone.threshold - currentProgress} more cards needed
                </div>
              </div>
            )
          )}

          {/* How to Achieve */}
          {!achieved && (
            <div className="p-3 bg-blue-50 dark:bg-blue-950 rounded-lg border border-blue-200 dark:border-blue-800">
              <div className="text-sm font-medium text-blue-900 dark:text-blue-100">How to Achieve</div>
              <div className="text-xs text-blue-700 dark:text-blue-300 mt-1">
                {milestone.type === 'SENT' && 'Send thank you cards to your colleagues to unlock this achievement.'}
                {milestone.type === 'RECEIVED' &&
                  'Receive thank you cards from your colleagues to unlock this achievement.'}
                {milestone.type === 'TOTAL' &&
                  'Send or receive thank you cards to unlock this achievement. Both sent and received cards count towards this milestone.'}
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default MilestoneDetailDialog

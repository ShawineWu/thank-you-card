import { Milestone } from '@/types/milestone'
import { Trophy, Medal, Award } from 'lucide-react'
import { cn } from '@/lib/utils'

interface MilestoneBadgeProps {
  milestone: Milestone
  achieved?: boolean
  achievedAt?: string
  onClick?: () => void
  currentProgress?: number // Current count for progress display
}

const MilestoneBadge = ({
  milestone,
  achieved = false,
  achievedAt,
  onClick,
  currentProgress,
}: MilestoneBadgeProps) => {
  // Get badge color based on type and threshold
  const getBadgeColors = () => {
    if (!achieved) {
      return {
        bg: 'bg-gradient-to-br from-gray-100 to-gray-200',
        border: 'border-gray-300',
        text: 'text-gray-600',
        icon: 'text-gray-400',
        shadow: 'shadow-md',
        glow: '',
      }
    }

    // Different colors for different thresholds - more vibrant and metallic
    if (milestone.threshold >= 100) {
      return {
        bg: 'bg-gradient-to-br from-purple-500 via-purple-600 to-purple-700',
        border: 'border-purple-400',
        text: 'text-white',
        icon: 'text-purple-100',
        shadow: 'shadow-xl shadow-purple-500/50',
        glow: 'ring-2 ring-purple-400/30',
      }
    } else if (milestone.threshold >= 50) {
      return {
        bg: 'bg-gradient-to-br from-blue-500 via-blue-600 to-blue-700',
        border: 'border-blue-400',
        text: 'text-white',
        icon: 'text-blue-100',
        shadow: 'shadow-xl shadow-blue-500/50',
        glow: 'ring-2 ring-blue-400/30',
      }
    } else if (milestone.threshold >= 10) {
      return {
        bg: 'bg-gradient-to-br from-emerald-500 via-emerald-600 to-emerald-700',
        border: 'border-emerald-400',
        text: 'text-white',
        icon: 'text-emerald-100',
        shadow: 'shadow-xl shadow-emerald-500/50',
        glow: 'ring-2 ring-emerald-400/30',
      }
    } else {
      return {
        bg: 'bg-gradient-to-br from-amber-400 via-amber-500 to-amber-600',
        border: 'border-amber-300',
        text: 'text-white',
        icon: 'text-amber-50',
        shadow: 'shadow-xl shadow-amber-500/50',
        glow: 'ring-2 ring-amber-400/30',
      }
    }
  }

  const getIcon = () => {
    if (milestone.threshold >= 100) {
      return <Trophy className="h-10 w-10" fill="currentColor" />
    } else if (milestone.threshold >= 50) {
      return <Medal className="h-10 w-10" fill="currentColor" />
    } else {
      return <Award className="h-10 w-10" fill="currentColor" />
    }
  }

  const colors = getBadgeColors()
  const progress = currentProgress !== undefined ? Math.min((currentProgress / milestone.threshold) * 100, 100) : undefined

  return (
    <div
      className={cn(
        'relative cursor-pointer transform transition-all duration-300 hover:scale-105 hover:z-10',
        'flex flex-col items-center'
      )}
      onClick={onClick}
    >
      {/* Badge container - improved octagonal shape */}
      <div
        className={cn(
          'relative w-28 h-28 flex flex-col items-center justify-center',
          'border-2 shadow-lg transition-all duration-300',
          colors.bg,
          colors.border,
          colors.shadow,
          colors.glow,
          achieved ? 'opacity-100' : 'opacity-50',
          'overflow-hidden'
        )}
        style={{
          clipPath: 'polygon(30% 0%, 70% 0%, 100% 30%, 100% 70%, 70% 100%, 30% 100%, 0% 70%, 0% 30%)',
        }}
      >
        {/* Enhanced shine effect */}
        <div className="absolute inset-0 bg-gradient-to-br from-white/30 via-white/10 to-transparent pointer-events-none" />
        <div className="absolute top-0 left-0 right-0 h-1/3 bg-gradient-to-b from-white/20 to-transparent pointer-events-none" />

        {/* Inner glow for achieved badges */}
        {achieved && (
          <div className="absolute inset-2 bg-gradient-to-br from-white/10 to-transparent rounded-lg pointer-events-none" />
        )}

        {/* Icon */}
        <div className={cn('mb-1.5 drop-shadow-lg', colors.icon)}>{getIcon()}</div>

        {/* Threshold number - larger and bolder */}
        <div className={cn('text-base font-extrabold drop-shadow-md', colors.text)}>
          {milestone.threshold}
        </div>

        {/* Progress bar for unearned badges */}
        {!achieved && progress !== undefined && (
          <div className="absolute bottom-0 left-0 right-0 h-1.5 bg-black/30">
            <div
              className="h-full bg-white/70 transition-all duration-500 shadow-sm"
              style={{ width: `${progress}%` }}
            />
          </div>
        )}
      </div>

      {/* Badge label - fixed truncation issue */}
      <div className="mt-3 text-center w-full px-1">
        <div
          className={cn(
            'text-xs font-semibold leading-tight',
            achieved ? 'text-foreground' : 'text-muted-foreground'
          )}
          style={{
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
            wordBreak: 'break-word',
          }}
        >
          {milestone.name}
        </div>
      </div>
    </div>
  )
}

export default MilestoneBadge

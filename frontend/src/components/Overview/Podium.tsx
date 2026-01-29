import { TopEmployee } from '@/types/stats'
import { Badge } from '@/components/ui/badge'
import { Trophy, Medal } from 'lucide-react'

interface PodiumProps {
  topThree: TopEmployee[]
  onEmployeeClick?: (employeeId: number, employeeName: string) => void
}

const Podium = ({ topThree, onEmployeeClick }: PodiumProps) => {
  if (topThree.length < 3) {
    return null
  }

  const [first, second, third] = topThree // 第1名、第2名、第3名

  const PodiumItem = ({
    employee,
    rank,
    podiumHeight,
    position,
  }: {
    employee: TopEmployee
    rank: number
    podiumHeight: number // 颁奖台高度（px）
    position: 'left' | 'center' | 'right'
  }) => {
    let rankColorClass = ''
    let rankShadowClass = ''
    let medalIcon = null
    let medalColor = ''

    if (rank === 1) {
      // 金牌 - 金色渐变
      rankColorClass = 'bg-gradient-to-b from-yellow-400 via-yellow-500 to-yellow-600'
      rankShadowClass = 'shadow-xl shadow-yellow-500/40'
      medalIcon = <Trophy className="h-6 w-6 text-yellow-300 drop-shadow-lg" fill="currentColor" />
      medalColor = 'text-yellow-300'
    } else if (rank === 2) {
      // 银牌 - 银色渐变
      rankColorClass = 'bg-gradient-to-b from-gray-200 via-gray-300 to-gray-400'
      rankShadowClass = 'shadow-lg shadow-gray-400/30'
      medalIcon = <Medal className="h-5 w-5 text-gray-100 drop-shadow-md" fill="currentColor" />
      medalColor = 'text-gray-100'
    } else if (rank === 3) {
      // 铜牌 - 铜色渐变
      rankColorClass = 'bg-gradient-to-b from-orange-500 via-amber-600 to-orange-700'
      rankShadowClass = 'shadow-md shadow-orange-600/30'
      medalIcon = <Medal className="h-5 w-5 text-orange-200 drop-shadow-md" fill="currentColor" />
      medalColor = 'text-orange-200'
    }

    return (
      <div
        className={`flex flex-col items-end justify-end relative ${position === 'center' ? 'order-2 z-10' : position === 'left' ? 'order-1 z-0' : 'order-3 z-0'}`}
        style={{ height: `${podiumHeight + 40}px` }}
      >
        {/* 员工信息卡片 - 悬浮在颁奖台上方 */}
        <div
          className={`absolute w-full ${rankColorClass} ${rankShadowClass} rounded-t-xl border-2 border-white/30 transition-all duration-300 hover:scale-105 hover:shadow-2xl cursor-pointer overflow-visible`}
          style={{ 
            height: `${podiumHeight}px`,
            bottom: '40px'
          }}
          onClick={() => onEmployeeClick?.(employee.employeeId, employee.name)}
        >
          {/* 背景光效 */}
          <div className="absolute inset-0 bg-gradient-to-b from-white/10 to-transparent rounded-t-xl"></div>
          
          <div className="relative flex flex-col items-center justify-center h-full text-center text-white p-2 overflow-visible">
            <div className="mb-0.5 transform hover:scale-110 transition-transform flex-shrink-0">{medalIcon}</div>
            <div className="text-lg font-bold mb-0.5 drop-shadow-lg flex-shrink-0">{rank}</div>
            <div className="font-bold text-xs mb-0.5 drop-shadow-md line-clamp-2 flex-shrink-0 min-h-[1.25rem] flex items-center justify-center px-1">
              {employee.name || 'Unknown'}
            </div>
            {employee.department && (
              <div className="text-[10px] opacity-95 drop-shadow line-clamp-1 flex-shrink-0 px-1">
                {employee.department}
              </div>
            )}
          </div>
        </div>

        {/* 颁奖台底座 - 更美观的设计 */}
        <div
          className={`absolute bottom-0 w-full ${rankColorClass} rounded-lg ${rankShadowClass} relative overflow-hidden border-t-2 border-white/20 cursor-pointer`}
          style={{ height: '40px' }}
          onClick={() => onEmployeeClick?.(employee.employeeId, employee.name)}
        >
          {/* 底座渐变效果 */}
          <div className="absolute inset-0 bg-gradient-to-b from-white/5 to-transparent"></div>
          {/* 底座装饰线 */}
          <div className="absolute top-0 left-0 right-0 h-0.5 bg-white/30"></div>
          {/* 卡片数量显示 */}
          <div className="relative flex items-center justify-center h-full">
            <Badge
              variant="secondary"
              className="bg-white/25 text-white border-white/40 hover:bg-white/35 cursor-pointer backdrop-blur-sm text-xs px-2 py-0.5 font-semibold"
              onClick={(e) => {
                e.stopPropagation()
                onEmployeeClick?.(employee.employeeId, employee.name)
              }}
            >
              {employee.count} {employee.count === 1 ? 'card' : 'cards'}
            </Badge>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full py-6 px-4 bg-gradient-to-b from-background to-muted/20 rounded-lg">
      <div className="relative flex items-end justify-center gap-2 sm:gap-3 md:gap-4 max-w-3xl mx-auto" style={{ minHeight: '240px' }}>
        {/* 第2名 - 左侧，中等高度 */}
        <div className="flex-1 max-w-[140px] sm:max-w-[160px]">
          <PodiumItem
            employee={second}
            rank={2}
            podiumHeight={140}
            position="left"
          />
        </div>
        
        {/* 第1名 - 居中，最高 */}
        <div className="flex-1 max-w-[160px] sm:max-w-[180px]">
          <PodiumItem
            employee={first}
            rank={1}
            podiumHeight={180}
            position="center"
          />
        </div>
        
        {/* 第3名 - 右侧，最低 */}
        <div className="flex-1 max-w-[140px] sm:max-w-[160px]">
          <PodiumItem
            employee={third}
            rank={3}
            podiumHeight={120}
            position="right"
          />
        </div>
      </div>
    </div>
  )
}

export default Podium

import { useMemo } from 'react'

interface WordCloudItem {
  text: string
  value: number
  percentage: number
}

interface WordCloudProps {
  data: WordCloudItem[]
  onClick?: (item: WordCloudItem) => void
}

const WordCloud = ({ data, onClick }: WordCloudProps) => {
  // Calculate min and max values for scaling
  const { minValue, maxValue } = useMemo(() => {
    if (data.length === 0) return { minValue: 0, maxValue: 0 }
    const values = data.map(item => item.value)
    return {
      minValue: Math.min(...values),
      maxValue: Math.max(...values),
    }
  }, [data])

  // Calculate font size based on value (scaled between 14px and 48px)
  const getFontSize = (value: number) => {
    if (maxValue === minValue) return 24
    const ratio = (value - minValue) / (maxValue - minValue)
    return 14 + ratio * 34 // Range: 14px to 48px
  }

  // Calculate opacity based on value (scaled between 0.6 and 1.0)
  const getOpacity = (value: number) => {
    if (maxValue === minValue) return 0.8
    const ratio = (value - minValue) / (maxValue - minValue)
    return 0.6 + ratio * 0.4 // Range: 0.6 to 1.0
  }

  // Color palette
  const colors = [
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

  const getColor = (index: number) => colors[index % colors.length]

  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-96 text-muted-foreground">
        No data available
      </div>
    )
  }

  return (
    <div className="flex flex-wrap items-center justify-center gap-4 p-8 min-h-[400px]">
      {data.map((item, index) => {
        const fontSize = getFontSize(item.value)
        const opacity = getOpacity(item.value)
        const color = getColor(index)

        return (
          <button
            key={index}
            onClick={() => onClick?.(item)}
            className="transition-all duration-200 hover:scale-110 hover:opacity-100"
            style={{
              fontSize: `${fontSize}px`,
              opacity,
              color,
              fontWeight: item.value === maxValue ? 700 : item.value > (minValue + maxValue) / 2 ? 600 : 500,
              cursor: onClick ? 'pointer' : 'default',
            }}
            title={`${item.text}: ${item.value} cards (${item.percentage.toFixed(1)}%)`}
          >
            {item.text}
          </button>
        )
      })}
    </div>
  )
}

export default WordCloud

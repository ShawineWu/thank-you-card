import { useMemo } from 'react'
// @ts-ignore - adorable-word-cloud types may not be perfect
import { AdorableWordCloud, CloudWord, Options, Callbacks } from 'adorable-word-cloud'

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
  // Convert data to CloudWord format
  const words: CloudWord[] = useMemo(() => {
    return data.map(item => ({
      text: item.text,
      value: item.value,
    }))
  }, [data])

  // Color palette matching the original design
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

  // Calculate min and max values for font size scaling
  const { minValue, maxValue } = useMemo(() => {
    if (data.length === 0) return { minValue: 0, maxValue: 0 }
    const values = data.map(item => item.value)
    return {
      minValue: Math.min(...values),
      maxValue: Math.max(...values),
    }
  }, [data])

  // Configure options for the word cloud
  const options: Options = useMemo(() => {
    // Calculate font size range based on data
    const baseMinSize = 20
    const baseMaxSize = 60
    const sizeRange = maxValue > minValue 
      ? [baseMinSize, baseMaxSize] 
      : [baseMinSize, baseMinSize + 20]

    return {
      colors: colors,
      fontFamily: 'system-ui, -apple-system, sans-serif',
      fontSizeRange: sizeRange as [number, number],
      rotationDivision: 0.5, // Allow random rotations (0-0.5 means -90 to 90 degrees)
      spiral: 'archimedean', // Use archimedean spiral for better layout
      padding: 5,
    }
  }, [maxValue, minValue])

  // Configure callbacks
  const callbacks: Callbacks = useMemo(() => {
    if (!onClick) return {}
    
    return {
      onWordClick: (word: CloudWord) => {
        const item = data.find(d => d.text === word.text)
        if (item) {
          onClick(item)
        }
      },
    }
  }, [onClick, data])

  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-96 text-muted-foreground">
        No data available
      </div>
    )
  }

  return (
    <div className="w-full" style={{ height: '500px' }}>
      <AdorableWordCloud 
        words={words} 
        options={options} 
        callbacks={callbacks}
      />
    </div>
  )
}

export default WordCloud

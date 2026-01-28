import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { reactionsApi } from '@/services/reactions'
import { EmojiSummary } from '@/types/card'
import { ThumbsUp, Heart, PartyPopper, Star } from 'lucide-react'

interface CardReactionsProps {
  cardId: number
  reactions: EmojiSummary[]
  onReactionChange?: () => void
  readOnly?: boolean // If true, reactions are display-only (e.g., in Analytics pages)
}

const EMOJI_OPTIONS = [
  { code: '👍', icon: ThumbsUp, label: 'Like' },
  { code: '❤️', icon: Heart, label: 'Love' },
  { code: '🎉', icon: PartyPopper, label: 'Celebrate' },
  { code: '🙌', icon: ThumbsUp, label: 'Praise' },
  { code: '⭐', icon: Star, label: 'Star' },
]

const CardReactions = ({ cardId, reactions, onReactionChange, readOnly = false }: CardReactionsProps) => {
  const queryClient = useQueryClient()
  // Ensure reactions is an array, default to empty array if null/undefined
  const safeReactions = reactions || []
  const [selectedEmoji, setSelectedEmoji] = useState<string | null>(
    safeReactions.find(r => r.count > 0 && r.userIds.length > 0)?.emojiCode || null
  )

  const addReactionMutation = useMutation({
    mutationFn: (emojiCode: string) => reactionsApi.addReaction(cardId, emojiCode),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cards'] })
      onReactionChange?.()
    },
  })

  const removeReactionMutation = useMutation({
    mutationFn: () => reactionsApi.removeReaction(cardId),
    onSuccess: () => {
      setSelectedEmoji(null)
      queryClient.invalidateQueries({ queryKey: ['cards'] })
      onReactionChange?.()
    },
  })

  const handleReactionClick = (emojiCode: string) => {
    if (readOnly) return // Do nothing in read-only mode
    if (selectedEmoji === emojiCode) {
      removeReactionMutation.mutate()
    } else {
      setSelectedEmoji(emojiCode)
      addReactionMutation.mutate(emojiCode)
    }
  }

  const getReactionCount = (emojiCode: string): number => {
    return safeReactions.find(r => r.emojiCode === emojiCode)?.count || 0
  }

  // In read-only mode, just display reactions without interactive buttons
  if (readOnly) {
    return (
      <div className="flex items-center space-x-2 pt-2 border-t">
        {EMOJI_OPTIONS.map(({ code }) => {
          const count = getReactionCount(code)
          // Only show reactions that have counts > 0
          if (count === 0) return null

          return (
            <div
              key={code}
              className="flex items-center space-x-1 px-2 py-1 rounded-md bg-muted/50"
            >
              <span className="text-base">{code}</span>
              <span className="text-xs font-medium text-muted-foreground">{count}</span>
            </div>
          )
        })}
        {/* Show message if no reactions */}
        {safeReactions.every(r => r.count === 0) && (
          <span className="text-xs text-muted-foreground">No reactions</span>
        )}
      </div>
    )
  }

  // Interactive mode (default)
  return (
    <div className="flex items-center space-x-2 pt-2 border-t">
      {EMOJI_OPTIONS.map(({ code }) => {
        const count = getReactionCount(code)
        const isSelected = selectedEmoji === code

        return (
          <Button
            key={code}
            variant={isSelected ? 'default' : 'outline'}
            size="sm"
            className="flex items-center space-x-1 h-8"
            onClick={() => handleReactionClick(code)}
            disabled={addReactionMutation.isPending || removeReactionMutation.isPending}
          >
            <span className="text-base">{code}</span>
            {count > 0 && (
              <span className="text-xs font-medium">{count}</span>
            )}
          </Button>
        )
      })}
    </div>
  )
}

export default CardReactions

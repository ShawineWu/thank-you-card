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
}

const EMOJI_OPTIONS = [
  { code: '👍', icon: ThumbsUp, label: 'Like' },
  { code: '❤️', icon: Heart, label: 'Love' },
  { code: '🎉', icon: PartyPopper, label: 'Celebrate' },
  { code: '🙌', icon: ThumbsUp, label: 'Praise' },
  { code: '⭐', icon: Star, label: 'Star' },
]

const CardReactions = ({ cardId, reactions, onReactionChange }: CardReactionsProps) => {
  const queryClient = useQueryClient()
  const [selectedEmoji, setSelectedEmoji] = useState<string | null>(
    reactions.find(r => r.count > 0 && r.userIds.length > 0)?.emojiCode || null
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
    if (selectedEmoji === emojiCode) {
      removeReactionMutation.mutate()
    } else {
      setSelectedEmoji(emojiCode)
      addReactionMutation.mutate(emojiCode)
    }
  }

  const getReactionCount = (emojiCode: string): number => {
    return reactions.find(r => r.emojiCode === emojiCode)?.count || 0
  }

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

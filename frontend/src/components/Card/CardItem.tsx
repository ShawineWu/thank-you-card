import { CardResponse } from '@/types/card'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { formatRelativeTime } from '@/utils/date'
import CardReactions from './CardReactions'
import ValueBadge from './ValueBadge'

interface CardItemProps {
  card: CardResponse
  onReactionChange?: () => void
}

const CardItem = ({ card, onReactionChange }: CardItemProps) => {
  return (
    <Card className="mb-4">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              <span className="font-semibold text-primary">{card.sender.name}</span>
              <span className="text-muted-foreground">→</span>
              <div className="flex flex-wrap gap-1">
                {card.recipients.map((recipient, idx) => (
                  <span key={recipient.id} className="font-medium">
                    {recipient.name}
                    {idx < card.recipients.length - 1 && ','}
                  </span>
                ))}
              </div>
            </div>
            <div className="text-sm text-muted-foreground">
              {card.sender.department} • {formatRelativeTime(card.createdAt)}
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <p className="mb-4 text-base leading-relaxed whitespace-pre-wrap">{card.reason}</p>
        
        <div className="flex flex-wrap gap-2 mb-4">
          {card.values.map((value) => (
            <ValueBadge key={value.id} value={value} />
          ))}
        </div>

        <CardReactions 
          cardId={card.id} 
          reactions={card.reactions}
          onReactionChange={onReactionChange}
        />
      </CardContent>
    </Card>
  )
}

export default CardItem

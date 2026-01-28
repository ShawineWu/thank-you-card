import { CardResponse } from '@/types/card'
import CardItem from './CardItem'

interface CardListProps {
  cards: CardResponse[]
  onReactionChange?: () => void
  emptyMessage?: string
  readOnly?: boolean // If true, reactions are display-only (e.g., in Analytics pages)
}

const CardList = ({ cards, onReactionChange, emptyMessage = 'No cards found', readOnly = false }: CardListProps) => {
  if (cards.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p className="text-lg">{emptyMessage}</p>
      </div>
    )
  }

  return (
    <div>
      {cards.map((card) => (
        <CardItem 
          key={card.id} 
          card={card} 
          onReactionChange={onReactionChange}
          readOnly={readOnly}
        />
      ))}
    </div>
  )
}

export default CardList

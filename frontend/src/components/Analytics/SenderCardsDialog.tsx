import { useQuery } from '@tanstack/react-query'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Loader2 } from 'lucide-react'
import { cardsApi } from '@/services/cards'
import { CardResponse, CompanyValueSummary } from '@/types/card'
import CardItem from '@/components/Card/CardItem'
import { format } from 'date-fns'

interface SenderCardsDialogProps {
  senderId: number
  senderName: string
  open: boolean
  onOpenChange: (open: boolean) => void
  from?: string
  to?: string
  onValueClick?: (value: CompanyValueSummary) => void
}

const SenderCardsDialog = ({
  senderId,
  senderName,
  open,
  onOpenChange,
  from,
  to,
  onValueClick,
}: SenderCardsDialogProps) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['sender-cards', senderId, from, to],
    queryFn: () =>
      cardsApi.getFeed({
        senderIds: [senderId],
        from: from ? `${from}T00:00:00Z` : undefined,
        to: to ? `${to}T23:59:59Z` : undefined,
        pageSize: 100, // Get all cards sent by this user
      }),
    enabled: open, // Only fetch when dialog is open
  })

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            Cards Sent by {senderName}
            {from && to && (
              <span className="text-sm font-normal text-muted-foreground ml-2">
                ({format(new Date(from), 'MMM d')} - {format(new Date(to), 'MMM d')})
              </span>
            )}
          </DialogTitle>
        </DialogHeader>
        <div className="mt-4">
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
          ) : error ? (
            <div className="text-center py-12">
              <p className="text-destructive">Failed to load cards. Please try again.</p>
            </div>
          ) : data && data.items.length > 0 ? (
            <div className="space-y-4">
              {data.items.map((card: CardResponse) => (
                <CardItem 
                  key={card.id} 
                  card={card} 
                  readOnly 
                  onValueClick={onValueClick}
                />
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No cards found sent by this user.</p>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default SenderCardsDialog

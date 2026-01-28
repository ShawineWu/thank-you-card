import { useState, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Loader2, X, User } from 'lucide-react'
import { cardsApi } from '@/services/cards'
import { CardResponse, CompanyValueSummary, EmployeeSummary } from '@/types/card'
import CardItem from '@/components/Card/CardItem'
import { Badge } from '@/components/ui/badge'
import { format } from 'date-fns'

interface EmployeeCardsDialogProps {
  employeeId: number
  employeeName: string
  open: boolean
  onOpenChange: (open: boolean) => void
  from?: string
  to?: string
  onValueClick?: (value: CompanyValueSummary) => void
}

const EmployeeCardsDialog = ({
  employeeId,
  employeeName,
  open,
  onOpenChange,
  from,
  to,
  onValueClick,
}: EmployeeCardsDialogProps) => {
  const [selectedValueIds, setSelectedValueIds] = useState<number[]>([])
  const [selectedSenderIds, setSelectedSenderIds] = useState<number[]>([])

  const { data, isLoading, error } = useQuery({
    queryKey: ['employee-cards', employeeId, from, to],
    queryFn: () =>
      cardsApi.getFeed({
        recipientIds: [employeeId],
        from: from ? `${from}T00:00:00Z` : undefined,
        to: to ? `${to}T23:59:59Z` : undefined,
        pageSize: 100, // Get all cards for this employee
      }),
    enabled: open, // Only fetch when dialog is open
  })

  // Calculate value counts from all cards
  const valueCounts = useMemo(() => {
    if (!data?.items) return new Map<number, { value: CompanyValueSummary; count: number }>()

    const counts = new Map<number, { value: CompanyValueSummary; count: number }>()
    
    data.items.forEach((card) => {
      card.values.forEach((value) => {
        const existing = counts.get(value.id)
        if (existing) {
          existing.count += 1
        } else {
          counts.set(value.id, { value, count: 1 })
        }
      })
    })

    return counts
  }, [data?.items])

  // Calculate sender counts from all cards
  const senderCounts = useMemo(() => {
    if (!data?.items) return new Map<number, { sender: EmployeeSummary; count: number }>()

    const counts = new Map<number, { sender: EmployeeSummary; count: number }>()
    
    data.items.forEach((card) => {
      const sender = card.sender
      const existing = counts.get(sender.id)
      if (existing) {
        existing.count += 1
      } else {
        counts.set(sender.id, { sender, count: 1 })
      }
    })

    return counts
  }, [data?.items])

  // Filter cards based on selected values and senders
  const filteredCards = useMemo(() => {
    if (!data?.items) return []
    
    let filtered = data.items

    // Filter by values
    if (selectedValueIds.length > 0) {
      filtered = filtered.filter((card) =>
        card.values.some((value) => selectedValueIds.includes(value.id))
      )
    }

    // Filter by senders
    if (selectedSenderIds.length > 0) {
      filtered = filtered.filter((card) =>
        selectedSenderIds.includes(card.sender.id)
      )
    }

    return filtered
  }, [data?.items, selectedValueIds, selectedSenderIds])

  const handleValueToggle = (valueId: number) => {
    setSelectedValueIds((prev) =>
      prev.includes(valueId) ? prev.filter((id) => id !== valueId) : [...prev, valueId]
    )
  }

  const handleSenderToggle = (senderId: number) => {
    setSelectedSenderIds((prev) =>
      prev.includes(senderId) ? prev.filter((id) => id !== senderId) : [...prev, senderId]
    )
  }

  const handleClearFilters = () => {
    setSelectedValueIds([])
    setSelectedSenderIds([])
  }

  const hasActiveFilters = selectedValueIds.length > 0 || selectedSenderIds.length > 0

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            Cards Received by {employeeName}
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
            <>
              {/* Filters */}
              <div className="mb-6 space-y-6">
                {/* Value Statistics and Filters */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-foreground">Filter by Values/Credos</h3>
                  </div>
                  <div className="flex flex-wrap gap-2.5">
                    {Array.from(valueCounts.values())
                      .sort((a, b) => b.count - a.count) // Sort by count descending
                      .map(({ value, count }) => {
                        const isSelected = selectedValueIds.includes(value.id)
                        const variant = value.type === 'VALUE' ? 'default' : 'success'
                        
                        return (
                          <button
                            key={value.id}
                            onClick={() => handleValueToggle(value.id)}
                            className={`group relative inline-flex items-center gap-2 px-3 py-1.5 rounded-full border transition-all duration-200 ${
                              isSelected
                                ? variant === 'default'
                                  ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                                  : 'bg-green-600 text-white border-green-600 shadow-sm'
                                : 'bg-background hover:bg-accent border-border hover:border-primary/50'
                            }`}
                          >
                            <span className={`text-sm font-medium ${isSelected ? 'text-inherit' : 'text-foreground'}`}>
                              {value.name}
                            </span>
                            <span
                              className={`inline-flex items-center justify-center min-w-[20px] h-5 px-1.5 rounded-full text-xs font-semibold ${
                                isSelected
                                  ? 'bg-white/20 text-inherit'
                                  : 'bg-muted text-muted-foreground group-hover:bg-primary/10 group-hover:text-primary'
                              }`}
                            >
                              {count}
                            </span>
                          </button>
                        )
                      })}
                  </div>
                </div>

                {/* User Statistics and Filters */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-foreground">Filter by Sender</h3>
                    {hasActiveFilters && (
                      <button
                        onClick={handleClearFilters}
                        className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1.5 transition-colors"
                      >
                        <X className="h-3.5 w-3.5" />
                        Clear filters
                      </button>
                    )}
                  </div>
                  <div className="flex flex-wrap gap-2.5">
                    {Array.from(senderCounts.values())
                      .sort((a, b) => b.count - a.count) // Sort by count descending
                      .map(({ sender, count }) => {
                        const isSelected = selectedSenderIds.includes(sender.id)
                        
                        return (
                          <button
                            key={sender.id}
                            onClick={() => handleSenderToggle(sender.id)}
                            className={`group relative inline-flex items-center gap-2 px-3 py-1.5 rounded-full border transition-all duration-200 ${
                              isSelected
                                ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                                : 'bg-background hover:bg-accent border-border hover:border-primary/50'
                            }`}
                          >
                            <User className={`h-3.5 w-3.5 ${isSelected ? 'text-inherit' : 'text-muted-foreground'}`} />
                            <span className={`text-sm font-medium ${isSelected ? 'text-inherit' : 'text-foreground'}`}>
                              {sender.name}
                            </span>
                            <span
                              className={`inline-flex items-center justify-center min-w-[20px] h-5 px-1.5 rounded-full text-xs font-semibold ${
                                isSelected
                                  ? 'bg-white/20 text-inherit'
                                  : 'bg-muted text-muted-foreground group-hover:bg-primary/10 group-hover:text-primary'
                              }`}
                            >
                              {count}
                            </span>
                          </button>
                        )
                      })}
                  </div>
                </div>

                {hasActiveFilters && (
                  <div className="flex items-center gap-2 pt-2 border-t">
                    <p className="text-xs text-muted-foreground">
                      Showing <span className="font-semibold text-foreground">{filteredCards.length}</span> of{' '}
                      <span className="font-semibold text-foreground">{data.items.length}</span> cards
                    </p>
                  </div>
                )}
              </div>

              {/* Filtered Cards List */}
              <div className="space-y-4">
                {filteredCards.length > 0 ? (
                  filteredCards.map((card: CardResponse) => (
                    <CardItem 
                      key={card.id} 
                      card={card} 
                      readOnly 
                      onValueClick={onValueClick}
                    />
                  ))
                ) : (
                  <div className="text-center py-8">
                    <p className="text-muted-foreground">
                      No cards match the selected filters.
                    </p>
                  </div>
                )}
              </div>
            </>
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground">No cards found for this employee.</p>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default EmployeeCardsDialog

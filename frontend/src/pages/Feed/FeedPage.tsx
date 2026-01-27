import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { cardsApi } from '@/services/cards'
import CardList from '@/components/Card/CardList'
import FeedFilters from '@/components/Filters/FeedFilters'
import { CardFilters } from '@/types/card'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'

const FeedPage = () => {
  const [filters, setFilters] = useState<CardFilters>({
    page: 1,
    pageSize: 20,
  })

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['cards', 'feed', filters],
    queryFn: () => cardsApi.getFeed(filters),
  })

  const handleFiltersChange = (newFilters: CardFilters) => {
    setFilters({ ...newFilters, page: 1 })
  }

  const handleLoadMore = () => {
    setFilters(prev => ({ ...prev, page: (prev.page || 1) + 1 }))
  }

  if (isLoading && !data) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-destructive">Failed to load feed. Please try again.</p>
        <Button onClick={() => refetch()} className="mt-4">
          Retry
        </Button>
      </div>
    )
  }

  const hasMore = data ? (data.items.length < data.total) : false

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Company Feed</h1>
        <p className="text-muted-foreground">See all thank you cards from across the company</p>
      </div>

      <FeedFilters filters={filters} onFiltersChange={handleFiltersChange} />

      <div className="mt-6">
        <CardList 
          cards={data?.items || []} 
          onReactionChange={() => refetch()}
          emptyMessage="No cards found. Be the first to send a thank you card!"
        />
      </div>

      {hasMore && (
        <div className="text-center mt-6">
          <Button 
            onClick={handleLoadMore} 
            variant="outline"
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Loading...
              </>
            ) : (
              'Load More'
            )}
          </Button>
        </div>
      )}
    </div>
  )
}

export default FeedPage

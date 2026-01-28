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

  const totalPages = data ? Math.ceil(data.total / (filters.pageSize || 20)) : 0
  const currentPage = filters.page || 1
  const hasMore = data ? currentPage < totalPages : false

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setFilters(prev => ({ ...prev, page: newPage }))
      window.scrollTo({ top: 0, behavior: 'smooth' })
    }
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Company Feed</h1>
        <p className="text-muted-foreground">See all thank you cards from across the company</p>
      </div>

      <FeedFilters filters={filters} onFiltersChange={handleFiltersChange} />

      <div className="mt-6">
        {data && data.total > 0 && (
          <div className="mb-4 text-sm text-muted-foreground">
            Showing {((currentPage - 1) * (filters.pageSize || 20)) + 1} to {Math.min(currentPage * (filters.pageSize || 20), data.total)} of {data.total} cards
          </div>
        )}
        <CardList 
          cards={data?.items || []} 
          onReactionChange={() => refetch()}
          emptyMessage="No cards found. Be the first to send a thank you card!"
          readOnly={true}
        />
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2 mt-6">
          <Button
            variant="outline"
            size="sm"
            onClick={() => handlePageChange(currentPage - 1)}
            disabled={currentPage === 1 || isLoading}
          >
            Previous
          </Button>
          
          {/* Page numbers */}
          <div className="flex items-center gap-1">
            {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
              let pageNum: number
              if (totalPages <= 5) {
                pageNum = i + 1
              } else if (currentPage <= 3) {
                pageNum = i + 1
              } else if (currentPage >= totalPages - 2) {
                pageNum = totalPages - 4 + i
              } else {
                pageNum = currentPage - 2 + i
              }
              
              return (
                <Button
                  key={pageNum}
                  variant={currentPage === pageNum ? "default" : "outline"}
                  size="sm"
                  onClick={() => handlePageChange(pageNum)}
                  disabled={isLoading}
                  className="min-w-[40px]"
                >
                  {pageNum}
                </Button>
              )
            })}
          </div>

          <Button
            variant="outline"
            size="sm"
            onClick={() => handlePageChange(currentPage + 1)}
            disabled={currentPage === totalPages || isLoading}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  )
}

export default FeedPage

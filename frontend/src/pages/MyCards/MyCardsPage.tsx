import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { cardsApi } from '@/services/cards'
import { CardFilters } from '@/types/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import CardList from '@/components/Card/CardList'
import { Loader2, Inbox, Send } from 'lucide-react'

const MyCardsPage = () => {
  const [activeTab, setActiveTab] = useState<'received' | 'sent'>('received')
  const [filters] = useState<CardFilters>({
    page: 1,
    pageSize: 20,
  })

  const { data: receivedData, isLoading: isLoadingReceived, refetch: refetchReceived } = useQuery({
    queryKey: ['cards', 'received', filters],
    queryFn: () => cardsApi.getMyReceived(filters),
    enabled: activeTab === 'received',
  })

  const { data: sentData, isLoading: isLoadingSent, refetch: refetchSent } = useQuery({
    queryKey: ['cards', 'sent', filters],
    queryFn: () => cardsApi.getMySent(filters),
    enabled: activeTab === 'sent',
  })

  const currentData = activeTab === 'received' ? receivedData : sentData
  const isLoading = activeTab === 'received' ? isLoadingReceived : isLoadingSent

  if (isLoading && !currentData) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">My Cards</h1>
        <p className="text-muted-foreground">View cards you've received and sent</p>
      </div>

      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'received' | 'sent')}>
        <TabsList>
          <TabsTrigger value="received" className="flex items-center space-x-2">
            <Inbox className="h-4 w-4" />
            <span>Received ({receivedData?.total || 0})</span>
          </TabsTrigger>
          <TabsTrigger value="sent" className="flex items-center space-x-2">
            <Send className="h-4 w-4" />
            <span>Sent ({sentData?.total || 0})</span>
          </TabsTrigger>
        </TabsList>

        <TabsContent value="received" className="mt-6">
          <CardList
            cards={receivedData?.items || []}
            onReactionChange={() => refetchReceived()}
            emptyMessage="You haven't received any cards yet."
          />
        </TabsContent>

        <TabsContent value="sent" className="mt-6">
          <CardList
            cards={sentData?.items || []}
            onReactionChange={() => refetchSent()}
            emptyMessage="You haven't sent any cards yet. Create your first one!"
          />
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default MyCardsPage

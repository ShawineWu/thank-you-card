import React, { useEffect, useState, useCallback, useRef } from "react";
import { api } from "@/services/api";
import type { Card as RecognitionCard, ValueResponse } from "@/services/api";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Loader2 } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { getValueColor, getValueBadgeClasses } from "@/constants/valueColors";

export const Feed: React.FC = () => {
  const [cards, setCards] = useState<RecognitionCard[]>([]);
  const [values, setValues] = useState<ValueResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [totalItems, setTotalItems] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [page, setPage] = useState(1);
  const [selectedValues, setSelectedValues] = useState<Set<string>>(new Set());
  const loadMoreRef = useRef<HTMLDivElement>(null);
  const pageSize = 20;

  useEffect(() => {
    const fetchValues = async () => {
      try {
        const res = await api.getValues();
        setValues(res.data.data);
      } catch (error) {
        console.error("Failed to fetch values:", error);
      }
    };
    fetchValues();
  }, []);

  // 初始加载或筛选条件变化时重新加载
  const fetchCards = useCallback(async (currentPage: number, append = false) => {
    try {
      if (!append) {
        setLoading(true);
      } else {
        setLoadingMore(true);
      }
      
      const res = await api.getCards({
        page: currentPage,
        pageSize,
        valueIds: Array.from(selectedValues),
      });
      
      const newCards = res.data.data.data;
      const pagination = res.data.data.pagination;
      
      if (append) {
        setCards(prev => [...prev, ...newCards]);
      } else {
        setCards(newCards);
      }
      
      setTotalItems(pagination.totalItems);
      setHasMore(pagination.hasNext);
    } catch (error) {
      console.error("Failed to fetch cards:", error);
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }, [selectedValues]);

  // 筛选条件变化时重新加载
  useEffect(() => {
    setPage(1);
    fetchCards(1, false);
  }, [fetchCards]);

  // 加载更多
  const loadMore = useCallback(() => {
    if (!loadingMore && hasMore && !loading) {
      const nextPage = page + 1;
      setPage(nextPage);
      fetchCards(nextPage, true);
    }
  }, [loadingMore, hasMore, loading, page, fetchCards]);

  // Intersection Observer 监听滚动到底部
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading && !loadingMore) {
          loadMore();
        }
      },
      { threshold: 0.1 }
    );

    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current);
    }

    return () => observer.disconnect();
  }, [hasMore, loading, loadingMore, loadMore]);

  const toggleValue = (valueId: string) => {
    setSelectedValues((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(valueId)) {
        newSet.delete(valueId);
      } else {
        newSet.add(valueId);
      }
      return newSet;
    });
  };

  const getRecipientNames = (card: RecognitionCard) => {
    return card.recipients.map((r) => r.name).join(", ");
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Company Feed</h1>
        <p className="text-muted-foreground">
          See all thank you cards from across the company
        </p>
      </div>

      {/* Search */}
      {/* <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search cards..."
          value={filters.search}
          onChange={handleSearchChange}
          className="pl-10 bg-background border-border"
        />
      </div> */}

      {/* Sender and Recipient Filters */}
      {/* <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="text-sm font-medium mb-2 block">Filter by Sender</label>
          <Input
            placeholder="Search sender..."
            value={filters.senderId}
            onChange={handleSenderChange}
            className="bg-background border-border"
          />
        </div>
        <div>
          <label className="text-sm font-medium mb-2 block">Filter by Recipient</label>
          <Input
            placeholder="Search recipient..."
            value={filters.recipientId}
            onChange={handleRecipientChange}
            className="bg-background border-border"
          />
        </div>
      </div> */}

      {/* Value Tags Filter - 分组显示 Credo 和 Value */}
      <div className="space-y-3">
        {/* Credo */}
        {values.filter(v => v.type === 'Credo').length > 0 && (
          <div>
            <label className="text-sm font-medium mb-2 block">Filter by Credo</label>
            <div className="flex flex-wrap gap-2">
              {values.filter(v => v.type === 'Credo').map((value) => {
                const color = getValueColor(value.name);
                return (
                  <div key={value.id} className="relative group">
                    <Badge
                      variant={selectedValues.has(value.id) ? "default" : "outline"}
                      className={`cursor-pointer transition-colors ${
                        selectedValues.has(value.id)
                          ? `${color.bg} ${color.text} border ${color.border}`
                          : `border ${color.border} ${color.text} ${color.bgHover}`
                      }`}
                      onClick={() => toggleValue(value.id)}
                    >
                      {value.name}
                    </Badge>
                    {/* Tooltip */}
                    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-popover border rounded-md shadow-lg text-sm w-64 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50 pointer-events-none">
                      <p className="text-foreground">{value.description}</p>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
        
        {/* Value */}
        {values.filter(v => v.type === 'Value').length > 0 && (
          <div>
            <label className="text-sm font-medium mb-2 block">Filter by Value</label>
            <div className="flex flex-wrap gap-2">
              {values.filter(v => v.type === 'Value').map((value) => {
                const color = getValueColor(value.name);
                return (
                  <div key={value.id} className="relative group">
                    <Badge
                      variant={selectedValues.has(value.id) ? "default" : "outline"}
                      className={`cursor-pointer transition-colors ${
                        selectedValues.has(value.id)
                          ? `${color.bg} ${color.text} border ${color.border}`
                          : `border ${color.border} ${color.text} ${color.bgHover}`
                      }`}
                      onClick={() => toggleValue(value.id)}
                    >
                      {value.name}
                    </Badge>
                    {/* Tooltip */}
                    <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-popover border rounded-md shadow-lg text-sm w-64 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50 pointer-events-none">
                      <p className="text-foreground">{value.description}</p>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>

      {/* Results Count */}
      <p className="text-sm text-muted-foreground">
        Showing 1 to {cards.length} of {totalItems} cards
      </p>

      {/* Cards List */}
      <div className="space-y-4">
        {loading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="p-6 border border-border rounded-lg bg-card">
              <Skeleton className="h-5 w-48 mb-2" />
              <Skeleton className="h-4 w-32 mb-4" />
              <Skeleton className="h-20 w-full" />
            </div>
          ))
        ) : cards.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground">
            No cards found matching your filters.
          </div>
        ) : (
          cards.map((card) => (
            <div
              key={card.id}
              className="p-6 border border-border rounded-lg bg-card hover:shadow-md transition-shadow"
            >
              <div className="flex items-center gap-2 mb-1">
                <span className="font-semibold text-primary">
                  {card.senderName || card.senderId}
                </span>
                <span className="text-muted-foreground">→</span>
                <span className="font-semibold">{getRecipientNames(card)}</span>
              </div>
              <p className="text-xs text-muted-foreground mb-3">
                {formatDistanceToNow(new Date(card.createdAt), { addSuffix: true })}
              </p>
              <div className="text-sm text-foreground whitespace-pre-wrap mb-4">
                {card.recognitionReason}
              </div>
              {card.selectedValues.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {card.selectedValues.map((value) => (
                    <div key={value.id} className="relative group">
                      <Badge
                        className={`text-xs font-medium px-3 py-1 border ${getValueBadgeClasses(value.name)}`}
                      >
                        {value.name}
                      </Badge>
                      {/* Tooltip */}
                      <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-popover border rounded-md shadow-lg text-sm w-64 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50 pointer-events-none">
                        <p className="text-foreground">{value.description}</p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))
        )}
        
        {/* 加载更多触发器 */}
        <div ref={loadMoreRef} className="py-4 flex justify-center">
          {loadingMore && (
            <div className="flex items-center gap-2 text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span>加载更多...</span>
            </div>
          )}
          {!hasMore && cards.length > 0 && (
            <p className="text-sm text-muted-foreground">已加载全部卡片</p>
          )}
        </div>
      </div>
    </div>
  );
};

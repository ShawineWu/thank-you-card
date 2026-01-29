import React, { useEffect, useState, useCallback } from "react";
import { api } from "@/services/api";
import type { Card as RecognitionCard, ValueResponse, CardFilter } from "@/services/api";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Search } from "lucide-react";
import { formatDistanceToNow } from "date-fns";

export const Feed: React.FC = () => {
  const [cards, setCards] = useState<RecognitionCard[]>([]);
  const [values, setValues] = useState<ValueResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [totalItems, setTotalItems] = useState(0);
  const [filters, setFilters] = useState<CardFilter>({
    page: 1,
    pageSize: 20,
    search: "",
    senderId: "",
    recipientId: "",
    valueIds: [],
  });
  const [selectedValues, setSelectedValues] = useState<Set<string>>(new Set());

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

  const fetchCards = useCallback(async () => {
    try {
      setLoading(true);
      const res = await api.getCards({
        ...filters,
        valueIds: Array.from(selectedValues),
      });
      setCards(res.data.data.data);
      setTotalItems(res.data.data.pagination.totalItems);
    } catch (error) {
      console.error("Failed to fetch cards:", error);
    } finally {
      setLoading(false);
    }
  }, [filters, selectedValues]);

  useEffect(() => {
    fetchCards();
  }, [fetchCards]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilters((prev) => ({ ...prev, search: e.target.value, page: 1 }));
  };

  const handleSenderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilters((prev) => ({ ...prev, senderId: e.target.value, page: 1 }));
  };

  const handleRecipientChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFilters((prev) => ({ ...prev, recipientId: e.target.value, page: 1 }));
  };

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
    setFilters((prev) => ({ ...prev, page: 1 }));
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
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search cards..."
          value={filters.search}
          onChange={handleSearchChange}
          className="pl-10 bg-background border-border"
        />
      </div>

      {/* Sender and Recipient Filters */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
      </div>

      {/* Value Tags Filter */}
      <div>
        <label className="text-sm font-medium mb-2 block">Filter by Values/Credos</label>
        <div className="flex flex-wrap gap-2">
          {values.map((value) => (
            <Badge
              key={value.id}
              variant={selectedValues.has(value.id) ? "default" : "outline"}
              className={`cursor-pointer transition-colors ${
                selectedValues.has(value.id)
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-accent"
              }`}
              onClick={() => toggleValue(value.id)}
            >
              {value.name}
            </Badge>
          ))}
        </div>
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
                  {card.selectedValues.map((value, index) => (
                    <Badge
                      key={value.id}
                      className={`text-white text-xs font-medium px-3 py-1 ${
                        index % 2 === 0
                          ? "bg-blue-600 hover:bg-blue-700"
                          : "bg-green-600 hover:bg-green-700"
                      }`}
                    >
                      {value.name}
                    </Badge>
                  ))}
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

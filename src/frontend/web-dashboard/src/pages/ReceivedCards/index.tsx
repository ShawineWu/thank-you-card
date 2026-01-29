import React, { useEffect, useState, useCallback } from "react";
import { useAuthStore } from "@/store/authStore";
import { api } from "@/services/api";
import type {
  Card as CardType,
  PaginatedResponse,
  CardFilter,
} from "@/services/api";
import { CardList } from "@/components/Cards/CardList";
import { CardFilters } from "@/components/Cards/CardFilters";
import { Inbox } from "lucide-react";

export const ReceivedCards: React.FC = () => {
  const { user } = useAuthStore();
  const [data, setData] = useState<PaginatedResponse<CardType> | null>(null);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<Partial<CardFilter>>({});

  const employeeId = user?.profile?.sub || "";

  const fetchReceivedCards = useCallback(
    async (pageNum: number, currentFilters: Partial<CardFilter>) => {
      if (!employeeId) return;

      try {
        setLoading(true);
        const res = await api.getCards({
          ...currentFilters,
          page: pageNum,
          pageSize: 10,
          recipientId: employeeId,
        });
        setData(res.data.data);
      } catch (error) {
        console.error("Failed to fetch received cards:", error);
      } finally {
        setLoading(false);
      }
    },
    [employeeId],
  );

  // Debounce filter changes
  useEffect(() => {
    const timer = setTimeout(() => {
      fetchReceivedCards(page, filters);
    }, 400);

    return () => clearTimeout(timer);
  }, [page, filters, fetchReceivedCards]);

  const handleFilterChange = (newFilters: Partial<CardFilter>) => {
    setFilters(newFilters);
    setPage(1); // Reset to first page on filter change
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <div className="p-3 rounded-2xl bg-primary/10 text-primary">
          <Inbox size={24} />
        </div>
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Received Cards</h1>
          <p className="text-muted-foreground">
            Appreciation and feedback shared with you by your colleagues.
          </p>
        </div>
      </div>

      <CardFilters
        onFilterChange={handleFilterChange}
        showRecipientFilter={false}
      />

      <CardList
        data={data}
        loading={loading}
        onPageChange={setPage}
        showRecipient={false} // Since we know we are the recipient
      />
    </div>
  );
};

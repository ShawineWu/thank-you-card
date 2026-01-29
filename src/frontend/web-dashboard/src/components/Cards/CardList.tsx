import React, { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import type { Card as CardType, PaginatedResponse } from "@/services/api";
import { format } from "date-fns";
import { Eye, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { CardDetail } from "./CardDetail";
import { Skeleton } from "@/components/ui/skeleton";

interface CardListProps {
  data: PaginatedResponse<CardType> | null;
  loading: boolean;
  onPageChange: (page: number) => void;
  showRecipient?: boolean;
  showSender?: boolean;
}

export const CardList: React.FC<CardListProps> = ({
  data,
  loading,
  onPageChange,
  showRecipient = true,
  showSender = true,
}) => {
  const [selectedCard, setSelectedCard] = useState<CardType | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);

  const handleViewDetails = (card: CardType) => {
    setSelectedCard(card);
    setIsDetailOpen(true);
  };

  const getInitials = (id: string) => id.substring(0, 2).toUpperCase();

  if (loading) {
    return (
      <div className="space-y-4">
        {[1, 2, 3, 4, 5].map((i) => (
          <Skeleton key={i} className="h-16 w-full rounded-xl" />
        ))}
      </div>
    );
  }

  if (!data || data.data.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20 bg-slate-900/30 rounded-3xl border border-dashed border-slate-800">
        <div className="p-4 rounded-full bg-slate-900 text-slate-500 mb-4">
          <ExternalLink size={32} />
        </div>
        <h3 className="text-xl font-semibold text-slate-300">No cards found</h3>
        <p className="text-slate-500 mt-2">
          When recognitions are shared, they'll appear here.
        </p>
      </div>
    );
  }

  const { pagination } = data;

  return (
    <div className="space-y-6">
      <div className="rounded-2xl border border-slate-800 bg-slate-950/50 overflow-hidden">
        <Table>
          <TableHeader className="bg-slate-900/50">
            <TableRow className="border-slate-800 hover:bg-transparent">
              {showSender && (
                <TableHead className="w-[150px]">Sender</TableHead>
              )}
              {showRecipient && <TableHead>Recipients</TableHead>}
              <TableHead className="max-w-[300px]">Reason</TableHead>
              <TableHead>Values</TableHead>
              <TableHead className="text-right">Date</TableHead>
              <TableHead className="w-[100px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((card) => (
              <TableRow
                key={card.id}
                className="border-slate-800 hover:bg-slate-900/30 transition-colors"
              >
                {showSender && (
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Avatar className="h-7 w-7 border border-slate-800">
                        <AvatarFallback className="text-[10px] bg-primary/10 text-primary">
                          {getInitials(card.senderId)}
                        </AvatarFallback>
                      </Avatar>
                      <span className="text-sm font-medium">
                        {card.senderId}
                      </span>
                    </div>
                  </TableCell>
                )}
                {showRecipient && (
                  <TableCell>
                    <div className="flex flex-wrap gap-1">
                      {card.recipients.slice(0, 2).map((r) => (
                        <Badge
                          key={r}
                          variant="outline"
                          className="text-[10px] py-0 px-2 border-slate-700 bg-slate-900"
                        >
                          {r}
                        </Badge>
                      ))}
                      {card.recipients.length > 2 && (
                        <Badge
                          variant="outline"
                          className="text-[10px] py-0 px-2 border-slate-700"
                        >
                          +{card.recipients.length - 2}
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                )}
                <TableCell className="max-w-[300px]">
                  <p className="text-sm text-slate-300 truncate">
                    {card.recognitionReason}
                  </p>
                </TableCell>
                <TableCell>
                  <div className="flex gap-1">
                    {card.selectedValues.slice(0, 1).map((v) => (
                      <Badge
                        key={v.id}
                        className="bg-primary/10 text-primary border-none text-[10px] px-2"
                      >
                        {v.name}
                      </Badge>
                    ))}
                    {card.selectedValues.length > 1 && (
                      <span className="text-[10px] text-slate-500">
                        +{card.selectedValues.length - 1}
                      </span>
                    )}
                  </div>
                </TableCell>
                <TableCell className="text-right text-xs text-slate-500">
                  {format(new Date(card.createdAt), "MMM d, yyyy")}
                </TableCell>
                <TableCell className="text-right">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-slate-500 hover:text-primary hover:bg-primary/10"
                    onClick={() => handleViewDetails(card)}
                  >
                    <Eye size={16} />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <Button
              variant="ghost"
              disabled={!pagination.hasPrevious}
              onClick={() => onPageChange(pagination.page - 1)}
              className="gap-1 pl-2.5"
            >
              <PaginationPrevious />
            </Button>
          </PaginationItem>

          {Array.from({ length: pagination.totalPages }, (_, i) => i + 1).map(
            (p) => (
              <PaginationItem key={p}>
                <PaginationLink
                  isActive={p === pagination.page}
                  onClick={() => onPageChange(p)}
                  className={
                    p === pagination.page
                      ? "bg-primary text-primary-foreground"
                      : "text-slate-400 hover:bg-slate-900"
                  }
                >
                  {p}
                </PaginationLink>
              </PaginationItem>
            ),
          )}

          <PaginationItem>
            <Button
              variant="ghost"
              disabled={!pagination.hasNext}
              onClick={() => onPageChange(pagination.page + 1)}
              className="gap-1 pr-2.5"
            >
              <PaginationNext />
            </Button>
          </PaginationItem>
        </PaginationContent>
      </Pagination>

      <CardDetail
        card={selectedCard}
        isOpen={isDetailOpen}
        onOpenChange={setIsDetailOpen}
      />
    </div>
  );
};

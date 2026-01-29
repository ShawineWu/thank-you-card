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
import { getValueBadgeClasses } from "@/constants/valueColors";

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
      <div className="flex flex-col items-center justify-center py-20 bg-white rounded-xl border border-gray-100 shadow-sm">
        <div className="p-4 rounded-full bg-gray-50 text-gray-400 mb-4">
          <ExternalLink size={32} />
        </div>
        <h3 className="text-lg font-semibold text-gray-700">No cards found</h3>
        <p className="text-gray-400 mt-2 text-sm">
          When recognitions are shared, they'll appear here.
        </p>
      </div>
    );
  }

  const { pagination } = data;

  return (
    <div className="space-y-6">
      <div className="rounded-xl border border-gray-100 bg-white shadow-sm overflow-hidden">
        <Table>
          <TableHeader className="bg-gray-50/80">
            <TableRow className="border-gray-100 hover:bg-transparent">
              {showSender && (
                <TableHead className="w-[150px] text-gray-500 font-medium">Sender</TableHead>
              )}
              {showRecipient && <TableHead className="text-gray-500 font-medium">Recipients</TableHead>}
              <TableHead className="max-w-[300px] text-gray-500 font-medium">Reason</TableHead>
              <TableHead className="text-gray-500 font-medium">Values</TableHead>
              <TableHead className="text-right text-gray-500 font-medium">Date</TableHead>
              <TableHead className="w-[100px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((card) => (
              <TableRow
                key={card.id}
                className="border-gray-100 hover:bg-gray-50/50 transition-colors"
              >
                {showSender && (
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Avatar className="h-7 w-7 border border-gray-200">
                        <AvatarFallback className="text-[10px] bg-indigo-50 text-indigo-600">
                          {getInitials(card.senderName || card.senderId)}
                        </AvatarFallback>
                      </Avatar>
                      <span className="text-sm font-medium text-gray-700">
                        {card.senderName || card.senderId}
                      </span>
                    </div>
                  </TableCell>
                )}
                {showRecipient && (
                  <TableCell>
                    <div className="flex flex-col gap-1">
                      {card.recipients.slice(0, 2).map((r) => (
                        <div key={r.id} className="flex items-center gap-2">
                          <Avatar className="h-7 w-7 border border-gray-200">
                            <AvatarFallback className="text-[10px] bg-pink-50 text-pink-600">
                              {getInitials(r.name || r.id)}
                            </AvatarFallback>
                          </Avatar>
                          <span className="text-sm font-medium text-gray-700">
                            {r.name}
                          </span>
                        </div>
                      ))}
                      {card.recipients.length > 2 && (
                        <span className="text-xs text-gray-400 ml-9">
                          +{card.recipients.length - 2} more
                        </span>
                      )}
                    </div>
                  </TableCell>
                )}
                <TableCell className="max-w-[300px]">
                  <p className="text-sm text-gray-600 truncate">
                    {card.recognitionReason}
                  </p>
                </TableCell>
                <TableCell>
                  <div className="flex gap-1">
                    {card.selectedValues.slice(0, 1).map((v) => (
                      <Badge
                        key={v.id}
                        className={`border-none text-[10px] px-2 font-medium ${getValueBadgeClasses(v.name)}`}
                      >
                        {v.name}
                      </Badge>
                    ))}
                    {card.selectedValues.length > 1 && (
                      <span className="text-[10px] text-gray-400">
                        +{card.selectedValues.length - 1}
                      </span>
                    )}
                  </div>
                </TableCell>
                <TableCell className="text-right text-xs text-gray-400">
                  {format(new Date(card.createdAt), "MMM d, yyyy")}
                </TableCell>
                <TableCell className="text-right">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg"
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
              className="gap-1 pl-2.5 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg"
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
                      ? "bg-indigo-500 text-white rounded-lg"
                      : "text-gray-500 hover:bg-gray-100 rounded-lg"
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
              className="gap-1 pr-2.5 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg"
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

import React, { useEffect, useState } from "react";
import { api, type ValueResponse } from "@/services/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { Eye } from "lucide-react";

type FilterType = "all" | "value" | "credo";

export const CompanyValues: React.FC = () => {
  const [values, setValues] = useState<ValueResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailValue, setDetailValue] = useState<ValueResponse | null>(null);
  const [detailOpen, setDetailOpen] = useState(false);
  const [filter, setFilter] = useState<FilterType>("all");

  useEffect(() => {
    const fetchValues = async () => {
      try {
        setLoading(true);
        const res = await api.getValues();
        setValues(res.data.data ?? []);
      } catch (error) {
        console.error("Failed to fetch company values:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchValues();
  }, []);

  const valueItems = values.filter((v) => v.type === "Value");
  const credoItems = values.filter((v) => v.type === "Credo");

  const filteredItems =
    filter === "all"
      ? values
      : filter === "value"
        ? valueItems
        : credoItems;

  const openDetails = (v: ValueResponse) => {
    setDetailValue(v);
    setDetailOpen(true);
  };

  const TabButton: React.FC<{
    active: boolean;
    onClick: () => void;
    children: React.ReactNode;
  }> = ({ active, onClick, children }) => (
    <button
      onClick={onClick}
      className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
        active
          ? "border-primary text-foreground"
          : "border-transparent text-muted-foreground hover:text-foreground"
      }`}
    >
      {children}
    </button>
  );

  const ValueCard: React.FC<{
    item: ValueResponse;
    onViewDetails: () => void;
  }> = ({ item, onViewDetails }) => (
    <div className="p-5 border border-border rounded-lg bg-card hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between gap-3 mb-3">
        <h3 className="font-semibold text-base">{item.name}</h3>
        <Badge
          variant="outline"
          className={`shrink-0 text-xs ${
            item.type === "Credo"
              ? "bg-green-500 text-white border-green-500"
              : "bg-blue-500 text-white border-blue-500"
          }`}
        >
          {item.type.toUpperCase()}
        </Badge>
      </div>
      <p className="text-sm text-muted-foreground line-clamp-3 mb-4">
        {item.description}
      </p>
      <Button
        variant="outline"
        size="sm"
        onClick={onViewDetails}
        className="w-full justify-center gap-2"
      >
        <Eye className="h-4 w-4" />
        View Details
      </Button>
    </div>
  );

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          Company Values & Credos
        </h1>
        <p className="text-muted-foreground">
          Explore our company values and credos that guide our work and culture.
        </p>
      </div>

      {/* Filter Tabs */}
      <div className="border-b border-border">
        <div className="flex gap-2">
          <TabButton active={filter === "all"} onClick={() => setFilter("all")}>
            All ({values.length})
          </TabButton>
          <TabButton
            active={filter === "value"}
            onClick={() => setFilter("value")}
          >
            Values ({valueItems.length})
          </TabButton>
          <TabButton
            active={filter === "credo"}
            onClick={() => setFilter("credo")}
          >
            Credos ({credoItems.length})
          </TabButton>
        </div>
      </div>

      {/* Cards Grid */}
      {loading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="p-5 border border-border rounded-lg">
              <div className="flex justify-between mb-3">
                <Skeleton className="h-5 w-32" />
                <Skeleton className="h-5 w-16" />
              </div>
              <Skeleton className="h-4 w-full mb-2" />
              <Skeleton className="h-4 w-3/4 mb-4" />
              <Skeleton className="h-9 w-full" />
            </div>
          ))}
        </div>
      ) : filteredItems.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">
          No items found.
        </div>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filteredItems.map((v) => (
            <ValueCard key={v.id} item={v} onViewDetails={() => openDetails(v)} />
          ))}
        </div>
      )}

      {/* Detail Dialog */}
      <Dialog open={detailOpen} onOpenChange={setDetailOpen}>
        <DialogContent className="max-w-md sm:max-w-lg">
          <DialogHeader>
            <div className="flex items-center gap-3">
              <DialogTitle>{detailValue?.name}</DialogTitle>
              {detailValue && (
                <Badge
                  variant="outline"
                  className={`text-xs ${
                    detailValue.type === "Credo"
                      ? "bg-green-500 text-white border-green-500"
                      : "bg-blue-500 text-white border-blue-500"
                  }`}
                >
                  {detailValue.type.toUpperCase()}
                </Badge>
              )}
            </div>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <p className="text-sm font-medium text-muted-foreground mb-1">
                Description
              </p>
              <p className="text-sm">{detailValue?.description}</p>
            </div>
            {detailValue?.examples && detailValue.examples.length > 0 && (
              <div>
                <p className="text-sm font-medium text-muted-foreground mb-2">
                  Examples
                </p>
                <ul className="list-disc list-inside space-y-1 text-sm">
                  {detailValue.examples.map((ex, i) => (
                    <li key={i}>{ex}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

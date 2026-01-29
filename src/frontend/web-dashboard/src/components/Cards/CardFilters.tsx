import React, { useEffect, useState } from "react";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { api } from "@/services/api";
import type { ValueResponse, CardFilter } from "@/services/api";
import { Search, Filter, X, User } from "lucide-react";
import { Button } from "@/components/ui/button";

interface CardFiltersProps {
  onFilterChange: (filters: Partial<CardFilter>) => void;
  showSenderFilter?: boolean;
  showRecipientFilter?: boolean;
}

export const CardFilters: React.FC<CardFiltersProps> = ({
  onFilterChange,
  showSenderFilter = true,
  showRecipientFilter = true,
}) => {
  const [values, setValues] = useState<ValueResponse[]>([]);
  const [filters, setFilters] = useState<Partial<CardFilter>>({
    search: "",
    valueIds: [],
    senderId: "",
    recipientId: "",
  });

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

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newFilters = { ...filters, search: e.target.value };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleValueChange = (val: string) => {
    const newFilters = {
      ...filters,
      valueIds: val === "all" ? [] : [val],
    };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleParticipantChange = (
    field: "senderId" | "recipientId",
    val: string,
  ) => {
    const newFilters = { ...filters, [field]: val };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const clearFilters = () => {
    const reset = {
      search: "",
      valueIds: [],
      senderId: "",
      recipientId: "",
    };
    setFilters(reset);
    onFilterChange(reset);
  };

  const hasActiveFilters =
    filters.search ||
    filters.valueIds?.length ||
    filters.senderId ||
    filters.recipientId;

  return (
    <div className="bg-slate-900/40 p-4 rounded-2xl border border-slate-800 space-y-4">
      <div className="flex flex-wrap gap-4 items-center">
        {/* Search */}
        <div className="relative flex-1 min-w-[240px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
          <Input
            placeholder="Search recognition reason..."
            value={filters.search}
            onChange={handleSearchChange}
            className="pl-10 bg-slate-950 border-slate-800 focus-visible:ring-primary shadow-inner"
          />
        </div>

        {/* Value Filter */}
        <div className="w-[200px]">
          <Select
            value={filters.valueIds?.[0] || "all"}
            onValueChange={handleValueChange}
          >
            <SelectTrigger className="bg-slate-950 border-slate-800 focus:ring-primary">
              <div className="flex items-center gap-2">
                <Filter className="h-3.5 w-3.5 text-slate-500" />
                <SelectValue placeholder="All Values" />
              </div>
            </SelectTrigger>
            <SelectContent className="bg-slate-950 border-slate-800 text-slate-200">
              <SelectItem value="all">All Values</SelectItem>
              {values.map((v) => (
                <SelectItem key={v.id} value={v.id}>
                  {v.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Participant Filters */}
        {showSenderFilter && (
          <div className="relative w-[180px]">
            <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
            <Input
              placeholder="Sender ID"
              value={filters.senderId}
              onChange={(e) =>
                handleParticipantChange("senderId", e.target.value)
              }
              className="pl-10 bg-slate-950 border-slate-800 focus-visible:ring-primary"
            />
          </div>
        )}

        {showRecipientFilter && (
          <div className="relative w-[180px]">
            <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-500" />
            <Input
              placeholder="Recipient ID"
              value={filters.recipientId}
              onChange={(e) =>
                handleParticipantChange("recipientId", e.target.value)
              }
              className="pl-10 bg-slate-950 border-slate-800 focus-visible:ring-primary"
            />
          </div>
        )}

        {/* Date Range */}
        <div className="flex items-center gap-2">
          <Input
            type="date"
            value={filters.startDate || ""}
            onChange={(e) => {
              const nf = { ...filters, startDate: e.target.value };
              setFilters(nf);
              onFilterChange(nf);
            }}
            className="w-[150px] bg-slate-950 border-slate-800 text-slate-400 text-xs"
          />
          <span className="text-slate-600">-</span>
          <Input
            type="date"
            value={filters.endDate || ""}
            onChange={(e) => {
              const nf = { ...filters, endDate: e.target.value };
              setFilters(nf);
              onFilterChange(nf);
            }}
            className="w-[150px] bg-slate-950 border-slate-800 text-slate-400 text-xs"
          />
        </div>

        {/* Clear Button */}
        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="sm"
            onClick={clearFilters}
            className="text-slate-500 hover:text-slate-200 gap-1 px-2"
          >
            <X size={14} />
            Reset
          </Button>
        )}
      </div>
    </div>
  );
};

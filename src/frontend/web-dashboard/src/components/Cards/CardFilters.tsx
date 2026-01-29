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
    <div className="flex flex-wrap gap-3 items-center">
      {/* Search */}
      <div className="relative flex-1 min-w-[240px]">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          placeholder="Search recognition reason..."
          value={filters.search}
          onChange={handleSearchChange}
          className="pl-10 bg-white border-gray-200 focus-visible:ring-indigo-500 focus-visible:border-indigo-500 rounded-lg text-gray-700 placeholder:text-gray-400"
        />
      </div>

      {/* Value Filter */}
      <div className="w-[160px]">
        <Select
          value={filters.valueIds?.[0] || "all"}
          onValueChange={handleValueChange}
        >
          <SelectTrigger className="bg-white border-gray-200 focus:ring-indigo-500 rounded-lg text-gray-700">
            <div className="flex items-center gap-2">
              <Filter className="h-3.5 w-3.5 text-gray-400" />
              <SelectValue placeholder="All Values" />
            </div>
          </SelectTrigger>
          <SelectContent className="bg-white border-gray-200 text-gray-700 rounded-lg shadow-lg">
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
        <div className="relative w-[160px]">
          <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Sender ID"
            value={filters.senderId}
            onChange={(e) =>
              handleParticipantChange("senderId", e.target.value)
            }
            className="pl-10 bg-white border-gray-200 focus-visible:ring-indigo-500 rounded-lg text-gray-700 placeholder:text-gray-400"
          />
        </div>
      )}

      {showRecipientFilter && (
        <div className="relative w-[160px]">
          <User className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Recipient ID"
            value={filters.recipientId}
            onChange={(e) =>
              handleParticipantChange("recipientId", e.target.value)
            }
            className="pl-10 bg-white border-gray-200 focus-visible:ring-indigo-500 rounded-lg text-gray-700 placeholder:text-gray-400"
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
          className="w-[140px] bg-white border-gray-200 text-gray-500 text-sm rounded-lg"
        />
        <span className="text-gray-300">-</span>
        <Input
          type="date"
          value={filters.endDate || ""}
          onChange={(e) => {
            const nf = { ...filters, endDate: e.target.value };
            setFilters(nf);
            onFilterChange(nf);
          }}
          className="w-[140px] bg-white border-gray-200 text-gray-500 text-sm rounded-lg"
        />
      </div>

      {/* Clear Button */}
      {hasActiveFilters && (
        <Button
          variant="ghost"
          size="sm"
          onClick={clearFilters}
          className="text-gray-400 hover:text-gray-600 hover:bg-gray-100 gap-1 px-2 rounded-lg"
        >
          <X size={14} />
          Reset
        </Button>
      )}
    </div>
  );
};

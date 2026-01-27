import { useState } from 'react'
import { CardFilters } from '@/types/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { getCompanyValues } from '@/services/companyValues'
import { X, Search } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

interface FeedFiltersProps {
  filters: CardFilters
  onFiltersChange: (filters: CardFilters) => void
}

const FeedFilters = ({ filters, onFiltersChange }: FeedFiltersProps) => {
  const [searchQuery, setSearchQuery] = useState(filters.q || '')
  const companyValues = getCompanyValues()

  const handleSearch = (value: string) => {
    setSearchQuery(value)
    onFiltersChange({ ...filters, q: value || undefined, page: 1 })
  }

  const toggleValue = (valueId: number) => {
    const currentValues = filters.values || []
    const newValues = currentValues.includes(valueId)
      ? currentValues.filter(id => id !== valueId)
      : [...currentValues, valueId]
    onFiltersChange({ ...filters, values: newValues.length > 0 ? newValues : undefined, page: 1 })
  }

  const clearFilters = () => {
    setSearchQuery('')
    onFiltersChange({ page: 1, pageSize: 20 })
  }

  const hasActiveFilters = !!(filters.q || (filters.values && filters.values.length > 0))

  return (
    <Card className="p-4">
      <div className="space-y-4">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search cards..."
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
            className="pl-10"
          />
        </div>

        {/* Values Filter */}
        <div>
          <label className="text-sm font-medium mb-2 block">Filter by Values/Credos</label>
          <div className="flex flex-wrap gap-2">
            {companyValues.map((value) => {
              const isSelected = filters.values?.includes(value.id)
              return (
                <Badge
                  key={value.id}
                  variant={isSelected ? 'default' : 'outline'}
                  className="cursor-pointer"
                  onClick={() => toggleValue(value.id)}
                >
                  {value.name}
                </Badge>
              )
            })}
          </div>
        </div>

        {/* Clear Filters */}
        {hasActiveFilters && (
          <Button variant="ghost" size="sm" onClick={clearFilters} className="w-full">
            <X className="mr-2 h-4 w-4" />
            Clear Filters
          </Button>
        )}
      </div>
    </Card>
  )
}

export default FeedFilters

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { CardFilters } from '@/types/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { companyValuesApi } from '@/services/companyValues'
import { employeesApi } from '@/services/employees'
import { EmployeeSummary } from '@/types/card'
import { X, Search } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

interface FeedFiltersProps {
  filters: CardFilters
  onFiltersChange: (filters: CardFilters) => void
}

const FeedFilters = ({ filters, onFiltersChange }: FeedFiltersProps) => {
  const [searchQuery, setSearchQuery] = useState(filters.q || '')
  const [senderSearchQuery, setSenderSearchQuery] = useState('')
  const [recipientSearchQuery, setRecipientSearchQuery] = useState('')
  const [senderSearchResults, setSenderSearchResults] = useState<EmployeeSummary[]>([])
  const [recipientSearchResults, setRecipientSearchResults] = useState<EmployeeSummary[]>([])
  const [showSenderResults, setShowSenderResults] = useState(false)
  const [showRecipientResults, setShowRecipientResults] = useState(false)
  
  const { data: companyValuesData } = useQuery({
    queryKey: ['company-values', 'all'],
    queryFn: () => companyValuesApi.getAll(),
  })
  
  const { data: employeesData } = useQuery({
    queryKey: ['employees', 'all'],
    queryFn: () => employeesApi.getAll(),
  })
  
  const companyValues = companyValuesData?.items || []
  const allEmployees = employeesData?.items || []

  // Calculate selected employees
  const selectedSender = filters.senderIds && filters.senderIds.length > 0
    ? allEmployees.find(emp => emp.id === filters.senderIds![0])
    : null

  const selectedRecipient = filters.recipientIds && filters.recipientIds.length > 0
    ? allEmployees.find(emp => emp.id === filters.recipientIds![0])
    : null

  // Search sender employees
  useEffect(() => {
    const search = async () => {
      // 如果已经选择了发送者，不显示搜索结果
      if (selectedSender) {
        setSenderSearchResults([])
        setShowSenderResults(false)
        return
      }
      if (senderSearchQuery.trim().length > 0) {
        const results = await employeesApi.searchEmployees(senderSearchQuery)
        setSenderSearchResults(results)
        setShowSenderResults(true)
      } else {
        setSenderSearchResults([])
        setShowSenderResults(false)
      }
    }
    const timeoutId = setTimeout(search, 300)
    return () => clearTimeout(timeoutId)
  }, [senderSearchQuery, selectedSender])

  // Search recipient employees
  useEffect(() => {
    const search = async () => {
      // 如果已经选择了接收者，不显示搜索结果
      if (selectedRecipient) {
        setRecipientSearchResults([])
        setShowRecipientResults(false)
        return
      }
      if (recipientSearchQuery.trim().length > 0) {
        const results = await employeesApi.searchEmployees(recipientSearchQuery)
        setRecipientSearchResults(results)
        setShowRecipientResults(true)
      } else {
        setRecipientSearchResults([])
        setShowRecipientResults(false)
      }
    }
    const timeoutId = setTimeout(search, 300)
    return () => clearTimeout(timeoutId)
  }, [recipientSearchQuery, selectedRecipient])

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

  const handleSenderSelect = (employee: EmployeeSummary) => {
    onFiltersChange({ ...filters, senderIds: [employee.id], page: 1 })
    setSenderSearchQuery('')
    setShowSenderResults(false)
  }

  const handleRecipientSelect = (employee: EmployeeSummary) => {
    onFiltersChange({ ...filters, recipientIds: [employee.id], page: 1 })
    setRecipientSearchQuery('')
    setShowRecipientResults(false)
  }

  const clearSender = () => {
    onFiltersChange({ ...filters, senderIds: undefined, page: 1 })
    setSenderSearchQuery('')
    setShowSenderResults(false)
  }

  const clearRecipient = () => {
    onFiltersChange({ ...filters, recipientIds: undefined, page: 1 })
    setRecipientSearchQuery('')
    setShowRecipientResults(false)
  }

  const clearFilters = () => {
    setSearchQuery('')
    setSenderSearchQuery('')
    setRecipientSearchQuery('')
    setShowSenderResults(false)
    setShowRecipientResults(false)
    onFiltersChange({ page: 1, pageSize: 20 })
  }

  const hasActiveFilters = !!(
    filters.q || 
    (filters.values && filters.values.length > 0) ||
    (filters.senderIds && filters.senderIds.length > 0) ||
    (filters.recipientIds && filters.recipientIds.length > 0)
  )

  // Update search query when filter changes from outside (only if not manually editing)
  useEffect(() => {
    if (!selectedSender && senderSearchQuery) {
      setSenderSearchQuery('')
    }
  }, [selectedSender])

  useEffect(() => {
    if (!selectedRecipient && recipientSearchQuery) {
      setRecipientSearchQuery('')
    }
  }, [selectedRecipient])

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

        {/* User Filters */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Sender Filter */}
          <div className="space-y-2">
            <label className="text-sm font-medium block">Filter by Sender</label>
            <div className="relative">
              <Input
                placeholder="Search sender..."
                value={senderSearchQuery}
                onChange={(e) => {
                  setSenderSearchQuery(e.target.value)
                  // 如果开始输入，清除已选择的发送者
                  if (selectedSender && e.target.value !== selectedSender.name) {
                    onFiltersChange({ ...filters, senderIds: undefined, page: 1 })
                  }
                }}
                onFocus={() => {
                  if (senderSearchQuery.trim().length > 0 && !selectedSender) {
                    setShowSenderResults(true)
                  }
                }}
                className="pr-8"
              />
              {selectedSender && (
                <button
                  onClick={clearSender}
                  className="absolute right-2 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
              {showSenderResults && senderSearchResults.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-popover border border-border rounded-md shadow-lg max-h-60 overflow-auto">
                  {senderSearchResults.map((emp) => (
                    <div
                      key={emp.id}
                      onClick={() => handleSenderSelect(emp)}
                      className="px-4 py-2 hover:bg-accent cursor-pointer"
                    >
                      <div className="font-medium">{emp.name}</div>
                      <div className="text-sm text-muted-foreground">{emp.department}</div>
                    </div>
                  ))}
                </div>
              )}
            </div>
            {selectedSender && (
              <Badge variant="default" className="cursor-pointer" onClick={clearSender}>
                {selectedSender.name} ({selectedSender.department})
                <X className="ml-1 h-3 w-3" />
              </Badge>
            )}
          </div>

          {/* Recipient Filter */}
          <div className="space-y-2">
            <label className="text-sm font-medium block">Filter by Recipient</label>
            <div className="relative">
              <Input
                placeholder="Search recipient..."
                value={recipientSearchQuery}
                onChange={(e) => {
                  setRecipientSearchQuery(e.target.value)
                  // 如果开始输入，清除已选择的接收者
                  if (selectedRecipient && e.target.value !== selectedRecipient.name) {
                    onFiltersChange({ ...filters, recipientIds: undefined, page: 1 })
                  }
                }}
                onFocus={() => {
                  if (recipientSearchQuery.trim().length > 0 && !selectedRecipient) {
                    setShowRecipientResults(true)
                  }
                }}
                className="pr-8"
              />
              {selectedRecipient && (
                <button
                  onClick={clearRecipient}
                  className="absolute right-2 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
              {showRecipientResults && recipientSearchResults.length > 0 && (
                <div className="absolute z-10 w-full mt-1 bg-popover border border-border rounded-md shadow-lg max-h-60 overflow-auto">
                  {recipientSearchResults.map((emp) => (
                    <div
                      key={emp.id}
                      onClick={() => handleRecipientSelect(emp)}
                      className="px-4 py-2 hover:bg-accent cursor-pointer"
                    >
                      <div className="font-medium">{emp.name}</div>
                      <div className="text-sm text-muted-foreground">{emp.department}</div>
                    </div>
                  ))}
                </div>
              )}
            </div>
            {selectedRecipient && (
              <Badge variant="default" className="cursor-pointer" onClick={clearRecipient}>
                {selectedRecipient.name} ({selectedRecipient.department})
                <X className="ml-1 h-3 w-3" />
              </Badge>
            )}
          </div>
        </div>

        {/* Values Filter */}
        <div>
          <label className="text-sm font-medium mb-2 block">Filter by Values/Credos</label>
          <div className="flex flex-wrap gap-2">
            {companyValues.map((value) => {
              const isSelected = filters.values?.includes(value.id)
              const variant = isSelected 
                ? (value.type === 'VALUE' ? 'default' : 'success')
                : 'outline'
              return (
                <Badge
                  key={value.id}
                  variant={variant}
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

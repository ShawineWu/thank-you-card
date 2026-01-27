import { useState, useEffect } from 'react'
import { EmployeeSummary } from '@/types/card'
import { employeesApi } from '@/services/employees'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import { X, Search, User } from 'lucide-react'

interface RecipientSelectorProps {
  selected: EmployeeSummary[]
  onChange: (recipients: EmployeeSummary[]) => void
  error?: string
}

const RecipientSelector = ({ selected, onChange, error }: RecipientSelectorProps) => {
  const [searchQuery, setSearchQuery] = useState('')
  const [searchResults, setSearchResults] = useState<EmployeeSummary[]>([])
  const [isSearching, setIsSearching] = useState(false)
  const [showResults, setShowResults] = useState(false)

  useEffect(() => {
    const search = async () => {
      if (searchQuery.trim().length > 0) {
        setIsSearching(true)
        const results = await employeesApi.searchEmployees(searchQuery)
        // Filter out already selected employees
        const filtered = results.filter(
          emp => !selected.find(s => s.id === emp.id)
        )
        setSearchResults(filtered)
        setShowResults(true)
        setIsSearching(false)
      } else {
        setSearchResults([])
        setShowResults(false)
      }
    }

    const timeoutId = setTimeout(search, 300)
    return () => clearTimeout(timeoutId)
  }, [searchQuery, selected])

  const handleSelect = (employee: EmployeeSummary) => {
    onChange([...selected, employee])
    setSearchQuery('')
    setShowResults(false)
  }

  const handleRemove = (id: number) => {
    onChange(selected.filter(emp => emp.id !== id))
  }

  return (
    <div className="space-y-2">
      <label className="text-sm font-medium">Recipients *</label>
      
      {/* Selected Recipients */}
      {selected.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-2">
          {selected.map((emp) => (
            <Badge key={emp.id} variant="default" className="flex items-center gap-1">
              {emp.name}
              <button
                onClick={() => handleRemove(emp.id)}
                className="ml-1 hover:bg-primary/80 rounded-full p-0.5"
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}

      {/* Search Input */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search employees by name or department..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          onFocus={() => searchQuery && setShowResults(true)}
          className="pl-10"
        />

        {/* Search Results */}
        {showResults && (
          <Card className="absolute z-10 w-full mt-1 max-h-60 overflow-y-auto">
            {isSearching ? (
              <div className="p-4 text-center text-sm text-muted-foreground">
                Searching...
              </div>
            ) : searchResults.length > 0 ? (
              <div className="p-2">
                {searchResults.map((emp) => (
                  <button
                    key={emp.id}
                    onClick={() => handleSelect(emp)}
                    className="w-full text-left p-2 hover:bg-accent rounded-md flex items-center space-x-2"
                  >
                    <User className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <div className="font-medium">{emp.name}</div>
                      <div className="text-xs text-muted-foreground">{emp.department}</div>
                    </div>
                  </button>
                ))}
              </div>
            ) : (
              <div className="p-4 text-center text-sm text-muted-foreground">
                No employees found
              </div>
            )}
          </Card>
        )}
      </div>

      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}
      <p className="text-xs text-muted-foreground">
        Select at least one recipient
      </p>
    </div>
  )
}

export default RecipientSelector

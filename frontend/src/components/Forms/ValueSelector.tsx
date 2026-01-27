import { getValues, getCredos } from '@/services/companyValues'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface ValueSelectorProps {
  selected: number[]
  onChange: (valueIds: number[]) => void
  error?: string
}

const ValueSelector = ({ selected, onChange, error }: ValueSelectorProps) => {
  const values = getValues()
  const credos = getCredos()
  const maxSelection = 3

  const handleToggle = (valueId: number) => {
    if (selected.includes(valueId)) {
      onChange(selected.filter(id => id !== valueId))
    } else {
      if (selected.length < maxSelection) {
        onChange([...selected, valueId])
      }
    }
  }

  return (
    <div className="space-y-4">
      <div>
        <label className="text-sm font-medium">Company Values & Credos *</label>
        <p className="text-xs text-muted-foreground mb-2">
          Select 1-3 values/credos that this recognition aligns with ({selected.length}/{maxSelection} selected)
        </p>
      </div>

      {/* Values */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Values</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {values.map((value) => {
              const isSelected = selected.includes(value.id)
              const isDisabled = !isSelected && selected.length >= maxSelection
              return (
                <Badge
                  key={value.id}
                  variant={isSelected ? 'default' : 'outline'}
                  className={isDisabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
                  onClick={() => !isDisabled && handleToggle(value.id)}
                >
                  {value.name}
                </Badge>
              )
            })}
          </div>
        </CardContent>
      </Card>

      {/* Credos */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Credos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            {credos.map((credo) => {
              const isSelected = selected.includes(credo.id)
              const isDisabled = !isSelected && selected.length >= maxSelection
              return (
                <Badge
                  key={credo.id}
                  variant={isSelected ? 'secondary' : 'outline'}
                  className={isDisabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
                  onClick={() => !isDisabled && handleToggle(credo.id)}
                >
                  {credo.name}
                </Badge>
              )
            })}
          </div>
        </CardContent>
      </Card>

      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}
    </div>
  )
}

export default ValueSelector

import { Badge } from '@/components/ui/badge'
import { CompanyValueSummary } from '@/types/card'

interface ValueBadgeProps {
  value: CompanyValueSummary
}

const ValueBadge = ({ value }: ValueBadgeProps) => {
  const variant = value.type === 'VALUE' ? 'default' : 'secondary'
  
  return (
    <Badge variant={variant} className="text-xs">
      {value.name}
    </Badge>
  )
}

export default ValueBadge

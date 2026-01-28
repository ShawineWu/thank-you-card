import { Badge } from '@/components/ui/badge'
import { CompanyValueSummary } from '@/types/card'

interface ValueBadgeProps {
  value: CompanyValueSummary
  onClick?: (value: CompanyValueSummary) => void
}

const ValueBadge = ({ value, onClick }: ValueBadgeProps) => {
  const variant = value.type === 'VALUE' ? 'default' : 'success'
  
  return (
    <Badge 
      variant={variant} 
      className={`text-xs ${onClick ? 'cursor-pointer hover:opacity-80 transition-opacity' : ''}`}
      onClick={() => onClick?.(value)}
    >
      {value.name}
    </Badge>
  )
}

export default ValueBadge

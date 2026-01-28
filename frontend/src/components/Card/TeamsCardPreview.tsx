import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { EmployeeSummary, CompanyValueSummary } from '@/types/card'
import { useQuery } from '@tanstack/react-query'
import { companyValuesApi } from '@/services/companyValues'
import { Heart } from 'lucide-react'

interface TeamsCardPreviewProps {
  sender: {
    name: string
    department: string
  }
  recipients: EmployeeSummary[]
  reason: string
  selectedValueIds: number[]
}

const TeamsCardPreview = ({ sender, recipients, reason, selectedValueIds }: TeamsCardPreviewProps) => {
  // Fetch all company values to get names for selected IDs
  const { data: allValuesData } = useQuery({
    queryKey: ['company-values', 'all'],
    queryFn: () => companyValuesApi.getAll(),
  })

  const allValues = allValuesData?.items || []
  const selectedValues = allValues.filter(v => selectedValueIds.includes(v.id))

  // Format recipients text
  const recipientsText = recipients.length === 0
    ? '...'
    : recipients.length === 1
    ? recipients[0].name
    : recipients.length === 2
    ? `${recipients[0].name} and ${recipients[1].name}`
    : `${recipients[0].name} and ${recipients.length - 1} others`

  return (
    <Card className="border border-border shadow-lg overflow-hidden bg-white">
      <div className="p-6 space-y-5">
        {/* Title section - centered */}
        <div className="text-center space-y-2">
          <h2 className="text-3xl font-bold text-foreground">✨ Thank You Card</h2>
          <p className="text-base text-muted-foreground">
            Some help may be routine, but it's worth being sincerely thanked ❤️
          </p>
        </div>

        {/* To section - small font at the beginning */}
        <div className="space-y-2">
          <p className="text-sm text-foreground">
            <span className="font-semibold">To:</span> {recipientsText}
          </p>
        </div>

        {/* Reason/Message section - before values */}
        <div className="space-y-2">
          <div className="bg-muted/20 rounded-lg p-4 min-h-[120px]">
            <p className="text-base text-foreground whitespace-pre-wrap leading-relaxed">
              {reason || (
                <span className="text-muted-foreground italic">Your message will appear here...</span>
              )}
            </p>
          </div>
        </div>

        {/* Values section - as tags without "Values:" prefix */}
        {selectedValues.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {selectedValues.map((value) => (
              <Badge
                key={value.id}
                variant={value.type === 'VALUE' ? 'default' : 'success'}
                className="text-xs px-3 py-1"
              >
                {value.name}
              </Badge>
            ))}
          </div>
        )}

        {/* Closing sentiment */}
        <div className="pt-4">
          <p className="text-center text-base font-semibold text-foreground">
            🌈 Working with you is a lucky thing!
          </p>
        </div>

        {/* Sender signature at the end */}
        <div className="pt-2">
          <p className="text-sm text-muted-foreground">
            {sender.name || '...'} ({sender.department || '...'})
          </p>
        </div>
      </div>
    </Card>
  )
}

export default TeamsCardPreview

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { CompanyValueDetail } from '@/types/card'

interface ValueDetailDialogProps {
  value: CompanyValueDetail | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

const ValueDetailDialog = ({ value, open, onOpenChange }: ValueDetailDialogProps) => {
  if (!value) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <div className="flex items-start justify-between pr-8">
            <div className="flex-1">
              <DialogTitle className="text-2xl mb-2">{value.name}</DialogTitle>
              <Badge variant={value.type === 'VALUE' ? 'default' : 'success'}>
                {value.type}
              </Badge>
            </div>
          </div>
        </DialogHeader>
        <div className="mt-4 space-y-6">
          <div>
            <h3 className="text-lg font-semibold mb-2">Description</h3>
            <p className="text-muted-foreground whitespace-pre-wrap leading-relaxed">
              {value.description}
            </p>
          </div>
          
          {value.examples && value.examples.length > 0 && (
            <div>
              <h3 className="text-lg font-semibold mb-3">Examples</h3>
              <ul className="space-y-2">
                {value.examples.map((example, index) => (
                  <li key={index} className="flex items-start gap-3">
                    <span className="text-primary font-semibold mt-0.5">•</span>
                    <span className="text-muted-foreground leading-relaxed flex-1">
                      {example}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default ValueDetailDialog

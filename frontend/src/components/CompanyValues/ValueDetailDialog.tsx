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
          <div className="flex items-center justify-between">
            <DialogTitle className="text-2xl">{value.name}</DialogTitle>
            <Badge variant={value.type === 'VALUE' ? 'default' : 'success'}>
              {value.type}
            </Badge>
          </div>
          <DialogDescription className="text-sm text-muted-foreground mt-2">
            Code: {value.code}
          </DialogDescription>
        </DialogHeader>
        <div className="mt-4">
          <h3 className="text-lg font-semibold mb-2">Description</h3>
          <p className="text-muted-foreground whitespace-pre-wrap leading-relaxed">
            {value.description}
          </p>
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default ValueDetailDialog

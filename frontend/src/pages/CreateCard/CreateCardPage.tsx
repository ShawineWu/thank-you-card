import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { cardsApi } from '@/services/cards'
import { CreateCardRequest, EmployeeSummary } from '@/types/card'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import RecipientSelector from '@/components/Forms/RecipientSelector'
import ValueSelector from '@/components/Forms/ValueSelector'
import TeamsCardPreview from '@/components/Card/TeamsCardPreview'
import { Loader2, Send, Eye } from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAuth } from '@/contexts/AuthContext'

const MAX_REASON_LENGTH = 2000

const CreateCardPage = () => {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuth()

  const [recipients, setRecipients] = useState<EmployeeSummary[]>([])
  const [valueIds, setValueIds] = useState<number[]>([])
  const [reason, setReason] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [previewOpen, setPreviewOpen] = useState(false)

  const createCardMutation = useMutation({
    mutationFn: (data: CreateCardRequest) => cardsApi.createCard(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cards'] })
      navigate('/')
    },
    onError: (error: any) => {
      setErrors({ submit: error.response?.data?.error || 'Failed to create card' })
    },
  })

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (recipients.length === 0) {
      newErrors.recipients = 'Please select at least one recipient'
    }

    if (valueIds.length === 0) {
      newErrors.values = 'Please select at least one value or credo'
    } else if (valueIds.length > 3) {
      newErrors.values = 'Please select no more than 3 values/credos'
    }

    if (!reason.trim()) {
      newErrors.reason = 'Please enter a reason for recognition'
    } else if (reason.length > MAX_REASON_LENGTH) {
      newErrors.reason = `Reason must be no more than ${MAX_REASON_LENGTH} characters`
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) {
      return
    }

    const data: CreateCardRequest = {
      recipientIds: recipients.map(r => r.id),
      valueIds,
      reason: reason.trim(),
    }

    createCardMutation.mutate(data)
  }

  const remainingChars = MAX_REASON_LENGTH - reason.length

  // Get sender info from auth context
  const senderName = user?.employee?.name || user?.username || '...'
  const senderDepartment = user?.employee?.department || '...'

  return (
    <div className="max-w-3xl mx-auto">
      <div className="mb-6">
        <h1 className="text-3xl font-bold mb-2">Create Thank You Card</h1>
        <p className="text-muted-foreground">
          Recognize a colleague for their great work
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle>Card Details</CardTitle>
            <CardDescription>
              Fill in the details to create a thank you card
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Recipients */}
            <RecipientSelector
              selected={recipients}
              onChange={setRecipients}
              error={errors.recipients}
            />

            {/* Values */}
            <ValueSelector
              selected={valueIds}
              onChange={setValueIds}
              error={errors.values}
            />

            {/* Reason */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                Reason for Recognition *
              </label>
              <textarea
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder="Describe why you're recognizing this person..."
                className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                maxLength={MAX_REASON_LENGTH}
              />
              <div className="flex justify-between items-center">
                {errors.reason && (
                  <p className="text-sm text-destructive">{errors.reason}</p>
                )}
                <p className={cn(
                  "text-xs ml-auto",
                  remainingChars < 50 ? "text-destructive" : "text-muted-foreground"
                )}>
                  {remainingChars} characters remaining
                </p>
              </div>
            </div>

            {/* Submit Error */}
            {errors.submit && (
              <div className="p-3 rounded-md bg-destructive/10 text-destructive text-sm">
                {errors.submit}
              </div>
            )}

            {/* Actions */}
            <div className="flex justify-end space-x-4 pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => navigate('/')}
                disabled={createCardMutation.isPending}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => setPreviewOpen(true)}
                disabled={createCardMutation.isPending}
              >
                <Eye className="mr-2 h-4 w-4" />
                Preview
              </Button>
              <Button
                type="submit"
                disabled={createCardMutation.isPending}
              >
                {createCardMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Creating...
                  </>
                ) : (
                  <>
                    <Send className="mr-2 h-4 w-4" />
                    Send Card
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>
      </form>

      {/* Preview Dialog */}
      <Dialog open={previewOpen} onOpenChange={setPreviewOpen}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Teams Card Preview</DialogTitle>
            <DialogDescription>
              This is how your card will appear in Teams
            </DialogDescription>
          </DialogHeader>
          <div className="mt-4">
            <TeamsCardPreview
              sender={{
                name: senderName,
                department: senderDepartment,
              }}
              recipients={recipients}
              reason={reason}
              selectedValueIds={valueIds}
            />
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default CreateCardPage

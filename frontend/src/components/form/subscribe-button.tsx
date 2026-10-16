import { Loader2 } from 'lucide-react'
import { Button } from '../ui/button'
import { useFormContext } from './form-context'

type SubscribeButtonProps = {
  label: string
  loadingText?: string
}

export const SubscribeButton = ({ label, loadingText }: SubscribeButtonProps) => {
  const form = useFormContext()
  return (
    <form.Subscribe selector={({ canSubmit, isSubmitting }) => ({ canSubmit, isSubmitting })}>
      {({ canSubmit, isSubmitting }) => (
        <Button type="submit" disabled={!canSubmit || isSubmitting}>
          {isSubmitting && <Loader2 className="animate-spin" aria-hidden="true" />}
          {isSubmitting && loadingText ? loadingText : label}
        </Button>
      )}
    </form.Subscribe>
  )
}

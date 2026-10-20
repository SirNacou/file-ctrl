import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { useState, type ReactElement } from 'react'
import LoaderCircleIcon from '~icons/lucide/loader-circle'

type ConfirmDeleteProps = {
  onConfirm: () => void | Promise<void>
  trigger?: ReactElement
  open?: boolean
  onOpenChanged?: (open: boolean) => void
  title?: string
  description?: string
}

export function ConfirmDeleteDialog({
  onConfirm,
  open,
  onOpenChanged,
  trigger,
  title = 'Are you sure?',
  description = 'This action cannot be undone.',
}: ConfirmDeleteProps) {
  const [internalOpen, setInternalOpen] = useState(false)
  const [loading, setLoading] = useState(false)

  const isControlled = open !== undefined
  const actualOpen = isControlled ? open : internalOpen
  const setActualOpen = isControlled ? (onOpenChanged ?? (() => {})) : setInternalOpen

  async function handleConfirm() {
    setLoading(true)
    try {
      await onConfirm()
      setActualOpen(false)
    } catch (err) {
      console.log(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <AlertDialog open={actualOpen} onOpenChange={setActualOpen}>
      {trigger && <AlertDialogTrigger render={trigger} />}
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription>{description}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction variant={'destructive'} onClick={handleConfirm} disabled={loading}>
            {loading && <LoaderCircleIcon className="mr-2 h-4 w-4 animate-spin" />}
            {loading ? 'Deleting...' : 'Delete'}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

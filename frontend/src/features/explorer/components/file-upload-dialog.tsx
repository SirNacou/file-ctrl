import { type UploadItem as UploadItemData } from '@uploadcn/core'
import { useUploadContext, useUploadProgress, useUploadSelector } from '@uploadcn/react'
import { atom, useAtom } from 'jotai'
import UploadIcon from '~icons/lucide/upload'

import { Button } from '#/components/ui/button.tsx'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog.tsx'
import {
  UploadDropzone,
  UploadDropzoneDescription,
  UploadDropzoneHeader,
  UploadDropzoneMedia,
  UploadDropzoneTitle,
  UploadModalFooter,
  UploadSummary,
} from '#/components/ui/upload.tsx'
import { UploadQueue } from '#/components/upload-queue.tsx'

export const uploadDialogOpenAtom = atom(false)

type Props = {
  onOpenChange?: (open: boolean) => void
}
export const FileUploadDialog = ({ onOpenChange }: Props) => {
  const [open, setOpen] = useAtom(uploadDialogOpenAtom)

  function handleOnOpenChange(open: boolean) {
    setOpen(open)
    onOpenChange?.(open)
  }

  return (
    <Dialog open={open} onOpenChange={handleOnOpenChange}>
      <DialogTrigger
        render={
          <Button>
            <UploadIcon />
            Upload files
          </Button>
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upload files</DialogTitle>
          <DialogDescription>Add files, review them, then upload.</DialogDescription>
        </DialogHeader>
        <DialogBody />
        <DialogActions
          onClose={() => handleOnOpenChange(false)}
          onDone={() => handleOnOpenChange(false)}
        />
      </DialogContent>
    </Dialog>
  )
}

function DialogActions({
  onClose,
  onDone,
}: {
  onClose: () => void
  onDone?: (items: UploadItemData[]) => void
}) {
  const { uploader } = useUploadContext()
  const summary = useUploadProgress(uploader)
  const idle = useUploadSelector(
    uploader,
    (state) => state.items.filter((item) => item.status === 'idle').length,
  )
  const finished = summary.total > 0 && idle === 0 && !summary.isUploading

  const done = () => {
    onDone?.(uploader.getState().items.filter((item) => item.status === 'success'))
    uploader.clearCompleted()
    onClose()
  }

  return (
    <UploadModalFooter className="items-center sm:justify-between">
      <UploadSummary className="max-sm:order-last" />
      <div className="flex gap-2 sm:ms-auto">
        {finished ? (
          <Button onClick={done}>Done</Button>
        ) : (
          <>
            <Button variant="outline" onClick={onClose}>
              {summary.isUploading ? 'Hide' : 'Cancel'}
            </Button>
            <Button disabled={idle === 0 || summary.isUploading} onClick={() => uploader.start()}>
              {summary.isUploading
                ? `Uploading… ${Math.round(summary.percent)}%`
                : idle > 0
                  ? `Upload ${idle} file${idle === 1 ? '' : 's'}`
                  : 'Upload'}
            </Button>
          </>
        )}
      </div>
    </UploadModalFooter>
  )
}

function DialogBody() {
  return (
    <div className="flex min-w-0 flex-col gap-3">
      <UploadDropzone size="sm" className="py-8">
        <UploadDropzoneHeader>
          <UploadDropzoneMedia variant="icon">
            <UploadIcon strokeWidth={2} />
          </UploadDropzoneMedia>
          <UploadDropzoneTitle>Drop files or click to browse</UploadDropzoneTitle>
          <UploadDropzoneDescription>Nothing uploads until you confirm</UploadDropzoneDescription>
        </UploadDropzoneHeader>
      </UploadDropzone>
      <UploadQueue size="sm" variant="muted" className="max-h-64 overflow-y-auto" />
    </div>
  )
}

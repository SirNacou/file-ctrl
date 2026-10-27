import { useMutation, useQueryClient } from '@tanstack/react-query'

import {
  renameFileMutation,
  storageListDirectoryQueryKey,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import type { DirectoryEntry } from '#/client/gen/index.ts'
import { zRenameFileBody, type RenameFileBodyZodInput } from '#/client/gen/zod.gen.ts'
import { Button } from '#/components/ui/button.tsx'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog.tsx'
import { Field, FieldGroup } from '#/components/ui/field.tsx'
import { useAppForm } from '#/hooks/use-form.ts'
import { notify } from '#/lib/notify.ts'

type Props = {
  storagePath: string
  subpath?: string
  entry: DirectoryEntry
  open: boolean
  openChanged: (isOpen: boolean) => void
}
export const RenameFileDialog = ({ storagePath, subpath, entry, open, openChanged }: Props) => {
  const defaultValues: RenameFileBodyZodInput = {
    old_name: entry.name,
    parent_path: subpath || '',
    new_name: entry.name,
  }

  const form = useAppForm({
    defaultValues,
    validators: {
      onBlur: zRenameFileBody,
      onSubmit: zRenameFileBody,
    },
    onSubmit: async ({ value }) => {
      await renameFile.mutateAsync({ path: { storage_path: storagePath }, body: value })
    },
  })

  const queryClient = useQueryClient()

  const renameFile = useMutation({
    ...renameFileMutation(),
    onError: (e) => {
      notify.error(e.title, e.errors?.at(0)?.message)
    },
    onSuccess: async () => {
      notify.success('Create successfully!')
      form.reset()
      await queryClient.invalidateQueries({
        queryKey: storageListDirectoryQueryKey({
          path: { storage_path: storagePath },
          query: {
            subpath: subpath,
          },
        }),
      })
      openChanged(false)
    },
  })

  function handleCancel() {
    form.reset()
    openChanged(false)
  }

  return (
    <Dialog open={open} onOpenChange={openChanged}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Rename</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={async (e) => {
            e.stopPropagation()
            e.preventDefault()
            await form.handleSubmit()
          }}
          className="flex flex-col gap-4"
        >
          <FieldGroup>
            <form.AppField name="new_name" children={(field) => <field.TextField label="Name" />} />
            <DialogFooter>
              <form.AppForm>
                <Field orientation={'horizontal'}>
                  <Button type="button" variant={'outline'} onClick={handleCancel}>
                    Cancel
                  </Button>
                  <form.SubscribeButton label="Confirm" />
                </Field>
              </form.AppForm>
            </DialogFooter>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  )
}

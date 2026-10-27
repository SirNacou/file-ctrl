import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import PlusIcon from '~icons/lucide/plus'

import {
  createStorageFolderMutation,
  storageListDirectoryQueryKey,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import {
  zCreateStorageFolderBody,
  type CreateStorageFolderBodyZodInput,
} from '#/client/gen/zod.gen.ts'
import { Button } from '#/components/ui/button.tsx'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog.tsx'
import { Field, FieldGroup } from '#/components/ui/field.tsx'
import { useAppForm } from '#/hooks/use-form.ts'
import { notify } from '#/lib/notify.ts'

type CreateFolderDialogProps = {
  storagePath: string
  subpath?: string
}

export function CreateFolderDialog({ subpath, storagePath }: CreateFolderDialogProps) {
  const [open, setOpen] = useState(false)

  const defaultValues: CreateStorageFolderBodyZodInput = {
    name: '',
    parent_path: subpath,
  }

  const form = useAppForm({
    defaultValues,
    validators: {
      onBlur: zCreateStorageFolderBody,
      onSubmit: zCreateStorageFolderBody,
    },
    onSubmit: async ({ value }) => {
      await createFolder.mutateAsync({
        body: { name: value.name, parent_path: value.parent_path ?? '' },
        path: {
          storage_path: storagePath,
        },
      })
    },
  })

  const queryClient = useQueryClient()

  const createFolder = useMutation({
    ...createStorageFolderMutation(),
    onSuccess: async () => {
      form.reset()
      await queryClient.invalidateQueries({
        queryKey: storageListDirectoryQueryKey({
          path: { storage_path: storagePath },
          query: { subpath: subpath },
        }),
      })
      setOpen(false)
    },
    onError: (e) => {
      notify.error(e.title, e.errors?.at(0)?.message)
    },
  })

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger
        render={
          <Button>
            <PlusIcon />
            Create folder
          </Button>
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Folder</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={async (e) => {
            e.preventDefault()
            e.stopPropagation()
            await form.handleSubmit()
          }}
        >
          <FieldGroup>
            <form.AppField
              name="name"
              children={(field) => <field.TextField label="Folder name" />}
            />

            <form.AppForm>
              <Field orientation={'horizontal'}>
                <Button type="button" variant={'outline'} onClick={() => setOpen(false)}>
                  Cancel
                </Button>
                <form.SubscribeButton label="Create" />
              </Field>
            </form.AppForm>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  )
}

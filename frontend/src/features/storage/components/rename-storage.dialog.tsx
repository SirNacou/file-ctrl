import {
  getApiStorageQueryKey,
  storageRenameMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import {
  zRenameStorageRequestBody,
  type RenameStorageRequestBodyZodInput,
} from '#/client/gen/zod.gen.ts'
import { fieldContext, formContext } from '#/components/form/form-context.ts'
import { SubscribeButton } from '#/components/form/subscribe-button.tsx'
import { TextField } from '#/components/form/text-field.tsx'
import { Button } from '#/components/ui/button.tsx'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog.tsx'
import { Field, FieldGroup } from '#/components/ui/field.tsx'
import { toast } from '#/components/ui/toast.tsx'
import { createFormHook } from '@tanstack/react-form'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import EditIcon from '~icons/lucide/edit'

const { useAppForm } = createFormHook({
  fieldContext,
  formContext,
  fieldComponents: {
    TextField,
  },
  formComponents: {
    SubscribeButton,
  },
})

type RenameStorageDialogProps = {
  subpath: string
  cur_name: string
}

export function RenameStorageDialog({ subpath, cur_name }: RenameStorageDialogProps) {
  const [open, setOpen] = useState(false)

  const defaultValues: RenameStorageRequestBodyZodInput = {
    new_name: cur_name,
  }

  const form = useAppForm({
    defaultValues,
    validators: {
      onBlur: zRenameStorageRequestBody,
      onSubmit: zRenameStorageRequestBody,
    },
    onSubmit: async ({ value }) => {
      await renameStorage.mutateAsync({
        path: {
          subpath,
        },
        body: {
          new_name: value.new_name,
        },
      })
    },
  })

  const queryClient = useQueryClient()

  const renameStorage = useMutation({
    ...storageRenameMutation(),
    onError: (error) => {
      toast.add({ type: 'error', title: error.title, description: error.errors?.at(0)?.message })
    },
    onSuccess: async () => {
      toast.add({ type: 'success', title: 'Rename successfully' })
      await queryClient.invalidateQueries({
        queryKey: getApiStorageQueryKey(),
      })
      form.reset()
      setOpen(false)
    },
  })

  function handleCancel() {
    form.reset()
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger
        render={
          <Button>
            <EditIcon strokeWidth={2} />
            <span>Storage</span>
          </Button>
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="font-bold text-xl">Add Storage</DialogTitle>
        </DialogHeader>

        <form
          onSubmit={async (event) => {
            event.preventDefault()
            event.stopPropagation()
            await form.handleSubmit()
          }}
        >
          <form.AppForm>
            <FieldGroup>
              <form.AppField
                name="new_name"
                children={(field) => <field.TextField label="Storage name" />}
              />
              <Field className="flex justify-end" orientation={'horizontal'}>
                <Button type="button" variant={'outline'} onClick={handleCancel}>
                  Cancel
                </Button>
                <form.SubscribeButton label="Submit" loadingText="Loading..." />
              </Field>
            </FieldGroup>
          </form.AppForm>
        </form>
      </DialogContent>
    </Dialog>
  )
}

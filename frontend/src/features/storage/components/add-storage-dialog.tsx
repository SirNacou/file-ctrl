import { createFormHook, formOptions } from '@tanstack/react-form'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import LucidePlusIcon from '~icons/lucide/plus'

import {
  listStorageSourcesQueryKey,
  listUnregisteredSourcesOptions,
  listUnregisteredSourcesQueryKey,
  registerStorageMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { type RegisterStorageSourceRequestBody } from '#/client/gen/index.ts'
import { zRegisterStorageSourceRequestBody } from '#/client/gen/zod.gen.ts'
import { fieldContext, formContext } from '#/components/form/form-context.ts'
import { SelectField, type ItemType } from '#/components/form/select-field.tsx'
import { SubscribeButton } from '#/components/form/subscribe-button.tsx'
import { TextField } from '#/components/form/text-field.tsx'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog.tsx'
import { Field, FieldGroup } from '#/components/ui/field.tsx'
import { toast } from '#/components/ui/toast.tsx'
import { Button } from '@/components/ui/button'

const defaultValues: RegisterStorageSourceRequestBody = {
  name: '',
  path: '',
}

const formOpts = formOptions({
  defaultValues,
  validators: {
    onBlur: zRegisterStorageSourceRequestBody,
    onSubmit: zRegisterStorageSourceRequestBody,
  },
})

const { useAppForm } = createFormHook({
  fieldContext,
  formContext,
  fieldComponents: {
    TextField,
    SelectField,
  },
  formComponents: {
    SubscribeButton,
  },
})

export const AddStorageDialog = () => {
  const [open, setOpen] = useState(false)

  const form = useAppForm({
    ...formOpts,
    onSubmit: async ({ value }) => {
      await registerStorageSource.mutateAsync({
        body: value,
      })
    },
  })

  const queryClient = useQueryClient()

  const registerStorageSource = useMutation({
    ...registerStorageMutation(),
    onError: (e) => {
      toast.add({
        type: 'error',
        title: e.title,
        description: e.errors?.at(0)?.message,
      })
    },
    onSuccess: async () => {
      toast.add({
        type: 'success',
        title: 'Success',
        description: 'Add storage successfully!',
      })

      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: listStorageSourcesQueryKey(),
        }),
        queryClient.invalidateQueries({
          queryKey: listUnregisteredSourcesQueryKey(),
        }),
      ])

      form.reset()
      setOpen(false)
    },
  })

  const { data } = useQuery(listUnregisteredSourcesOptions())

  const selectItems = useMemo(
    () =>
      data?.dirs?.map((dir): ItemType<string> => ({
        label: dir,
        value: dir,
      })),
    [data],
  )

  const handleCancel = () => {
    form.reset()
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger
        render={
          <Button>
            <LucidePlusIcon strokeWidth={2} />
            <span>Storage</span>
          </Button>
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="text-xl font-bold">Add Storage</DialogTitle>
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
                name="name"
                children={(field) => <field.TextField label="Storage name" />}
              />
              <form.AppField
                name="path"
                children={(field) => (
                  <field.SelectField label="Directory" items={selectItems} disabled={!data} />
                )}
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

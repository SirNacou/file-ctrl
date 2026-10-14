import {
  getApiStorageListUnregisteredSourcesOptions,
  getApiStorageQueryKey,
  postApiStorageMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { type RegisterStorageSourceRequestBody } from '#/client/gen/index.ts'
import { zRegisterStorageSourceRequestBody } from '#/client/gen/zod.gen.ts'
import { fieldContext, formContext } from '#/components/form/form-context.ts'
import { SelectField, type ItemType } from '#/components/form/select-field.tsx'
import { TextField } from '#/components/form/text-field.tsx'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '#/components/ui/dialog.tsx'
import { toast } from '#/components/ui/toast.tsx'
import { Button } from '@/components/ui/button'
import { createFormHook, formOptions } from '@tanstack/react-form'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useMemo } from 'react'
import LucidePlusIcon from '~icons/lucide/plus'

const defaultValues: RegisterStorageSourceRequestBody = {
  name: '',
  sub_path: '',
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
  formComponents: {},
})

export const AddStorageDialog = () => {
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
    ...postApiStorageMutation(),
    onError: (e) => {
      toast.add({
        title: e.title,
        description: e.errors?.at(0)?.message ?? '',
      })
    },
    onSuccess: async () => {
      toast.add({
        title: 'Success',
        description: 'Add storage successfully!',
      })

      await queryClient.invalidateQueries({
        queryKey: getApiStorageQueryKey(),
      })
    },
  })

  const { data } = useQuery(getApiStorageListUnregisteredSourcesOptions())

  const selectItems = useMemo(
    () =>
      data?.dirs?.map((dir): ItemType<string> => ({
        label: dir,
        value: dir,
      })),
    [data],
  )

  return (
    <Dialog>
      <DialogTrigger
        render={
          <Button className="w-full">
            <LucidePlusIcon className="size-7" />
            <span className="font-bold text-lg uppercase">Storage</span>
          </Button>
        }
      />
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="font-bold text-xl">Add Storage</DialogTitle>
        </DialogHeader>

        <form.AppForm>
          <form.AppField
            name="name"
            children={(field) => <field.TextField label="Storage name" />}
          />
          <form.AppField
            name="sub_path"
            children={(field) => (
              <field.SelectField label="Directory" items={selectItems} disabled={!data} />
            )}
          />
        </form.AppForm>
      </DialogContent>
    </Dialog>
  )
}

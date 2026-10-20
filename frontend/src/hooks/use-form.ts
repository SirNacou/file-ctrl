import { fieldContext, formContext } from '#/components/form/form-context.ts'
import { SelectField } from '#/components/form/select-field.tsx'
import { SubscribeButton } from '#/components/form/subscribe-button.tsx'
import { TextField } from '#/components/form/text-field.tsx'
import { createFormHook } from '@tanstack/react-form'

export const { useAppForm, extendForm } = createFormHook({
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

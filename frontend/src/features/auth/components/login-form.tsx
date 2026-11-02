import { createFormHook } from '@tanstack/react-form'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { AlertCircle, KeyRound } from 'lucide-react'

import { authGetStatusQueryKey, authLoginMutation } from '#/client/gen/@tanstack/react-query.gen.ts'
import { zAuthLoginBody, type AuthLoginBodyZodInput } from '#/client/gen/zod.gen.ts'
import { fieldContext, formContext } from '#/components/form/form-context.ts'
import { InputField } from '#/components/ui/shuip/tanstack-form/input-field.tsx'
import { PasswordField } from '#/components/ui/shuip/tanstack-form/password-field.tsx'
import { SubmitButton } from '#/components/ui/shuip/tanstack-form/submit-button.tsx'

type LoginFormProps = {
  onSuccess: () => void
}

const { useAppForm } = createFormHook({
  fieldContext,
  formContext,
  fieldComponents: {
    InputField,
    PasswordField,
  },
  formComponents: {
    SubmitButton,
  },
})

const defaultValues: AuthLoginBodyZodInput = {
  username: '',
  password: '',
}

export function LoginForm({ onSuccess }: LoginFormProps) {
  const queryClient = useQueryClient()
  const loginMutation = useMutation({
    ...authLoginMutation(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: authGetStatusQueryKey(),
      })
    },
  })

  const form = useAppForm({
    defaultValues,
    validators: {
      onSubmit: zAuthLoginBody,
      onChange: zAuthLoginBody,
    },
    onSubmit: async ({ value }) => {
      await loginMutation.mutateAsync({
        body: value,
      })
      onSuccess()
    },
  })

  const displayError = loginMutation.error
    ? (loginMutation.error as { detail?: string; message?: string }).detail ||
      (loginMutation.error as { message?: string }).message ||
      'Invalid username or password.'
    : null

  return (
    <div className="bg-background flex min-h-screen items-center justify-center p-4">
      <div className="bg-card text-card-foreground w-full max-w-md space-y-6 rounded-xl border p-6 shadow-sm sm:p-8">
        {/* Header */}
        <div className="flex flex-col items-center space-y-2 text-center">
          <div className="bg-primary/10 text-primary flex h-12 w-12 items-center justify-center rounded-full">
            <KeyRound className="h-6 w-6" />
          </div>
          <h1 className="text-2xl font-semibold tracking-tight">Sign in to File Ctrl</h1>
          <p className="text-muted-foreground text-sm">
            Enter your credentials to access your files.
          </p>
        </div>

        {/* Error Alert Banner */}
        {displayError && (
          <div
            role="alert"
            className="border-destructive/30 bg-destructive/10 text-destructive flex items-start gap-2.5 rounded-md border p-3 text-sm"
          >
            <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
            <span>{displayError}</span>
          </div>
        )}

        {/* Form */}
        <form
          className="flex flex-col items-stretch gap-4"
          noValidate
          onSubmit={async (e) => {
            e.preventDefault()
            e.stopPropagation()
            await form.handleSubmit()
          }}
        >
          <form.AppField name="username" listeners={{ onChange: () => loginMutation.reset() }}>
            {(field) => (
              <field.InputField props={{ placeholder: 'Enter Username' }} label="Username" />
            )}
          </form.AppField>

          <form.AppField name="password" listeners={{ onChange: () => loginMutation.reset() }}>
            {(field) => <field.PasswordField label="Password" />}
          </form.AppField>

          <form.AppForm>
            <form.SubmitButton pendingLabel="Signing in...">Sign In</form.SubmitButton>
          </form.AppForm>
        </form>
      </div>
    </div>
  )
}

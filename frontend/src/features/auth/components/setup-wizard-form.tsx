import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import AlertCircle from '~icons/lucide/circle-alert'
import CheckCircle2 from '~icons/lucide/circle-check-big'
import Eye from '~icons/lucide/eye'
import EyeOff from '~icons/lucide/eye-off'
import Loader2 from '~icons/lucide/loader-circle'
import ShieldCheck from '~icons/lucide/shield-check'

import {
  authGetStatusQueryKey,
  authSetupAdminMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

import { OidcLoginButton, oidcLoginUrl } from './oidc-login-button'

interface SetupWizardFormProps {
  onSuccess: () => void
  oidcEnabled: boolean
}

export function SetupWizardForm({ onSuccess, oidcEnabled }: SetupWizardFormProps) {
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [validationError, setValidationError] = useState<string | null>(null)

  const queryClient = useQueryClient()
  const enabledOidcLoginUrl = oidcEnabled ? oidcLoginUrl : undefined
  const setupMutation = useMutation({
    ...authSetupAdminMutation(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: authGetStatusQueryKey(),
      })
    },
  })

  const isUsernameValid = username.trim().length >= 3
  const isPasswordLengthValid = password.length >= 8
  const doPasswordsMatch = password.length > 0 && password === confirmPassword

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault()
    setValidationError(null)

    const trimmedUsername = username.trim()
    if (trimmedUsername.length < 3) {
      setValidationError('Username must be at least 3 characters.')
      return
    }

    if (password.length < 8) {
      setValidationError('Password must be at least 8 characters.')
      return
    }

    if (password !== confirmPassword) {
      setValidationError('Passwords do not match.')
      return
    }

    setupMutation.mutate(
      {
        body: {
          username: trimmedUsername,
          password,
        },
      },
      {
        onSuccess: () => {
          onSuccess()
        },
      },
    )
  }

  // Extract Huma problem+json error detail if present
  const serverError = setupMutation.error
    ? (setupMutation.error as { detail?: string; message?: string })?.detail ||
      (setupMutation.error as { message?: string })?.message ||
      'Failed to create administrator account.'
    : null

  const displayError = validationError || serverError

  return (
    <div className="bg-background flex min-h-screen items-center justify-center p-4">
      <div className="bg-card text-card-foreground w-full max-w-md space-y-6 rounded-xl border p-6 shadow-sm sm:p-8">
        {/* Header */}
        <div className="flex flex-col items-center space-y-2 text-center">
          <div className="bg-primary/10 text-primary flex h-12 w-12 items-center justify-center rounded-full">
            <ShieldCheck className="h-6 w-6" />
          </div>
          <h1 className="text-2xl font-semibold tracking-tight">Welcome to File Ctrl</h1>
          <p className="text-muted-foreground text-sm">
            Create your administrator account to secure this instance.
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
        <form onSubmit={handleSubmit} className="space-y-4" noValidate>
          <div className="space-y-2">
            <Label htmlFor="setup-username">Administrator Username</Label>
            <Input
              id="setup-username"
              type="text"
              autoComplete="username"
              placeholder="admin"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value)
                setValidationError(null)
              }}
              disabled={setupMutation.isPending}
              required
              autoFocus
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="setup-password">Password</Label>
            <div className="relative">
              <Input
                id="setup-password"
                type={showPassword ? 'text' : 'password'}
                autoComplete="new-password"
                placeholder="At least 8 characters"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value)
                  setValidationError(null)
                }}
                disabled={setupMutation.isPending}
                className="pr-10"
                required
              />
              <button
                type="button"
                onClick={() => setShowPassword((prev) => !prev)}
                className="text-muted-foreground hover:text-foreground absolute top-1/2 right-3 -translate-y-1/2"
                tabIndex={-1}
                aria-label={showPassword ? 'Hide password' : 'Show password'}
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="setup-confirm-password">Confirm Password</Label>
            <Input
              id="setup-confirm-password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="new-password"
              placeholder="Re-enter your password"
              value={confirmPassword}
              onChange={(e) => {
                setConfirmPassword(e.target.value)
                setValidationError(null)
              }}
              disabled={setupMutation.isPending}
              required
            />
          </div>

          {/* Live Requirements Checklist */}
          <div className="bg-muted/50 text-muted-foreground space-y-1.5 rounded-md p-3 text-xs">
            <RequirementItem met={isUsernameValid} label="Username is at least 3 characters" />
            <RequirementItem
              met={isPasswordLengthValid}
              label="Password is at least 8 characters"
            />
            <RequirementItem met={doPasswordsMatch} label="Passwords match" />
          </div>

          <Button
            type="submit"
            className="w-full"
            disabled={
              setupMutation.isPending ||
              !isUsernameValid ||
              !isPasswordLengthValid ||
              !doPasswordsMatch
            }
          >
            {setupMutation.isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Creating Administrator...
              </>
            ) : (
              'Complete Setup'
            )}
          </Button>
        </form>

        {enabledOidcLoginUrl && (
          <OidcLoginButton href={enabledOidcLoginUrl} label="Set up with SSO" />
        )}
      </div>
    </div>
  )
}

function RequirementItem({ met, label }: { met: boolean; label: string }) {
  return (
    <div className="flex items-center gap-2">
      <CheckCircle2
        className={`h-3.5 w-3.5 shrink-0 transition-colors ${
          met ? 'text-emerald-500' : 'text-muted-foreground/40'
        }`}
      />
      <span className={met ? 'text-foreground' : 'text-muted-foreground'}>{label}</span>
    </div>
  )
}

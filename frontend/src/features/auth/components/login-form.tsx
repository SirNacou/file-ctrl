import { Button, Input } from '@base-ui/react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { AlertCircle, Eye, EyeOff, KeyRound, Loader2 } from 'lucide-react'
import { useState } from 'react'

import { authGetStatusQueryKey, authLoginMutation } from '#/client/gen/@tanstack/react-query.gen.ts'
import { Label } from '#/components/ui/label.tsx'

interface LoginFormProps {
  onSuccess: () => void
}

export function LoginForm({ onSuccess }: LoginFormProps) {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [validationError, setValidationError] = useState<string | null>(null)

  const queryClient = useQueryClient()
  const loginMutation = useMutation({
    ...authLoginMutation(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: authGetStatusQueryKey(),
      })
    },
  })

  const isFormFilled = username.trim().length > 0 && password.length > 0

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault()
    setValidationError(null)

    const trimmedUsername = username.trim()
    if (!trimmedUsername || !password) {
      setValidationError('Please enter both your username and password.')
      return
    }

    loginMutation.mutate(
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
  const serverError = loginMutation.error
    ? (loginMutation.error as { detail?: string; message?: string })?.detail ||
      (loginMutation.error as { message?: string })?.message ||
      'Invalid username or password.'
    : null

  const displayError = validationError || serverError

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
        <form onSubmit={handleSubmit} className="space-y-4" noValidate>
          <div className="space-y-2">
            <Label htmlFor="login-username">Username</Label>
            <Input
              id="login-username"
              type="text"
              autoComplete="username"
              placeholder="admin"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value)
                setValidationError(null)
              }}
              disabled={loginMutation.isPending}
              required
              autoFocus
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="login-password">Password</Label>
            <div className="relative">
              <Input
                id="login-password"
                type={showPassword ? 'text' : 'password'}
                autoComplete="current-password"
                placeholder="Enter your password"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value)
                  setValidationError(null)
                }}
                disabled={loginMutation.isPending}
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

          <Button
            type="submit"
            className="w-full"
            disabled={loginMutation.isPending || !isFormFilled}
          >
            {loginMutation.isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Signing in...
              </>
            ) : (
              'Sign In'
            )}
          </Button>
        </form>
      </div>
    </div>
  )
}

import { useState } from 'react'
import EyeIcon from '~icons/lucide/eye'
import EyeOffIcon from '~icons/lucide/eye-off'

import { Field, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'

import { useFieldContext } from './form-context'

type PasswordFieldProps = {
  label: string
  autoComplete?: 'current-password' | 'new-password'
}

export const PasswordField = ({ label, autoComplete = 'current-password' }: PasswordFieldProps) => {
  const field = useFieldContext<string>()
  const [showPassword, setShowPassword] = useState(false)

  return (
    <Field>
      <FieldLabel htmlFor={field.name}>{label}</FieldLabel>
      <div className="relative">
        <Input
          id={field.name}
          name={field.name}
          type={showPassword ? 'text' : 'password'}
          autoComplete={autoComplete}
          value={field.state.value}
          onChange={(e) => field.handleChange(e.target.value)}
          onBlur={field.handleBlur}
          className="pr-10"
        />
        <button
          type="button"
          onClick={() => setShowPassword((prev) => !prev)}
          className="text-muted-foreground hover:text-foreground absolute top-1/2 right-3 -translate-y-1/2"
          tabIndex={-1}
          aria-label={showPassword ? 'Hide password' : 'Show password'}
        >
          {showPassword ? <EyeOffIcon className="h-4 w-4" /> : <EyeIcon className="h-4 w-4" />}
        </button>
      </div>
    </Field>
  )
}

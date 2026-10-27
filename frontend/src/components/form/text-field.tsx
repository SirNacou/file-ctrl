import { Field, FieldLabel } from '#/components/ui/field'
import { Input } from '#/components/ui/input'

import { useFieldContext } from './form-context'

type TextFieldProps = {
  label: string
}

export const TextField = ({ label }: TextFieldProps) => {
  const field = useFieldContext<string>()
  return (
    <Field>
      <FieldLabel htmlFor={field.name}>{label}</FieldLabel>
      <Input
        id={field.name}
        name={field.name}
        value={field.state.value}
        onChange={(e) => field.handleChange(e.target.value)}
        onBlur={field.handleBlur}
      />
    </Field>
  )
}

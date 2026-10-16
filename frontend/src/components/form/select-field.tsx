import type { Group } from '@base-ui/react/internals/resolveValueLabel'
import type { Key } from 'react'
import { Field, FieldLabel } from '../ui/field'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../ui/select'
import { useFieldContext } from './form-context'

type SelectFieldProps<T extends Key> = {
  label: string
  placeholder?: string
  items: readonly Group<T>[] | readonly ItemType<T>[] | Record<string, React.ReactNode> | undefined
  disabled?: boolean
  alignItemWithTrigger?: boolean
}

export type ItemType<T extends Key> = {
  label: React.ReactNode
  value: T
}

export const SelectField = <T extends Key>({
  label,
  placeholder,
  items,
  disabled,
  alignItemWithTrigger = false,
}: SelectFieldProps<T>) => {
  const field = useFieldContext<T>()

  function isItemArray<T extends Key>(items: any): items is readonly ItemType<T>[] {
    if (!Array.isArray(items) || items.length <= 0) return false

    const firstItem = items.at(0)
    return 'value' in firstItem && 'label' in firstItem
  }
  return (
    <Field>
      <FieldLabel htmlFor={field.name}>{label}</FieldLabel>
      <Select
        value={field.state.value}
        onValueChange={(value) => value && field.handleChange(value)}
        disabled={disabled}
        id={field.name}
        name={field.name}
        items={items}
      >
        <SelectTrigger className="w-full">
          <SelectValue placeholder={placeholder} />
        </SelectTrigger>
        <SelectContent alignItemWithTrigger={alignItemWithTrigger}>
          <SelectGroup>
            {isItemArray(items) &&
              items.map((item) => (
                <SelectItem key={item.value} value={item.value}>
                  {item.label}
                </SelectItem>
              ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </Field>
  )
}

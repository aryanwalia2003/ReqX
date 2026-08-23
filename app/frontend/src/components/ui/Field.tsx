import { cloneElement, isValidElement, useId } from 'react'
import type { ReactElement } from 'react'

import { Label } from '@/components/ui/Label'
import { cn } from '@/lib/cn'

export interface FieldProps {
  label?: string
  hint?: string
  error?: string
  required?: boolean
  className?: string
  /** Ek single form control — Input/Textarea/Select/Checkbox waghera. */
  children: ReactElement
}

/**
 * Label + control + hint/error ko id/aria se jod deta hai — form banate
 * waqt manually htmlFor/aria-describedby likhna nahi padta.
 */
export function Field({ label, hint, error, required, className, children }: FieldProps) {
  const generatedId = useId()
  const controlProps = isValidElement(children) ? (children.props as Record<string, unknown>) : {}
  const id = (controlProps.id as string | undefined) ?? generatedId
  const describedBy = error ? `${id}-error` : hint ? `${id}-hint` : undefined

  const control = cloneElement(children, {
    id,
    'aria-describedby': describedBy,
    'aria-invalid': Boolean(error) || undefined,
    invalid: Boolean(error) || (controlProps as { invalid?: boolean }).invalid,
  } as Record<string, unknown>)

  return (
    <div className={cn('flex flex-col gap-1.5', className)}>
      {label && (
        <Label htmlFor={id} required={required}>
          {label}
        </Label>
      )}
      {control}
      {error ? (
        <p id={`${id}-error`} className="text-fg-muted text-xs">
          {error}
        </p>
      ) : hint ? (
        <p id={`${id}-hint`} className="text-fg-subtle text-xs">
          {hint}
        </p>
      ) : null}
    </div>
  )
}

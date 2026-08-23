import { createContext, useContext, useId, useState } from 'react'
import type { HTMLAttributes, ReactNode } from 'react'

import { cn } from '@/lib/cn'

interface TabsContextValue {
  value: string
  setValue: (value: string) => void
  baseId: string
}

const TabsContext = createContext<TabsContextValue | null>(null)

function useTabsContext(component: string) {
  const ctx = useContext(TabsContext)
  if (!ctx) throw new Error(`<${component}> must be used inside <Tabs>`)
  return ctx
}

export interface TabsProps {
  /** Uncontrolled default; ya value+onValueChange se controlled bana lo. */
  defaultValue?: string
  value?: string
  onValueChange?: (value: string) => void
  className?: string
  children: ReactNode
}

export function Tabs({ defaultValue, value, onValueChange, className, children }: TabsProps) {
  const baseId = useId()
  const [internal, setInternal] = useState(defaultValue ?? '')
  const active = value ?? internal

  function setValue(next: string) {
    setInternal(next)
    onValueChange?.(next)
  }

  return (
    <TabsContext.Provider value={{ value: active, setValue, baseId }}>
      <div className={className}>{children}</div>
    </TabsContext.Provider>
  )
}

export function TabsList({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      role="tablist"
      className={cn('border-border bg-surface inline-flex gap-1 rounded-md border p-1', className)}
      {...props}
    />
  )
}

export interface TabsTriggerProps extends HTMLAttributes<HTMLButtonElement> {
  value: string
}

export function TabsTrigger({ value, className, children, ...props }: TabsTriggerProps) {
  const ctx = useTabsContext('TabsTrigger')
  const selected = ctx.value === value
  return (
    <button
      type="button"
      role="tab"
      id={`${ctx.baseId}-tab-${value}`}
      aria-controls={`${ctx.baseId}-panel-${value}`}
      aria-selected={selected}
      onClick={() => ctx.setValue(value)}
      className={cn(
        'rounded-sm px-3 py-1.5 text-sm font-medium transition-colors',
        selected ? 'bg-surface-3 text-fg' : 'text-fg-muted hover:text-fg',
        className,
      )}
      {...props}
    >
      {children}
    </button>
  )
}

export interface TabsContentProps extends HTMLAttributes<HTMLDivElement> {
  value: string
}

export function TabsContent({ value, className, children, ...props }: TabsContentProps) {
  const ctx = useTabsContext('TabsContent')
  if (ctx.value !== value) return null
  return (
    <div
      role="tabpanel"
      id={`${ctx.baseId}-panel-${value}`}
      aria-labelledby={`${ctx.baseId}-tab-${value}`}
      className={className}
      {...props}
    >
      {children}
    </div>
  )
}

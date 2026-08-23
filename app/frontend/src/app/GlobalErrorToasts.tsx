import { useEffect } from 'react'

import { useToast } from '@/components/ui'
import { getErrorMessage } from '@/lib/errors'
import { onErrorReported } from '@/lib/reportError'

/** Kahin bhi `reportError()` chale to turant ek danger toast dikhata hai. */
export function GlobalErrorToasts() {
  const toast = useToast()

  useEffect(
    () =>
      onErrorReported(({ error }) => {
        toast({ title: getErrorMessage(error), variant: 'danger' })
      }),
    [toast],
  )

  return null
}

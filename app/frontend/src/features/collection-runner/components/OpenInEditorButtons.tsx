import { useState } from 'react'

import { Button } from '@/components/ui'
import { openInEditor } from '@/features/collection-runner/api'
import { getErrorMessage } from '@/lib/errors'

export interface OpenInEditorButtonsProps {
  /** Khaali/undefined ho to buttons disabled — abhi tak koi file path nahi. */
  path: string | undefined
  className?: string
}

/** VS Code / Vim / OS text editor me diya gaya file kholne ke teen chhote buttons. */
export function OpenInEditorButtons({ path, className }: OpenInEditorButtonsProps) {
  const [error, setError] = useState<string | undefined>()
  const disabled = !path?.trim()

  async function open(kind: 'vscode' | 'vim' | 'system', label: string) {
    if (!path) return
    setError(undefined)
    const result = await openInEditor(kind, path)
    if (!result.ok) {
      setError(`Could not open in ${label}: ${getErrorMessage(result.error)}`)
    }
  }

  return (
    <div
      className={className ? `flex flex-col items-end gap-1 ${className}` : 'flex flex-col gap-1'}
    >
      <div className="flex gap-1">
        <Button
          variant="ghost"
          size="sm"
          disabled={disabled}
          onClick={() => void open('vscode', 'VS Code')}
        >
          VS Code
        </Button>
        <Button
          variant="ghost"
          size="sm"
          disabled={disabled}
          onClick={() => void open('vim', 'Vim')}
        >
          Vim
        </Button>
        <Button
          variant="ghost"
          size="sm"
          disabled={disabled}
          onClick={() => void open('system', 'text editor')}
        >
          Text editor
        </Button>
      </div>
      {error && (
        <span role="alert" className="text-fg-muted text-xs">
          {error}
        </span>
      )}
    </div>
  )
}

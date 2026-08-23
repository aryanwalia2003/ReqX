import { render } from '@testing-library/react'
import { useRef } from 'react'
import { describe, expect, it, vi } from 'vitest'

import { useClickOutside } from '@/hooks/useClickOutside'

function TestBox({ onOutside }: { onOutside: () => void }) {
  const ref = useRef<HTMLDivElement>(null)
  useClickOutside(ref, onOutside)
  return (
    <div>
      <div data-testid="inside" ref={ref}>
        inside
      </div>
      <div data-testid="outside">outside</div>
    </div>
  )
}

describe('useClickOutside', () => {
  it('fires only for clicks outside the ref element', () => {
    const onOutside = vi.fn()
    const { getByTestId } = render(<TestBox onOutside={onOutside} />)

    getByTestId('inside').dispatchEvent(new MouseEvent('mousedown', { bubbles: true }))
    expect(onOutside).not.toHaveBeenCalled()

    getByTestId('outside').dispatchEvent(new MouseEvent('mousedown', { bubbles: true }))
    expect(onOutside).toHaveBeenCalledOnce()
  })
})

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import App from '@/App'
import { ToastProvider } from '@/components/ui'
import '@/styles/globals.css'

const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('#root element not found — check index.html')
}

createRoot(rootElement).render(
  <StrictMode>
    <ToastProvider>
      <App />
    </ToastProvider>
  </StrictMode>,
)

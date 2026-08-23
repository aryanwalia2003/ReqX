import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import App from '@/App'
import { GlobalErrorToasts } from '@/app/GlobalErrorToasts'
import { installGlobalErrorHandlers } from '@/app/installGlobalErrorHandlers'
import { ErrorBoundary, ToastProvider } from '@/components/ui'
import '@/styles/globals.css'

installGlobalErrorHandlers()

const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('#root element not found — check index.html')
}

createRoot(rootElement).render(
  <StrictMode>
    <ErrorBoundary context="root">
      <ToastProvider>
        <GlobalErrorToasts />
        <App />
      </ToastProvider>
    </ErrorBoundary>
  </StrictMode>,
)

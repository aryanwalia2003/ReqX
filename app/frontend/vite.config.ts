/// <reference types="vitest/config" />
import path from 'node:path'

import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      // Keep in sync with tsconfig.app.json's "paths" entries.
      '@': path.resolve(import.meta.dirname, './src'),
      // Wails auto-generates ./wailsjs (Go bindings + runtime) once `wails
      // dev`/`wails generate module` has run — this alias means feature code
      // never reaches for it with a relative "../../wailsjs/..." import.
      '@wails': path.resolve(import.meta.dirname, './wailsjs'),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
  },
})

import { defineConfig } from 'vite'
import path from 'path'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      bindings: path.resolve(__dirname, './bindings'),
    }
  },
  plugins: [
    react(),
  ],
  test: {
    globals: true, 
    environment: 'node',
    setupFiles: './vitest.setup.ts',
  },
})



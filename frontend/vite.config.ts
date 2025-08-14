import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'node:url'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    }
  },
  server: {
    port: 5173,
    proxy: {
      '/confessions': 'https://go-rest-api-469015.as.r.appspot.com/',
      '/tags': 'https://go-rest-api-469015.as.r.appspot.com/'
    }
  }
})

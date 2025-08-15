import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'node:url'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const backend = env.VITE_API_BASE_URL || 'https://go-rest-api-469015.as.r.appspot.com'

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      }
    },
    server: {
      port: 5173,
      proxy: {
        '/confessions': {
          target: backend,
          changeOrigin: true,
          secure: true
        },
        '/tags': {
          target: backend,
          changeOrigin: true,
          secure: true
        }
      }
    }
  }
})

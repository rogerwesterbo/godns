import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        // Split vendor code into stable chunks so it caches independently of app code.
        advancedChunks: {
          groups: [
            { name: 'react', test: /node_modules[\\/](react|react-dom|scheduler)[\\/]/, priority: 30 },
            { name: 'router', test: /node_modules[\\/](react-router|react-router-dom)[\\/]/, priority: 20 },
            { name: 'radix', test: /node_modules[\\/](@radix-ui|radix-ui)[\\/]/, priority: 10 },
            { name: 'vendor', test: /node_modules[\\/]/, priority: 0 },
          ],
        },
      },
    },
  },
})

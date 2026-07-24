import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],

  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  server: {
    // Listen on all interfaces (not just localhost) so the dev server is
    // reachable from other devices on the same network, e.g. a phone over Wi-Fi.
    host: true,

    // Proxy API calls through this same origin instead of hitting the
    // backend on a different port directly. Some mobile browsers apply
    // stricter same-site cookie rules across ports than desktop Chrome
    // does, which broke session cookies when the frontend and backend
    // were on different ports of the same LAN IP. Same-origin sidesteps
    // that entirely (and works for any host the phone reaches Vite on —
    // localhost, the LAN IP, whatever — with no hardcoded IP needed).
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
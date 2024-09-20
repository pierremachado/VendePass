import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({

  server: {
    host: '0.0.0.0',  // Isso permite que a aplicação seja acessível externamente
    port: 5173
  },
  
  plugins: [react()],
})

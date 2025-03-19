import { defineConfig } from 'vite'

export default defineConfig({
  base: "/movie_cat/",
  build: {
    rollupOptions: {
      external: ['register-service-worker']
    }
  },

})

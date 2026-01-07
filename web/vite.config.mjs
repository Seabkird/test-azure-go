import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './vitest.setup.mjs',
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
        // Si ton backend s'attend à recevoir "/ping" et non "/api/ping", décommente la ligne ci-dessous :
        // rewrite: (path) => path.replace(/^\/api/, ''),
      }
    },
  },
});
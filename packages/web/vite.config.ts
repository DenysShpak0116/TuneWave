import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  resolve: {
    alias: {
      '@assets': '/src/assets',
      '@modules': '/src/modules',
      '@ui': '/src/ui',
      '@components': '/src/components',
      '@api': '/src/api',
      '@consts': '/src/consts'
    }
  }
});

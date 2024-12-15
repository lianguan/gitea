import {defineConfig} from 'vitest/config';
import vuePlugin from '@vitejs/plugin-vue';
import {stringPlugin} from 'vite-string-plugin';

export default defineConfig({
  test: {
    include: ['bundles/**/*.test.ts'],
    setupFiles: ['bundles/js/vitest.setup.ts'],
    environment: 'happy-dom',
    testTimeout: 20000,
    open: false,
    allowOnly: true,
    passWithNoTests: true,
    globals: true,
    watch: false,
  },
  plugins: [
    stringPlugin(),
    vuePlugin(),
  ],
});

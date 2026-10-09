import { devtools } from '@tanstack/devtools-vite'
import { defineConfig, lazyPlugins } from 'vite-plus'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'

import babel from '@rolldown/plugin-babel'
import tailwindcss from '@tailwindcss/vite'
import viteReact, { reactCompilerPreset } from '@vitejs/plugin-react'
import { nitro } from 'nitro/vite'
import Icons from 'unplugin-icons/vite'

const config = defineConfig({
  fmt: {
    semi: false,
    singleQuote: true,
    trailingComma: 'all',
  },
  lint: {
    jsPlugins: [{ name: 'vite-plus', specifier: 'vite-plus/oxlint-plugin' }],
    rules: { 'vite-plus/prefer-vite-plus-imports': 'error' },
    options: { typeAware: true, typeCheck: true },
    ignorePatterns: ['*.gen.*'],
  },
  resolve: { tsconfigPaths: true },
  server: {
    host: true,
    port: 3000,
    watch: {
      usePolling: true,
    },
  },
  envPrefix: ['PUBLIC_', 'VITE_'],
  plugins: lazyPlugins(() => [
    devtools(),
    nitro({ rollupConfig: { external: [/^@sentry\//] } }),
    tailwindcss(),
    tanstackStart({
      // spa: {
      //   enabled: true,
      // },
    }),
    viteReact(),
    babel({ presets: [reactCompilerPreset()] }),
    Icons({
      autoInstall: true,
      compiler: 'jsx',
      jsx: 'react',
    }),
  ]),
})

export default config

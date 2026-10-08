import { devtools } from "@tanstack/devtools-vite";
import { defineConfig, lazyPlugins } from "vite-plus";

import { tanstackStart } from "@tanstack/react-start/plugin/vite";

import { heyApiPlugin } from "@hey-api/vite-plugin";
import babel from "@rolldown/plugin-babel";
import tailwindcss from "@tailwindcss/vite";
import viteReact, { reactCompilerPreset } from "@vitejs/plugin-react";
import { nitro } from "nitro/vite";

const config = defineConfig({
  fmt: {
    semi: false,
    singleQuote: true,
    trailingComma: "all",
  },
  lint: {
    jsPlugins: [{ name: "vite-plus", specifier: "vite-plus/oxlint-plugin" }],
    rules: { "vite-plus/prefer-vite-plus-imports": "error" },
    options: { typeAware: true, typeCheck: true },
  },
  resolve: { tsconfigPaths: true },
  server: {
    host: true,
    port: 3000,
    watch: {
      usePolling: true,
      interval: 300,
    },
  },
  envPrefix: ["PUBLIC_", "VITE_"],
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
    heyApiPlugin({
      config: {
        input: {
          path: "http://localhost:8080/openapi.json",
          watch: true,
        },
        output: { path: "src/client/gen", postProcess: ["oxlint", "oxfmt"] },
        plugins: [
          {
            name: "@hey-api/typescript",
            enums: "typescript-const",
          },
          {
            name: "@hey-api/sdk",
            validator: true,
            transformer: true,
          },
          {
            name: "@hey-api/client-fetch",
            runtimeConfigPath: "./src/client/hey-api.ts",
          },
          {
            name: "zod",
            requests: true,
            responses: true,
            definitions: true,
            dates: {
              offset: true,
            },
            types: {
              infer: true,
              input: true,
              output: true,
            },
          },
          {
            name: "@tanstack/react-query",
            queryKeys: {
              tags: true,
            },
            infiniteQueryKeys: {
              tags: true,
            },
          },
        ],
      },
      vite: {
        apply: "serve",
      },
    }),
  ]),
});

export default config;

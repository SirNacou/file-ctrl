import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: {
    path: 'http://localhost:8080/openapi.json',
    watch: true,
  },
  output: {
    path: 'src/client/gen',
    postProcess: ['oxlint', 'oxfmt'],
    clean: true,
  },
  logs: 'logs',
  plugins: [
    // 1. Hook up the SDK plugin specifically using Zod validation and transformation rules
    {
      name: '@hey-api/sdk',
      validator: 'zod',
      transformer: true,
    },
    {
      name: '@hey-api/transformers',
      dates: true,
      bigInt: true,
    },
    {
      name: '@hey-api/client-fetch',
      runtimeConfigPath: './src/client/hey-api.ts',
    },
    // 2. Configure Zod v4 to exclusively handle schemas + infer static TS types
    {
      name: 'zod',
      requests: {
        enabled: true,
      },
      responses: {
        enabled: true,
      },
      definitions: {
        enabled: true,
      },
      dates: {
        offset: true,
      },
      types: {
        input: true,
        infer: true,
        output: true,
      },
    },
    {
      name: '@tanstack/react-query',
      queryKeys: {
        tags: true,
      },
      infiniteQueryKeys: {
        tags: true,
      },
    },
  ],
})

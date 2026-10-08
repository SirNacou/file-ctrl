import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: {
    path: "http://localhost:8080/openapi.json",
    watch: true,
  },
  output: { path: "src/client/gen", postProcess: ["oxlint", "oxfmt"], clean: true },
  logs: "logs",
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
});

import { env } from '#/config/env.ts'
import type { CreateClientConfig } from './gen/client.gen'

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseUrl: env.PUBLIC_APP_URL,
})

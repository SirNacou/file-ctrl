import { defineEnv } from 'envin'
import { z } from 'zod'

export const env = defineEnv({
  clientPrefix: 'PUBLIC_',
  client: {
    PUBLIC_APP_URL: z.url(),
  },
  env: {
    ...process.env, // Server env variables
    ...import.meta.env, // Client env variables
  },
})

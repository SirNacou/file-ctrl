import { defineEnv } from "envin"
import { z } from "zod"

export const env = defineEnv({
  shared: {
    APP_URL: z.url(),
  },
  clientPrefix: "PUBLIC_",
  env: {
    ...process.env,
    ...import.meta.env,
  },
})

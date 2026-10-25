import { env } from '#/config/env.ts'
import { tusAdapter } from '@uploadcn/core/tus'

const FILE_UPLOAD_ENDPOINT = `${env.PUBLIC_APP_URL}/api/storage/files/uploads`

const MB = 1024 * 1024

type Props = {
  path: string
  chunkSize?: number
  subpath?: string
}

export const getAdapter = ({ path, chunkSize = 1 * MB, subpath = '' }: Props) =>
  tusAdapter({
    endpoint: FILE_UPLOAD_ENDPOINT,
    chunkSize: chunkSize,
    metadata: {
      storage_path: path,
      parent_path: subpath,
    },
  })

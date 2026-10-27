import type { UploadAdapter } from '@uploadcn/core'
import { tusAdapter } from '@uploadcn/core/tus'

import { env } from '#/config/env.ts'

const FILE_UPLOAD_ENDPOINT = `${env.PUBLIC_APP_URL}/api/storage/files/uploads`

const MB = 1024 * 1024

type Props = {
  path: string
  chunkSize?: number
  subpath?: string
}

export const getAdapter = ({ path, chunkSize = 1 * MB, subpath = '' }: Props): UploadAdapter => {
  const buildAdapter = (relativePath: string) => {
    const lastIndex = relativePath.lastIndexOf('/')
    const directoryPath = lastIndex !== -1 ? relativePath.slice(0, lastIndex) : ''
    return tusAdapter({
      endpoint: FILE_UPLOAD_ENDPOINT,
      chunkSize: chunkSize,
      metadata: {
        storage_path: path,
        storage_subpath: subpath,
        directory_path: directoryPath,
      },
    })
  }

  return {
    name: 'tus',
    resumable: true,
    upload: async (context) => buildAdapter(context.file.webkitRelativePath).upload(context),
    abort: async (context) => buildAdapter(context.item.file.webkitRelativePath).abort?.(context),
  }
}

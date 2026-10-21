import { client } from '#/client/gen/client.gen.ts'
import { type StorageFileDownloadData } from '#/client/gen/index.ts'
import type { AnchorHTMLAttributes } from 'react'

type Props = Omit<AnchorHTMLAttributes<HTMLAnchorElement>, 'href' | 'download'> & {
  storagePath: string
  filePath: string
}

export function DownloadLink({ storagePath, filePath, children, ...rest }: Props) {
  const url = client.buildUrl<StorageFileDownloadData>({
    path: {
      storage_path: storagePath,
    },
    query: {
      file_path: filePath,
    },
    url: '/api/storage/{storage_path}/files/download',
  })
  return (
    <a {...rest} href={url} download target="_blank" rel="noopener">
      {children}
    </a>
  )
}

import type { AnchorHTMLAttributes } from 'react'

import { client } from '#/client/gen/client.gen.ts'
import {
  type StorageDownloadArchiveData,
  type StorageFileDownloadData,
} from '#/client/gen/index.ts'

type BaseProps = Omit<AnchorHTMLAttributes<HTMLAnchorElement>, 'href' | 'download'>

function BaseLink({ url, children, ...rest }: BaseProps & { url: string }) {
  return (
    <a {...rest} href={url} download target="_blank" rel="noopener">
      {children}
    </a>
  )
}

export function FileDownloadLink({
  storagePath,
  filePath,
  ...rest
}: BaseProps & { storagePath: string; filePath: string }) {
  const url = client.buildUrl<StorageFileDownloadData>({
    path: {
      storage_path: storagePath,
    },
    query: {
      file_path: filePath,
    },
    querySerializer: {
      array: { explode: false, style: 'form' },
    },
    url: '/api/storage/{storage_path}/files/download',
  })

  return <BaseLink url={url} {...rest} />
}

export function ArchiveDownloadLink({
  storagePath,
  filePaths,
  ...rest
}: BaseProps & { storagePath: string; filePaths: string[] }) {
  const url = client.buildUrl<StorageDownloadArchiveData>({
    path: {
      storage_path: storagePath,
    },
    query: {
      file_paths: filePaths,
    },
    querySerializer: {
      array: { explode: false, style: 'form' },
    },
    url: '/api/storage/{storage_path}/files/download-archive',
  })
  console.log(url)

  return <BaseLink url={url} {...rest} />
}

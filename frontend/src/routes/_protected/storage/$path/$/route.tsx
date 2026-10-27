import { useQueryClient } from '@tanstack/react-query'
import { createFileRoute, Outlet } from '@tanstack/react-router'
import { createIndexedDBPersistence, type UploadError, type UploadItem } from '@uploadcn/core'
import { UploadConfigProvider } from '@uploadcn/react'
import { useAtomValue } from 'jotai'

import { storageListDirectoryQueryKey } from '#/client/gen/@tanstack/react-query.gen.ts'
import { uploadDialogOpenAtom } from '#/features/explorer/components/file-upload-dialog.tsx'
import { getAdapter } from '#/features/explorer/lib/tus-adapter.ts'
import { notify } from '#/lib/notify.ts'
import { GlobalDropzone, UploadProvider, UploadQueuePanel } from '@/components/global-upload'

export const Route = createFileRoute('/_protected/storage/$path/$')({
  component: RouteComponent,
})

function RouteComponent() {
  const { path, _splat: subpath } = Route.useParams()
  const queryClient = useQueryClient()

  const dialogOpen = useAtomValue(uploadDialogOpenAtom)

  function handleError(_: UploadItem, error: UploadError) {
    notify.error('Upload failed', error.message)
  }

  async function handleSuccess() {
    await queryClient.invalidateQueries({
      queryKey: storageListDirectoryQueryKey({
        path: {
          storage_path: path,
        },
        query: {
          subpath,
        },
      }),
    })
  }

  return (
    <UploadConfigProvider
      adapter={getAdapter({
        path,
        subpath,
      })}
      autoUpload={!dialogOpen}
      persistence={createIndexedDBPersistence()}
    >
      <UploadProvider duplicates="replace" persist onSuccess={handleSuccess} onError={handleError}>
        <Outlet />

        {!dialogOpen && (
          <>
            <GlobalDropzone />
            <UploadQueuePanel />
          </>
        )}
      </UploadProvider>
    </UploadConfigProvider>
  )
}

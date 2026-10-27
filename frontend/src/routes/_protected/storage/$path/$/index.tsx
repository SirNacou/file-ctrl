import { useMutation, useQueryClient, useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useAtom } from '@tanstack/react-store'
import { useMemo } from 'react'
import DownloadIcon from '~icons/lucide/download'
import TrashIcon from '~icons/lucide/trash'

import {
  storageDeleteFilesMutation,
  storageListDirectoryOptions,
  storageListDirectoryQueryKey,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { type DirectoryEntry } from '#/client/gen/index.ts'
import { ConfirmDeleteDialog } from '#/components/confirm-detele-dialog.tsx'
import { PageHeader } from '#/components/page-header.tsx'
import { Button } from '#/components/ui/button.tsx'
import { CreateFolderDialog } from '#/features/explorer/components/create-folder-dialog.tsx'
import {
  ArchiveDownloadLink,
  FileDownloadLink,
} from '#/features/explorer/components/download-link.tsx'
import { FileBreadcrumbs } from '#/features/explorer/components/file-breadcrumb.tsx'
import { FileTable, rowSelectionAtom } from '#/features/explorer/components/file-table.tsx'
import { FileUploadDialog } from '#/features/explorer/components/file-upload-dialog.tsx'
import { notify } from '#/lib/notify.ts'

export const Route = createFileRoute('/_protected/storage/$path/$/')({
  loader: async ({ params: { path, _splat: subpath }, context: { queryClient } }) =>
    await queryClient.query({
      ...storageListDirectoryOptions({ path: { storage_path: path }, query: { subpath: subpath } }),
      staleTime: 'static',
    }),
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = Route.useNavigate()
  const { path, _splat: subpath } = Route.useParams()

  const { data } = useSuspenseQuery(
    storageListDirectoryOptions({ path: { storage_path: path }, query: { subpath: subpath } }),
  )

  const [rowSelection, setRowSelection] = useAtom(rowSelectionAtom)
  const { emptySelection, rowIds } = useMemo(() => {
    const rowIds = Object.keys(rowSelection)
    const emptySelection = rowIds.length === 0
    return { rowIds, emptySelection }
  }, [rowSelection])

  async function handleItemSelected(item: DirectoryEntry) {
    if (item.is_dir) {
      await navigate({
        params: {
          path,
          _splat: item.path,
        },
      })
    }
  }

  const queryClient = useQueryClient()
  const deleteFiles = useMutation({
    ...storageDeleteFilesMutation(),
    onSuccess: async () => {
      setRowSelection(() => ({}))
      await queryClient.invalidateQueries({
        queryKey: storageListDirectoryQueryKey({
          path: { storage_path: path },
          query: { subpath: subpath },
        }),
      })
    },
    onError: (error) => {
      notify.error(error.title, error.errors?.at(0)?.message)
    },
  })

  async function handleBulkDelete() {
    if (emptySelection) return

    await deleteFiles.mutateAsync({
      path: {
        storage_path: path,
      },
      body: {
        file_paths: rowIds,
      },
    })
  }

  return (
    <>
      <PageHeader
        title={`Storage: ${data.storage_name}`}
        actions={
          <>
            <ConfirmDeleteDialog
              onConfirm={handleBulkDelete}
              trigger={
                <Button variant={'destructive'} disabled={emptySelection}>
                  <TrashIcon />
                  Delete
                  {!emptySelection && ` (${rowIds.length})`}
                </Button>
              }
            />
            <Button
              disabled={emptySelection}
              render={
                rowIds.length > 1 ? (
                  <ArchiveDownloadLink storagePath={path} filePaths={rowIds} />
                ) : rowIds.length === 1 ? (
                  <FileDownloadLink storagePath={path} filePath={rowIds[0]} />
                ) : undefined
              }
            >
              <DownloadIcon />
              Download
              {!emptySelection && ` (${rowIds.length})`}
            </Button>
            <FileUploadDialog />
            <CreateFolderDialog storagePath={path} subpath={subpath} />
          </>
        }
      />
      <FileBreadcrumbs />
      <FileTable
        storagePath={path}
        subpath={subpath}
        data={data.items}
        itemSelected={handleItemSelected}
      />
    </>
  )
}

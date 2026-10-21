import {
  listDirectoryOptions,
  listDirectoryQueryKey,
  storageDeleteFilesMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import type { DirectoryEntry } from '#/client/gen/index.ts'
import { ConfirmDeleteDialog } from '#/components/confirm-detele-dialog.tsx'
import { PageHeader } from '#/components/page-header.tsx'
import { Button } from '#/components/ui/button.tsx'
import { CreateFolderDialog } from '#/features/explorer/components/create-folder-dialog.tsx'
import { FileBreadcrumbs } from '#/features/explorer/components/file-breadcrumb.tsx'
import { FileTable, rowSelectionAtom } from '#/features/explorer/components/file-table.tsx'
import { notify } from '#/lib/notify.ts'
import { useMutation, useQueryClient, useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useSelector } from '@tanstack/react-store'
import { useMemo } from 'react'
import TrashIcon from '~icons/lucide/trash'

export const Route = createFileRoute('/_protected/storage/$path/$/')({
  loader: async ({ params: { path }, context: { queryClient } }) =>
    await queryClient.query({
      ...listDirectoryOptions({ path: { storage_path: path } }),
      staleTime: 'static',
    }),
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = Route.useNavigate()
  const { path, _splat: subpath } = Route.useParams()
  const { data } = useSuspenseQuery(
    listDirectoryOptions({ path: { storage_path: path }, query: { subpath: subpath } }),
  )

  const rowSelection = useSelector(rowSelectionAtom)
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
      notify.success('Delete successfully!')
      await queryClient.invalidateQueries({
        queryKey: listDirectoryQueryKey({ path: { storage_path: path } }),
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
            <CreateFolderDialog storagePath={path} parentPath={subpath} />
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

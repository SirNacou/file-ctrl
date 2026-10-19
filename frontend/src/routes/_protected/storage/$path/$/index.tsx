import { listDirectoryOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import type { DirectoryEntry } from '#/client/gen/index.ts'
import { PageHeader } from '#/components/page-header.tsx'
import { CreateFolderDialog } from '#/features/explorer/components/create-folder-dialog.tsx'
import { FileBreadcrumbs } from '#/features/explorer/components/file-breadcrumb.tsx'
import { FileTable } from '#/features/explorer/components/file-table.tsx'
import { useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/storage/$path/$/')({
  loader: ({ params: { path }, context: { queryClient } }) =>
    queryClient.query({
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

  return (
    <>
      <PageHeader
        title={`Storage: ${data.storage_name}`}
        actions={<CreateFolderDialog storagePath={path} parentPath={subpath} />}
      />
      <FileBreadcrumbs />
      <FileTable data={data.items ?? []} itemSelected={handleItemSelected} />
    </>
  )
}

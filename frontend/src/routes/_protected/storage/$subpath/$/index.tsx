import { listDirectoryOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { PageHeader } from '#/components/page-header.tsx'
import { Breadcrumb, BreadcrumbItem, BreadcrumbList } from '#/components/ui/breadcrumb.tsx'
import { FileTable } from '#/features/explorer/components/file-table.tsx'
import { useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/storage/$subpath/$/')({
  loader: ({ params: { subpath }, context: { queryClient } }) =>
    queryClient.query({
      ...listDirectoryOptions({ path: { storage_subpath: subpath } }),
      staleTime: 'static',
    }),
  component: RouteComponent,
})

function RouteComponent() {
  const { subpath, _splat } = Route.useParams()
  const { data } = useSuspenseQuery(
    listDirectoryOptions({ path: { storage_subpath: subpath }, query: { path: _splat } }),
  )

  return (
    <>
      <PageHeader title={`Storage: ${data.storage_name}`} />
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem></BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <FileTable data={data.items ?? []} />
    </>
  )
}

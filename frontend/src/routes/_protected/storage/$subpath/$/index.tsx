import { listDirectoryOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import type { DirectoryEntry } from '#/client/gen/index.ts'
import { PageHeader } from '#/components/page-header.tsx'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '#/components/ui/breadcrumb.tsx'
import { FileTable } from '#/features/explorer/components/file-table.tsx'
import { useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute, Link } from '@tanstack/react-router'
import { Fragment, useMemo } from 'react'
export const Route = createFileRoute('/_protected/storage/$subpath/$/')({
  loader: ({ params: { subpath }, context: { queryClient } }) =>
    queryClient.query({
      ...listDirectoryOptions({ path: { storage_subpath: subpath } }),
      staleTime: 'static',
    }),
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = Route.useNavigate()
  const { subpath, _splat } = Route.useParams()
  const { data } = useSuspenseQuery(
    listDirectoryOptions({ path: { storage_subpath: subpath }, query: { path: _splat } }),
  )

  async function handleItemSelected(item: DirectoryEntry) {
    if (item.is_dir) {
      await navigate({
        params: {
          subpath,
          _splat: item.path,
        },
      })
    }
  }

  const breadcrumbItems = useMemo(() => getPathItems(data.current_path), [data.current_path])

  return (
    <>
      <PageHeader title={`Storage: ${data.storage_name}`} />
      <Breadcrumb>
        <BreadcrumbList className="text-base">
          <BreadcrumbItem>
            <BreadcrumbLink
              render={
                <Link to={Route.path} params={{ subpath, _splat: undefined }}>
                  {subpath}
                </Link>
              }
            />
          </BreadcrumbItem>
          {breadcrumbItems.map((item) => (
            <Fragment key={item.path}>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink
                  render={
                    <Link to={Route.path} params={{ subpath, _splat: item.path }}>
                      {item.name}
                    </Link>
                  }
                />
              </BreadcrumbItem>
            </Fragment>
          ))}
        </BreadcrumbList>
      </Breadcrumb>
      <FileTable data={data.items ?? []} itemSelected={handleItemSelected} />
    </>
  )
}
type PathItem = {
  name: string
  path: string
}

function getPathItems(path: string): PathItem[] {
  const segments = path.split('/').filter(Boolean)

  return segments.map((name, index) => ({
    name,
    path: segments.slice(0, index + 1).join('/'),
  }))
}

import { listDirectoryOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { PageHeader } from '#/components/page-header.tsx'
import { useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/storage/$sourceId/')({
  loader: ({ params: { sourceId }, context: { queryClient } }) =>
    queryClient.query({
      ...listDirectoryOptions({ path: { id: sourceId } }),
      staleTime: 'static',
    }),
  component: RouteComponent,
})

function RouteComponent() {
  const { sourceId } = Route.useParams()
  const { data } = useSuspenseQuery(listDirectoryOptions({ path: { id: sourceId } }))

  return (
    <>
      <PageHeader title={`Storage: ${data.storage_name}`} />
    </>
  )
}

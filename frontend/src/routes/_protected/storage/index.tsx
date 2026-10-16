import { getApiStorageOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { PageHeader } from '#/components/page-header.tsx'
import { AddStorageDialog } from '#/features/storage/components/add-storage-dialog.tsx'
import { StorageTable } from '#/features/storage/components/storage-table.tsx'
import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/storage/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { data } = useQuery(getApiStorageOptions())

  return (
    <>
      <PageHeader title="Storage" actions={<AddStorageDialog />} />
      <StorageTable data={data?.sources ?? []} />
    </>
  )
}

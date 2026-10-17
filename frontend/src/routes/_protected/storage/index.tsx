import { getApiStorageOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { PageHeader } from '#/components/page-header.tsx'
import { AddStorageDialog } from '#/features/storage/components/add-storage-dialog.tsx'
import { StorageCard } from '#/features/storage/components/storage-card.tsx'
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
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {data?.sources?.map((s) => (
          <StorageCard key={s.subpath} storage={s} onEdit={() => {}} onRemove={() => {}} />
        ))}
      </div>
    </>
  )
}

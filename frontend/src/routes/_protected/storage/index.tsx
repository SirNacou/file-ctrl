import { useSuspenseQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

import { listStorageSourcesOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { PageHeader } from '#/components/page-header.tsx'
import { AddStorageDialog } from '#/features/storage/components/add-storage-dialog.tsx'
import { StorageCard } from '#/features/storage/components/storage-card.tsx'

export const Route = createFileRoute('/_protected/storage/')({
  loader: async ({ context: { queryClient } }) => {
    await queryClient.query({
      ...listStorageSourcesOptions(),
      staleTime: 'static',
    })
  },
  component: RouteComponent,
})
function RouteComponent() {
  const { data: sourcesRes } = useSuspenseQuery(listStorageSourcesOptions())

  return (
    <>
      <PageHeader title="Storage" actions={<AddStorageDialog />} />
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {sourcesRes.sources.map((s) => (
          <StorageCard key={s.path} storage={s} onEdit={() => {}} onRemove={() => {}} />
        ))}
      </div>
    </>
  )
}

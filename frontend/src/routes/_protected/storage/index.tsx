import { PageHeader } from '#/components/page-header.tsx'
import { AddStorageDialog } from '#/features/storage/components/add-storage-dialog.tsx'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/storage/')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div>
      <PageHeader title="Storage" actions={<AddStorageDialog />} />
    </div>
  )
}

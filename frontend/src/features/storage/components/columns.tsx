import type { Source } from '#/client/gen/index.ts'
import { ConfirmDelete } from '#/components/confirm-detele.tsx'
import { Button } from '#/components/ui/button.tsx'
import { Link } from '@tanstack/react-router'
import { createColumnHelper } from '@tanstack/react-table'
import type { Features } from './features'

const columnHelper = createColumnHelper<Features, Source>()

export const columns = columnHelper.columns([
  columnHelper.accessor('name', {
    header: () => <span>Name</span>,
  }),
  columnHelper.accessor('subpath', {
    header: () => <span>Path</span>,
  }),
  columnHelper.accessor('created_at', {
    header: () => <span>Created At</span>,
    cell: ({ row }) => {
      const date = new Date(row.original.created_at)
      return (
        <span>
          {date.toLocaleDateString(undefined, {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
          })}
        </span>
      )
    },
  }),
  columnHelper.display({
    id: 'actions',
    header: () => <span>Actions</span>,
    cell: ({ row, table }) => (
      <div className="flex gap-2">
        <Button
          render={
            <Link
              to="/storage/$storage_subpath"
              params={{ storage_subpath: row.original.subpath }}
            />
          }
        >
          Open
        </Button>
        <ConfirmDelete
          buttonText="Remove"
          onConfirm={() => table.options.meta?.handleRemove(row.original.subpath)}
        />
      </div>
    ),
    enablePinning: true,
    size: 20,
  }),
])

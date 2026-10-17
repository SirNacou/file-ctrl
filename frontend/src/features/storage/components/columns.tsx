import type { Source } from '#/client/gen/index.ts'
import { ConfirmDeleteDialog } from '#/components/confirm-detele-dialog.tsx'
import { Button } from '#/components/ui/button.tsx'
import { Link } from '@tanstack/react-router'
import { createColumnHelper } from '@tanstack/react-table'
import type { Features } from './features'
import { RenameStorageDialog } from './rename-storage-dialog'

const columnHelper = createColumnHelper<Features, Source>()

export const columns = columnHelper.columns([
  columnHelper.accessor('name', {
    header: () => <span>Name</span>,
    cell: ({ row }) => (
      <Link
        to="/storage/$subpath/$"
        params={{ subpath: row.original.subpath }}
        className={`block px-3 py-2 text-inherit no-underline hover:underline focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-blue-500`}
      >
        {row.original.name}
      </Link>
    ),
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
          render={<Link to="/storage/$subpath/$" params={{ subpath: row.original.subpath }} />}
        >
          Open
        </Button>
        <RenameStorageDialog storage={row.original} />
        <ConfirmDeleteDialog
          buttonText="Remove"
          onConfirm={() => table.options.meta?.handleRemove(row.original.subpath)}
        />
      </div>
    ),
    enablePinning: true,
    size: 20,
  }),
])

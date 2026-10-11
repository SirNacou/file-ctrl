import type { FileItem } from '#/client/gen/index.ts'
import { Checkbox } from '#/components/ui/checkbox.tsx'
import { createColumnHelper } from '@tanstack/react-table'
import { formatBytes, getFileIcon } from '../types'
import type { Features } from './features'

const columnHelper = createColumnHelper<Features, FileItem>()

export const columns = columnHelper.columns([
  columnHelper.display({
    id: 'select',
    header: ({ table }) => {
      return (
        <Checkbox
          className="size-5"
          checked={table.getIsAllRowsSelected()}
          indeterminate={!table.getIsAllRowsSelected() && table.getIsSomeRowsSelected()}
          onCheckedChange={(checked) => table.toggleAllRowsSelected(checked)}
        />
      )
    },
    cell: ({ row }) => (
      <Checkbox
        className="size-5"
        checked={row.getIsSelected()}
        disabled={!row.getCanSelect()}
        indeterminate={row.getIsSomeSelected()}
        onCheckedChange={(checked) => row.toggleSelected(checked)}
      />
    ),
    enableSorting: false,
    size: 40,
  }),
  columnHelper.accessor('name', {
    header: () => <span>Name</span>,
    cell: ({ row }) => {
      const item = row.original
      const Icon = getFileIcon(item.name, item.is_dir)

      return (
        <div
          className="group flex items-center gap-2.5 py-0.5 cursor-pointer select-none"
          onDoubleClick={() => {
            row.table.options.meta?.itemSelected(item)
          }}
        >
          <Icon
            className={`h-4 w-4 shrink-0 transition-transform group-hover:scale-110 ${
              item.is_dir ? 'text-amber-500 fill-amber-500/20' : 'text-muted-foreground'
            }`}
          />
          <span
            className={`font-medium truncate ${
              item.is_dir ? 'hover:underline text-foreground' : 'text-foreground/90 font-normal'
            }`}
          >
            {item.name}
          </span>
        </div>
      )
    },
    sortFn: (rowA, rowB) => {
      const a = rowA.original
      const b = rowB.original
      if (a.is_dir && !b.is_dir) return -1
      if (!a.is_dir && b.is_dir) return 1
      return a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' })
    },
    size: undefined,
  }),
  columnHelper.accessor('size', {
    header: () => <span>Size</span>,
    cell: ({ row }) => {
      const item = row.original
      return <span className="font-mono">{item.is_dir ? 'Folder' : formatBytes(item.size)}</span>
    },
    size: 300,
  }),
  columnHelper.accessor('mod_time', {
    header: () => <span>Modified At</span>,
    cell: ({ row }) => {
      const date = new Date(row.original.mod_time)
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
    size: 300,
  }),
])

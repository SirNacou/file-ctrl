import type { DirectoryEntry } from '#/client/gen/index.ts'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '#/components/ui/context-menu.tsx'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table.tsx'
import { useTable } from '@tanstack/react-table'
import { columns } from './columns'
import { features } from './features'

type FileTableProps = {
  data: DirectoryEntry[]
  itemSelected?: (item: DirectoryEntry) => void
}

export const FileTable = ({ data, itemSelected }: FileTableProps) => {
  const table = useTable({
    key: 'file-table',
    features,
    columns,
    data,
    meta: {
      itemSelected: (item) => {
        itemSelected?.(item)
      },
    },
  })

  function handleRename(path: string) {
    console.log(path)
  }

  return (
    <div
      className="bg-card shadow border-2 overflow-auto"
      style={{
        height: 'calc(100vh - 200px)',
      }}
    >
      <Table noWrapper className="text-base">
        <TableHeader className="top-0 z-10 sticky bg-card">
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                const hasExplicitSize = header.column.columnDef.size !== undefined
                return (
                  <TableHead
                    className="font-bold"
                    key={header.id}
                    style={{
                      width: hasExplicitSize ? header.column.getSize() : 'auto',
                    }}
                  >
                    {header.isPlaceholder ? null : <table.FlexRender header={header} />}
                  </TableHead>
                )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <ContextMenu key={row.id}>
                <ContextMenuTrigger
                  render={
                    <TableRow data-state={row.getIsSelected() && 'selected'}>
                      {row.getAllCells().map((cell) => {
                        const hasExplicitSize = cell.column.columnDef.size !== undefined
                        return (
                          <TableCell
                            key={cell.id}
                            style={{
                              width: hasExplicitSize ? cell.column.getSize() : 'auto',
                            }}
                          >
                            <table.FlexRender cell={cell} />
                          </TableCell>
                        )
                      })}
                    </TableRow>
                  }
                />
                <ContextMenuContent>
                  <ContextMenuItem onClick={() => handleRename(row.original.path)}>
                    Rename
                  </ContextMenuItem>
                </ContextMenuContent>
              </ContextMenu>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                This folder is empty.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}

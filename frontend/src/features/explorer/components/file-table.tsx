import { createAtom } from '@tanstack/react-store'
import { useTable, type RowSelectionState } from '@tanstack/react-table'
import { useEffect } from 'react'

import type { DirectoryEntry } from '#/client/gen/index.ts'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table.tsx'

import { columns } from './columns'
import { features } from './features'
import { FileTableContextMenu } from './file-table-context-menu'

type FileTableProps = {
  storagePath: string
  subpath?: string
  data: DirectoryEntry[]
  itemSelected?: (item: DirectoryEntry) => void
}

export const rowSelectionAtom = createAtom<RowSelectionState>({})

export const FileTable = ({ storagePath, subpath, data, itemSelected }: FileTableProps) => {
  const table = useTable({
    key: 'file-table',
    features,
    columns,
    data,
    getRowId: (row) => row.path,
    meta: {
      itemSelected: (item) => {
        itemSelected?.(item)
      },
    },
    atoms: {
      rowSelection: rowSelectionAtom,
    },
  })

  useEffect(() => {
    rowSelectionAtom.set(() => ({}))
  }, [storagePath, subpath])

  return (
    <div
      className="bg-card overflow-auto border-2 shadow"
      style={{
        height: 'calc(100vh - 300px)',
      }}
    >
      <Table noWrapper className="text-base">
        <TableHeader className="bg-card sticky top-0 z-10">
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
              <FileTableContextMenu
                key={row.id}
                entry={row.original}
                storagePath={storagePath}
                subpath={subpath}
              >
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
              </FileTableContextMenu>
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

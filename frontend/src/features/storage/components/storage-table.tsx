import {
  getApiStorageListUnregisteredSourcesQueryKey,
  getApiStorageQueryKey,
  removeStorageMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { type Source } from '#/client/gen/index.ts'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table.tsx'
import { toast } from '#/components/ui/toast.tsx'
import { getPinningProps } from '#/lib/table-utils.ts'
import { cn } from '#/lib/utils.ts'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useTable } from '@tanstack/react-table'
import { columns } from './columns'
import { features } from './features'

type StorageTableProps = {
  data: Source[]
}

export const StorageTable = ({ data }: StorageTableProps) => {
  const queryClient = useQueryClient()

  const removeStorage = useMutation({
    ...removeStorageMutation(),
    onError: (error) =>
      toast.add({ title: error.title, description: error.errors?.at(0)?.message }),
    onSuccess: async () => {
      toast.add({ title: 'Success' })
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: getApiStorageQueryKey(),
        }),
        queryClient.invalidateQueries({
          queryKey: getApiStorageListUnregisteredSourcesQueryKey(),
        }),
      ])
    },
  })

  const table = useTable({
    features,
    columns,
    data: data,
    initialState: {
      columnPinning: {
        start: [],
        end: [columns[3].id!],
      },
    },
    meta: {
      handleRemove: async (subpath) => {
        await removeStorage.mutateAsync({ path: { subpath } })
      },
    },
  })

  return (
    <div
      className="bg-card shadow border-2 overflow-auto"
      style={{
        height: 'calc(100vh - 200px)',
      }}
    >
      <Table>
        <TableHeader className="top-0 z-10 sticky bg-card">
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                const { className, style } = getPinningProps(header.column)

                return (
                  <TableHead
                    key={header.id}
                    className={cn('w-auto font-bold text-base', className)}
                    style={style}
                  >
                    {header.isPlaceholder ? null : (
                      <span>
                        <table.FlexRender header={header} />
                      </span>
                    )}
                  </TableHead>
                )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows.length == 0 ? (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          ) : (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getAllCells().map((cell) => {
                  const { className, style } = getPinningProps(cell.column)
                  return (
                    <TableCell
                      key={cell.id}
                      className={cn('w-auto text-base', className)}
                      style={style}
                    >
                      <table.FlexRender cell={cell} />
                    </TableCell>
                  )
                })}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}

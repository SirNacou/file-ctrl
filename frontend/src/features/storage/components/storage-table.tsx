import { Table, TableHeader, TableRow } from '#/components/ui/table.tsx'
import { useTable } from '@tanstack/react-table'
import { columns } from './columns'
import { features } from './features'

export const StorageTable = () => {
  useTable({
    features,
    columns,
    data: [],
  })

  return (
    <div
      className="bg-card shadow border-2 overflow-auto"
      style={{
        height: 'calc(100vh - 200px)',
      }}
    >
      <Table>
        <TableHeader>
          <TableRow></TableRow>
        </TableHeader>
      </Table>
    </div>
  )
}

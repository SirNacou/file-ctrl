import type { Source } from '#/client/gen/index.ts'
import { createColumnHelper } from '@tanstack/react-table'
import type { Features } from './features'

const columnHelper = createColumnHelper<Features, Source>()

export const columns = columnHelper.columns([
  columnHelper.accessor('name', {
    header: () => <span>Name</span>,
  }),
  columnHelper.accessor('sub_path', {
    header: () => <span>Path</span>,
  }),
  columnHelper.accessor('created_at', {
    header: () => <span>Created At</span>,
  }),
])

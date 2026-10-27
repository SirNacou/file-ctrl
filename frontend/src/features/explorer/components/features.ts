import {
  columnSizingFeature,
  metaHelper,
  rowSelectionFeature,
  rowSortingFeature,
  tableFeatures,
} from '@tanstack/react-table'

import type { DirectoryEntry } from '#/client/gen/index.ts'

interface TableMeta {
  itemSelected: (item: DirectoryEntry) => void
}

export const features = tableFeatures({
  rowSelectionFeature,
  rowSortingFeature,
  columnSizingFeature,
  tableMeta: metaHelper<TableMeta>(),
})

export type Features = typeof features

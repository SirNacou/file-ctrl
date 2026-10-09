import type { FileItem } from '#/client/gen/index.ts'
import {
  columnSizingFeature,
  metaHelper,
  rowSelectionFeature,
  rowSortingFeature,
  tableFeatures,
} from '@tanstack/react-table'

interface TableMeta {
  itemSelected: (item: FileItem) => void
}

export const features = tableFeatures({
  rowSelectionFeature,
  rowSortingFeature,
  columnSizingFeature,
  tableMeta: metaHelper<TableMeta>(),
})

export type Features = typeof features

import {
  columnPinningFeature,
  columnSizingFeature,
  metaHelper,
  tableFeatures,
} from '@tanstack/react-table'

type TableMeta = {
  handleRemove: (subpath: string) => void
}

export const features = tableFeatures({
  columnPinningFeature,
  columnSizingFeature,
  tableMeta: metaHelper<TableMeta>(),
})

export type Features = typeof features

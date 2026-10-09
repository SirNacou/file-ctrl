import {
  columnSizingFeature,
  rowSelectionFeature,
  rowSortingFeature,
  tableFeatures,
} from "@tanstack/react-table";

export const features = tableFeatures({
  rowSelectionFeature,
  rowSortingFeature,
  columnSizingFeature,
});

export type Features = typeof features;

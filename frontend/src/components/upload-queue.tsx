import { type UploadItem as UploadItemData } from '@uploadcn/core'
import { UploadList } from '@uploadcn/react'
import { type VariantProps } from 'class-variance-authority'

import {
  UploadItem,
  UploadItemActions,
  UploadItemContent,
  UploadItemDescription,
  UploadItemMedia,
  UploadItemProgress,
  UploadItemTitle,
  uploadItemVariants,
} from './ui/upload'

export function UploadQueue({
  className,
  filter,
  variant = 'outline',
  size = 'default',
}: {
  className?: string
  filter?: (item: UploadItemData) => boolean
} & VariantProps<typeof uploadItemVariants>) {
  return (
    <UploadList className={className} filter={filter}>
      {(item) => (
        <UploadItem key={item.id} item={item} variant={variant} size={size}>
          <UploadItemMedia />
          <UploadItemContent>
            <UploadItemTitle>{item.file.webkitRelativePath}</UploadItemTitle>
            <UploadItemDescription />
          </UploadItemContent>
          <UploadItemActions />
          <UploadItemProgress />
        </UploadItem>
      )}
    </UploadList>
  )
}

'use client'

import {
  type UploadAdapter,
  type UploadPersistence,
  createIndexedDBPersistence,
} from '@uploadcn/core'
import {
  Upload as UploadPrimitive,
  type UseUploaderOptions,
  useNetworkStatus,
  useUploadContext,
  useUploadGuard,
  useUploadProgress,
  useUploader,
  useWindowDrop,
} from '@uploadcn/react'
import { cn } from 'cn'
import * as React from 'react'
import ChevronDownIcon from '~icons/lucide/chevron-down'
import CloudUploadIcon from '~icons/lucide/cloud-upload'
import WifiOffIcon from '~icons/lucide/wifi-off'
import XIcon from '~icons/lucide/x'

import { NativeButton, UploadProgress } from '@/components/ui/upload'

import { UploadQueue } from './upload-queue'

interface UploadProviderProps<TResult> extends Partial<UseUploaderOptions<TResult>> {
  adapter?: UploadAdapter<TResult>
  /**
   * Keep uploads across reloads in IndexedDB and resume them. `true` uses the
   * default store; pass your own `UploadPersistence` to customize.
   */
  persist?: boolean | UploadPersistence
  /** Warn before leaving the page while uploads run. Default `true`. */
  guard?: boolean
  children: React.ReactNode
}

/**
 * One uploader for the whole app. Any component inside can add files with
 * `useUploadContext().addFiles`, and `GlobalDropzone` / `UploadQueuePanel`
 * render app-wide UI for it, uploads keep running across client-side
 * navigation because the provider stays mounted.
 */
function UploadProvider<TResult = unknown>({
  adapter,
  persist = false,
  guard = true,
  children,
  ...options
}: UploadProviderProps<TResult>) {
  const [persistence] = React.useState(() =>
    persist === true ? createIndexedDBPersistence() : persist === false ? undefined : persist,
  )
  const uploader = useUploader<TResult>({
    adapter,
    persistence,
    restore: !!persistence,
    ...options,
  })
  useUploadGuard(uploader, { enabled: guard })
  return (
    <UploadPrimitive.Root uploader={uploader} className="contents">
      {children}
    </UploadPrimitive.Root>
  )
}

interface GlobalDropzoneProps {
  title?: React.ReactNode
  description?: React.ReactNode
  disabled?: boolean
  className?: string
}

/** A full-window overlay that appears when files are dragged over the page. */
function GlobalDropzone({
  title = 'Drop files anywhere to upload',
  description = "They'll keep uploading while you work",
  disabled = false,
  className,
}: GlobalDropzoneProps) {
  const { addFiles, disabled: uploadDisabled } = useUploadContext()
  const { isDragging } = useWindowDrop({
    onDrop: addFiles,
    enabled: !disabled && !uploadDisabled,
  })
  if (!isDragging) return null
  return (
    <div
      data-slot="global-dropzone"
      aria-hidden
      className={cn(
        'pointer-events-none fixed inset-0 z-50 flex animate-in items-center justify-center bg-background/80 p-6 backdrop-blur-sm duration-150 fade-in-0 motion-reduce:animate-none',
        className,
      )}
    >
      <div className="animate-in border-primary/60 zoom-in-95 flex size-full max-h-128 max-w-3xl flex-col items-center justify-center gap-3 rounded-2xl border-2 border-dashed text-center duration-200 motion-reduce:animate-none">
        <div className="bg-background flex size-12 items-center justify-center rounded-full border shadow-sm">
          <CloudUploadIcon strokeWidth={2} className="size-5" />
        </div>
        <p className="text-lg font-medium">{title}</p>
        <p className="text-muted-foreground text-sm">{description}</p>
      </div>
    </div>
  )
}

/** A collapsible panel showing every upload in the app, like a cloud drive. */
function UploadQueuePanel({ className }: { className?: string }) {
  const { uploader } = useUploadContext()
  const summary = useUploadProgress(uploader)
  const online = useNetworkStatus()
  const [collapsed, setCollapsed] = React.useState(false)
  const [dismissed, setDismissed] = React.useState(false)
  const panelId = React.useId()

  // A new upload brings the panel back.
  const [lastTotal, setLastTotal] = React.useState(summary.total)
  if (summary.total !== lastTotal) {
    setLastTotal(summary.total)
    if (summary.total > lastTotal) setDismissed(false)
  }

  if (summary.total === 0 || dismissed) return null

  const title = summary.isUploading
    ? `Uploading ${summary.counts.uploading + summary.counts.queued} file${summary.counts.uploading + summary.counts.queued === 1 ? '' : 's'}`
    : summary.counts.paused > 0
      ? `${summary.counts.paused} paused`
      : `${summary.counts.success} upload${summary.counts.success === 1 ? '' : 's'} complete`

  return (
    <section
      aria-label="Uploads"
      data-slot="upload-queue-panel"
      className={cn(
        'fixed inset-x-4 bottom-4 z-40 flex animate-in flex-col overflow-hidden rounded-xl border bg-popover text-popover-foreground shadow-lg duration-200 fade-in-0 slide-in-from-bottom-4 motion-reduce:animate-none sm:inset-x-auto sm:right-4 sm:w-96',
        className,
      )}
    >
      <header className="flex items-center gap-2 border-b px-3 py-2">
        <div className="flex min-w-0 flex-1 flex-col gap-1.5">
          <div className="flex items-center gap-2 text-sm font-medium">
            {!online ? (
              <WifiOffIcon
                strokeWidth={2}
                className="text-muted-foreground size-3.5"
                aria-label="Offline"
              />
            ) : null}
            <span className="truncate">{title}</span>
            {summary.isUploading ? (
              <span className="text-muted-foreground ms-auto text-xs font-normal tabular-nums">
                {Math.round(summary.percent)}%
              </span>
            ) : null}
          </div>
          {summary.isUploading ? <UploadProgress /> : null}
        </div>
        <NativeButton
          variant="ghost"
          size="icon-xs"
          aria-label={collapsed ? 'Expand uploads' : 'Collapse uploads'}
          aria-expanded={!collapsed}
          aria-controls={panelId}
          onClick={() => setCollapsed((value) => !value)}
        >
          <ChevronDownIcon
            strokeWidth={2}
            className={cn('transition-transform', collapsed && 'rotate-180')}
          />
        </NativeButton>
        {summary.isUploading ? null : (
          <NativeButton
            variant="ghost"
            size="icon-xs"
            aria-label="Close uploads"
            onClick={() => setDismissed(true)}
          >
            <XIcon strokeWidth={2} />
          </NativeButton>
        )}
      </header>
      <div id={panelId} hidden={collapsed} className="max-h-80 overflow-y-auto p-2">
        <UploadQueue variant="default" size="sm" />
      </div>
    </section>
  )
}

export { GlobalDropzone, UploadProvider, type UploadProviderProps, UploadQueuePanel }

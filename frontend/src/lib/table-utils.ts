// pinning.ts
import { cn } from '@/lib/utils'
import type { Column_ColumnPinning, Column_ColumnSizing } from '@tanstack/react-table'
import type { CSSProperties } from 'react'

export function getPinningProps(
  column: Column_ColumnPinning & Column_ColumnSizing,
  { header = false }: { header?: boolean } = {},
) {
  const pinned = column.getIsPinned()
  const isLastLeft = pinned === 'start'
  const isFirstRight = pinned === 'end'

  const style: CSSProperties = {
    width: column.getSize(),
    left: pinned === 'start' ? column.getStart('start') : undefined,
    right: pinned === 'end' ? column.getAfter('end') : undefined,
  }

  const className = cn(
    // pinned base: sticky + opaque so content doesn't bleed through
    pinned && 'sticky z-10 bg-background',
    // headers sit above body cells
    pinned && header && 'z-20',
    // keep the hover/selected tint on pinned cells (parent <tr> has `group/row`)
    pinned && !header && 'group-hover/row:bg-muted group-data-[state=selected]/row:bg-muted',
    // 1px edge line + soft shadow on the boundary with scrolling content
    isLastLeft &&
      'shadow-[inset_-1px_0_0_var(--border),6px_0_8px_-6px_rgb(0_0_0/0.15)] dark:shadow-[inset_-1px_0_0_var(--border),6px_0_8px_-6px_rgb(0_0_0/0.6)]',
    isFirstRight &&
      'shadow-[inset_1px_0_0_var(--border),-6px_0_8px_-6px_rgb(0_0_0/0.15)] dark:shadow-[inset_1px_0_0_var(--border),-6px_0_8px_-6px_rgb(0_0_0/0.6)]',
  )

  return { style, className }
}

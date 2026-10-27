'use client'

import { Toast as ToastPrimitive } from '@base-ui/react/toast'
import { cn } from 'cn'
import * as React from 'react'
import CircleCheckIcon from '~icons/lucide/circle-check'
import CircleXIcon from '~icons/lucide/circle-x'
import InfoIcon from '~icons/lucide/info'
import LoaderCircleIcon from '~icons/lucide/loader-circle'
import TriangleAlertIcon from '~icons/lucide/triangle-alert'
import XIcon from '~icons/lucide/x'

import { Button } from '@/components/ui/button'

export type ToastPosition =
  | 'top-left'
  | 'top-center'
  | 'top-right'
  | 'bottom-left'
  | 'bottom-center'
  | 'bottom-right'

const viewportPositionClasses: Record<ToastPosition, string> = {
  'top-left': 'top-4 sm:left-4 sm:right-auto [--y-dir:1]',
  'top-center': 'top-4 left-1/2 -translate-x-1/2 [--y-dir:1]',
  'top-right': 'top-4 sm:right-4 sm:left-auto [--y-dir:1]',
  'bottom-left': 'bottom-4 sm:left-4 sm:right-auto [--y-dir:-1]',
  'bottom-center': 'bottom-4 left-1/2 -translate-x-1/2 [--y-dir:-1]',
  'bottom-right': 'bottom-4 sm:right-4 sm:left-auto [--y-dir:-1]',
}

function getSwipeDirections(position: ToastPosition): ToastPrimitive.Root.Props['swipeDirection'] {
  const vertical = position.startsWith('top') ? 'up' : 'down'
  if (position.endsWith('left')) return [vertical, 'left']
  if (position.endsWith('right')) return [vertical, 'right']
  return [vertical, 'left', 'right']
}

const toast = ToastPrimitive.createToastManager()

function ToastProvider({ ...props }: ToastPrimitive.Provider.Props) {
  return <ToastPrimitive.Provider {...props} />
}

function ToastPortal({ ...props }: ToastPrimitive.Portal.Props) {
  return <ToastPrimitive.Portal data-slot="toast-portal" {...props} />
}

function ToastViewport({
  className,
  position = 'bottom-right',
  ...props
}: ToastPrimitive.Viewport.Props & { position?: ToastPosition }) {
  const isTop = position.startsWith('top')

  return (
    <ToastPrimitive.Viewport
      data-slot="toast-viewport"
      data-position={position}
      data-y-position={isTop ? 'top' : 'bottom'}
      className={cn(
        'group/viewport z-100 fixed inset-x-4 mx-auto sm:mx-0 outline-none w-auto sm:w-full max-w-sm pointer-events-none',
        viewportPositionClasses[position],
        className,
      )}
      {...props}
    />
  )
}

function Toast({ className, ...props }: ToastPrimitive.Root.Props) {
  return (
    <ToastPrimitive.Root
      data-slot="toast"
      className={cn(
        'group/toast right-0 left-0 z-[calc(1000-var(--toast-index))] absolute bg-popover shadow-lg border focus-visible:border-ring rounded-2xl outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 w-full text-popover-foreground pointer-events-auto select-none will-change-transform',
        'group-data-[y-position=bottom]/viewport:bottom-0 group-data-[y-position=bottom]/viewport:origin-bottom',
        'group-data-[y-position=top]/viewport:top-0 group-data-[y-position=top]/viewport:origin-top',
        '[--gap:0.75rem] [--peek:0.75rem]',
        '[--height:var(--toast-frontmost-height,var(--toast-height))]',
        '[--scale:calc(max(0,1-(var(--toast-index)*0.1)))]',
        '[--shrink:calc(1-var(--scale))]',
        '[--offset-y:calc(var(--toast-offset-y)*var(--y-dir)+calc(var(--toast-index)*var(--gap)*var(--y-dir))+var(--toast-swipe-movement-y))]',
        'h-(--height) [transform:translateX(var(--toast-swipe-movement-x))_translateY(calc(var(--toast-swipe-movement-y)+((var(--toast-index)*var(--peek))+(var(--shrink)*var(--height)))*var(--y-dir)))_scale(var(--scale))] [transition:transform_500ms_cubic-bezier(0.22,1,0.36,1),opacity_500ms,height_150ms]',
        "after:absolute after:left-0 after:h-[calc(var(--gap)+1px)] after:w-full after:content-['']",
        'group-data-[y-position=bottom]/viewport:after:top-full',
        'group-data-[y-position=top]/viewport:after:bottom-full',
        'data-expanded:h-(--toast-height) data-expanded:[transform:translateX(var(--toast-swipe-movement-x))_translateY(var(--offset-y))]',
        'data-limited:opacity-0 data-starting-style:[transform:translateY(calc(150%*var(--y-dir)*-1))]',
        '[&[data-ending-style]:not([data-limited]):not([data-swipe-direction])]:[transform:translateY(calc(150%*var(--y-dir)*-1))]',
        'data-ending-style:data-[swipe-direction=down]:[transform:translateY(calc(var(--toast-swipe-movement-y)+150%))]',
        'data-ending-style:data-[swipe-direction=up]:[transform:translateY(calc(var(--toast-swipe-movement-y)-150%))]',
        'data-ending-style:data-[swipe-direction=left]:[transform:translateX(calc(var(--toast-swipe-movement-x)-150%))_translateY(var(--offset-y))]',
        'data-ending-style:data-[swipe-direction=right]:[transform:translateX(calc(var(--toast-swipe-movement-x)+150%))_translateY(var(--offset-y))]',
        'data-expanded:data-ending-style:data-[swipe-direction=down]:[transform:translateY(calc(var(--toast-swipe-movement-y)+150%))]',
        'data-expanded:data-ending-style:data-[swipe-direction=up]:[transform:translateY(calc(var(--toast-swipe-movement-y)-150%))]',
        'data-expanded:data-ending-style:data-[swipe-direction=left]:[transform:translateX(calc(var(--toast-swipe-movement-x)-150%))_translateY(var(--offset-y))]',
        'data-expanded:data-ending-style:data-[swipe-direction=right]:[transform:translateX(calc(var(--toast-swipe-movement-x)+150%))_translateY(var(--offset-y))]',
        className,
      )}
      {...props}
    />
  )
}

function ToastContent({ className, ...props }: ToastPrimitive.Content.Props) {
  return (
    <ToastPrimitive.Content
      data-slot="toast-content"
      className={cn(
        'flex items-center gap-3 data-behind:opacity-0 data-expanded:opacity-100 p-4 h-full overflow-hidden transition-opacity duration-250 ease-[cubic-bezier(0.22,1,0.36,1)]',
        className,
      )}
      {...props}
    />
  )
}

function ToastTitle({ className, ...props }: ToastPrimitive.Title.Props) {
  return (
    <ToastPrimitive.Title
      data-slot="toast-title"
      className={cn('font-medium text-sm', className)}
      {...props}
    />
  )
}

function ToastDescription({ className, ...props }: ToastPrimitive.Description.Props) {
  return (
    <ToastPrimitive.Description
      data-slot="toast-description"
      className={cn('text-muted-foreground text-sm', className)}
      {...props}
    />
  )
}

function ToastAction({
  className,
  render = <Button variant="outline" size="sm" />,
  ...props
}: ToastPrimitive.Action.Props) {
  return (
    <ToastPrimitive.Action
      data-slot="toast-action"
      render={render}
      className={cn('shrink-0', className)}
      {...props}
    />
  )
}

function ToastClose({
  className,
  children,
  render = <Button variant="ghost" size="icon-sm" />,
  ...props
}: ToastPrimitive.Close.Props) {
  return (
    <ToastPrimitive.Close
      data-slot="toast-close"
      aria-label="Close toast"
      render={render}
      className={cn(
        "after:absolute relative after:-inset-2 text-muted-foreground hover:text-foreground after:content-[''] shrink-0",
        className,
      )}
      {...props}
    >
      {children ?? <XIcon strokeWidth={2} aria-hidden="true" />}
    </ToastPrimitive.Close>
  )
}

function ToastIcon({ type }: { type: string | undefined }) {
  let icon: React.ReactNode = null

  if (type === 'success') {
    icon = <CircleCheckIcon strokeWidth={2} aria-hidden="true" />
  }

  if (type === 'info') {
    icon = <InfoIcon strokeWidth={2} aria-hidden="true" />
  }

  if (type === 'warning') {
    icon = <TriangleAlertIcon strokeWidth={2} aria-hidden="true" />
  }

  if (type === 'error') {
    icon = <CircleXIcon strokeWidth={2} className="text-destructive" aria-hidden="true" />
  }

  if (type === 'loading') {
    icon = <LoaderCircleIcon strokeWidth={2} className="animate-spin" aria-hidden="true" />
  }

  if (!icon) {
    return null
  }

  return (
    <span
      data-slot="toast-icon"
      className="shrink-0 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
    >
      {icon}
    </span>
  )
}

function ToastList({ position = 'bottom-right' }: { position?: ToastPosition }) {
  const { toasts } = ToastPrimitive.useToastManager()
  const swipeDirection = getSwipeDirections(position)

  return toasts.map((toastItem) => (
    <Toast key={toastItem.id} toast={toastItem} swipeDirection={swipeDirection}>
      <ToastContent>
        <ToastIcon type={toastItem.type} />
        <div className="flex min-w-0 flex-1 flex-col gap-1">
          <ToastTitle />
          <ToastDescription />
        </div>
        <ToastAction />
        <ToastClose />
      </ToastContent>
    </Toast>
  ))
}

interface ToasterProps extends ToastPrimitive.Provider.Props {
  position?: ToastPosition
}

function Toaster({
  children,
  toastManager = toast,
  position = 'bottom-right',
  ...props
}: ToasterProps) {
  return (
    <ToastProvider toastManager={toastManager} {...props}>
      {children}
      <ToastPortal>
        <ToastViewport position={position}>
          <ToastList position={position} />
        </ToastViewport>
      </ToastPortal>
    </ToastProvider>
  )
}

const createToastManager = ToastPrimitive.createToastManager
const useToastManager = ToastPrimitive.useToastManager

export {
  createToastManager,
  Toast,
  toast,
  ToastAction,
  ToastClose,
  ToastContent,
  ToastDescription,
  Toaster,
  ToastPortal,
  ToastProvider,
  ToastTitle,
  ToastViewport,
  useToastManager,
  type ToasterProps,
}

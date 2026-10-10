'use client'

import { cn } from 'cn'
import * as React from 'react'

const RawTable = ({ className, ...props }: React.ComponentProps<'table'>) => {
  return (
    <table
      data-slot="table"
      className={cn('w-full text-sm caption-bottom', className)}
      {...props}
    />
  )
}

function Table({
  className,
  noWrapper,
  divClassname,
  ...props
}: React.ComponentProps<'table'> & {
  noWrapper?: boolean
  divClassname?: string
}) {
  if (noWrapper) {
    return <RawTable className={className} {...props} />
  }

  return (
    <div
      data-slot="table-container"
      className={cn('relative w-full overflow-x-auto', divClassname)}
    >
      <RawTable className={className} {...props} />;
    </div>
  )
}

function TableHeader({ className, ...props }: React.ComponentProps<'thead'>) {
  return <thead data-slot="table-header" className={cn('[&_tr]:border-b', className)} {...props} />
}

function TableBody({ className, ...props }: React.ComponentProps<'tbody'>) {
  return (
    <tbody
      data-slot="table-body"
      className={cn('[&_tr:last-child]:border-0', className)}
      {...props}
    />
  )
}

function TableFooter({ className, ...props }: React.ComponentProps<'tfoot'>) {
  return (
    <tfoot
      data-slot="table-footer"
      className={cn('bg-muted/50 border-t [&>tr]:last:border-b-0 font-medium', className)}
      {...props}
    />
  )
}

function TableRow({ className, ...props }: React.ComponentProps<'tr'>) {
  return (
    <tr
      data-slot="table-row"
      className={cn(
        'data-[state=selected]:bg-muted has-aria-expanded:bg-muted/50 hover:bg-muted/50 border-b transition-colors',
        className,
      )}
      {...props}
    />
  )
}

function TableHead({ className, ...props }: React.ComponentProps<'th'>) {
  return (
    <th
      data-slot="table-head"
      className={cn(
        'h-12 px-3 text-left align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pr-0',
        className,
      )}
      {...props}
    />
  )
}

function TableCell({ className, ...props }: React.ComponentProps<'td'>) {
  return (
    <td
      data-slot="table-cell"
      className={cn('p-3 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0', className)}
      {...props}
    />
  )
}

function TableCaption({ className, ...props }: React.ComponentProps<'caption'>) {
  return (
    <caption
      data-slot="table-caption"
      className={cn('mt-4 text-muted-foreground text-sm', className)}
      {...props}
    />
  )
}

export { Table, TableBody, TableCaption, TableCell, TableFooter, TableHead, TableHeader, TableRow }

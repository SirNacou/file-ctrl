import type { ReactNode } from 'react'

type PageHeaderProps = {
  title: ReactNode
  description?: string
  actions?: ReactNode
}

export function PageHeader({ title, description, actions }: PageHeaderProps) {
  return (
    <header className="flex items-center justify-between gap-4">
      <div className="space-y-1">
        <h1 className="text-3xl font-semibold">{title}</h1>

        {description && <p className="text-muted-foreground text-sm">{description}</p>}
      </div>

      {actions && <div className="flex items-center gap-2">{actions}</div>}
    </header>
  )
}

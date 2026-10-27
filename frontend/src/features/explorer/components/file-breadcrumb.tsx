import { Link, useParams } from '@tanstack/react-router'
import { Fragment, useMemo } from 'react'

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '#/components/ui/breadcrumb.tsx'

export function FileBreadcrumbs() {
  const { path, _splat: current_path } = useParams({ from: '/_protected/storage/$path/$/' })

  const breadcrumbItems = useMemo(() => getPathItems(current_path), [current_path])

  return (
    <Breadcrumb>
      <BreadcrumbList className="text-base">
        <BreadcrumbItem>
          <BreadcrumbLink
            render={
              <Link to="/storage/$path/$" params={{ path, _splat: undefined }}>
                {path}
              </Link>
            }
          />
        </BreadcrumbItem>
        {breadcrumbItems.map((item) => (
          <Fragment key={item.path}>
            <BreadcrumbSeparator />
            <BreadcrumbItem>
              <BreadcrumbLink
                render={
                  <Link to="/storage/$path/$" params={{ path, _splat: item.path }}>
                    {item.name}
                  </Link>
                }
              />
            </BreadcrumbItem>
          </Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  )
}

type PathItem = {
  name: string
  path: string
}

function getPathItems(path: string | undefined): PathItem[] {
  if (path === undefined) return []

  const segments = path.split('/').filter(Boolean)

  return segments.map((name, index) => ({
    name,
    path: segments.slice(0, index + 1).join('/'),
  }))
}

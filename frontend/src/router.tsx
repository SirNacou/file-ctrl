import { createRouter as createTanStackRouter, ErrorComponent } from '@tanstack/react-router'
import { setupRouterSsrQueryIntegration } from '@tanstack/react-router-ssr-query'

import { getContext } from './integrations/tanstack-query/root-provider'
import { Loading } from './routes/-components/loading'
import { NotFound } from './routes/-components/not-found'
import { routeTree } from './routeTree.gen'

export function getRouter() {
  const context = getContext()

  const router = createTanStackRouter({
    routeTree,
    context,
    scrollRestoration: true,
    defaultPreload: 'intent',
    defaultPreloadStaleTime: 0,

    defaultPendingComponent: Loading,
    defaultPendingMs: 0,
    defaultPendingMinMs: 300,
    defaultNotFoundComponent: NotFound,
    defaultErrorComponent: ErrorComponent,
  })

  setupRouterSsrQueryIntegration({ router, queryClient: context.queryClient })

  return router
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}

import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

import { authGetStatusOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { ThemeToggle } from '#/components/theme-toggle.tsx'
import { AppSidebar } from '#/routes/_protected/-components/app-sidebar.tsx'
import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'

export const Route = createFileRoute('/_protected')({
  beforeLoad: async ({ context }) => {
    const status = await context.queryClient.query({
      ...authGetStatusOptions(),
      staleTime: 60 * 1_000,
    })

    if (status.needs_setup) {
      throw redirect({ to: '/setup' })
    }
    if (!status.authenticated) {
      throw redirect({ to: '/login' })
    }

    return {
      user: {
        ...status,
      },
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <SidebarProvider>
      <AppSidebar />
      <div className="flex w-full flex-col gap-6 px-6 py-3">
        <header className="flex justify-between">
          <SidebarTrigger className="md:hidden" />
          <ThemeToggle className="ml-auto" />
        </header>
        <main className="flex flex-1 flex-col gap-4">
          <Outlet />
        </main>
      </div>
    </SidebarProvider>
  )
}

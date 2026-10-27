import { createFileRoute, Outlet } from '@tanstack/react-router'

import { AppSidebar } from '#/components/app-sidebar.tsx'
import { ThemeToggle } from '#/components/theme-toggle.tsx'
import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'

export const Route = createFileRoute('/_protected')({
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

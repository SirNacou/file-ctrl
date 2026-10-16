import { AppSidebar } from '#/components/app-sidebar.tsx'
import { ThemeToggle } from '#/components/theme-toggle.tsx'
import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <SidebarProvider>
      <AppSidebar />
      <div className="flex flex-col gap-6 w-full py-3 px-6">
        <header className="flex justify-between">
          <SidebarTrigger className="md:hidden" />
          <ThemeToggle className="ml-auto" />
        </header>
        <main className="flex-1 flex flex-col gap-4">
          <Outlet />
        </main>
      </div>
    </SidebarProvider>
  )
}

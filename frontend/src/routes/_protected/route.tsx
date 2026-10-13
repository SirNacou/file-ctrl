import { AppSidebar } from '#/components/app-sidebar.tsx'
import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <SidebarProvider
      style={
        {
          '--sidebar-width': '20rem',
        } as React.CSSProperties
      }
    >
      <AppSidebar />
      <div className="flex flex-col px-4 md:px-6 py-2 md:py-4 w-full">
        <header>
          <SidebarTrigger className="md:hidden" />
        </header>
        <main className="flex-1">
          <Outlet />
        </main>
      </div>
    </SidebarProvider>
  )
}

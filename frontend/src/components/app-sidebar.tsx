import { Link, useMatchRoute, type LinkProps } from '@tanstack/react-router'

import { ThemeToggle } from './theme-toggle'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from './ui/sidebar'

type NavItem = {
  label: string
} & LinkProps

type NavGroup = {
  label?: string
  navItems: NavItem[]
}

const navGroups: NavGroup[] = [
  {
    navItems: [
      {
        to: '/storage',
        label: 'Storage',
      },
    ],
  },
]

export const AppSidebar = () => {
  const matchRoute = useMatchRoute()

  return (
    <Sidebar>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem className="self-center text-2xl font-bold">File Ctrl</SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          {navGroups.map(({ label, navItems }, i) => (
            <SidebarGroup key={i}>
              {label && <SidebarGroupLabel>{label}</SidebarGroupLabel>}
              <SidebarGroupContent>
                <SidebarMenu>
                  {navItems.map(({ label, to, ...props }) => (
                    <SidebarMenuItem key={label}>
                      <SidebarMenuButton
                        isActive={!!matchRoute({ to, fuzzy: true })}
                        render={
                          <Link className="text-base" to={to} {...props}>
                            {label}
                          </Link>
                        }
                      />
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          ))}
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <ThemeToggle />
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}

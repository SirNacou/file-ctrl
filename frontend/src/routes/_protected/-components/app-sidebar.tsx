import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Link, useMatchRoute, useNavigate, type LinkProps } from '@tanstack/react-router'
import LogoutIcon from '~icons/lucide/log-out'

import {
  authGetStatusQueryKey,
  authLogoutMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { Button } from '#/components/ui/button.tsx'

import { ThemeToggle } from '../../../components/theme-toggle'
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
} from '../../../components/ui/sidebar'

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
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const logoutMutation = useMutation({
    ...authLogoutMutation(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: authGetStatusQueryKey(),
      })
      await navigate({
        to: '/login',
      })
    },
  })

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
            <Button size={'icon'} onClick={() => logoutMutation.mutate({})}>
              <LogoutIcon />
            </Button>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}

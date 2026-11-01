import { createFileRoute, redirect } from '@tanstack/react-router'

import { authGetStatusOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { LoginForm } from '#/features/auth/components/login-form.tsx'

export const Route = createFileRoute('/login')({
  beforeLoad: async ({ context }) => {
    const status = await context.queryClient.query(authGetStatusOptions())
    if (status.needs_setup) {
      throw redirect({ to: '/setup' })
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = Route.useNavigate()
  return <LoginForm onSuccess={() => navigate({ to: '/' })} />
}

import { createFileRoute, redirect } from '@tanstack/react-router'

import { authGetStatusOptions } from '#/client/gen/@tanstack/react-query.gen.ts'
import { SetupWizardForm } from '#/features/auth/components/setup-wizard-form.tsx'

export const Route = createFileRoute('/setup')({
  beforeLoad: async ({ context }) => {
    const status = await context.queryClient.query(authGetStatusOptions())
    if (!status.needs_setup) {
      throw redirect({ to: status.authenticated ? '/' : '/login' })
    }

    return { oidcEnabled: status.oidc_enabled }
  },
  component: RouteComponent,
})

function RouteComponent() {
  const navigate = Route.useNavigate()
  const { oidcEnabled } = Route.useRouteContext()
  return <SetupWizardForm oidcEnabled={oidcEnabled} onSuccess={() => navigate({ to: '/' })} />
}

// features/auth/components/oidc-login-button.tsx
import LogIn from '~icons/lucide/log-in'

import type { AuthOidcLoginData } from '#/client/gen/types.gen.ts'
import { Button } from '#/components/ui/button.tsx'
import { env } from '#/config/env.ts'

const oidcLoginPath: AuthOidcLoginData['url'] = '/api/auth/oidc/login'

export const oidcLoginUrl = `${env.PUBLIC_APP_URL}${oidcLoginPath}`

type OidcLoginButtonProps = {
  href: string
  label?: string
}

export function OidcLoginButton({ href, label = 'Continue with SSO' }: OidcLoginButtonProps) {
  return (
    <Button
      nativeButton={false}
      render={
        <a href={href}>
          <LogIn />
          {label}
        </a>
      }
      variant="outline"
      className="w-full"
    />
  )
}

import { Link, useRouter, type ErrorComponentProps } from '@tanstack/react-router'
import { TriangleAlert } from 'lucide-react'

import { zErrorModel, type ErrorModelZodType } from '#/client/gen/zod.gen.ts'
import { Button } from '@/components/ui/button'
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty'

export function ErrorComponent({ error, reset }: ErrorComponentProps) {
  const router = useRouter()

  const handleRetry = async () => {
    // Re-run loaders, then reset the error boundary
    await router.invalidate()
    reset()
  }
  const { title, detail, errors } = getErrorInfo(error)

  return (
    <Empty className="min-h-svh">
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <TriangleAlert />
        </EmptyMedia>
        <EmptyTitle className="text-xl">{title ?? 'Something went wrong'}</EmptyTitle>
        <EmptyDescription className="text-lg">
          {detail ??
            'An unexpected error occurred. You can try again, or head back to the home page.'}
        </EmptyDescription>
      </EmptyHeader>

      <EmptyContent>
        {import.meta.env.DEV && errors.length > 0 && (
          <ul className="bg-muted text-muted-foreground max-h-40 w-full max-w-md space-y-1 overflow-auto rounded-md p-3 text-left text-base">
            {errors.map((e, i) => (
              <li key={i}>
                {e.location && <span className="font-mono">{e.location}: </span>}
                {e.message}
              </li>
            ))}
          </ul>
        )}
        <div className="flex justify-center gap-2">
          <Button onClick={handleRetry}>Try again</Button>
          <Button nativeButton={false} variant="outline" render={<Link to="/">Go to home</Link>} />
        </div>
      </EmptyContent>
    </Empty>
  )
}

export default ErrorComponent

type ErrorInfo = {
  title?: string
  detail?: string
  errors: { location?: string; message?: string }[]
}

function getErrorInfo(error: unknown): ErrorInfo {
  if (error instanceof Error) {
    return { detail: error.message, errors: [] }
  }

  const parsed = zErrorModel.safeParse(error)
  if (parsed.success) {
    const { title, detail, errors } = parsed.data as ErrorModelZodType
    if (title || detail || errors?.length) {
      return { title, detail, errors: errors ?? [] }
    }
  }

  return { detail: 'Unknown error', errors: [] }
}

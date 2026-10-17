import { Button } from '@/components/ui/button'
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty'
import { Link, useRouter } from '@tanstack/react-router'
import FileQuestionMarkIcon from '~icons/lucide/file-question-mark'

export const NotFound = () => {
  const router = useRouter()
  return (
    <Empty className="min-h-svh">
      <EmptyHeader>
        <EmptyMedia className="size-15 rounded-full" variant="icon">
          <FileQuestionMarkIcon className="size-8" />
        </EmptyMedia>
        <EmptyTitle className="text-2xl">Page not found</EmptyTitle>
        <EmptyDescription className="text-md">
          We couldn't find the page you're looking for. It may have moved, or the link may be
          incorrect.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent className="flex-row justify-center">
        <Button nativeButton={false} render={<Link to="/">Go to home</Link>} />
        {router.history.canGoBack() && (
          <Button variant="outline" onClick={() => router.history.back()}>
            Go back
          </Button>
        )}
      </EmptyContent>
    </Empty>
  )
}

export default NotFound

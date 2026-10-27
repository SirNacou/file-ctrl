import LoaderCircleIcon from '~icons/lucide/loader-circle'

type Props = {
  message?: string
}
export const Loading = ({ message = 'Loading...' }: Props) => {
  return (
    <div
      className="flex h-screen flex-col items-center justify-center gap-3"
      role="status"
      aria-live="polite"
    >
      <LoaderCircleIcon className="text-muted-foreground size-6 animate-spin" aria-hidden="true" />

      <span className="text-muted-foreground text-sm">{message}</span>
    </div>
  )
}

import {
  listStorageSourcesQueryKey,
  removeStorageMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import { type Source } from '#/client/gen/index.ts'
import { ConfirmDeleteDialog } from '#/components/confirm-detele-dialog.tsx'
import {
  Card,
  CardAction,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '#/components/ui/card.tsx'
import { toast } from '#/components/ui/toast.tsx'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'
import ArchiveIcon from '~icons/lucide/archive'
import { RenameStorageDialog } from './rename-storage-dialog'

type StorageCardProps = {
  storage: Source
  onEdit: () => void
  onRemove: () => void
}

export function StorageCard({ storage }: StorageCardProps) {
  const created = new Date(storage.created_at).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })

  const queryClient = useQueryClient()

  const removeStorage = useMutation({
    ...removeStorageMutation(),
    onError: (error) => {
      toast.add({ type: 'error', title: error.title, description: error.errors?.at(0)?.message })
    },
    onSuccess: async () => {
      toast.add({ type: 'success', title: 'Remove successfully.' })
      await queryClient.invalidateQueries({
        queryKey: listStorageSourcesQueryKey(),
      })
    },
  })

  async function handleRemove() {
    await removeStorage.mutateAsync({ path: { path: storage.path } })
  }

  return (
    <Card className="relative transition hover:border-muted-foreground/40 hover:shadow-md active:bg-muted">
      <CardHeader>
        <CardTitle>
          <Link
            to="/storage/$path/$"
            params={{ path: storage.path }}
            className="after:absolute after:inset-0 after:rounded-[inherit] after:content-[''] focus-visible:outline-none focus-visible:after:ring-2 focus-visible:after:ring-primary"
          >
            {storage.name}
          </Link>
        </CardTitle>
        <CardDescription className="inline-flex gap-1">
          <ArchiveIcon />
          {storage.path}
        </CardDescription>
        <CardAction className="relative z-10 flex gap-2">
          <RenameStorageDialog storage={storage} />
          <ConfirmDeleteDialog buttonText="Remove" onConfirm={handleRemove} />
        </CardAction>
      </CardHeader>
      <CardFooter>{created}</CardFooter>
    </Card>
  )
}

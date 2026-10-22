import {
  listDirectoryQueryKey,
  storageDeleteFilesMutation,
} from '#/client/gen/@tanstack/react-query.gen.ts'
import type { DirectoryEntry } from '#/client/gen/index.ts'
import { ConfirmDeleteDialog } from '#/components/confirm-detele-dialog.tsx'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuGroup,
  ContextMenuItem,
  ContextMenuLabel,
  ContextMenuTrigger,
} from '#/components/ui/context-menu.tsx'
import { notify } from '#/lib/notify.ts'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState, type ReactElement } from 'react'
import DownloadIcon from '~icons/lucide/download'
import TextCursorInputIcon from '~icons/lucide/text-cursor-input'
import TrashIcon from '~icons/lucide/trash'
import { ArchiveDownloadLink, FileDownloadLink } from './download-link'
import { RenameFileDialog } from './rename-file-dialog'

type Props = {
  children: ReactElement
  entry: DirectoryEntry
  storagePath: string
  subpath?: string
}

type Dialog = 'rename' | 'delete' | null

export const FileTableContextMenu = ({ children, entry, storagePath, subpath = '' }: Props) => {
  const [dialog, setDialog] = useState<Dialog>(null)

  const queryClient = useQueryClient()
  const deleteFiles = useMutation({
    ...storageDeleteFilesMutation(),
    onError: (error) => {
      notify.error(error.title, error.errors?.at(0)?.message)
    },
    onSuccess: async () => {
      notify.success('Delete successfully!')
      await queryClient.invalidateQueries({
        queryKey: listDirectoryQueryKey({ path: { storage_path: storagePath } }),
      })
    },
  })

  async function handleDelete() {
    await deleteFiles.mutateAsync({
      path: { storage_path: storagePath },
      body: {
        file_paths: [entry.path],
      },
    })
  }

  return (
    <>
      <ContextMenu>
        <ContextMenuTrigger render={children} />
        <ContextMenuContent>
          <ContextMenuGroup>
            <ContextMenuLabel>{entry.name}</ContextMenuLabel>
            <ContextMenuItem
              render={
                entry.is_dir ? (
                  <ArchiveDownloadLink storagePath={storagePath} filePaths={[entry.path]} />
                ) : (
                  <FileDownloadLink storagePath={storagePath} filePath={entry.path} />
                )
              }
            >
              <DownloadIcon />
              Download
            </ContextMenuItem>
            <ContextMenuItem variant="default" onClick={() => setDialog('rename')}>
              <TextCursorInputIcon />
              Rename
            </ContextMenuItem>
            <ContextMenuItem variant="destructive" onClick={() => setDialog('delete')}>
              <TrashIcon />
              Delete
            </ContextMenuItem>
          </ContextMenuGroup>
        </ContextMenuContent>
      </ContextMenu>
      <RenameFileDialog
        open={dialog === 'rename'}
        openChanged={(open) => setDialog(open ? 'rename' : null)}
        entry={entry}
        storagePath={storagePath}
        subpath={subpath}
      />
      <ConfirmDeleteDialog
        open={dialog === 'delete'}
        onOpenChanged={(open) => setDialog(open ? 'delete' : null)}
        onConfirm={handleDelete}
      />
    </>
  )
}

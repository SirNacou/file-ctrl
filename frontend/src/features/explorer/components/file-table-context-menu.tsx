import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '#/components/ui/context-menu.tsx'
import { useState, type ReactElement } from 'react'
import { RenameFileDialog } from './rename-file-dialog'

type Props = {
  children: ReactElement
  fileName: string
  storagePath: string
  subpath?: string
}
export const FileTableContextMenu = ({ children, fileName, storagePath, subpath }: Props) => {
  const [open, setOpen] = useState(false)

  return (
    <>
      <ContextMenu>
        <ContextMenuTrigger render={children} />
        <ContextMenuContent>
          <ContextMenuItem onPointerDown={() => setOpen(true)}>Rename</ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
      <RenameFileDialog
        open={open}
        openChanged={setOpen}
        fileName={fileName}
        storagePath={storagePath}
        subpath={subpath}
      />
    </>
  )
}

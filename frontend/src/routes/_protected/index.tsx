import { FileTable } from '#/features/explorer/components/file-table.tsx'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/')({
  component: Home,
})

const items = Array.from({ length: 1_0 }).flatMap(() => [
  {
    name: 'Test',
    is_dir: true,
    mod_time: new Date().toISOString(),
    path: '/',
    size: 0,
  },
  {
    name: 'file.txt',
    is_dir: false,
    mod_time: new Date().toISOString(),
    path: '/',
    size: 100,
  },
])

function Home() {
  return (
    <div className="h-full">
      <FileTable data={items} />
    </div>
  )
}

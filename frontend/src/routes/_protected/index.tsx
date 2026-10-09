import type { FileItem } from '#/client/gen/index.ts'
import { FileTable } from '#/features/explorer/components/file-table.tsx'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/')({
  loader: ({ context: { queryClient } }) => [],
  component: Home,
})

const items: FileItem[] = Array.from({ length: 1_000 }).flatMap(() => [
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
  const data = Route.useLoaderData()
  return (
    <div className="h-full">
      <FileTable data={[]} />
    </div>
  )
}

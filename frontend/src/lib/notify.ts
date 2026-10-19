import { toast } from '@/components/ui/toast'

type Opts = Omit<Parameters<typeof toast.add>[0], 'type' | 'title'>

export const notify = {
  success: (title: string, description?: string, opts?: Opts) =>
    toast.add({ type: 'success', title, description, timeout: 3000, ...opts }),

  error: (title?: string, description?: string, opts?: Opts) =>
    toast.add({
      type: 'error',
      title: title || 'Something went wrong',
      description,
      timeout: 5000,
      ...opts,
    }),

  info: (title: string, description?: string, opts?: Opts) =>
    toast.add({ type: 'info', title, description, ...opts }),

  warning: (title: string, description?: string, opts?: Opts) =>
    toast.add({ type: 'warning', title, description, ...opts }),
}

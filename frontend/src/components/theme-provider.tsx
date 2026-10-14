import { ScriptOnce } from '@tanstack/react-router'
import { createContext, useCallback, useContext, useEffect, useSyncExternalStore } from 'react'

export type Theme = 'dark' | 'light' | 'system'

type ThemeProviderProps = {
  children: React.ReactNode
  defaultTheme?: Theme
  storageKey?: string
}

type ThemeProviderState = {
  theme: Theme
  setTheme: (theme: Theme) => void
}

const isTheme = (v: unknown): v is Theme => v === 'light' || v === 'dark' || v === 'system'

function getThemeScript(storageKey: string, defaultTheme: Theme) {
  const key = JSON.stringify(storageKey)
  const fallback = JSON.stringify(defaultTheme)

  return `(function(){try{var t=localStorage.getItem(${key});if(t!=='light'&&t!=='dark'&&t!=='system'){t=${fallback}}var d=matchMedia('(prefers-color-scheme: dark)').matches;var r=t==='system'?(d?'dark':'light'):t;var e=document.documentElement;e.classList.add(r);e.style.colorScheme=r}catch(e){}})();`
}

const ThemeProviderContext = createContext<ThemeProviderState>({
  theme: 'system',
  setTheme: () => {},
})

function applyTheme(theme: Theme) {
  const root = document.documentElement
  root.classList.remove('light', 'dark')

  const resolved =
    theme === 'system'
      ? window.matchMedia('(prefers-color-scheme: dark)').matches
        ? 'dark'
        : 'light'
      : theme

  root.classList.add(resolved)
  root.style.colorScheme = resolved
}

const listeners = new Set<() => void>()
const notify = () => listeners.forEach((l) => l())

function readStoredTheme(storageKey: string, defaultTheme: Theme): Theme {
  try {
    const stored = localStorage.getItem(storageKey)
    return isTheme(stored) ? stored : defaultTheme
  } catch {
    return defaultTheme
  }
}

export function ThemeProvider({
  children,
  defaultTheme = 'system',
  storageKey = 'theme',
}: ThemeProviderProps) {
  const subscribe = useCallback(
    (callback: () => void) => {
      const onStorage = (e: StorageEvent) => {
        if (e.key === storageKey || e.key === null) callback()
      }
      window.addEventListener('storage', onStorage)
      listeners.add(callback)
      return () => {
        window.removeEventListener('storage', onStorage)
        listeners.delete(callback)
      }
    },
    [storageKey],
  )

  const theme = useSyncExternalStore(
    subscribe,
    () => readStoredTheme(storageKey, defaultTheme), // client
    () => defaultTheme, // server / hydration
  )
  useEffect(() => {
    applyTheme(theme)
  }, [theme])

  useEffect(() => {
    if (theme !== 'system') return

    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => applyTheme('system')
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [theme])

  const setTheme = useCallback(
    (next: Theme) => {
      try {
        localStorage.setItem(storageKey, next)
      } catch {}
      notify()
    },
    [storageKey],
  )

  return (
    <ThemeProviderContext value={{ theme, setTheme }}>
      <ScriptOnce>{getThemeScript(storageKey, defaultTheme)}</ScriptOnce>
      {children}
    </ThemeProviderContext>
  )
}

export function useTheme() {
  const context = useContext(ThemeProviderContext)
  if (context === undefined) throw new Error('useTheme must be used within a ThemeProvider')
  return context
}

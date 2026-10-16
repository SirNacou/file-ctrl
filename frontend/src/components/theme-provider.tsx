import { ScriptOnce } from '@tanstack/react-router'
import { useAtomValue, useSetAtom } from 'jotai'
import { atomWithStorage } from 'jotai/utils'
import { useEffect } from 'react'

export type Theme = 'dark' | 'light' | 'system'

const STORAGE_KEY = 'theme'
const DEFAULT_THEME: Theme = 'system'

const isTheme = (v: unknown): v is Theme => v === 'light' || v === 'dark' || v === 'system'

// Plain-string storage so the inline script can read the raw value
const themeStorage = {
  getItem: (key: string, initial: Theme) => {
    try {
      const v = localStorage.getItem(key)
      return isTheme(v) ? v : initial
    } catch {
      return initial
    }
  },
  setItem: (key: string, value: Theme) => {
    try {
      localStorage.setItem(key, value)
    } catch {}
  },
  removeItem: (key: string) => {
    try {
      localStorage.removeItem(key)
    } catch {}
  },
  subscribe: (key: string, callback: (value: Theme) => void) => {
    const onStorage = (e: StorageEvent) => {
      if (e.key === key && isTheme(e.newValue)) callback(e.newValue)
    }
    window.addEventListener('storage', onStorage)
    return () => window.removeEventListener('storage', onStorage)
  },
}

const themeAtom = atomWithStorage<Theme>(STORAGE_KEY, DEFAULT_THEME, themeStorage)

const themeScript = `(function(){try{var t=localStorage.getItem(${JSON.stringify(STORAGE_KEY)});if(t!=='light'&&t!=='dark'&&t!=='system')t=${JSON.stringify(DEFAULT_THEME)};var r=t==='system'?(matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'):t;var e=document.documentElement;e.classList.add(r);e.style.colorScheme=r}catch(e){}})()`

function applyTheme(theme: Theme) {
  const dark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const resolved = theme === 'system' ? (dark ? 'dark' : 'light') : theme
  const root = document.documentElement
  root.classList.remove('light', 'dark')
  root.classList.add(resolved)
  root.style.colorScheme = resolved
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const theme = useAtomValue(themeAtom)

  useEffect(() => {
    applyTheme(theme)
    if (theme !== 'system') return
    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => applyTheme('system')
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [theme])

  return (
    <>
      <ScriptOnce>{themeScript}</ScriptOnce>
      {children}
    </>
  )
}

export function useTheme() {
  return { theme: useAtomValue(themeAtom), setTheme: useSetAtom(themeAtom) }
}

import { acceptHMRUpdate, defineStore } from 'pinia'

type ThemeMode = 'light' | 'dark'

const STORAGE_KEY = 'theme'

function applyToDocument(mode: ThemeMode) {
  document.documentElement.classList.toggle('dark', mode === 'dark')
}

function getInitialTheme(): ThemeMode {
  const stored = localStorage.getItem(STORAGE_KEY)

  if (stored === 'light' || stored === 'dark') {
    return stored
  }

  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export const useThemeStore = defineStore('theme', {
  state: () => ({
    mode: getInitialTheme() as ThemeMode,
  }),

  actions: {
    // Applies the current mode to <html>. Call once on app startup, before
    // the first paint, so there's no flash of the wrong theme.
    init() {
      applyToDocument(this.mode)
    },

    toggle() {
      this.set(this.mode === 'dark' ? 'light' : 'dark')
    },

    set(mode: ThemeMode) {
      this.mode = mode
      localStorage.setItem(STORAGE_KEY, mode)
      applyToDocument(mode)
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useThemeStore, import.meta.hot))
}

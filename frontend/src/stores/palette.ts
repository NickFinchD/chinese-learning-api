import { acceptHMRUpdate, defineStore } from 'pinia'

export type PaletteId = 'default' | 'teal-honey' | 'emerald'

interface PaletteColors {
  primary: string
  secondary: string
  accent: string
  mauve: string
  mint: string
}

// Five roles, reused everywhere via CSS custom properties (see applyPalette):
// primary — main buttons/active nav/focus rings/progress bars
// secondary — the logout button (low-stakes, only a couple of uses)
// accent — favorites/star/sparkles/highlights
// mauve — purely decorative (background blobs)
// mint — lighter dark-mode-friendly variant of primary
export const PALETTES: Record<PaletteId, { label: string, colors: PaletteColors }> = {
  'default': {
    label: 'Стандартная',
    colors: {
      primary: '#41b3a3',
      secondary: '#e27d60',
      accent: '#e8a87c',
      mauve: '#c38d9e',
      mint: '#85dcba',
    },
  },
  'teal-honey': {
    label: 'Бирюза и мёд',
    colors: {
      primary: '#026670',
      secondary: '#fef9c7',
      accent: '#fce181',
      mauve: '#edeae5',
      mint: '#9fedd7',
    },
  },
  'emerald': {
    label: 'Изумруд',
    colors: {
      primary: '#479761',
      secondary: '#19181a',
      accent: '#a16e83',
      mauve: '#b19f9e',
      mint: '#cebc81',
    },
  },
}

const STORAGE_KEY = 'wojiao-palette'

// WCAG relative luminance: picks black or white text so it stays readable
// against a background color of any lightness, since palette colors (and any
// added later) aren't guaranteed to be dark enough for white text.
function readableTextColor(hex: string): string {
  const r = parseInt(hex.slice(1, 3), 16) / 255
  const g = parseInt(hex.slice(3, 5), 16) / 255
  const b = parseInt(hex.slice(5, 7), 16) / 255

  const [rl, gl, bl] = [r, g, b].map(c => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
  const luminance = 0.2126 * rl + 0.7152 * gl + 0.0722 * bl

  return luminance > 0.5 ? '#1f2937' : '#ffffff'
}

function applyPalette(id: PaletteId) {
  const { colors } = PALETTES[id]
  const root = document.documentElement

  root.style.setProperty('--color-primary', colors.primary)
  root.style.setProperty('--color-secondary', colors.secondary)
  root.style.setProperty('--color-secondary-text', readableTextColor(colors.secondary))
  root.style.setProperty('--color-accent', colors.accent)
  root.style.setProperty('--color-mauve', colors.mauve)
  root.style.setProperty('--color-mint', colors.mint)
}

export const usePaletteStore = defineStore('palette', {
  state: () => ({
    paletteId: (localStorage.getItem(STORAGE_KEY) as PaletteId | null) ?? 'default',
  }),

  actions: {
    init() {
      applyPalette(this.paletteId)
    },

    setPalette(id: PaletteId) {
      this.paletteId = id
      localStorage.setItem(STORAGE_KEY, id)
      applyPalette(id)
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(usePaletteStore, import.meta.hot))
}

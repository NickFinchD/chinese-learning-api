import { acceptHMRUpdate, defineStore } from 'pinia'

export type HanziFont = 'system' | 'noto-sans' | 'noto-serif' | 'calligraphy'

const FONT_STACKS: Record<HanziFont, string> = {
  system: '"Manrope", ui-sans-serif, system-ui, sans-serif',
  'noto-sans': '"Noto Sans SC", "Manrope", ui-sans-serif, sans-serif',
  'noto-serif': '"Noto Serif SC", serif',
  calligraphy: '"Ma Shan Zheng", "Noto Sans SC", cursive',
}

export const HANZI_FONT_OPTIONS: { value: HanziFont, label: string }[] = [
  { value: 'system', label: 'Стандартный' },
  { value: 'noto-sans', label: 'Noto Sans SC (гладкий)' },
  { value: 'noto-serif', label: 'Noto Serif SC (с засечками)' },
  { value: 'calligraphy', label: 'Ma Shan Zheng (каллиграфия)' },
]

const STORAGE_KEY = 'wojiao-hanzi-font'

function applyFont(font: HanziFont) {
  document.documentElement.style.setProperty('--font-hanzi', FONT_STACKS[font])
}

export const useFontsStore = defineStore('fonts', {
  state: () => ({
    hanziFont: (localStorage.getItem(STORAGE_KEY) as HanziFont | null) ?? 'noto-sans',
  }),

  actions: {
    init() {
      applyFont(this.hanziFont)
    },

    setHanziFont(font: HanziFont) {
      this.hanziFont = font
      localStorage.setItem(STORAGE_KEY, font)
      applyFont(font)
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useFontsStore, import.meta.hot))
}

import { defineStore } from 'pinia'

import { getText, getTexts } from '@/services/texts'

import type { TextItem } from '@/types/text'

export const useTextsStore = defineStore('texts', {
  state: () => ({
    items: [] as TextItem[],
    current: null as TextItem | null,
    hsk: 0,
    loading: false,
  }),

  actions: {
    async loadTexts() {
      this.loading = true

      try {
        const response = await getTexts(this.hsk || undefined)

        this.items = response.data ?? []
      } finally {
        this.loading = false
      }
    },

    async loadText(id: number) {
      this.loading = true

      try {
        const response = await getText(id)

        this.current = response.data
      } finally {
        this.loading = false
      }
    },
  },
})

import { acceptHMRUpdate, defineStore } from 'pinia'

import { getText, getTexts, markTextRead, markTextUnread } from '@/services/texts'

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

    async markRead(id: number) {
      await markTextRead(id)

      this.applyStatus(id, 'completed')
    },

    async markUnread(id: number) {
      await markTextUnread(id)

      this.applyStatus(id, 'in_progress')
    },

    applyStatus(id: number, status: TextItem['status']) {
      if (this.current?.id === id) {
        this.current.status = status
      }

      const item = this.items.find(text => text.id === id)

      if (item) {
        item.status = status
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useTextsStore, import.meta.hot))
}

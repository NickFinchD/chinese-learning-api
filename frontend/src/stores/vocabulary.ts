import { acceptHMRUpdate, defineStore } from 'pinia'

import { getWords } from '@/services/words'

import type { Word } from '@/types/word'

export const useVocabularyStore = defineStore('vocabulary', {
  state: () => ({
    items: [] as Word[],
    search: '',
    hsk: 0,
    loading: false,
  }),

  actions: {
    async loadWords() {
      this.loading = true

      try {
        const response = await getWords({
          search: this.search,
          hsk: this.hsk,
        })

        this.items = response.data ?? []
      } finally {
        this.loading = false
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useVocabularyStore, import.meta.hot))
}

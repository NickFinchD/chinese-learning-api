import { acceptHMRUpdate, defineStore } from 'pinia'

import { getSavedWords, saveWord, unsaveWord } from '@/services/savedWords'

import type { Word } from '@/types/word'

export const useSavedWordsStore = defineStore('savedWords', {
  state: () => ({
    items: [] as Word[],
    loading: false,
  }),

  actions: {
    async loadSavedWords() {
      this.loading = true

      try {
        const response = await getSavedWords()

        this.items = response.data ?? []
      } finally {
        this.loading = false
      }
    },

    async removeWord(wordId: number) {
      await unsaveWord(wordId)

      this.items = this.items.filter(word => word.id !== wordId)
    },

    async addWord(word: Word) {
      await saveWord(word.id)

      if (!this.items.some(item => item.id === word.id)) {
        this.items.push(word)
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSavedWordsStore, import.meta.hot))
}

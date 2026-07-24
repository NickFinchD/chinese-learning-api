import { acceptHMRUpdate, defineStore } from 'pinia'

import {
  addWordToCollection,
  createCollection,
  deleteCollection,
  getCollection,
  getCollections,
  removeWordFromCollection,
  renameCollection,
} from '@/services/collections'

import type { Collection, CollectionDetail } from '@/types/collection'
import type { Word } from '@/types/word'

export const useCollectionsStore = defineStore('collections', {
  state: () => ({
    items: [] as Collection[],
    current: null as CollectionDetail | null,
    loading: false,
    loadingCurrent: false,
  }),

  actions: {
    async loadCollections() {
      this.loading = true

      try {
        const response = await getCollections()

        this.items = response.data
      } finally {
        this.loading = false
      }
    },

    async loadCollection(id: number) {
      this.loadingCurrent = true

      try {
        const response = await getCollection(id)

        this.current = response.data
      } finally {
        this.loadingCurrent = false
      }
    },

    async create(name: string) {
      const response = await createCollection(name)

      this.items.unshift(response.data)

      return response.data
    },

    async rename(id: number, name: string) {
      await renameCollection(id, name)

      const item = this.items.find(c => c.id === id)

      if (item) {
        item.name = name
      }

      if (this.current?.id === id) {
        this.current.name = name
      }
    },

    async remove(id: number) {
      await deleteCollection(id)

      this.items = this.items.filter(c => c.id !== id)

      if (this.current?.id === id) {
        this.current = null
      }
    },

    async addWord(collectionId: number, word: Word) {
      await addWordToCollection(collectionId, word.id)

      const item = this.items.find(c => c.id === collectionId)

      if (item) {
        item.word_count++
      }

      if (this.current?.id === collectionId && !this.current.words.some(w => w.id === word.id)) {
        this.current.words.unshift(word)
        this.current.word_count++
      }
    },

    async removeWord(collectionId: number, wordId: number) {
      await removeWordFromCollection(collectionId, wordId)

      const item = this.items.find(c => c.id === collectionId)

      if (item) {
        item.word_count--
      }

      if (this.current?.id === collectionId) {
        this.current.words = this.current.words.filter(w => w.id !== wordId)
        this.current.word_count--
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useCollectionsStore, import.meta.hot))
}

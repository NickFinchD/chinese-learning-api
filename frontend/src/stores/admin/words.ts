import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListWords, createWord, deleteWord, updateWord } from '@/services/admin/words'
import type { WordPayload } from '@/services/admin/words'

import type { Word } from '@/types/word'

export const useAdminWordsStore = defineStore('adminWords', {
  state: () => ({
    items: [] as Word[],
    total: 0,
    page: 1,
    limit: 20,
    search: '',
    loading: false,
    error: '',
  }),

  getters: {
    totalPages: state => Math.max(1, Math.ceil(state.total / state.limit)),
  },

  actions: {
    async fetchPage() {
      this.loading = true
      this.error = ''

      try {
        const response = await adminListWords({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить слова'
      } finally {
        this.loading = false
      }
    },

    setSearch(search: string) {
      this.search = search
      this.page = 1
      this.fetchPage()
    },

    setPage(page: number) {
      this.page = page
      this.fetchPage()
    },

    async create(payload: WordPayload) {
      await createWord(payload)
      await this.fetchPage()
    },

    async update(id: number, payload: WordPayload) {
      await updateWord(id, payload)
      await this.fetchPage()
    },

    async remove(id: number) {
      await deleteWord(id)
      await this.fetchPage()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminWordsStore, import.meta.hot))
}

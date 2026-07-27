import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListSentences, createSentence, deleteSentence, updateSentence } from '@/services/admin/sentences'
import type { SentenceExercisePayload } from '@/services/admin/sentences'

import type { SentenceExercise } from '@/types/lesson'

export const useAdminSentencesStore = defineStore('adminSentences', {
  state: () => ({
    items: [] as SentenceExercise[],
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
        const response = await adminListSentences({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить упражнения'
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

    async create(payload: SentenceExercisePayload) {
      await createSentence(payload)
      await this.fetchPage()
    },

    async update(id: number, payload: SentenceExercisePayload) {
      await updateSentence(id, payload)
      await this.fetchPage()
    },

    async remove(id: number) {
      await deleteSentence(id)
      await this.fetchPage()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminSentencesStore, import.meta.hot))
}

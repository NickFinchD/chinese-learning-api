import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListQuizzes, createQuiz, deleteQuiz, updateQuiz } from '@/services/admin/quizzes'
import type { AdminQuiz, QuizPayload } from '@/services/admin/quizzes'

export const useAdminQuizzesStore = defineStore('adminQuizzes', {
  state: () => ({
    items: [] as AdminQuiz[],
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
        const response = await adminListQuizzes({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить тесты'
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

    async create(payload: QuizPayload) {
      await createQuiz(payload)
      await this.fetchPage()
    },

    async update(id: number, payload: QuizPayload) {
      await updateQuiz(id, payload)
      await this.fetchPage()
    },

    async remove(id: number) {
      await deleteQuiz(id)
      await this.fetchPage()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminQuizzesStore, import.meta.hot))
}

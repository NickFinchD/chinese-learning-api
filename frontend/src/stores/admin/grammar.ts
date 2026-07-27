import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListGrammarNotes, createGrammarNote, deleteGrammarNote, updateGrammarNote } from '@/services/admin/grammar'
import type { GrammarNotePayload } from '@/services/admin/grammar'

import type { GrammarNote } from '@/types/lesson'

export const useAdminGrammarStore = defineStore('adminGrammar', {
  state: () => ({
    items: [] as GrammarNote[],
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
        const response = await adminListGrammarNotes({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить грамматические заметки'
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

    async create(payload: GrammarNotePayload) {
      await createGrammarNote(payload)
      await this.fetchPage()
    },

    async update(id: number, payload: GrammarNotePayload) {
      await updateGrammarNote(id, payload)
      await this.fetchPage()
    },

    async remove(id: number) {
      await deleteGrammarNote(id)
      await this.fetchPage()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminGrammarStore, import.meta.hot))
}

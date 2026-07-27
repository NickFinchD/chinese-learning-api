import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListTexts, createText, deleteText, updateText } from '@/services/admin/texts'
import type { TextPayload } from '@/services/admin/texts'

import type { TextItem } from '@/types/text'

export const useAdminTextsStore = defineStore('adminTexts', {
  state: () => ({
    items: [] as TextItem[],
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
        const response = await adminListTexts({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить тексты'
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

    async create(payload: TextPayload) {
      await createText(payload)
      await this.fetchPage()
    },

    async update(id: number, payload: TextPayload) {
      await updateText(id, payload)
      await this.fetchPage()
    },

    async remove(id: number) {
      await deleteText(id)
      await this.fetchPage()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminTextsStore, import.meta.hot))
}

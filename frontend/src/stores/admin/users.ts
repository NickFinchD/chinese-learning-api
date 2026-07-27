import { acceptHMRUpdate, defineStore } from 'pinia'

import { adminListUsers, setUserAdmin } from '@/services/admin/users'
import type { AdminUser } from '@/services/admin/users'

export const useAdminUsersStore = defineStore('adminUsers', {
  state: () => ({
    items: [] as AdminUser[],
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
        const response = await adminListUsers({
          search: this.search || undefined,
          page: this.page,
          limit: this.limit,
        })

        this.items = response.data.items
        this.total = response.data.total
      } catch {
        this.error = 'Не удалось загрузить пользователей'
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

    async toggleAdmin(id: number, isAdmin: boolean) {
      this.error = ''

      try {
        await setUserAdmin(id, isAdmin)
        await this.fetchPage()
      } catch {
        this.error = isAdmin
          ? 'Не удалось выдать права администратора'
          : 'Не удалось отозвать права администратора (возможно, это последний админ)'
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAdminUsersStore, import.meta.hot))
}

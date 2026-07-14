import { defineStore } from 'pinia'
import { api } from '@/services/client'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as any,
    loading: false,
  }),

  getters: {
    isAuthenticated: (state) => state.user !== null,
  },

  actions: {
    async loadUser() {
      this.loading = true

      try {
        const response = await api.get('/me')

        this.user = response.data.data
      } catch {
        this.user = null
      } finally {
        this.loading = false
      }
    },
  },
})
import { defineStore } from 'pinia'

import { api } from '@/services/client'

import type { ApiResponse } from '@/types/api'
import type { User } from '@/types/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    loading: false,
  }),

  getters: {
    isAuthenticated: (state) => state.user !== null,
  },

  actions: {
    async loadUser() {
      this.loading = true

      try {
        const response = await api.get<ApiResponse<User>>('/me')

        this.user = response.data.data
      } catch {
        this.user = null
      } finally {
        this.loading = false
      }
    },

    logout() {
      this.user = null
    },
  },
})
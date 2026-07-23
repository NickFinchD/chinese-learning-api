import { defineStore } from 'pinia'

import { getAchievements, getProgress, sendHeartbeat } from '@/services/gamification'

import type { Achievement, GamificationProgress } from '@/types/gamification'

export const useGamificationStore = defineStore('gamification', {
  state: () => ({
    progress: null as GamificationProgress | null,
    achievements: [] as Achievement[],
    loading: false,
  }),

  actions: {
    async loadProgress() {
      const response = await getProgress()

      this.progress = response.data
    },

    async loadAchievements() {
      this.loading = true

      try {
        const response = await getAchievements()

        this.achievements = response.data ?? []
      } finally {
        this.loading = false
      }
    },

    async heartbeat() {
      try {
        const response = await sendHeartbeat()

        this.progress = response.data
      } catch (error) {
        console.error('Failed to send activity heartbeat:', error)
      }
    },
  },
})

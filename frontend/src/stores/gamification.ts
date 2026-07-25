import { acceptHMRUpdate, defineStore } from 'pinia'

import { getAchievements, getProgress, sendHeartbeat } from '@/services/gamification'

import type { Achievement, GamificationProgress } from '@/types/gamification'

export const useGamificationStore = defineStore('gamification', {
  state: () => ({
    progress: null as GamificationProgress | null,
    achievements: [] as Achievement[],
    loading: false,
    // Set right after `progress` jumps to a higher level, so a toast can
    // react to it; null the rest of the time (including on first load).
    leveledUpTo: null as number | null,
  }),

  actions: {
    async loadProgress() {
      const response = await getProgress()

      this.applyProgress(response.data)
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

        this.applyProgress(response.data)
      } catch (error) {
        console.error('Failed to send activity heartbeat:', error)
      }
    },

    applyProgress(next: GamificationProgress) {
      const previousLevel = this.progress?.level

      this.progress = next

      if (previousLevel !== undefined && next.level > previousLevel) {
        this.leveledUpTo = next.level
      }
    },

    dismissLevelUp() {
      this.leveledUpTo = null
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useGamificationStore, import.meta.hot))
}

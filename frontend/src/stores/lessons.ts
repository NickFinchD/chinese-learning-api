import { defineStore } from 'pinia'

import { getLesson } from '@/services/lessons'

import type { Lesson } from '@/types/lesson'

export const useLessonsStore = defineStore('lessons', {
  state: () => ({
    current: null as Lesson | null,
    loading: false,
  }),

  actions: {
    async loadLesson(id: number) {
      this.loading = true

      try {
        const response = await getLesson(id)

        this.current = response.data
      } finally {
        this.loading = false
      }
    },
  },
})
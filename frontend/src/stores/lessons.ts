import { defineStore } from 'pinia'

import { getLesson } from '@/services/lessons'
import { completeLesson, getLessonProgress, startLesson, updateLessonStep } from '@/services/progress'

import type { Lesson } from '@/types/lesson'
import type { LessonProgress } from '@/types/progress'

export const useLessonsStore = defineStore('lessons', {
  state: () => ({
    current: null as Lesson | null,
    progress: null as LessonProgress | null,
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

    // Loads saved progress for the lesson, starting it on the backend
    // the first time it's opened. Returns the step index (0-based) to resume from.
    async resumeOrStart(lessonId: number): Promise<number> {
      const progressResponse = await getLessonProgress(lessonId)

      if (progressResponse.data.status === 'not_started') {
        await startLesson(lessonId)

        this.progress = {
          status: 'in_progress',
          current_step: 1,
          score: 0,
        }

        return 0
      }

      this.progress = progressResponse.data

      return Math.max(0, progressResponse.data.current_step - 1)
    },

    async saveStep(lessonId: number, stepIndex: number) {
      const currentStep = stepIndex + 1

      await updateLessonStep(lessonId, currentStep)

      if (this.progress) {
        this.progress.current_step = currentStep
      }
    },

    async finishLesson(lessonId: number, score: number) {
      const response = await completeLesson(lessonId, score)

      this.progress = {
        status: 'completed',
        current_step: this.progress?.current_step ?? 0,
        score: response.data.score,
      }
    },
  },
})

import { acceptHMRUpdate, defineStore } from 'pinia'

import { getSentencesByHSK } from '@/services/sentences'

import type { SentenceExercise } from '@/types/lesson'

function shuffle<T>(items: T[]): T[] {
  const result = [...items]

  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))

    ;[result[i], result[j]] = [result[j], result[i]]
  }

  return result
}

export const useSentenceTestStore = defineStore('sentenceTest', {
  state: () => ({
    hsk: 1,
    exercises: [] as SentenceExercise[],
    currentIndex: 0,
    correctCount: 0,
    answeredResult: null as boolean | null,
    loading: false,
    started: false,
  }),

  getters: {
    currentExercise: (state): SentenceExercise | null => state.exercises[state.currentIndex] ?? null,

    isFinished: (state) => state.started && state.currentIndex >= state.exercises.length,
  },

  actions: {
    async start(hsk: number) {
      this.hsk = hsk
      this.loading = true

      try {
        const response = await getSentencesByHSK(hsk)

        this.exercises = shuffle(response.data)
        this.currentIndex = 0
        this.correctCount = 0
        this.answeredResult = null
        this.started = true
      } finally {
        this.loading = false
      }
    },

    answer(orderedChunks: string[]) {
      const exercise = this.currentExercise

      if (!exercise || this.answeredResult !== null) {
        return
      }

      const correct = orderedChunks.every((chunk, index) => chunk === exercise.chunks[index])
        && orderedChunks.length === exercise.chunks.length

      this.answeredResult = correct

      if (correct) {
        this.correctCount++
      }
    },

    next() {
      this.currentIndex++
      this.answeredResult = null
    },

    reset() {
      this.exercises = []
      this.currentIndex = 0
      this.correctCount = 0
      this.answeredResult = null
      this.started = false
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSentenceTestStore, import.meta.hot))
}

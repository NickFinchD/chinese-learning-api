import { defineStore } from 'pinia'

import { checkAnswer, getQuizzesByHSK } from '@/services/quizzes'

import type { Quiz } from '@/types/quiz'

export const useGrammarTestStore = defineStore('grammarTest', {
  state: () => ({
    hsk: 1,
    quizzes: [] as Quiz[],
    currentIndex: 0,
    correctCount: 0,
    answeredResult: null as boolean | null,
    loading: false,
    started: false,
  }),

  getters: {
    currentQuiz: (state): Quiz | null => state.quizzes[state.currentIndex] ?? null,

    isFinished: (state) => state.started && state.currentIndex >= state.quizzes.length,
  },

  actions: {
    async start(hsk: number) {
      this.hsk = hsk
      this.loading = true

      try {
        const response = await getQuizzesByHSK(hsk)

        this.quizzes = response.data
        this.currentIndex = 0
        this.correctCount = 0
        this.answeredResult = null
        this.started = true
      } finally {
        this.loading = false
      }
    },

    async answer(optionId: number) {
      const quiz = this.currentQuiz

      if (!quiz || this.answeredResult !== null) {
        return
      }

      const result = await checkAnswer(quiz.id, optionId)

      this.answeredResult = result.correct

      if (result.correct) {
        this.correctCount++
      }
    },

    next() {
      this.currentIndex++
      this.answeredResult = null
    },

    reset() {
      this.quizzes = []
      this.currentIndex = 0
      this.correctCount = 0
      this.answeredResult = null
      this.started = false
    },
  },
})

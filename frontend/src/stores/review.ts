import { defineStore } from 'pinia'

import { answerReview, getReviewStatistics, startReviewSession } from '@/services/review'

import type { ReviewStatistics, ReviewWord } from '@/types/review'

export const useReviewStore = defineStore('review', {
  state: () => ({
    statistics: null as ReviewStatistics | null,
    words: [] as ReviewWord[],
    currentIndex: 0,
    showAnswer: false,
    loading: false,
    sessionStarted: false,
  }),

  getters: {
    currentWord: (state): ReviewWord | null => state.words[state.currentIndex] ?? null,

    isFinished: (state) => state.sessionStarted && state.currentIndex >= state.words.length,
  },

  actions: {
    async loadStatistics() {
      this.loading = true

      try {
        const response = await getReviewStatistics()

        this.statistics = response.data
      } finally {
        this.loading = false
      }
    },

    async startSession() {
      this.loading = true

      try {
        const response = await startReviewSession()

        this.words = response.data.words
        this.currentIndex = 0
        this.showAnswer = false
        this.sessionStarted = true
      } finally {
        this.loading = false
      }
    },

    reveal() {
      this.showAnswer = true
    },

    async answer(correct: boolean) {
      const word = this.currentWord

      if (!word) {
        return
      }

      await answerReview(word.word_id, correct)

      this.currentIndex++
      this.showAnswer = false

      await this.loadStatistics()
    },

    resetSession() {
      this.words = []
      this.currentIndex = 0
      this.showAnswer = false
      this.sessionStarted = false
    },
  },
})

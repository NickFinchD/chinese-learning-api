import { defineStore } from 'pinia'

import { getInProgressWords, getLearnedWords } from '@/services/learning'

import type { InProgressWord } from '@/types/learning'
import type { Word } from '@/types/word'

export const useLearningStore = defineStore('learning', {
  state: () => ({
    learnedWords: [] as Word[],
    loading: false,

    inProgressWords: [] as InProgressWord[],
    loadingInProgress: false,
  }),

  actions: {
    async loadLearned() {
      this.loading = true

      try {
        const response = await getLearnedWords()

        this.learnedWords = response.data ?? []
      } finally {
        this.loading = false
      }
    },

    async loadInProgress() {
      this.loadingInProgress = true

      try {
        const response = await getInProgressWords()

        this.inProgressWords = response.data ?? []
      } finally {
        this.loadingInProgress = false
      }
    },
  },
})

import { defineStore } from 'pinia'

import { getLearnedWords, recordWordAnswer } from '@/services/learning'
import { getSavedWords } from '@/services/savedWords'
import { getWords } from '@/services/words'

import type { WordProgress } from '@/types/learning'
import type { Word } from '@/types/word'

interface TrainingQuestion {
  word: Word
  options: Word[]
}

function shuffle<T>(items: T[]): T[] {
  const copy = [...items]

  for (let i = copy.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[copy[i], copy[j]] = [copy[j], copy[i]]
  }

  return copy
}

export const useWordTrainingStore = defineStore('wordTraining', {
  state: () => ({
    questions: [] as TrainingQuestion[],
    currentIndex: 0,
    correctCount: 0,
    answeredWordId: null as number | null,
    lastProgress: null as WordProgress | null,
    loading: false,
    started: false,
    notEnoughWords: false,
    allLearned: false,
  }),

  getters: {
    currentQuestion: (state): TrainingQuestion | null => state.questions[state.currentIndex] ?? null,

    isFinished: (state) => state.started && state.currentIndex >= state.questions.length,
  },

  actions: {
    async start() {
      this.loading = true

      try {
        const [savedResponse, allResponse, learnedResponse] = await Promise.all([
          getSavedWords(),
          getWords(),
          getLearnedWords(),
        ])

        const saved = savedResponse.data ?? []
        const all = allResponse.data ?? []
        const learnedIds = new Set((learnedResponse.data ?? []).map(word => word.id))

        this.notEnoughWords = false
        this.allLearned = false

        if (saved.length < 2) {
          this.notEnoughWords = true
          this.started = false
          return
        }

        const trainable = saved.filter(word => !learnedIds.has(word.id))

        if (trainable.length < 2) {
          this.allLearned = true
          this.started = false
          return
        }

        const distractorPool = all.length > 0 ? all : saved

        this.questions = shuffle(trainable).map(word => {
          const distractors = shuffle(
            distractorPool.filter(candidate => candidate.id !== word.id),
          ).slice(0, 3)

          return {
            word,
            options: shuffle([word, ...distractors]),
          }
        })

        this.currentIndex = 0
        this.correctCount = 0
        this.answeredWordId = null
        this.lastProgress = null
        this.started = true
      } finally {
        this.loading = false
      }
    },

    async answer(wordId: number) {
      if (this.answeredWordId !== null) {
        return
      }

      this.answeredWordId = wordId

      const question = this.currentQuestion

      if (!question) {
        return
      }

      const correct = wordId === question.word.id

      if (correct) {
        this.correctCount++
      }

      try {
        const response = await recordWordAnswer(question.word.id, correct)

        this.lastProgress = response.data
      } catch (error) {
        console.error('Failed to record learning progress:', error)
      }
    },

    next() {
      this.currentIndex++
      this.answeredWordId = null
      this.lastProgress = null
    },

    reset() {
      this.questions = []
      this.currentIndex = 0
      this.correctCount = 0
      this.answeredWordId = null
      this.lastProgress = null
      this.started = false
      this.notEnoughWords = false
      this.allLearned = false
    },
  },
})

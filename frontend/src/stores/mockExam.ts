import { acceptHMRUpdate, defineStore } from 'pinia'

import { getMockExamHistory, getMockExamPaper, submitMockExam } from '@/services/mockExam'

import type { MockExamAttempt, MockExamQuestion, MockExamResult } from '@/types/mockExam'

export const useMockExamStore = defineStore('mockExam', {
  state: () => ({
    hsk: 1,
    questions: [] as MockExamQuestion[],
    timeLimitSeconds: 0,
    currentIndex: 0,
    quizAnswers: [] as { quiz_id: number, option_id: number }[],
    sentenceAnswers: [] as { exercise_id: number, chunks: string[] }[],
    startedAt: 0,
    loading: false,
    submitting: false,
    started: false,
    result: null as MockExamResult | null,
    history: [] as MockExamAttempt[],
    loadingHistory: false,
  }),

  getters: {
    currentQuestion: (state): MockExamQuestion | null => state.questions[state.currentIndex] ?? null,

    isFinished: (state) => state.started && state.currentIndex >= state.questions.length,
  },

  actions: {
    async start(hsk: number) {
      this.hsk = hsk
      this.loading = true
      this.result = null

      try {
        const response = await getMockExamPaper(hsk)

        this.questions = response.data.questions
        this.timeLimitSeconds = response.data.time_limit_seconds
        this.currentIndex = 0
        this.quizAnswers = []
        this.sentenceAnswers = []
        this.startedAt = Date.now()
        this.started = true
      } finally {
        this.loading = false
      }
    },

    answerQuiz(quizId: number, optionId: number) {
      if (this.currentQuestion?.type !== 'quiz') {
        return
      }

      this.quizAnswers.push({ quiz_id: quizId, option_id: optionId })
      this.currentIndex++
    },

    answerSentence(exerciseId: number, chunks: string[]) {
      if (this.currentQuestion?.type !== 'sentence') {
        return
      }

      this.sentenceAnswers.push({ exercise_id: exerciseId, chunks })
      this.currentIndex++
    },

    async submit() {
      if (this.submitting || !this.started) {
        return
      }

      this.submitting = true

      try {
        const durationSeconds = Math.round((Date.now() - this.startedAt) / 1000)

        const response = await submitMockExam(this.hsk, {
          quiz_answers: this.quizAnswers,
          sentence_answers: this.sentenceAnswers,
          duration_seconds: durationSeconds,
        })

        this.result = response.data
        this.started = false
      } finally {
        this.submitting = false
      }
    },

    async loadHistory() {
      this.loadingHistory = true

      try {
        const response = await getMockExamHistory()

        this.history = response.data
      } finally {
        this.loadingHistory = false
      }
    },

    reset() {
      this.questions = []
      this.currentIndex = 0
      this.quizAnswers = []
      this.sentenceAnswers = []
      this.started = false
      this.result = null
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useMockExamStore, import.meta.hot))
}

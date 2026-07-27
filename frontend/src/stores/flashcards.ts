import { acceptHMRUpdate, defineStore } from 'pinia'

import { getInProgressWords, recordWordAnswer } from '@/services/learning'

import type { InProgressWord } from '@/types/learning'

function shuffle<T>(items: T[]): T[] {
  const copy = [...items]

  for (let i = copy.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[copy[i], copy[j]] = [copy[j], copy[i]]
  }

  return copy
}

export const useFlashcardsStore = defineStore('flashcards', {
  state: () => ({
    queue: [] as InProgressWord[],
    totalCount: 0,
    flipped: false,
    loading: false,
    started: false,
    noWords: false,
  }),

  getters: {
    current: (state): InProgressWord | null => state.queue[0] ?? null,
    isFinished: (state) => state.started && state.queue.length === 0,
  },

  actions: {
    async start() {
      this.loading = true

      try {
        const response = await getInProgressWords()
        const words = response.data ?? []

        this.noWords = words.length === 0

        if (this.noWords) {
          this.started = false
          return
        }

        this.queue = shuffle(words)
        this.totalCount = words.length
        this.flipped = false
        this.started = true
      } finally {
        this.loading = false
      }
    },

    flip() {
      this.flipped = !this.flipped
    },

    // "Знаю" — the card leaves the session for good. Reuses the same
    // spaced-repetition signal the quiz-based word training sends, so a
    // flashcard pass counts toward the word's real learning progress
    // (respecting its cooldown — spamming "знаю" won't fast-track it).
    async know() {
      const word = this.current

      if (!word) return

      this.queue.shift()
      this.flipped = false

      try {
        await recordWordAnswer(word.word_id, true)
      } catch (error) {
        console.error('Failed to record learning progress:', error)
      }
    },

    // "Не знаю" — no server call. The card just goes to the back of this
    // session's queue and will come back around until answered "знаю".
    dontKnow() {
      const word = this.current

      if (!word) return

      this.queue.shift()
      this.queue.push(word)
      this.flipped = false
    },

    reset() {
      this.queue = []
      this.totalCount = 0
      this.flipped = false
      this.started = false
      this.noWords = false
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useFlashcardsStore, import.meta.hot))
}

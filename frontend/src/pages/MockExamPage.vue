<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Пробный экзамен HSK
    </h1>

    <div
      v-if="!test.started && !test.result"
      class="max-w-md"
    >
      <p class="mb-4 text-gray-600 dark:text-gray-400">
        Проверьте себя в формате, близком к настоящему экзамену: лексика, грамматика и
        сборка предложений одним таймером. Проходной балл — 60%.
      </p>

      <div class="mb-8 flex flex-wrap gap-2">
        <button
          v-for="level in [1, 2]"
          :key="level"
          class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] dark:border-white/10 dark:bg-white/5 dark:text-gray-300"
          :disabled="test.loading"
          @click="onStart(level)"
        >
          HSK {{ level }}
        </button>
      </div>

      <div
        v-if="test.loading"
        class="mb-6 flex items-center gap-2 text-gray-500 dark:text-gray-400"
      >
        <BaseSpinner />
        Загрузка...
      </div>

      <h2 class="mb-3 text-lg font-semibold text-gray-900 dark:text-white">
        История попыток
      </h2>

      <div
        v-if="test.loadingHistory"
        class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
      >
        <BaseSpinner />
        Загрузка...
      </div>

      <div
        v-else-if="test.history.length === 0"
        class="text-gray-500 dark:text-gray-400"
      >
        Вы ещё не проходили пробный экзамен.
      </div>

      <div
        v-else
        class="space-y-2"
      >
        <div
          v-for="attempt in test.history"
          :key="attempt.id"
          class="flex items-center justify-between rounded-xl border border-white/50 bg-white/30 px-4 py-3 backdrop-blur-md dark:border-white/10 dark:bg-white/5"
        >
          <div class="flex items-center gap-3">
            <AppIcon
              :name="attempt.passed ? 'check-circle' : 'x-circle'"
              :class="attempt.passed ? 'text-green-600 dark:text-green-400' : 'text-red-500 dark:text-red-400'"
            />

            <div>
              <div class="font-medium text-gray-900 dark:text-white">
                HSK {{ attempt.hsk_level }} · {{ attempt.score_percent }}%
              </div>

              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ formatDate(attempt.created_at) }}
              </div>
            </div>
          </div>

          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ attempt.correct_count }}/{{ attempt.total_questions }}
          </div>
        </div>
      </div>
    </div>

    <div
      v-else-if="test.submitting"
      class="flex max-w-md items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Проверяем результаты...
    </div>

    <div
      v-else-if="test.result"
      class="max-w-md"
    >
      <div class="animate-pop-in rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div
          class="mb-2 flex items-center justify-center gap-2 text-xl font-semibold"
          :class="test.result.attempt.passed ? 'text-green-600 dark:text-green-400' : 'text-red-500 dark:text-red-400'"
        >
          <AppIcon :name="test.result.attempt.passed ? 'check-circle' : 'x-circle'" />
          {{ test.result.attempt.passed ? 'Экзамен сдан!' : 'Экзамен не сдан' }}
        </div>

        <p class="mb-1 text-gray-600 dark:text-gray-400">
          Правильных ответов: {{ test.result.attempt.correct_count }} из {{ test.result.attempt.total_questions }}
          ({{ test.result.attempt.score_percent }}%)
        </p>

        <p
          v-if="test.result.xp_awarded > 0"
          class="mb-6 flex items-center justify-center gap-1 font-semibold text-[var(--color-accent)]"
        >
          <AppIcon
            name="sparkles"
            :size="18"
          />
          +{{ test.result.xp_awarded }} XP
        </p>

        <p
          v-else
          class="mb-6"
        />

        <button
          class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
          @click="onRestart"
        >
          Выбрать другой уровень
        </button>
      </div>
    </div>

    <div
      v-else-if="test.currentQuestion"
      class="max-w-md"
    >
      <div class="mb-2 flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
        <span>Вопрос {{ test.currentIndex + 1 }} из {{ test.questions.length }} · HSK {{ test.hsk }}</span>

        <span
          class="flex items-center gap-1 font-semibold"
          :class="secondsLeft <= 60 ? 'text-red-500 dark:text-red-400' : ''"
        >
          <AppIcon
            name="clock"
            :size="16"
          />
          {{ formatTime(secondsLeft) }}
        </span>
      </div>

      <div
        :key="test.currentIndex"
        class="animate-fade-in-up rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
      >
        <template v-if="test.currentQuestion.type === 'quiz'">
          <h2 class="font-hanzi mb-4 text-lg font-semibold text-gray-900 dark:text-white">
            {{ test.currentQuestion.quiz.question }}
          </h2>

          <div class="space-y-2">
            <button
              v-for="option in test.currentQuestion.quiz.options"
              :key="option.id"
              class="w-full rounded-lg border border-white/50 bg-white/20 p-3 text-left text-gray-800 backdrop-blur-md transition hover:bg-white/40 dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
              @click="onAnswerQuiz(option.id)"
            >
              {{ option.text }}
            </button>
          </div>
        </template>

        <template v-else-if="test.currentQuestion.type === 'sentence'">
          <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
            {{ test.currentQuestion.sentence.translation }}
          </h2>

          <div class="mb-4 flex min-h-16 flex-wrap gap-2 rounded-xl border border-dashed border-white/50 bg-white/10 p-3 dark:border-white/10">
            <button
              v-for="chunk in selected"
              :key="chunk.key"
              type="button"
              class="font-hanzi rounded-lg bg-[var(--color-primary)] px-3 py-2 text-lg text-white shadow transition hover:bg-[var(--color-primary)]/90"
              @click="unselect(chunk)"
            >
              {{ chunk.text }}
            </button>

            <span
              v-if="selected.length === 0"
              class="self-center text-sm text-gray-400 dark:text-gray-500"
            >
              Нажимайте на иероглифы ниже, чтобы собрать предложение
            </span>
          </div>

          <div class="mb-4 flex flex-wrap gap-2">
            <button
              v-for="chunk in pool"
              :key="chunk.key"
              type="button"
              class="font-hanzi rounded-lg border border-white/50 bg-white/20 px-3 py-2 text-lg text-gray-800 backdrop-blur-md transition hover:bg-white/40 dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
              @click="select(chunk)"
            >
              {{ chunk.text }}
            </button>
          </div>

          <button
            class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="pool.length > 0"
            @click="onAnswerSentence"
          >
            Ответить
          </button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'

import { useMockExamStore } from '@/stores/mockExam'

import AppIcon from '@/components/base/AppIcon.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

interface Chunk {
  text: string
  key: number
}

const test = useMockExamStore()

onMounted(() => {
  test.reset()
  test.loadHistory()
})

async function onStart(level: number) {
  await test.start(level)
}

function onRestart() {
  test.reset()
  test.loadHistory()
}

function onAnswerQuiz(optionId: number) {
  if (test.currentQuestion?.type !== 'quiz') {
    return
  }

  test.answerQuiz(test.currentQuestion.quiz.id, optionId)
}

function shuffle(chunks: string[]): Chunk[] {
  const items = chunks.map((text, key) => ({ text, key }))

  for (let i = items.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))

    ;[items[i], items[j]] = [items[j], items[i]]
  }

  return items
}

const pool = ref<Chunk[]>([])
const selected = ref<Chunk[]>([])

watch(
  () => test.currentQuestion,
  (question) => {
    if (question?.type === 'sentence') {
      pool.value = shuffle(question.sentence.chunks)
      selected.value = []
    }
  },
  { immediate: true },
)

function select(chunk: Chunk) {
  pool.value = pool.value.filter(c => c.key !== chunk.key)
  selected.value.push(chunk)
}

function unselect(chunk: Chunk) {
  selected.value = selected.value.filter(c => c.key !== chunk.key)
  pool.value.push(chunk)
}

function onAnswerSentence() {
  if (pool.value.length > 0 || test.currentQuestion?.type !== 'sentence') {
    return
  }

  test.answerSentence(test.currentQuestion.sentence.id, selected.value.map(chunk => chunk.text))
}

// Timer: ticks locally in the component (not the store) so it stops
// cleanly on unmount; hitting zero force-submits whatever's answered so
// far, same as a proctor calling time on a real exam.
const secondsLeft = ref(0)
let timerInterval: ReturnType<typeof setInterval> | undefined

function stopTimer() {
  clearInterval(timerInterval)
  timerInterval = undefined
}

function startTimer() {
  stopTimer()
  secondsLeft.value = test.timeLimitSeconds

  timerInterval = setInterval(() => {
    secondsLeft.value--

    if (secondsLeft.value <= 0) {
      stopTimer()
      test.submit()
    }
  }, 1000)
}

watch(
  () => test.started,
  (started) => {
    if (started) {
      startTimer()
    } else {
      stopTimer()
    }
  },
)

watch(
  () => test.isFinished,
  (finished) => {
    if (finished) {
      stopTimer()
      test.submit()
    }
  },
)

onUnmounted(stopTimer)

function formatTime(totalSeconds: number): string {
  const s = Math.max(0, totalSeconds)
  const minutes = Math.floor(s / 60)
  const seconds = s % 60

  return `${minutes}:${String(seconds).padStart(2, '0')}`
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

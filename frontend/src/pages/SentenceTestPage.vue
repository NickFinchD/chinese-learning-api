<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Конструктор предложений
    </h1>

    <div
      v-if="!test.started"
      class="max-w-md"
    >
      <p class="mb-4 text-gray-600 dark:text-gray-400">
        Выберите уровень HSK:
      </p>

      <div class="mb-6 flex flex-wrap gap-2">
        <button
          v-for="level in 6"
          :key="level"
          class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] dark:border-white/10 dark:bg-white/5 dark:text-gray-300"
          :disabled="test.loading"
          @click="test.start(level)"
        >
          HSK {{ level }}
        </button>
      </div>

      <div
        v-if="test.loading"
        class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
      >
        <BaseSpinner />
        Загрузка...
      </div>
    </div>

    <div
      v-else-if="test.exercises.length === 0"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500 dark:text-gray-400">
        Для уровня HSK {{ test.hsk }} пока нет упражнений.
      </p>

      <button
        class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
        @click="test.reset()"
      >
        Выбрать другой уровень
      </button>
    </div>

    <div
      v-else-if="test.isFinished"
      class="max-w-md"
    >
      <div class="animate-pop-in rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-2 flex items-center justify-center gap-2 text-xl font-semibold text-gray-900 dark:text-white">
          <AppIcon
            name="sparkles"
            :size="22"
            class="text-[var(--color-accent)]"
          />
          Тест завершён
        </div>

        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Правильных ответов: {{ test.correctCount }} из {{ test.exercises.length }}
        </p>

        <button
          class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
          @click="test.reset()"
        >
          Выбрать другой уровень
        </button>
      </div>
    </div>

    <div
      v-else-if="test.currentExercise"
      class="max-w-md"
    >
      <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">
        Предложение {{ test.currentIndex + 1 }} из {{ test.exercises.length }} · HSK {{ test.hsk }}
      </div>

      <div
        :key="test.currentIndex"
        class="animate-fade-in-up rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
          {{ test.currentExercise.translation }}
        </h2>

        <div class="mb-4 flex min-h-16 flex-wrap gap-2 rounded-xl border border-dashed border-white/50 bg-white/10 p-3 dark:border-white/10">
          <button
            v-for="chunk in selected"
            :key="chunk.key"
            type="button"
            class="font-hanzi rounded-lg bg-[var(--color-primary)] px-3 py-2 text-lg text-white shadow transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed"
            :disabled="test.answeredResult !== null"
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

        <div class="flex flex-wrap gap-2">
          <button
            v-for="chunk in pool"
            :key="chunk.key"
            type="button"
            class="font-hanzi rounded-lg border border-white/50 bg-white/20 px-3 py-2 text-lg text-gray-800 backdrop-blur-md transition hover:bg-white/40 disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
            :disabled="test.answeredResult !== null"
            @click="select(chunk)"
          >
            {{ chunk.text }}
          </button>
        </div>

        <div
          v-if="test.answeredResult !== null"
          class="animate-pop-in mt-4 text-lg font-semibold"
        >
          <span
            v-if="test.answeredResult"
            class="flex items-center gap-2 text-green-600 dark:text-green-400"
          >
            <AppIcon name="check-circle" />
            Верно!
          </span>

          <span
            v-else
            class="flex items-center gap-2 text-red-600 dark:text-red-400"
          >
            <AppIcon name="x-circle" />
            Неверно
          </span>

          <span
            v-if="!test.answeredResult"
            class="font-hanzi mt-1 block text-base font-normal text-gray-500 dark:text-gray-400"
          >
            Правильно: {{ test.currentExercise.chunks.join('') }}
          </span>

          <span
            v-if="test.currentExercise.pinyin"
            class="mt-1 block text-base font-normal text-gray-500 dark:text-gray-400"
          >
            {{ test.currentExercise.pinyin }}
          </span>
        </div>

        <button
          v-if="test.answeredResult === null"
          class="mt-6 w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="pool.length > 0"
          @click="check"
        >
          Проверить
        </button>
      </div>

      <button
        v-if="test.answeredResult !== null"
        class="mt-6 w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
        @click="test.next()"
      >
        Далее
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'

import { useSentenceTestStore } from '@/stores/sentenceTest'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

interface Chunk {
  text: string
  key: number
}

const test = useSentenceTestStore()

// The store survives navigation, so returning to this page after leaving
// mid-test would otherwise resume the old level instead of showing the
// picker — start every visit fresh.
onMounted(() => {
  test.reset()
})

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
  () => test.currentExercise,
  (exercise) => {
    pool.value = exercise ? shuffle(exercise.chunks) : []
    selected.value = []
  },
  { immediate: true },
)

function select(chunk: Chunk) {
  if (test.answeredResult !== null) {
    return
  }

  pool.value = pool.value.filter(c => c.key !== chunk.key)
  selected.value.push(chunk)
}

function unselect(chunk: Chunk) {
  if (test.answeredResult !== null) {
    return
  }

  selected.value = selected.value.filter(c => c.key !== chunk.key)
  pool.value.push(chunk)
}

function check() {
  if (pool.value.length > 0) {
    return
  }

  test.answer(selected.value.map(chunk => chunk.text))
}
</script>

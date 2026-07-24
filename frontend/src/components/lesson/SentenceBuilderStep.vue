<template>
  <div class="animate-fade-in-up mx-auto max-w-xl rounded-2xl border border-white/50 bg-white/30 p-8 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <h2 class="mb-1 text-sm font-medium text-gray-500 dark:text-gray-400">
      Собери предложение
    </h2>

    <p class="mb-6 text-xl font-bold text-gray-900 dark:text-white">
      {{ exercise.translation }}
    </p>

    <div class="mb-4 flex min-h-16 flex-wrap gap-2 rounded-xl border border-dashed border-white/50 bg-white/10 p-3 dark:border-white/10">
      <button
        v-for="chunk in selected"
        :key="chunk.key"
        type="button"
        class="font-hanzi rounded-lg bg-[var(--color-primary)] px-3 py-2 text-lg text-white shadow transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed"
        :disabled="result !== null"
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

    <div class="mb-6 flex flex-wrap gap-2">
      <button
        v-for="chunk in pool"
        :key="chunk.key"
        type="button"
        class="font-hanzi rounded-lg border border-white/50 bg-white/20 px-3 py-2 text-lg text-gray-800 backdrop-blur-md transition hover:bg-white/40 disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
        :disabled="result !== null"
        @click="select(chunk)"
      >
        {{ chunk.text }}
      </button>
    </div>

    <div
      v-if="result !== null"
      class="animate-pop-in mb-5 text-lg font-semibold"
    >
      <span
        v-if="result"
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
        v-if="!result"
        class="font-hanzi mt-1 block text-base font-normal text-gray-500 dark:text-gray-400"
      >
        Правильно: {{ exercise.chunks.join('') }}
      </span>

      <span
        v-if="exercise.pinyin"
        class="mt-1 block text-base font-normal text-gray-500 dark:text-gray-400"
      >
        {{ exercise.pinyin }}
      </span>
    </div>

    <button
      v-if="result === null"
      class="rounded-full bg-[var(--color-primary)] px-5 py-2 text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
      :disabled="pool.length > 0"
      @click="check"
    >
      Проверить
    </button>

    <button
      v-else
      class="rounded-full bg-[var(--color-primary)] px-5 py-2 text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
      @click="emit('next')"
    >
      {{ isLastStep ? 'Завершить урок' : 'Далее' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import type { SentenceExercise } from '@/types/lesson'
import AppIcon from '@/components/base/AppIcon.vue'

interface Chunk {
  text: string
  key: number
}

const props = defineProps<{
  exercise: SentenceExercise
  isLastStep: boolean
}>()

const emit = defineEmits<{
  (e: 'answered', correct: boolean): void
  (e: 'next'): void
}>()

function shuffle(chunks: string[]): Chunk[] {
  const items = chunks.map((text, key) => ({ text, key }))

  for (let i = items.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))

    ;[items[i], items[j]] = [items[j], items[i]]
  }

  return items
}

const pool = ref<Chunk[]>(shuffle(props.exercise.chunks))
const selected = ref<Chunk[]>([])
const result = ref<boolean | null>(null)

function select(chunk: Chunk) {
  if (result.value !== null) {
    return
  }

  pool.value = pool.value.filter(c => c.key !== chunk.key)
  selected.value.push(chunk)
}

function unselect(chunk: Chunk) {
  if (result.value !== null) {
    return
  }

  selected.value = selected.value.filter(c => c.key !== chunk.key)
  pool.value.push(chunk)
}

function check() {
  if (pool.value.length > 0 || result.value !== null) {
    return
  }

  const correct = selected.value.every(
    (chunk, index) => chunk.text === props.exercise.chunks[index],
  )

  result.value = correct

  emit('answered', correct)
}
</script>

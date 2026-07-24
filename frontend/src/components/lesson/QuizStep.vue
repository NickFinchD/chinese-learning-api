<template>
  <div class="animate-fade-in-up mx-auto max-w-xl rounded-2xl border border-white/50 bg-white/30 p-8 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <h2 class="font-hanzi mb-6 text-2xl font-bold text-gray-900 dark:text-white">
      {{ quiz.question }}
    </h2>

    <div
      v-for="option in quiz.options"
      :key="option.id"
      class="mb-3"
    >
      <label
        class="flex cursor-pointer items-center gap-3 rounded-lg border border-white/50 bg-white/20 p-3 text-gray-800 backdrop-blur-md transition hover:bg-white/40 dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
        :class="{ 'cursor-not-allowed opacity-60': result !== null }"
      >
        <input
          v-model="selected"
          type="radio"
          :value="option.id"
          :disabled="result !== null"
        >

        {{ option.text }}
      </label>
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
        v-if="pinyin"
        class="mt-1 block text-base font-normal text-gray-500 dark:text-gray-400"
      >
        {{ pinyin }}
      </span>
    </div>

    <button
      v-if="result === null"
      class="rounded-full bg-[var(--color-primary)] px-5 py-2 text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
      :disabled="selected === null"
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
import { computed, ref } from 'vue'

import type { Quiz } from '@/types/lesson'
import { checkAnswer } from '@/services/quizzes'
import AppIcon from '@/components/base/AppIcon.vue'

const props = defineProps<{
  quiz: Quiz
  isLastStep: boolean
  pinyinByHanzi: Record<string, string>
}>()

const emit = defineEmits<{
  (e: 'answered', correct: boolean): void
  (e: 'next'): void
}>()

const selected = ref<number | null>(null)
const result = ref<boolean | null>(null)

// Auto-generated vocab quizzes always follow this exact question format —
// pull the hanzi back out so the pinyin can be shown alongside the result.
const pinyin = computed(() => {
  const match = props.quiz.question.match(/Как переводится (.+)\?/)

  return match ? props.pinyinByHanzi[match[1]] : undefined
})

async function check() {
  if (selected.value === null || result.value !== null) {
    return
  }

  try {
    const response = await checkAnswer(
      props.quiz.id,
      selected.value,
    )

    result.value = response.correct

    emit('answered', response.correct)
  } catch (error) {
    console.error('Failed to check answer:', error)
  }
}
</script>

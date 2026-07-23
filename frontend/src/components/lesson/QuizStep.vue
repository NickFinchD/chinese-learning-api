<template>
  <div class="mx-auto max-w-xl rounded-2xl border border-white/50 bg-white/30 p-8 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">
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

    <button
      class="mt-6 rounded-full bg-[#41b3a3] px-5 py-2 text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90 disabled:cursor-not-allowed disabled:opacity-50"
      :disabled="selected === null || result !== null"
      @click="check"
    >
      Проверить
    </button>

    <div
      v-if="result !== null"
      class="mt-5 text-lg font-semibold"
    >
      <span
        v-if="result"
        class="text-green-600 dark:text-green-400"
      >
        ✅ Верно!
      </span>

      <span
        v-else
        class="text-red-600 dark:text-red-400"
      >
        ❌ Неверно
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import type { Quiz } from '@/types/lesson'
import { checkAnswer } from '@/services/quizzes'

const props = defineProps<{
  quiz: Quiz
}>()

const emit = defineEmits<{
  (e: 'answered', correct: boolean): void
}>()

const selected = ref<number | null>(null)
const result = ref<boolean | null>(null)

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

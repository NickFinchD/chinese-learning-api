<template>
  <div class="mx-auto max-w-xl rounded-2xl border bg-white p-8 shadow">
    <h2 class="mb-6 text-2xl font-bold">
      {{ quiz.question }}
    </h2>

    <div
      v-for="option in quiz.options"
      :key="option.id"
      class="mb-3"
    >
      <label
        class="flex cursor-pointer items-center gap-3 rounded-lg border p-3 hover:bg-gray-50"
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
      class="mt-6 rounded bg-blue-600 px-5 py-2 text-white disabled:cursor-not-allowed disabled:opacity-50"
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
        class="text-green-600"
      >
        ✅ Верно!
      </span>

      <span
        v-else
        class="text-red-600"
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

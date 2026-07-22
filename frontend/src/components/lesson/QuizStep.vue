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
      >
        <input
          v-model="selected"
          type="radio"
          :value="option.id"
        >

        {{ option.text }}
      </label>
    </div>

    <button
      class="mt-6 rounded bg-blue-600 px-5 py-2 text-white"
      @click="check"
    >
      Check
    </button>

    <div
      v-if="result !== null"
      class="mt-5 text-lg font-semibold"
    >
      <span
        v-if="result"
        class="text-green-600"
      >
        ✅ Correct!
      </span>

      <span
        v-else
        class="text-red-600"
      >
        ❌ Wrong
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

const selected = ref<number | null>(null)
const result = ref<boolean | null>(null)

async function check() {
  if (selected.value === null) {
    return
  }

  try {
    const response = await checkAnswer(
      props.quiz.id,
      selected.value,
    )

    result.value = response.correct
  } catch (error) {
    console.error('Failed to check answer:', error)
  }
}
</script>
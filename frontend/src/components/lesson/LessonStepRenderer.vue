<template>
  <WordStep
    v-if="step.step_type === 'word' && step.data"
    :key="step.id"
    :word="step.data"
  />

  <QuizStep
    v-else-if="step.step_type === 'quiz' && step.data"
    :key="step.id"
    :quiz="step.data"
    @answered="$emit('answered', $event)"
  />

  <div
    v-else
    class="rounded-xl border border-red-200 bg-red-50 p-6 text-red-700 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-400"
  >
    Неизвестный тип шага: {{ step.step_type }}
  </div>
</template>

<script setup lang="ts">
import WordStep from './WordStep.vue'
import QuizStep from './QuizStep.vue'

import type { LessonStep } from '@/types/lesson'

defineProps<{
  step: LessonStep
}>()

defineEmits<{
  (e: 'answered', correct: boolean): void
}>()
</script>

<template>
  <WordStep
    v-if="step.step_type === 'word' && step.data"
    :word="step.data"
  />

<QuizStep
  v-else-if="step.step_type === 'quiz' && step.data"
  :quiz="step.data"
/>

  <div
    v-else
    class="rounded-xl border border-red-200 bg-red-50 p-6 text-red-700"
  >
    Unknown lesson step: {{ step.step_type }}
  </div>
</template>

<script setup lang="ts">
import WordStep from './WordStep.vue'
import QuizStep from './QuizStep.vue'

import type { LessonStep } from '@/types/lesson'

defineProps<{
  step: LessonStep
}>()
</script>
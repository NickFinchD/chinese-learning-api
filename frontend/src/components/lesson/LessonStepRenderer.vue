<template>
  <WordStep
    v-if="step.step_type === 'word' && step.data"
    :key="stepKey"
    :word="step.data"
  />

  <QuizStep
    v-else-if="step.step_type === 'quiz' && step.data"
    :key="stepKey"
    :quiz="step.data"
    :is-last-step="isLastStep"
    @answered="$emit('answered', $event)"
    @next="$emit('next')"
  />

  <GrammarStep
    v-else-if="step.step_type === 'grammar' && step.data"
    :key="stepKey"
    :note="step.data"
  />

  <SentenceBuilderStep
    v-else-if="step.step_type === 'sentence_builder' && step.data"
    :key="stepKey"
    :exercise="step.data"
    :is-last-step="isLastStep"
    @answered="$emit('answered', $event)"
    @next="$emit('next')"
  />

  <div
    v-else
    class="rounded-xl border border-red-200 bg-red-50 p-6 text-red-700 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-400"
  >
    Неизвестный тип шага: {{ step.step_type }}
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import WordStep from './WordStep.vue'
import QuizStep from './QuizStep.vue'
import GrammarStep from './GrammarStep.vue'
import SentenceBuilderStep from './SentenceBuilderStep.vue'

import type { LessonStep } from '@/types/lesson'

const props = defineProps<{
  step: LessonStep
  attempt: number
  isLastStep: boolean
}>()

defineEmits<{
  (e: 'answered', correct: boolean): void
  (e: 'next'): void
}>()

// Folds in the retry attempt number so a step being re-queued after a wrong
// answer always gets a fresh component instance (reset local state), even
// when it lands back at the same position in the queue — see LessonPage's
// retryAttempts for why that position-only case needs this.
const stepKey = computed(() => `${props.step.id}-${props.attempt}`)
</script>

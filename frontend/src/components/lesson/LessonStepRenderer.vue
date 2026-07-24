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
    :is-last-step="isLastStep"
    :pinyin-by-hanzi="pinyinByHanzi"
    @answered="$emit('answered', $event)"
    @next="$emit('next')"
  />

  <GrammarStep
    v-else-if="step.step_type === 'grammar' && step.data"
    :key="step.id"
    :note="step.data"
  />

  <SentenceBuilderStep
    v-else-if="step.step_type === 'sentence_builder' && step.data"
    :key="step.id"
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
import WordStep from './WordStep.vue'
import QuizStep from './QuizStep.vue'
import GrammarStep from './GrammarStep.vue'
import SentenceBuilderStep from './SentenceBuilderStep.vue'

import type { LessonStep } from '@/types/lesson'

defineProps<{
  step: LessonStep
  isLastStep: boolean
  pinyinByHanzi: Record<string, string>
}>()

defineEmits<{
  (e: 'answered', correct: boolean): void
  (e: 'next'): void
}>()
</script>

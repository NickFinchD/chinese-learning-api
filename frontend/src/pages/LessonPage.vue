<template>
  <div v-if="lessons.loading || resuming">
    Загрузка...
  </div>

  <div
    v-else-if="finished"
    class="max-w-xl"
  >
    <div class="rounded-2xl border border-white/50 bg-white/30 p-8 text-center shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
      <div class="mb-2 text-2xl font-bold text-gray-900 dark:text-white">
        Урок пройден! 🎉
      </div>

      <p class="mb-6 text-gray-600 dark:text-gray-400">
        Результат: {{ finalScore }}%
      </p>

      <RouterLink
        v-if="lessons.current"
        :to="{ name: 'course', params: { id: lessons.current.course_id } }"
        class="inline-block rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
      >
        Вернуться к курсу
      </RouterLink>
    </div>
  </div>

  <div v-else-if="lessons.current">
    <h1 class="mb-2 text-3xl font-bold text-gray-900 dark:text-white">
      {{ lessons.current.title }}
    </h1>

    <p class="mb-2 text-gray-700 dark:text-gray-300">
      {{ lessons.current.description }}
    </p>

    <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">
      Шаг {{ currentStepIndex + 1 }} из {{ lessons.current.steps.length }}
    </div>

    <div class="mb-6 h-2 overflow-hidden rounded-full bg-gray-200/50 dark:bg-white/10">
      <div
        class="h-full bg-[#41b3a3] transition-all duration-300"
        :style="{
          width: `${((currentStepIndex + 1) / lessons.current.steps.length) * 100}%`
        }"
      />
    </div>

    <div class="mb-6 flex justify-between">
      <button
        class="rounded-full border border-white/50 bg-white/30 px-4 py-2 text-gray-700 backdrop-blur-md transition hover:bg-white/50 disabled:opacity-50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
        :disabled="currentStepIndex === 0"
        @click="previousStep"
      >
        Назад
      </button>

      <button
        class="rounded-full bg-[#41b3a3] px-4 py-2 text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
        @click="nextStep"
      >
        {{ isLastStep ? 'Завершить урок' : 'Далее' }}
      </button>
    </div>

    <div
      v-if="currentStep"
      class="mb-4"
    >
      <LessonStepRenderer
        :step="currentStep"
        @answered="onAnswered"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { useLessonsStore } from '@/stores/lessons'
import LessonStepRenderer from '@/components/lesson/LessonStepRenderer.vue'

const lessons = useLessonsStore()
const route = useRoute()

const lessonId = Number(route.params.id)

const currentStepIndex = ref(0)
const resuming = ref(true)
const finished = ref(false)
const finalScore = ref(0)

const quizResults = new Map<number, boolean>()

const currentStep = computed(() => {
  return lessons.current?.steps[currentStepIndex.value] ?? null
})

const isLastStep = computed(() => {
  if (!lessons.current) {
    return false
  }

  return currentStepIndex.value === lessons.current.steps.length - 1
})

function onAnswered(correct: boolean) {
  const stepId = currentStep.value?.id

  if (stepId !== undefined) {
    quizResults.set(stepId, correct)
  }
}

function computeScore(): number {
  if (!lessons.current) {
    return 0
  }

  const quizSteps = lessons.current.steps.filter(step => step.step_type === 'quiz')

  if (quizSteps.length === 0) {
    return 100
  }

  const correctCount = quizSteps.filter(step => quizResults.get(step.id)).length

  return Math.round((correctCount / quizSteps.length) * 100)
}

async function nextStep() {
  if (!lessons.current) {
    return
  }

  if (currentStepIndex.value < lessons.current.steps.length - 1) {
    currentStepIndex.value++

    lessons.saveStep(lessonId, currentStepIndex.value).catch(error => {
      console.error('Failed to save lesson step:', error)
    })

    return
  }

  finalScore.value = computeScore()

  try {
    await lessons.finishLesson(lessonId, finalScore.value)

    finished.value = true
  } catch (error) {
    console.error('Failed to complete the lesson:', error)
  }
}

function previousStep() {
  if (currentStepIndex.value > 0) {
    currentStepIndex.value--

    lessons.saveStep(lessonId, currentStepIndex.value).catch(error => {
      console.error('Failed to save lesson step:', error)
    })
  }
}

onMounted(async () => {
  await lessons.loadLesson(lessonId)

  try {
    const resumeStepIndex = await lessons.resumeOrStart(lessonId)

    currentStepIndex.value = lessons.current
      ? Math.min(resumeStepIndex, lessons.current.steps.length - 1)
      : 0
  } catch (error) {
    console.error('Failed to load lesson progress:', error)
  } finally {
    resuming.value = false
  }
})
</script>

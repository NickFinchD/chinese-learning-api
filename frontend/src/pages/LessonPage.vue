<template>
  <div
    v-if="lessons.loading || resuming"
    class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
  >
    <BaseSpinner />
    Загрузка...
  </div>

  <div
    v-else-if="finished"
    class="max-w-xl"
  >
    <div class="animate-pop-in rounded-2xl border border-white/50 bg-white/30 p-8 text-center shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
      <div class="mb-2 flex items-center justify-center gap-2 text-2xl font-bold text-gray-900 dark:text-white">
        <AppIcon
          name="sparkles"
          :size="24"
          class="text-[var(--color-accent)]"
        />
        Урок пройден!
      </div>

      <p class="mb-6 text-gray-600 dark:text-gray-400">
        Результат: {{ finalScore }}%

        <span
          v-if="xpAwarded > 0"
          class="ml-2 inline-flex items-center rounded-full bg-[var(--color-accent)]/15 px-3 py-1 text-sm font-semibold text-[var(--color-accent)] dark:bg-[var(--color-accent)]/20"
        >
          +{{ xpAwarded }} XP
        </span>
      </p>

      <div class="flex flex-wrap items-center justify-center gap-3">
        <RouterLink
          v-if="lessons.current"
          :to="{ name: 'course', params: { id: lessons.current.course_id } }"
          class="inline-block rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
        >
          Вернуться к курсу
        </RouterLink>

        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="restarting"
          @click="retakeLesson"
        >
          <AppIcon
            name="refresh"
            :size="16"
          />
          Пройти ещё раз
        </button>
      </div>
    </div>
  </div>

  <div v-else-if="lessons.current">
    <h1 class="mb-2 text-3xl font-bold text-gray-900 dark:text-white">
      {{ lessons.current.title }}
    </h1>

    <p class="font-hanzi mb-2 text-gray-700 dark:text-gray-300">
      {{ lessons.current.description }}
    </p>

    <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">
      Шаг {{ currentStepIndex + 1 }} из {{ stepsQueue.length }}
    </div>

    <div class="mb-6 h-2 overflow-hidden rounded-full bg-gray-200/50 dark:bg-white/10">
      <div
        class="h-full bg-[var(--color-primary)] transition-all duration-300"
        :style="{
          width: `${((currentStepIndex + 1) / stepsQueue.length) * 100}%`
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
        v-if="currentStep?.step_type !== 'quiz' && currentStep?.step_type !== 'sentence_builder'"
        class="rounded-full bg-[var(--color-primary)] px-4 py-2 text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
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
        :is-last-step="isLastStep"
        :pinyin-by-hanzi="pinyinByHanzi"
        @answered="onAnswered"
        @next="nextStep"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { useLessonsStore } from '@/stores/lessons'
import { useGamificationStore } from '@/stores/gamification'
import LessonStepRenderer from '@/components/lesson/LessonStepRenderer.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

import type { LessonStep } from '@/types/lesson'

const lessons = useLessonsStore()
const gamification = useGamificationStore()
const route = useRoute()

const lessonId = Number(route.params.id)

const currentStepIndex = ref(0)
const resuming = ref(true)
const restarting = ref(false)
const finished = ref(false)
const finalScore = ref(0)
const xpAwarded = ref(0)

// Local, mutable copy of the lesson's steps: wrong quiz answers get moved to
// the end of this queue so the learner sees them again before finishing.
const stepsQueue = ref<LessonStep[]>([])

const quizResults = new Map<number, boolean>()

const currentStep = computed(() => {
  return stepsQueue.value[currentStepIndex.value] ?? null
})

const isLastStep = computed(() => {
  return currentStepIndex.value === stepsQueue.value.length - 1
})

// Lets quiz steps show the pinyin for the word they're testing, since the
// same word always also appears as a 'word' step somewhere in the lesson.
const pinyinByHanzi = computed(() => {
  const map: Record<string, string> = {}

  for (const step of stepsQueue.value) {
    if (step.step_type === 'word') {
      map[step.data.hanzi] = step.data.pinyin
    }
  }

  return map
})

function onAnswered(correct: boolean) {
  const stepId = currentStep.value?.id

  if (stepId !== undefined) {
    quizResults.set(stepId, correct)
  }
}

function isScoredStep(step: LessonStep): boolean {
  return step.step_type === 'quiz' || step.step_type === 'sentence_builder'
}

function computeScore(): number {
  const scoredSteps = stepsQueue.value.filter(isScoredStep)

  if (scoredSteps.length === 0) {
    return 100
  }

  const correctCount = scoredSteps.filter(step => quizResults.get(step.id)).length

  return Math.round((correctCount / scoredSteps.length) * 100)
}

async function nextStep() {
  const step = currentStep.value

  if (!step) {
    return
  }

  if (isScoredStep(step) && quizResults.get(step.id) === false) {
    stepsQueue.value.splice(currentStepIndex.value, 1)
    stepsQueue.value.push(step)
  } else if (currentStepIndex.value < stepsQueue.value.length - 1) {
    currentStepIndex.value++
  } else {
    finalScore.value = computeScore()

    try {
      xpAwarded.value = await lessons.finishLesson(lessonId, finalScore.value)

      finished.value = true

      if (xpAwarded.value > 0) {
        gamification.loadProgress().catch(error => {
          console.error('Failed to refresh XP after finishing a lesson:', error)
        })
      }
    } catch (error) {
      console.error('Failed to complete the lesson:', error)
    }

    return
  }

  lessons.saveStep(lessonId, currentStepIndex.value).catch(error => {
    console.error('Failed to save lesson step:', error)
  })
}

function previousStep() {
  if (currentStepIndex.value > 0) {
    currentStepIndex.value--

    lessons.saveStep(lessonId, currentStepIndex.value).catch(error => {
      console.error('Failed to save lesson step:', error)
    })
  }
}

async function retakeLesson() {
  restarting.value = true

  try {
    await lessons.restart(lessonId)

    stepsQueue.value = lessons.current ? [...lessons.current.steps] : []
    currentStepIndex.value = 0
    quizResults.clear()
    finalScore.value = 0
    xpAwarded.value = 0
    finished.value = false
  } catch (error) {
    console.error('Failed to restart the lesson:', error)
  } finally {
    restarting.value = false
  }
}

onMounted(async () => {
  await lessons.loadLesson(lessonId)

  stepsQueue.value = lessons.current ? [...lessons.current.steps] : []

  try {
    const resumeStepIndex = await lessons.resumeOrStart(lessonId)

    if (lessons.progress?.status === 'completed') {
      finalScore.value = lessons.progress.score
      finished.value = true
    } else {
      currentStepIndex.value = Math.min(resumeStepIndex, Math.max(0, stepsQueue.value.length - 1))
    }
  } catch (error) {
    console.error('Failed to load lesson progress:', error)
  } finally {
    resuming.value = false
  }
})
</script>

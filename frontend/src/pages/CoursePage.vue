<template>
  <div>
    <div
      v-if="courses.loading"
      class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Загрузка...
    </div>

    <div v-else-if="courses.current">
      <h1 class="mb-2 text-3xl font-bold text-gray-900 dark:text-white">
        {{ courses.current.title }}
      </h1>

      <p class="mb-6 text-gray-600 dark:text-gray-400">
        {{ courses.current.description }}
      </p>

      <h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
        Уроки
      </h2>

      <div class="space-y-3">
        <RouterLink
          v-for="lesson in courses.current.lessons"
          :key="lesson.id"
          :to="{ name: 'lesson', params: { id: lesson.id } }"
          class="flex items-center justify-between rounded-lg border border-white/50 bg-white/30 p-4 backdrop-blur-xl transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10"
        >
          <div>
            <div class="font-medium text-gray-900 dark:text-white">
              Урок {{ lesson.lesson_number }}
            </div>

            <div class="text-gray-600 dark:text-gray-400">
              {{ lesson.title }}
            </div>
          </div>

          <span
            v-if="statusLabel(lesson.id)"
            class="inline-flex items-center gap-1 rounded-full px-3 py-1 text-xs font-medium"
            :class="statusClass(lesson.id)"
          >
            <AppIcon
              v-if="isCompleted(lesson.id)"
              name="check"
              :size="12"
            />
            {{ statusLabel(lesson.id) }}
          </span>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { getLessonProgress } from '@/services/progress'
import { useCoursesStore } from '@/stores/courses'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

import type { LessonProgress } from '@/types/progress'

const route = useRoute()
const courses = useCoursesStore()

const progressByLesson = ref<Record<number, LessonProgress>>({})

function isCompleted(lessonId: number) {
  return progressByLesson.value[lessonId]?.status === 'completed'
}

function statusLabel(lessonId: number) {
  const status = progressByLesson.value[lessonId]?.status

  if (status === 'completed') {
    return 'Пройден'
  }

  if (status === 'in_progress') {
    return 'В процессе'
  }

  return null
}

function statusClass(lessonId: number) {
  const status = progressByLesson.value[lessonId]?.status

  if (status === 'completed') {
    return 'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-400'
  }

  return 'bg-[var(--color-primary)]/15 text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]'
}

async function loadProgress() {
  if (!courses.current) {
    return
  }

  const entries = await Promise.all(
    courses.current.lessons.map(async lesson => {
      try {
        const response = await getLessonProgress(lesson.id)

        return [lesson.id, response.data] as const
      } catch (error) {
        console.error(`Failed to load progress for lesson ${lesson.id}:`, error)

        return null
      }
    }),
  )

  progressByLesson.value = Object.fromEntries(
    entries.filter((entry): entry is [number, LessonProgress] => entry !== null),
  )
}

onMounted(async () => {
  await courses.loadCourse(Number(route.params.id))
  await loadProgress()
})
</script>

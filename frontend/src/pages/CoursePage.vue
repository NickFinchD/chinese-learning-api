<template>
  <div>
    <div v-if="courses.loading">
      Загрузка...
    </div>

    <div v-else-if="courses.current">
      <h1 class="mb-2 text-3xl font-bold">
        {{ courses.current.title }}
      </h1>

      <p class="mb-6 text-gray-600">
        {{ courses.current.description }}
      </p>

      <h2 class="mb-4 text-xl font-semibold">
        Уроки
      </h2>

      <div class="space-y-3">
        <RouterLink
          v-for="lesson in courses.current.lessons"
          :key="lesson.id"
          :to="{ name: 'lesson', params: { id: lesson.id } }"
          class="flex items-center justify-between rounded-lg border p-4 transition hover:bg-gray-50"
        >
          <div>
            <div class="font-medium">
              Урок {{ lesson.lesson_number }}
            </div>

            <div class="text-gray-600">
              {{ lesson.title }}
            </div>
          </div>

          <span
            v-if="statusLabel(lesson.id)"
            class="inline-flex rounded-full px-3 py-1 text-xs font-medium"
            :class="statusClass(lesson.id)"
          >
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

import type { LessonProgress } from '@/types/progress'

const route = useRoute()
const courses = useCoursesStore()

const progressByLesson = ref<Record<number, LessonProgress>>({})

function statusLabel(lessonId: number) {
  const status = progressByLesson.value[lessonId]?.status

  if (status === 'completed') {
    return '✓ Пройден'
  }

  if (status === 'in_progress') {
    return 'В процессе'
  }

  return null
}

function statusClass(lessonId: number) {
  const status = progressByLesson.value[lessonId]?.status

  if (status === 'completed') {
    return 'bg-green-100 text-green-700'
  }

  return 'bg-blue-100 text-blue-700'
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

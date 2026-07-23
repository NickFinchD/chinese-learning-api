<template>
  <div>
    <h1 class="mb-2 text-3xl font-bold">
      С возвращением{{ auth.user ? `, ${auth.user.username}` : '' }} 👋
    </h1>

    <p class="mb-8 text-gray-600">
      Вот на чём вы остановились.
    </p>

    <div class="mb-8 grid gap-6 sm:grid-cols-3">
      <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div class="text-3xl font-bold text-blue-600">
          {{ courses.items.length }}
        </div>
        <div class="text-gray-500">
          Курсы
        </div>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div class="text-3xl font-bold text-blue-600">
          {{ review.statistics?.ready_for_review ?? 0 }}
        </div>
        <div class="text-gray-500">
          Слов готово к повторению
        </div>
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div class="text-3xl font-bold text-blue-600">
          {{ review.statistics?.total_words ?? 0 }}
        </div>
        <div class="text-gray-500">
          Слов на повторении
        </div>
      </div>
    </div>

    <div class="flex flex-wrap gap-4">
      <RouterLink
        to="/app/courses"
        class="rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
      >
        Продолжить обучение
      </RouterLink>

      <RouterLink
        to="/app/review"
        class="rounded-xl bg-gray-100 px-4 py-3 font-semibold text-gray-700 transition hover:bg-gray-200"
      >
        Повторить слова
      </RouterLink>

      <RouterLink
        to="/app/vocabulary"
        class="rounded-xl bg-gray-100 px-4 py-3 font-semibold text-gray-700 transition hover:bg-gray-200"
      >
        Словарь
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useCoursesStore } from '@/stores/courses'
import { useReviewStore } from '@/stores/review'

const auth = useAuthStore()
const courses = useCoursesStore()
const review = useReviewStore()

onMounted(() => {
  courses.loadCourses()
  review.loadStatistics()
})
</script>

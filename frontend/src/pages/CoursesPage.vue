<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Курсы
    </h1>

    <div
      v-if="courses.loading"
      class="text-gray-500 dark:text-gray-400"
    >
      Загрузка курсов...
    </div>

    <div
      v-else
      class="grid gap-6 md:grid-cols-2 xl:grid-cols-3"
    >
      <RouterLink
        v-for="course in courses.items"
        :key="course.id"
        :to="{ name: 'course', params: { id: course.id } }"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          {{ course.title }}
        </h2>

        <p class="mb-4 text-gray-600 dark:text-gray-400">
          {{ course.description }}
        </p>

        <span
          class="inline-flex rounded-full bg-[var(--color-primary)]/15 px-3 py-1 text-sm font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        >
          HSK {{ course.hsk_level }}
        </span>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'

import { useCoursesStore } from '@/stores/courses'

const courses = useCoursesStore()

onMounted(() => {
  courses.loadCourses()
})
</script>

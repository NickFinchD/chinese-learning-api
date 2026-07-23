<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold">
      Курсы
    </h1>

    <div
      v-if="courses.loading"
      class="text-gray-500"
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
        class="block rounded-xl border border-gray-200 bg-white p-6 shadow-sm transition hover:shadow-md"
      >
        <h2 class="mb-2 text-xl font-semibold">
          {{ course.title }}
        </h2>

        <p class="mb-4 text-gray-600">
          {{ course.description }}
        </p>

        <span
          class="inline-flex rounded-full bg-blue-100 px-3 py-1 text-sm font-medium text-blue-700"
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
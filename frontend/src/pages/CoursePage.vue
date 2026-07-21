<template>
  <div>
    <div v-if="courses.loading">
      Loading...
    </div>

    <div v-else-if="courses.current">
      <h1 class="mb-2 text-3xl font-bold">
        {{ courses.current.title }}
      </h1>

      <p class="mb-6 text-gray-600">
        {{ courses.current.description }}
      </p>

      <h2 class="mb-4 text-xl font-semibold">
        Lessons
      </h2>

      <div class="space-y-3">
        <div
          v-for="lesson in courses.current.lessons"
          :key="lesson.id"
          class="rounded-lg border p-4"
        >
          <div class="font-medium">
            Lesson {{ lesson.lesson_number }}
          </div>

          <div class="text-gray-600">
            {{ lesson.title }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

import { useCoursesStore } from '@/stores/courses'

const route = useRoute()
const courses = useCoursesStore()

onMounted(() => {
  courses.loadCourse(Number(route.params.id))
})
</script>
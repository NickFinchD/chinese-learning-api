<template>
  <div>
    <h1 class="mb-2 text-3xl font-bold text-gray-900 dark:text-white">
      С возвращением{{ auth.user ? `, ${auth.user.username}` : '' }} 👋
    </h1>

    <p class="mb-8 text-gray-600 dark:text-gray-400">
      Вот на чём вы остановились.
    </p>

    <div class="mb-8 grid gap-6 sm:grid-cols-3">
      <div class="rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="text-3xl font-bold text-[#41b3a3]">
          {{ courses.items.length }}
        </div>
        <div class="text-gray-500 dark:text-gray-400">
          Курсы
        </div>
      </div>

      <div class="rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="text-3xl font-bold text-[#41b3a3]">
          {{ learning.learnedWords.length }}
        </div>
        <div class="text-gray-500 dark:text-gray-400">
          Слов изучено
        </div>
      </div>

      <div class="rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="text-3xl font-bold text-[#41b3a3]">
          {{ texts.items.length }}
        </div>
        <div class="text-gray-500 dark:text-gray-400">
          Текстов для чтения
        </div>
      </div>
    </div>

    <div class="flex flex-wrap gap-4">
      <RouterLink
        to="/app/courses"
        class="rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
      >
        Продолжить обучение
      </RouterLink>

      <RouterLink
        to="/app/texts"
        class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
      >
        Читать тексты
      </RouterLink>

      <RouterLink
        to="/app/vocabulary"
        class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
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
import { useLearningStore } from '@/stores/learning'
import { useTextsStore } from '@/stores/texts'

const auth = useAuthStore()
const courses = useCoursesStore()
const learning = useLearningStore()
const texts = useTextsStore()

onMounted(() => {
  courses.loadCourses()
  learning.loadLearned()
  texts.loadTexts()
})
</script>

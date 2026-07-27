<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Админка
    </h1>

    <div class="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
      <component
        :is="section.to ? RouterLink : 'div'"
        v-for="section in sections"
        :key="section.label"
        :to="section.to"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
        :class="section.to ? 'transition hover:shadow-md' : 'opacity-50'"
      >
        <div class="mb-3 flex items-center gap-3">
          <AppIcon :name="section.icon" />
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
            {{ section.label }}
          </h2>
        </div>

        <p class="text-gray-600 dark:text-gray-400">
          {{ section.to ? section.description : `${section.description} (скоро)` }}
        </p>
      </component>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

import type { IconName } from '@/components/base/AppIcon.vue'
import AppIcon from '@/components/base/AppIcon.vue'

interface AdminSection {
  label: string
  description: string
  icon: IconName
  to?: string
}

// Sections go from disabled to active as their CRUD screens ship — this
// shell is intentionally built before any of them exist.
const sections: AdminSection[] = [
  { label: 'Слова', description: 'Ханьцзы, пиньинь, перевод, HSK-уровень', icon: 'book-open', to: '/app/admin/words' },
  { label: 'Тексты', description: 'Учебные тексты для чтения', icon: 'file-text', to: '/app/admin/texts' },
  { label: 'Грамматика', description: 'Грамматические заметки', icon: 'graduation-cap', to: '/app/admin/grammar' },
  { label: 'Тесты', description: 'Квизы и варианты ответов', icon: 'flask', to: '/app/admin/quizzes' },
  { label: 'Предложения', description: 'Упражнения на составление предложений', icon: 'sort', to: '/app/admin/sentences' },
  { label: 'Курсы и уроки', description: 'Структура курсов и шаги уроков', icon: 'folder', to: '/app/admin/courses' },
  { label: 'Пользователи', description: 'Список пользователей, права администратора', icon: 'lock', to: '/app/admin/users' },
]
</script>

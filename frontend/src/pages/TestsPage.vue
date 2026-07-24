<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Тесты
    </h1>

    <div class="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
      <RouterLink
        to="/app/tests/mock-exam"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          Пробный экзамен
        </h2>

        <p class="text-gray-600 dark:text-gray-400">
          Таймер, лексика, грамматика и предложения — как на настоящем HSK 1/2.
        </p>
      </RouterLink>

      <RouterLink
        to="/app/tests/grammar"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          Грамматика по HSK
        </h2>

        <p class="text-gray-600 dark:text-gray-400">
          Тест с вариантами ответов по выбранному уровню HSK.
        </p>
      </RouterLink>

      <RouterLink
        to="/app/tests/words"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          Слово → перевод
        </h2>

        <p class="text-gray-600 dark:text-gray-400">
          Тренировка на основе ваших сохранённых слов.
        </p>
      </RouterLink>

      <RouterLink
        to="/app/tests/sentences"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          Конструктор предложений
        </h2>

        <p class="text-gray-600 dark:text-gray-400">
          Собери предложение из иероглифов по переводу.
        </p>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

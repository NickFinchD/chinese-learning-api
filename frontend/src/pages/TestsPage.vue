<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold">
      Тесты
    </h1>

    <div class="grid gap-6 md:grid-cols-2">
      <RouterLink
        to="/app/tests/grammar"
        class="block rounded-xl border border-gray-200 bg-white p-6 shadow-sm transition hover:shadow-md"
      >
        <h2 class="mb-2 text-xl font-semibold">
          Грамматика по HSK
        </h2>

        <p class="text-gray-600">
          Тест с вариантами ответов по выбранному уровню HSK.
        </p>
      </RouterLink>

      <RouterLink
        to="/app/tests/words"
        class="block rounded-xl border border-gray-200 bg-white p-6 shadow-sm transition hover:shadow-md"
      >
        <h2 class="mb-2 text-xl font-semibold">
          Слово → перевод
        </h2>

        <p class="text-gray-600">
          Тренировка на основе ваших сохранённых слов.
        </p>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

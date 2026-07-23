<template>
  <aside class="w-64 border-r border-gray-200 bg-white">
    <nav class="flex flex-col gap-2 p-4">

      <RouterLink
        to="/app"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        🏠 Главная
      </RouterLink>

      <RouterLink
        to="/app/courses"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        📚 Курсы
      </RouterLink>

      <RouterLink
        to="/app/vocabulary"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        📖 Словарь
      </RouterLink>

      <RouterLink
        to="/app/review"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        🔁 Повторение
      </RouterLink>

      <RouterLink
        to="/app/tests"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        🧪 Тесты
      </RouterLink>

      <RouterLink
        to="/app/settings"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        ⚙️ Настройки
      </RouterLink>

    </nav>
  </aside>
</template>

<script setup lang="ts">
</script>
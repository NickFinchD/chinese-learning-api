<template>
  <aside class="w-64 border-r border-white/50 bg-white/30 backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <nav class="flex flex-col gap-2 p-4">

      <RouterLink
        to="/app"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        🏠 Главная
      </RouterLink>

      <RouterLink
        to="/app/courses"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        📚 Курсы
      </RouterLink>

      <RouterLink
        to="/app/vocabulary"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        📖 Словарь
      </RouterLink>

      <RouterLink
        to="/app/texts"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        📄 Тексты
      </RouterLink>

      <RouterLink
        to="/app/tests"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        🧪 Тесты
      </RouterLink>

      <RouterLink
        to="/app/achievements"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        🏆 Достижения
      </RouterLink>

      <RouterLink
        to="/app/settings"
        class="rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[#41b3a3]/15 text-[#41b3a3] font-semibold dark:bg-[#41b3a3]/20 dark:text-[#85dcba]"
      >
        ⚙️ Настройки
      </RouterLink>

    </nav>
  </aside>
</template>

<script setup lang="ts">
</script>

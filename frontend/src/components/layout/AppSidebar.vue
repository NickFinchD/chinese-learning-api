<template>
  <aside class="w-64 border-r border-gray-200 bg-white">
    <nav class="flex flex-col gap-2 p-4">

      <RouterLink
        to="/app"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        🏠 Dashboard
      </RouterLink>

      <RouterLink
        to="/app/courses"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        📚 Courses
      </RouterLink>

      <RouterLink
        to="/app/saved"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        ⭐ Saved Words
      </RouterLink>

      <RouterLink
        to="/app/review"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        🔁 Review
      </RouterLink>

      <RouterLink
        to="/app/settings"
        class="rounded-lg px-4 py-3 transition hover:bg-blue-50"
        active-class="bg-blue-100 text-blue-700 font-semibold"
      >
        ⚙️ Settings
      </RouterLink>

    </nav>
  </aside>
</template>

<script setup lang="ts">
</script>
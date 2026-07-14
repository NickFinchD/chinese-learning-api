<template>
  <aside
    class="w-64 border-r border-gray-200 bg-white"
  >
    <nav class="flex flex-col p-4">

      <RouterLink
        to="/app"
        class="rounded-lg px-4 py-3 hover:bg-blue-50"
      >
        Dashboard
      </RouterLink>

      <RouterLink
        to="/app/courses"
        class="rounded-lg px-4 py-3 hover:bg-blue-50"
      >
        Courses
      </RouterLink>

      <RouterLink
        to="/app/saved"
        class="rounded-lg px-4 py-3 hover:bg-blue-50"
      >
        Saved Words
      </RouterLink>

      <RouterLink
        to="/app/review"
        class="rounded-lg px-4 py-3 hover:bg-blue-50"
      >
        Review
      </RouterLink>

    </nav>
  </aside>
</template>

<script setup lang="ts">
</script>
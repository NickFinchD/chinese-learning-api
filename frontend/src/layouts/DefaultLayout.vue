<template>
  <div class="flex h-screen flex-col bg-gray-100">

    <AppHeader />

    <div class="flex flex-1 overflow-hidden">

      <AppSidebar />

      <main class="flex-1 overflow-y-auto p-8">
        <RouterView />
      </main>

    </div>

  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
</script>
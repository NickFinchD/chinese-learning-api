<template>
  <main class="relative flex min-h-screen items-center justify-center px-4">
    <ThemeToggle class="absolute right-4 top-4" />

    <RouterView v-slot="{ Component, route }">
      <Transition
        name="page"
        mode="out-in"
      >
        <component
          :is="Component"
          :key="route.fullPath"
        />
      </Transition>
    </RouterView>
  </main>
</template>

<script setup lang="ts">
import ThemeToggle from '@/components/base/ThemeToggle.vue'
</script>

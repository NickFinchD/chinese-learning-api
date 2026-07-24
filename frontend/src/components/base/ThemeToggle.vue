<template>
  <button
    type="button"
    class="flex h-10 w-10 items-center justify-center rounded-full border border-white/50 bg-white/30 text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
    :title="theme.mode === 'dark' ? 'Светлая тема' : 'Тёмная тема'"
    @click="theme.toggle()"
  >
    <AppIcon :name="theme.mode === 'dark' ? 'sun' : 'moon'" />
  </button>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'
import AppIcon from './AppIcon.vue'

const theme = useThemeStore()
</script>

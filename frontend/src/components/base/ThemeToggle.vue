<template>
  <button
    type="button"
    class="flex h-10 w-10 items-center justify-center rounded-full border border-white/50 bg-white/30 text-lg backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10"
    :title="theme.mode === 'dark' ? 'Светлая тема' : 'Тёмная тема'"
    @click="theme.toggle()"
  >
    {{ theme.mode === 'dark' ? '☀️' : '🌙' }}
  </button>
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'

const theme = useThemeStore()
</script>

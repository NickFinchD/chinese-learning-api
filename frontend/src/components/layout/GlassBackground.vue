<template>
  <div class="fixed inset-0 -z-10 overflow-hidden bg-gradient-to-b from-gray-100 via-gray-200 to-gray-300 transition-colors duration-300 dark:from-slate-950 dark:via-slate-950 dark:to-slate-900">
    <div class="animate-blob absolute -left-40 -top-40 h-[28rem] w-[28rem] rounded-full bg-[var(--color-primary)] opacity-30 blur-3xl dark:opacity-25" />
    <div
      class="animate-blob absolute -right-32 top-0 h-[26rem] w-[26rem] rounded-full bg-[var(--color-mint)] opacity-25 blur-3xl dark:opacity-20"
      style="animation-delay: -6s; animation-duration: 24s"
    />

    <!-- Background artwork chosen in Settings (see stores/background.ts), each
         with separate day/night artwork. Both images stay mounted and are
         crossfaded via opacity (rather than display:none) so switching theme
         fades smoothly between them instead of jump-cutting. Blurred so the
         glass cards' own backdrop-blur has soft, painterly colour to pick up
         rather than a crisp graphic. -->
    <img
      :src="current.day"
      alt=""
      class="absolute inset-0 h-full w-full min-w-[1440px] object-cover object-bottom blur-sm transition-opacity duration-700 ease-in-out dark:opacity-0"
      aria-hidden="true"
    >
    <img
      :src="current.night"
      alt=""
      class="absolute inset-0 h-full w-full min-w-[1440px] object-cover object-bottom opacity-0 blur-sm transition-opacity duration-700 ease-in-out dark:opacity-100"
      aria-hidden="true"
    >
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { BACKGROUNDS, useBackgroundStore } from '@/stores/background'

const background = useBackgroundStore()
const current = computed(() => BACKGROUNDS[background.backgroundId])
</script>

<template>
  <div class="mx-auto max-w-xl rounded-2xl border border-white/50 bg-white/30 p-10 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <div class="mb-6 text-center text-6xl font-bold text-gray-900 dark:text-white">
      {{ word.hanzi }}
    </div>

    <div class="mb-3 text-center text-2xl text-gray-500 dark:text-gray-400">
      {{ word.pinyin }}
    </div>

    <div class="text-center text-3xl font-medium text-gray-800 dark:text-gray-200">
      {{ word.translation }}
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Word } from '@/types/lesson'

defineProps<{
  word: Word
}>()
</script>

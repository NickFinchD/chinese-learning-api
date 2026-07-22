<template>
  <div class="mx-auto max-w-xl rounded-2xl border bg-white p-10 shadow">
    <div class="mb-6 text-center text-6xl font-bold">
      {{ word.hanzi }}
    </div>

    <div class="mb-3 text-center text-2xl text-gray-500">
      {{ word.pinyin }}
    </div>

    <div class="text-center text-3xl font-medium">
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
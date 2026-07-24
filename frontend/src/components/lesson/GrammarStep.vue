<template>
  <div class="animate-fade-in-up mx-auto max-w-xl rounded-2xl border border-white/50 bg-white/30 p-8 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
    <div class="mb-4 flex items-center gap-2 text-sm font-semibold uppercase tracking-wide text-[var(--color-primary)]">
      <AppIcon
        name="book-open"
        :size="18"
      />
      Грамматика
    </div>

    <h2 class="font-hanzi mb-4 text-2xl font-bold text-gray-900 dark:text-white">
      {{ note.title }}
    </h2>

    <p class="mb-6 whitespace-pre-line text-gray-700 dark:text-gray-300">
      {{ note.explanation }}
    </p>

    <div
      v-if="note.example_hanzi"
      class="rounded-xl border border-white/50 bg-white/20 p-4 backdrop-blur-md dark:border-white/10 dark:bg-white/5"
    >
      <div class="mb-2 flex items-center gap-3">
        <div class="font-hanzi text-2xl font-bold text-gray-900 dark:text-white">
          {{ note.example_hanzi }}
        </div>

        <AudioButton :text="note.example_hanzi" />
      </div>

      <div
        v-if="note.example_pinyin"
        class="mb-1 text-gray-500 dark:text-gray-400"
      >
        {{ note.example_pinyin }}
      </div>

      <div
        v-if="note.example_translation"
        class="text-gray-800 dark:text-gray-200"
      >
        {{ note.example_translation }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GrammarNote } from '@/types/lesson'
import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'

defineProps<{
  note: GrammarNote
}>()
</script>

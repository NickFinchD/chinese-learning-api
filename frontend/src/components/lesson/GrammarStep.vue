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

    <div class="space-y-3">
      <div
        v-for="(example, index) in examples"
        :key="index"
        class="rounded-xl border border-white/50 bg-white/20 p-4 backdrop-blur-md dark:border-white/10 dark:bg-white/5"
      >
        <div class="mb-2 flex items-center gap-3">
          <div class="font-hanzi text-2xl font-bold text-gray-900 dark:text-white">
            {{ example.hanzi }}
          </div>

          <AudioButton :text="example.hanzi" />
        </div>

        <div
          v-if="example.pinyin"
          class="mb-1 text-gray-500 dark:text-gray-400"
        >
          {{ example.pinyin }}
        </div>

        <div
          v-if="example.translation"
          class="text-gray-800 dark:text-gray-200"
        >
          {{ example.translation }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { GrammarNote } from '@/types/lesson'
import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'

const props = defineProps<{
  note: GrammarNote
}>()

// Most notes carry a single illustrative example; dedicated "Примеры: ..."
// notes (shown right after a particle/auxiliary verb's principle note) use
// all three slots to show several usage sentences in one step.
interface ExampleSlot {
  hanzi?: string
  pinyin?: string
  translation?: string
}

const examples = computed(() => {
  const slots: ExampleSlot[] = [
    { hanzi: props.note.example_hanzi, pinyin: props.note.example_pinyin, translation: props.note.example_translation },
    { hanzi: props.note.example2_hanzi, pinyin: props.note.example2_pinyin, translation: props.note.example2_translation },
    { hanzi: props.note.example3_hanzi, pinyin: props.note.example3_pinyin, translation: props.note.example3_translation },
  ]

  return slots.filter((slot): slot is ExampleSlot & { hanzi: string } => Boolean(slot.hanzi))
})
</script>

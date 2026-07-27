<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-150"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="wordId !== null"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
        @click.self="$emit('close')"
      >
        <div class="max-h-[85vh] w-full max-w-lg overflow-y-auto rounded-2xl border border-white/50 bg-white/90 p-6 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/95">
          <div class="mb-4 flex items-start justify-between gap-4">
            <div
              v-if="word"
              class="flex items-center gap-3"
            >
              <div class="font-hanzi text-4xl font-bold text-gray-900 dark:text-white">
                {{ word.hanzi }}
              </div>

              <AudioButton :text="word.hanzi" />
            </div>

            <button
              type="button"
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-gray-500 transition hover:bg-white/60 dark:text-gray-400 dark:hover:bg-white/10"
              @click="$emit('close')"
            >
              <AppIcon name="x" />
            </button>
          </div>

          <div
            v-if="loading"
            class="flex justify-center py-10"
          >
            <BaseSpinner :size="28" />
          </div>

          <template v-else-if="word">
            <div class="mb-1 text-xl text-gray-500 dark:text-gray-400">
              {{ word.pinyin }}
            </div>

            <div class="mb-4 text-2xl font-medium text-gray-800 dark:text-gray-200">
              {{ word.translation }}
            </div>

            <div class="mb-6 flex flex-wrap gap-2">
              <span class="inline-flex rounded-full bg-[var(--color-primary)]/15 px-3 py-1 text-sm font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]">
                HSK {{ word.hsk_level }}
              </span>

              <span
                v-if="word.part_of_speech"
                class="inline-flex rounded-full border border-white/50 px-3 py-1 text-sm text-gray-600 dark:border-white/10 dark:text-gray-300"
              >
                {{ word.part_of_speech }}
              </span>
            </div>

            <h3 class="mb-2 text-sm font-semibold text-gray-700 dark:text-gray-300">
              Примеры
            </h3>

            <div
              v-if="word.examples && word.examples.length > 0"
              class="space-y-3"
            >
              <div
                v-for="example in word.examples"
                :key="example.id"
                class="rounded-lg border border-white/50 bg-white/20 p-3 dark:border-white/10 dark:bg-white/5"
              >
                <div class="mb-0.5 flex items-center gap-2">
                  <span class="font-hanzi text-lg text-gray-900 dark:text-white">{{ example.hanzi }}</span>
                  <AudioButton
                    :text="example.hanzi"
                    size="sm"
                  />
                </div>

                <div class="text-sm text-gray-500 dark:text-gray-400">
                  {{ example.pinyin }}
                </div>

                <div class="text-sm text-gray-700 dark:text-gray-300">
                  {{ example.translation }}
                </div>
              </div>
            </div>

            <p
              v-else
              class="text-sm text-gray-400 dark:text-gray-500"
            >
              Примеров пока нет.
            </p>
          </template>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

import { getWord } from '@/services/words'
import type { WordDetail } from '@/types/word'

import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const props = defineProps<{
  wordId: number | null
}>()

defineEmits<{
  (e: 'close'): void
}>()

const word = ref<WordDetail | null>(null)
const loading = ref(false)

watch(
  () => props.wordId,
  async (id) => {
    word.value = null

    if (id === null) return

    loading.value = true

    try {
      const response = await getWord(id)

      word.value = response.data
    } finally {
      loading.value = false
    }
  },
  { immediate: true },
)
</script>

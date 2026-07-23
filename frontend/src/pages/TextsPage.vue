<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Тексты
    </h1>

    <div class="mb-6 flex items-center gap-4">
      <select
        v-model.number="texts.hsk"
        class="rounded-xl border border-white/50 bg-white/40 px-4 py-3 text-gray-900 outline-none backdrop-blur-md transition focus:border-[#41b3a3] dark:border-white/10 dark:bg-white/5 dark:text-white"
        @change="texts.loadTexts()"
      >
        <option :value="0">
          Любой уровень HSK
        </option>
        <option
          v-for="level in 6"
          :key="level"
          :value="level"
        >
          HSK {{ level }}
        </option>
      </select>
    </div>

    <div
      v-if="texts.loading"
      class="text-gray-500 dark:text-gray-400"
    >
      Загрузка...
    </div>

    <div
      v-else-if="texts.items.length === 0"
      class="text-gray-500 dark:text-gray-400"
    >
      Текстов для этого уровня пока нет.
    </div>

    <div
      v-else
      class="grid gap-4 md:grid-cols-2 xl:grid-cols-3"
    >
      <RouterLink
        v-for="text in texts.items"
        :key="text.id"
        :to="{ name: 'text', params: { id: text.id } }"
        class="block rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
      >
        <div class="mb-2 flex items-center justify-between gap-2">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ text.title }}
          </h2>

          <span class="inline-flex shrink-0 rounded-full bg-[#41b3a3]/15 px-3 py-1 text-xs font-medium text-[#41b3a3] dark:bg-[#41b3a3]/20 dark:text-[#85dcba]">
            HSK {{ text.hsk_level }}
          </span>
        </div>

        <p class="text-gray-600 dark:text-gray-400">
          {{ preview(text.hanzi) }}
        </p>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'

import { useTextsStore } from '@/stores/texts'

const texts = useTextsStore()

function preview(hanzi: string) {
  return [...hanzi].slice(0, 5).join('') + '...'
}

onMounted(() => {
  texts.loadTexts()
})
</script>

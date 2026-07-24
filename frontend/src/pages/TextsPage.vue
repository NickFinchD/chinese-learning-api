<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Тексты
    </h1>

    <div class="mb-6 flex flex-wrap items-center gap-4">
      <BaseInput
        v-model="search"
        placeholder="Поиск по названию, иероглифу или переводу"
        class="max-w-xs"
      />

      <BaseSelect
        v-model="texts.hsk"
        class="w-56"
        :options="hskOptions"
        @update:model-value="texts.loadTexts()"
      />

      <button
        type="button"
        class="flex items-center gap-2 rounded-xl border border-white/50 bg-white/40 px-4 py-3 text-gray-700 outline-none backdrop-blur-md transition hover:bg-white/60 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
        @click="sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'"
      >
        <AppIcon
          name="sort"
          :size="16"
        />
        {{ sortDirection === 'asc' ? 'Сначала короткие' : 'Сначала длинные' }}
      </button>

      <div class="flex gap-2 rounded-xl border border-white/50 bg-white/30 p-1 backdrop-blur-md dark:border-white/10 dark:bg-white/5">
        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="statusTab === 'all' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="statusTab = 'all'"
        >
          Все
        </button>

        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="statusTab === 'completed' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="statusTab = 'completed'"
        >
          Прочитанные
        </button>

        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="statusTab === 'in_progress' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="statusTab = 'in_progress'"
        >
          На изучении
        </button>
      </div>
    </div>

    <div
      v-if="texts.loading"
      class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Загрузка...
    </div>

    <div
      v-else-if="texts.items.length === 0"
      class="text-gray-500 dark:text-gray-400"
    >
      Текстов для этого уровня пока нет.
    </div>

    <div
      v-else-if="filteredTexts.length === 0"
      class="text-gray-500 dark:text-gray-400"
    >
      Ничего не найдено.
    </div>

    <div
      v-else
      class="grid gap-4 md:grid-cols-2 xl:grid-cols-3"
    >
      <RouterLink
        v-for="(text, index) in filteredTexts"
        :key="text.id"
        :to="{ name: 'text', params: { id: text.id } }"
        class="animate-fade-in-up flex h-40 flex-col rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:-translate-y-0.5 hover:shadow-md dark:border-white/10 dark:bg-white/5"
        :style="{ animationDelay: `${Math.min(index * 40, 300)}ms` }"
      >
        <div class="mb-2 flex items-start justify-between gap-2">
          <h2 class="line-clamp-2 text-lg font-semibold text-gray-900 dark:text-white">
            {{ text.title }}
          </h2>

          <div class="flex shrink-0 items-center gap-1">
            <AppIcon
              v-if="text.status === 'completed'"
              name="check-circle"
              :size="16"
              class="text-[var(--color-primary)]"
            />

            <span class="inline-flex rounded-full bg-[var(--color-primary)]/15 px-3 py-1 text-xs font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]">
              HSK {{ text.hsk_level }}
            </span>
          </div>
        </div>

        <p class="font-hanzi line-clamp-2 text-gray-600 dark:text-gray-400">
          {{ preview(text.hanzi) }}
        </p>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { useTextsStore } from '@/stores/texts'

import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const texts = useTextsStore()

const hskOptions = [
  { value: 0, label: 'Любой уровень HSK' },
  ...Array.from({ length: 6 }, (_, index) => ({ value: index + 1, label: `HSK ${index + 1}` })),
]

const search = ref('')
const sortDirection = ref<'asc' | 'desc'>('asc')
const statusTab = ref<'all' | 'completed' | 'in_progress'>('all')

const sortedTexts = computed(() => {
  const sorted = [...texts.items].sort((a, b) => [...a.hanzi].length - [...b.hanzi].length)

  return sortDirection.value === 'asc' ? sorted : sorted.reverse()
})

const filteredTexts = computed(() => {
  const query = search.value.trim().toLowerCase()

  return sortedTexts.value.filter((text) => {
    if (statusTab.value !== 'all' && text.status !== statusTab.value) {
      return false
    }

    if (!query) {
      return true
    }

    return text.title.toLowerCase().includes(query)
      || text.hanzi.includes(query)
      || text.pinyin.toLowerCase().includes(query)
      || text.translation.toLowerCase().includes(query)
  })
})

function preview(hanzi: string) {
  return [...hanzi].slice(0, 5).join('') + '...'
}

onMounted(() => {
  texts.loadTexts()
})
</script>

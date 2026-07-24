<template>
  <div>
    <RouterLink
      to="/app/vocabulary?tab=collections"
      class="mb-6 inline-flex items-center gap-2 text-gray-500 transition hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
    >
      <AppIcon
        name="arrow-left"
        :size="18"
      />
      Все подборки
    </RouterLink>

    <div
      v-if="collections.loadingCurrent"
      class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Загрузка...
    </div>

    <template v-else-if="collections.current">
      <h1 class="mb-8 flex items-center gap-3 text-3xl font-bold text-gray-900 dark:text-white">
        <AppIcon
          name="folder"
          :size="28"
          class="text-[var(--color-primary)]"
        />
        {{ collections.current.name }}
      </h1>

      <div
        v-if="collections.current.words.length === 0"
        class="text-gray-500 dark:text-gray-400"
      >
        В этой подборке пока нет слов. Добавляйте их со страницы «Словарь».
      </div>

      <div
        v-else
        class="grid gap-4 md:grid-cols-2 xl:grid-cols-3"
      >
        <div
          v-for="(word, index) in collections.current.words"
          :key="word.id"
          class="animate-fade-in-up relative rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:-translate-y-0.5 hover:shadow-md dark:border-white/10 dark:bg-white/5"
          :style="{ animationDelay: `${Math.min(index * 40, 300)}ms` }"
        >
          <button
            type="button"
            title="Убрать из подборки"
            class="absolute right-4 top-4 text-gray-300 transition hover:text-red-500 dark:text-gray-600 dark:hover:text-red-400"
            @click="onRemoveWord(word.id)"
          >
            <AppIcon name="x" />
          </button>

          <div class="mb-2 flex items-center gap-2 pr-8">
            <div class="font-hanzi text-2xl font-semibold text-gray-900 dark:text-white">
              {{ word.hanzi }}
            </div>

            <AudioButton
              :text="word.hanzi"
              size="sm"
            />

            <span class="inline-flex rounded-full bg-[var(--color-primary)]/15 px-3 py-1 text-xs font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]">
              HSK {{ word.hsk_level }}
            </span>
          </div>

          <div class="mb-1 text-gray-500 dark:text-gray-400">
            {{ word.pinyin }}
          </div>

          <div class="text-gray-700 dark:text-gray-300">
            {{ word.translation }}
          </div>
        </div>
      </div>
    </template>

    <div
      v-else
      class="text-gray-500 dark:text-gray-400"
    >
      Подборка не найдена.
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { useCollectionsStore } from '@/stores/collections'

import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const route = useRoute()
const collections = useCollectionsStore()

function load() {
  const id = Number(route.params.id)

  if (!Number.isNaN(id)) {
    collections.loadCollection(id)
  }
}

async function onRemoveWord(wordId: number) {
  if (!collections.current) {
    return
  }

  await collections.removeWord(collections.current.id, wordId)
}

watch(() => route.params.id, load)

onMounted(load)
</script>

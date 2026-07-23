<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Словарь
    </h1>

    <div class="mb-6 flex flex-wrap items-center gap-4">
      <BaseInput
        v-model="vocabulary.search"
        placeholder="Поиск по иероглифу, пиньиню или переводу"
        class="max-w-xs"
        @update:model-value="onSearchChange"
      />

      <select
        v-model.number="vocabulary.hsk"
        class="rounded-xl border border-white/50 bg-white/40 px-4 py-3 text-gray-900 outline-none backdrop-blur-md transition focus:border-[#41b3a3] dark:border-white/10 dark:bg-white/5 dark:text-white"
        @change="vocabulary.loadWords()"
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

      <div class="ml-auto flex gap-2 rounded-xl border border-white/50 bg-white/30 p-1 backdrop-blur-md dark:border-white/10 dark:bg-white/5">
        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="tab === 'all' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="onShowAll"
        >
          Все слова
        </button>

        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="tab === 'saved' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="onShowSaved"
        >
          Сохранённые
        </button>

        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="tab === 'learned' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="onShowLearned"
        >
          Изучено
        </button>

        <button
          class="rounded-lg px-4 py-2 text-sm font-medium transition"
          :class="tab === 'progress' ? 'bg-white/80 shadow-sm dark:bg-white/15 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
          @click="onShowInProgress"
        >
          На изучении
        </button>
      </div>
    </div>

    <template v-if="tab === 'progress'">
      <div
        v-if="learning.loadingInProgress"
        class="text-gray-500 dark:text-gray-400"
      >
        Загрузка...
      </div>

      <div
        v-else-if="learning.inProgressWords.length === 0"
        class="text-gray-500 dark:text-gray-400"
      >
        Пока нет слов на изучении. Пройдите тренировку «Слово → перевод» в разделе «Тесты», чтобы начать.
      </div>

      <div
        v-else
        class="grid gap-4 md:grid-cols-2 xl:grid-cols-3"
      >
        <div
          v-for="word in learning.inProgressWords"
          :key="word.word_id"
          class="rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
        >
          <div class="mb-2 flex items-center gap-2">
            <div class="text-2xl font-semibold text-gray-900 dark:text-white">
              {{ word.hanzi }}
            </div>

            <span class="inline-flex rounded-full bg-[#41b3a3]/15 px-3 py-1 text-xs font-medium text-[#41b3a3] dark:bg-[#41b3a3]/20 dark:text-[#85dcba]">
              HSK {{ word.hsk_level }}
            </span>
          </div>

          <div class="mb-1 text-gray-500 dark:text-gray-400">
            {{ word.pinyin }}
          </div>

          <div class="mb-4 text-gray-700 dark:text-gray-300">
            {{ word.translation }}
          </div>

          <div class="mb-2 h-2 overflow-hidden rounded-full bg-gray-200/50 dark:bg-white/10">
            <div
              class="h-full bg-[#41b3a3] transition-all duration-300"
              :style="{ width: `${(word.stage / word.max_stage) * 100}%` }"
            />
          </div>

          <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
            <span>Осталось повторений: {{ word.repetitions_left }}</span>
            <span>{{ formatTimeLeft(word.next_eligible_at) }}</span>
          </div>
        </div>
      </div>
    </template>

    <template v-else>
      <div
        v-if="loading"
        class="text-gray-500 dark:text-gray-400"
      >
        Загрузка...
      </div>

      <div
        v-else-if="visibleWords.length === 0"
        class="text-gray-500 dark:text-gray-400"
      >
        {{ emptyMessage }}
      </div>

      <div
        v-else
        class="grid gap-4 md:grid-cols-2 xl:grid-cols-3"
      >
        <div
          v-for="word in visibleWords"
          :key="word.id"
          class="relative rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
        >
          <button
            class="absolute right-4 top-4 text-xl leading-none"
            :class="isSaved(word.id) ? 'text-[#e8a87c]' : 'text-gray-300 hover:text-gray-400 dark:text-gray-600 dark:hover:text-gray-400'"
            @click="toggleSaved(word)"
          >
            {{ isSaved(word.id) ? '★' : '☆' }}
          </button>

          <div class="mb-2 flex items-center gap-2 pr-8">
            <div class="text-2xl font-semibold text-gray-900 dark:text-white">
              {{ word.hanzi }}
            </div>

            <span class="inline-flex rounded-full bg-[#41b3a3]/15 px-3 py-1 text-xs font-medium text-[#41b3a3] dark:bg-[#41b3a3]/20 dark:text-[#85dcba]">
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
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

import { useLearningStore } from '@/stores/learning'
import { useSavedWordsStore } from '@/stores/savedWords'
import { useVocabularyStore } from '@/stores/vocabulary'

import BaseInput from '@/components/base/BaseInput.vue'

import type { Word } from '@/types/word'

const vocabulary = useVocabularyStore()
const savedWords = useSavedWordsStore()
const learning = useLearningStore()

const tab = ref<'all' | 'saved' | 'learned' | 'progress'>('all')

let searchTimeout: ReturnType<typeof setTimeout> | undefined

function onSearchChange() {
  clearTimeout(searchTimeout)

  searchTimeout = setTimeout(() => vocabulary.loadWords(), 300)
}

function onShowSaved() {
  tab.value = 'saved'

  savedWords.loadSavedWords()
}

function onShowAll() {
  tab.value = 'all'

  if (vocabulary.items.length === 0) {
    vocabulary.loadWords()
  }
}

function onShowLearned() {
  tab.value = 'learned'

  learning.loadLearned()
}

function onShowInProgress() {
  tab.value = 'progress'

  learning.loadInProgress()
}

const loading = computed(() => {
  if (tab.value === 'saved') {
    return savedWords.loading
  }

  if (tab.value === 'learned') {
    return learning.loading
  }

  return vocabulary.loading
})

const visibleWords = computed<Word[]>(() => {
  if (tab.value === 'saved') {
    return savedWords.items
  }

  if (tab.value === 'learned') {
    return learning.learnedWords
  }

  return vocabulary.items
})

const emptyMessage = computed(() => {
  if (tab.value === 'saved') {
    return 'Вы ещё не сохранили ни одного слова.'
  }

  if (tab.value === 'learned') {
    return 'Пока нет изученных слов. Попробуйте тренировку «Слово → перевод» в разделе «Тесты».'
  }

  return 'Ничего не найдено.'
})

function isSaved(wordId: number) {
  return savedWords.items.some(word => word.id === wordId)
}

async function toggleSaved(word: Word) {
  if (isSaved(word.id)) {
    await savedWords.removeWord(word.id)
  } else {
    await savedWords.addWord(word)
  }
}

const now = ref(Date.now())
let clockInterval: ReturnType<typeof setInterval> | undefined

function formatTimeLeft(nextEligibleAt?: string): string {
  if (!nextEligibleAt) {
    return 'Готово к повторению'
  }

  const diffMs = new Date(nextEligibleAt).getTime() - now.value

  if (diffMs <= 0) {
    return 'Готово к повторению'
  }

  const minutes = Math.floor(diffMs / 60000)
  const days = Math.floor(minutes / 1440)
  const hours = Math.floor((minutes % 1440) / 60)
  const mins = minutes % 60

  if (days > 0) {
    return `через ${days} дн ${hours} ч`
  }

  if (hours > 0) {
    return `через ${hours} ч ${mins} мин`
  }

  return `через ${mins} мин`
}

onMounted(() => {
  vocabulary.loadWords()
  savedWords.loadSavedWords()

  clockInterval = setInterval(() => {
    now.value = Date.now()
  }, 30000)
})

onUnmounted(() => {
  clearInterval(clockInterval)
})
</script>

<template>
  <div @click="openMenuWordId = null">
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

      <BaseSelect
        v-model="vocabulary.hsk"
        class="w-56"
        :options="hskOptions"
        @update:model-value="vocabulary.loadWords()"
      />

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
        class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
      >
        <BaseSpinner />
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
          v-for="(word, index) in learning.inProgressWords"
          :key="word.word_id"
          class="animate-fade-in-up rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:-translate-y-0.5 hover:shadow-md dark:border-white/10 dark:bg-white/5"
          :style="{ animationDelay: `${Math.min(index * 40, 300)}ms` }"
        >
          <div class="mb-2 flex items-center gap-2">
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

          <div class="mb-4 text-gray-700 dark:text-gray-300">
            {{ word.translation }}
          </div>

          <div class="mb-2 h-2 overflow-hidden rounded-full bg-gray-200/50 dark:bg-white/10">
            <div
              class="h-full bg-[var(--color-primary)] transition-all duration-300"
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
        class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
      >
        <BaseSpinner />
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
          v-for="(word, index) in visibleWords"
          :key="word.id"
          class="animate-fade-in-up relative rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl transition hover:-translate-y-0.5 hover:shadow-md dark:border-white/10 dark:bg-white/5"
          :style="{ animationDelay: `${Math.min(index * 40, 300)}ms` }"
        >
          <button
            class="absolute right-4 top-4 leading-none"
            :class="isSaved(word.id) ? 'text-[var(--color-accent)]' : 'text-gray-300 hover:text-gray-400 dark:text-gray-600 dark:hover:text-gray-400'"
            @click="toggleSaved(word)"
          >
            <AppIcon
              name="star"
              :filled="isSaved(word.id)"
            />
          </button>

          <div class="absolute right-11 top-4">
            <button
              type="button"
              title="Добавить в подборку"
              class="leading-none text-gray-300 hover:text-gray-400 dark:text-gray-600 dark:hover:text-gray-400"
              @click.stop="toggleCollectionMenu(word.id)"
            >
              <AppIcon name="folder" />
            </button>

            <div
              v-if="openMenuWordId === word.id"
              class="absolute right-0 top-7 z-10 w-56 rounded-xl border border-white/50 bg-white/90 p-2 shadow-lg backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/95"
              @click.stop
            >
              <div
                v-if="collections.items.length === 0"
                class="px-2 py-1.5 text-sm text-gray-500 dark:text-gray-400"
              >
                Нет подборок
              </div>

              <button
                v-for="collection in collections.items"
                :key="collection.id"
                type="button"
                class="block w-full truncate rounded-lg px-2 py-1.5 text-left text-sm text-gray-700 hover:bg-white/60 dark:text-gray-200 dark:hover:bg-white/10"
                @click="onAddToCollection(word, collection.id)"
              >
                {{ collection.name }}
              </button>

              <div class="mt-1 flex gap-1 border-t border-white/50 pt-1 dark:border-white/10">
                <input
                  v-model="newCollectionName"
                  placeholder="Новая подборка"
                  class="min-w-0 flex-1 rounded-lg border border-white/50 bg-white/60 px-2 py-1 text-sm text-gray-900 outline-none dark:border-white/10 dark:bg-white/5 dark:text-white"
                  @keyup.enter="onCreateAndAdd(word)"
                >

                <button
                  type="button"
                  class="shrink-0 rounded-lg bg-[var(--color-primary)] px-2 text-white"
                  @click="onCreateAndAdd(word)"
                >
                  <AppIcon
                    name="plus"
                    :size="14"
                  />
                </button>
              </div>
            </div>
          </div>

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
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'

import { useCollectionsStore } from '@/stores/collections'
import { useLearningStore } from '@/stores/learning'
import { useSavedWordsStore } from '@/stores/savedWords'
import { useVocabularyStore } from '@/stores/vocabulary'

import BaseInput from '@/components/base/BaseInput.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

import type { Word } from '@/types/word'

const vocabulary = useVocabularyStore()
const savedWords = useSavedWordsStore()
const learning = useLearningStore()
const collections = useCollectionsStore()

const openMenuWordId = ref<number | null>(null)
const newCollectionName = ref('')

function toggleCollectionMenu(wordId: number) {
  openMenuWordId.value = openMenuWordId.value === wordId ? null : wordId
}

async function onAddToCollection(word: Word, collectionId: number) {
  await collections.addWord(collectionId, word)
  openMenuWordId.value = null
}

async function onCreateAndAdd(word: Word) {
  const name = newCollectionName.value.trim()

  if (!name) {
    return
  }

  const collection = await collections.create(name)

  await collections.addWord(collection.id, word)

  newCollectionName.value = ''
  openMenuWordId.value = null
}

const tab = ref<'all' | 'saved' | 'learned' | 'progress'>('all')

const hskOptions = [
  { value: 0, label: 'Любой уровень HSK' },
  ...Array.from({ length: 6 }, (_, index) => ({ value: index + 1, label: `HSK ${index + 1}` })),
]

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
  collections.loadCollections()

  clockInterval = setInterval(() => {
    now.value = Date.now()
  }, 30000)
})

onUnmounted(() => {
  clearInterval(clockInterval)
})
</script>

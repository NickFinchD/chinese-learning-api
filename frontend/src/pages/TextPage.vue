<template>
  <div>
    <div
      v-if="texts.loading"
      class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
    >
      <BaseSpinner />
      Загрузка...
    </div>

    <div
      v-else-if="texts.current"
      class="max-w-2xl"
    >
      <div class="mb-4 flex items-center justify-between gap-2">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          {{ texts.current.title }}
        </h1>

        <span class="inline-flex shrink-0 rounded-full bg-[var(--color-primary)]/15 px-3 py-1 text-xs font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]">
          HSK {{ texts.current.hsk_level }}
        </span>
      </div>

      <div class="rounded-2xl border border-white/50 bg-white/30 p-8 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-6 flex flex-wrap gap-3">
          <button
            type="button"
            :disabled="!hasVoice"
            :title="hasVoice ? '' : 'На этом устройстве не установлен голос для китайского языка'"
            class="flex items-center gap-2 rounded-full bg-white/40 px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-white/60 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
            @click="isPlaying ? stopText() : playText()"
          >
            <AppIcon
              :name="isPlaying ? 'stop' : 'volume'"
              :size="16"
            />
            {{ isPlaying ? 'Остановить' : 'Слушать текст' }}
          </button>

          <button
            type="button"
            class="flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60"
            :class="texts.current.status === 'completed'
              ? 'bg-[var(--color-primary)]/15 text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]'
              : 'bg-white/40 text-gray-700 hover:bg-white/60 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10'"
            :disabled="markingRead"
            @click="toggleRead"
          >
            <AppIcon
              name="check-circle"
              :size="16"
              :filled="texts.current.status === 'completed'"
            />
            {{ texts.current.status === 'completed' ? 'Прочитано' : 'Отметить как прочитано' }}
          </button>
        </div>

        <p
          class="font-hanzi mb-6 text-2xl leading-loose text-gray-900 dark:text-white"
          @click="activeIndex = null"
        >
          <span
            v-for="(segment, index) in segments"
            :key="index"
            :class="segment.word ? 'relative inline-block cursor-pointer border-b border-dotted ' + (isSaved(segment.word.id) ? 'border-[var(--color-accent)]' : 'border-[var(--color-primary)]') : ''"
            @click="onSegmentClick(segment, index, $event)"
          >
            {{ segment.text }}

            <span
              v-if="segment.word && activeIndex === index"
              class="absolute bottom-full left-1/2 z-10 mb-2 flex -translate-x-1/2 items-center gap-2 whitespace-nowrap rounded-lg bg-gray-900 px-3 py-2 text-sm font-normal text-white shadow-lg"
            >
              <span>{{ segment.word.pinyin }} — {{ segment.word.translation }}</span>

              <AudioButton
                :text="segment.word.hanzi"
                size="sm"
                class="!text-[var(--color-mint)] hover:!bg-white/10"
              />

              <button
                type="button"
                class="leading-none"
                :class="isSaved(segment.word.id) ? 'text-[var(--color-accent)]' : 'text-gray-400 hover:text-white'"
                @click.stop="toggleSaved(segment.word)"
              >
                <AppIcon
                  name="star"
                  :size="16"
                  :filled="isSaved(segment.word.id)"
                />
              </button>
            </span>
          </span>
        </p>

        <div class="mb-4 flex gap-3">
          <button
            class="rounded-full px-4 py-2 text-sm font-medium transition"
            :class="showPinyin ? 'bg-[var(--color-primary)]/15 text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]' : 'bg-white/40 text-gray-700 hover:bg-white/60 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10'"
            @click="showPinyin = !showPinyin"
          >
            Пиньинь
          </button>

          <button
            class="rounded-full px-4 py-2 text-sm font-medium transition"
            :class="showTranslation ? 'bg-[var(--color-primary)]/15 text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]' : 'bg-white/40 text-gray-700 hover:bg-white/60 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10'"
            @click="showTranslation = !showTranslation"
          >
            Перевод
          </button>
        </div>

        <p
          v-if="showPinyin"
          class="mb-4 text-lg text-gray-500 dark:text-gray-400"
        >
          {{ texts.current.pinyin }}
        </p>

        <p
          v-if="showTranslation"
          class="text-lg text-gray-700 dark:text-gray-300"
        >
          {{ texts.current.translation }}
        </p>
      </div>

      <RouterLink
        to="/app/texts"
        class="mt-6 inline-block rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
      >
        Назад к текстам
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { getWords } from '@/services/words'
import { useSavedWordsStore } from '@/stores/savedWords'
import { useTextsStore } from '@/stores/texts'
import { isSpeechSupported, useChineseVoiceAvailable } from '@/utils/speech'
import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

import type { Word } from '@/types/word'

interface Segment {
  text: string
  word: Word | null
}

const route = useRoute()
const texts = useTextsStore()
const savedWords = useSavedWordsStore()

function isSaved(wordId: number) {
  return savedWords.items.some(word => word.id === wordId)
}

const activeIndex = ref<number | null>(null)

function onSegmentClick(segment: Segment, index: number, event: MouseEvent) {
  if (!segment.word) {
    return
  }

  // Stop the click from bubbling to the paragraph's "close on click elsewhere"
  // handler, but only for actual word segments — punctuation/plain text should
  // still bubble up and close whatever popover is open.
  event.stopPropagation()

  activeIndex.value = activeIndex.value === index ? null : index
}

async function toggleSaved(word: Word) {
  try {
    if (isSaved(word.id)) {
      await savedWords.removeWord(word.id)
    } else {
      await savedWords.addWord(word)
    }
  } catch (error) {
    console.error('Failed to toggle saved word:', error)
  }
}

const isPlaying = ref(false)
const hasVoice = useChineseVoiceAvailable()

function playText() {
  if (!isSpeechSupported() || !texts.current) {
    return
  }

  window.speechSynthesis.cancel()

  const utterance = new SpeechSynthesisUtterance(texts.current.hanzi)

  utterance.lang = 'zh-CN'
  utterance.rate = 0.85
  utterance.onend = () => { isPlaying.value = false }
  utterance.onerror = () => { isPlaying.value = false }

  window.speechSynthesis.speak(utterance)
  isPlaying.value = true
}

function stopText() {
  window.speechSynthesis.cancel()
  isPlaying.value = false
}

const markingRead = ref(false)

async function toggleRead() {
  if (!texts.current) {
    return
  }

  markingRead.value = true

  try {
    if (texts.current.status === 'completed') {
      await texts.markUnread(texts.current.id)
    } else {
      await texts.markRead(texts.current.id)
    }
  } catch (error) {
    console.error('Failed to update read status:', error)
  } finally {
    markingRead.value = false
  }
}

const showPinyin = ref(false)
const showTranslation = ref(false)

const wordMap = ref(new Map<string, Word>())
const maxWordLength = ref(1)

// Greedy longest-match tokenizer: at each position, try the longest
// substring first and fall back to shorter ones, since Chinese text has
// no spaces between words.
const segments = computed<Segment[]>(() => {
  const hanzi = texts.current?.hanzi

  if (!hanzi) {
    return []
  }

  const chars = [...hanzi]
  const result: Segment[] = []

  let i = 0

  while (i < chars.length) {

    let matched = false

    for (let len = Math.min(maxWordLength.value, chars.length - i); len >= 1; len--) {

      const candidate = chars.slice(i, i + len).join('')
      const word = wordMap.value.get(candidate)

      if (word) {
        result.push({ text: candidate, word })
        i += len
        matched = true
        break
      }
    }

    if (!matched) {
      result.push({ text: chars[i], word: null })
      i += 1
    }
  }

  return result
})

onMounted(async () => {
  showPinyin.value = false
  showTranslation.value = false

  texts.loadText(Number(route.params.id))
  savedWords.loadSavedWords()

  try {
    const response = await getWords()
    const words = response.data ?? []

    wordMap.value = new Map(words.map(word => [word.hanzi, word]))
    maxWordLength.value = words.reduce((max, word) => Math.max(max, [...word.hanzi].length), 1)
  } catch (error) {
    console.error('Failed to load dictionary for hover hints:', error)
  }
})

onUnmounted(() => {
  if (isSpeechSupported()) {
    window.speechSynthesis.cancel()
  }
})
</script>

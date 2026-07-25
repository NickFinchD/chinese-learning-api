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
            :class="segmentClass(segment)"
            @click="onSegmentClick(segment, index, $event)"
          >{{ segment.text }}</span>
        </p>

        <Teleport to="body">
          <div
            v-if="activeSegment"
            ref="tooltipRef"
            class="fixed z-50 flex items-center gap-2 whitespace-nowrap rounded-lg bg-gray-900 px-3 py-2 text-sm font-normal text-white shadow-lg"
            :style="{ top: `${tooltipPos.top}px`, left: `${tooltipPos.left}px` }"
            @click.stop
          >
            <span>{{ (activeSegment.word ?? activeSegment.name)!.pinyin }} — {{ (activeSegment.word ?? activeSegment.name)!.translation }}</span>

            <AudioButton
              :text="activeSegment.text"
              size="sm"
              class="!text-[var(--color-mint)] hover:!bg-white/10"
            />

            <button
              v-if="activeSegment.word"
              type="button"
              class="leading-none"
              :class="isSaved(activeSegment.word.id) ? 'text-[var(--color-accent)]' : 'text-gray-400 hover:text-white'"
              @click.stop="toggleSaved(activeSegment.word)"
            >
              <AppIcon
                name="star"
                :size="16"
                :filled="isSaved(activeSegment.word.id)"
              />
            </button>
          </div>
        </Teleport>

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
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { getWords } from '@/services/words'
import { getNames } from '@/services/names'
import { useSavedWordsStore } from '@/stores/savedWords'
import { useTextsStore } from '@/stores/texts'
import { isSpeechSupported, useChineseVoiceAvailable } from '@/utils/speech'
import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

import type { Word } from '@/types/word'
import type { ProperName } from '@/types/name'

interface Segment {
  text: string
  word: Word | null
  // Proper nouns (people's names, book titles): get the same hover hint as
  // a dictionary word, but deliberately aren't part of `words` — no save
  // button, never trainable/quizzable. See the `names` table/endpoint.
  name: ProperName | null
}

function segmentClass(segment: Segment): string {
  if (segment.word) {
    return 'relative inline-block cursor-pointer border-b border-dotted ' +
      (isSaved(segment.word.id) ? 'border-[var(--color-accent)]' : 'border-[var(--color-primary)]')
  }

  if (segment.name) {
    return 'relative inline-block cursor-pointer border-b border-dotted border-gray-400 dark:border-gray-500'
  }

  return ''
}

const route = useRoute()
const texts = useTextsStore()
const savedWords = useSavedWordsStore()

function isSaved(wordId: number) {
  return savedWords.items.some(word => word.id === wordId)
}

const activeIndex = ref<number | null>(null)
const activeSegment = computed(() => activeIndex.value !== null ? segments.value[activeIndex.value] ?? null : null)

// Teleported to <body> and positioned in viewport coordinates (rather than
// CSS-anchored to the word span) so it can never be clipped by <main>'s
// overflow-y-auto — which, per spec, forces overflow-x to auto too, cutting
// off anything that pokes out past main's left edge (right where the
// sidebar sits) when a word near the margin has a long pinyin/translation.
const tooltipRef = ref<HTMLElement | null>(null)
const tooltipPos = ref({ top: 0, left: 0 })

const TOOLTIP_MARGIN = 8

function positionTooltip(anchor: HTMLElement) {
  const rect = anchor.getBoundingClientRect()
  const width = tooltipRef.value?.offsetWidth ?? 0
  const height = tooltipRef.value?.offsetHeight ?? 0

  const left = Math.min(
    Math.max(TOOLTIP_MARGIN, rect.left + rect.width / 2 - width / 2),
    window.innerWidth - width - TOOLTIP_MARGIN,
  )

  tooltipPos.value = {
    left,
    top: Math.max(TOOLTIP_MARGIN, rect.top - height - TOOLTIP_MARGIN),
  }
}

async function onSegmentClick(segment: Segment, index: number, event: MouseEvent) {
  if (!segment.word && !segment.name) {
    return
  }

  // Stop the click from bubbling to the paragraph's "close on click elsewhere"
  // handler, but only for actual word segments — punctuation/plain text should
  // still bubble up and close whatever popover is open.
  event.stopPropagation()

  if (activeIndex.value === index) {
    activeIndex.value = null
    return
  }

  activeIndex.value = index

  await nextTick()
  positionTooltip(event.currentTarget as HTMLElement)
}

function closeTooltip() {
  activeIndex.value = null
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
const nameMap = ref(new Map<string, ProperName>())
const maxMatchLength = ref(1)

// Greedy longest-match tokenizer: at each position, try the longest
// substring first and fall back to shorter ones, since Chinese text has
// no spaces between words. Checks the real dictionary first, then the
// proper-noun list, so a hanzi that somehow existed in both would resolve
// to the trainable word.
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

    for (let len = Math.min(maxMatchLength.value, chars.length - i); len >= 1; len--) {

      const candidate = chars.slice(i, i + len).join('')
      const word = wordMap.value.get(candidate)
      const name = nameMap.value.get(candidate)

      if (word || name) {
        result.push({ text: candidate, word: word ?? null, name: word ? null : (name ?? null) })
        i += len
        matched = true
        break
      }
    }

    if (!matched) {
      result.push({ text: chars[i], word: null, name: null })
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
    const [wordsResponse, namesResponse] = await Promise.all([getWords(), getNames()])
    const words = wordsResponse.data ?? []
    const properNames = namesResponse.data ?? []

    wordMap.value = new Map(words.map(word => [word.hanzi, word]))
    nameMap.value = new Map(properNames.map(name => [name.hanzi, name]))

    maxMatchLength.value = [...words, ...properNames]
      .reduce((max, entry) => Math.max(max, [...entry.hanzi].length), 1)
  } catch (error) {
    console.error('Failed to load dictionary for hover hints:', error)
  }

  // The tooltip is positioned in viewport coordinates on open and doesn't
  // track the anchor afterwards, so close it rather than let it drift out
  // of place if the text scrolls or the window resizes.
  window.addEventListener('scroll', closeTooltip, true)
  window.addEventListener('resize', closeTooltip)
})

onUnmounted(() => {
  if (isSpeechSupported()) {
    window.speechSynthesis.cancel()
  }

  window.removeEventListener('scroll', closeTooltip, true)
  window.removeEventListener('resize', closeTooltip)
})
</script>

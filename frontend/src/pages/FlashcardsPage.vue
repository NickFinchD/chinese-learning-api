<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Карточки
    </h1>

    <div
      v-if="!flash.started && !flash.noWords"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-600 dark:text-gray-400">
        Повторение слов, которые вы сейчас изучаете. Иероглиф → переверните карточку, чтобы увидеть перевод → отметьте «Знаю» или «Не знаю». Неизвестные слова вернутся в конец очереди.
      </p>

      <button
        class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="flash.loading"
        @click="flash.start()"
      >
        Начать
      </button>
    </div>

    <div
      v-else-if="flash.noWords"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500 dark:text-gray-400">
        Нет слов «в изучении». Начните изучать слова в курсах или тренировках — тогда они появятся здесь.
      </p>

      <RouterLink
        to="/app/vocabulary"
        class="inline-block rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
      >
        Перейти в словарь
      </RouterLink>
    </div>

    <div
      v-else-if="flash.isFinished"
      class="max-w-md"
    >
      <div class="animate-pop-in rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-2 flex items-center justify-center gap-2 text-xl font-semibold text-gray-900 dark:text-white">
          <AppIcon
            name="sparkles"
            :size="22"
            class="text-[var(--color-accent)]"
          />
          Все карточки пройдены
        </div>

        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Слов повторено: {{ flash.totalCount }}
        </p>

        <button
          class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
          @click="flash.reset()"
        >
          Начать заново
        </button>
      </div>
    </div>

    <div
      v-else-if="flash.current"
      class="max-w-md"
    >
      <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">
        Осталось: {{ flash.queue.length }}
      </div>

      <div
        :key="flash.current.word_id"
        class="animate-fade-in-up cursor-pointer select-none rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl transition hover:shadow-md dark:border-white/10 dark:bg-white/5"
        @click="flash.flip()"
      >
        <div class="flex items-center justify-center gap-2">
          <div class="font-hanzi text-5xl font-bold text-gray-900 dark:text-white">
            {{ flash.current.hanzi }}
          </div>

          <AudioButton :text="flash.current.hanzi" />
        </div>

        <div
          v-if="flash.flipped"
          class="animate-pop-in mt-4"
        >
          <div class="text-lg text-gray-500 dark:text-gray-400">
            {{ flash.current.pinyin }}
          </div>

          <div class="mt-1 text-xl font-medium text-gray-800 dark:text-gray-200">
            {{ flash.current.translation }}
          </div>
        </div>

        <div
          v-else
          class="mt-4 text-sm text-gray-400 dark:text-gray-500"
        >
          Нажмите, чтобы увидеть перевод
        </div>
      </div>

      <div class="mt-6 flex gap-3">
        <button
          class="flex-1 rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
          @click="flash.dontKnow()"
        >
          Не знаю
        </button>

        <button
          class="flex-1 rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
          @click="flash.know()"
        >
          Знаю
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'

import { useFlashcardsStore } from '@/stores/flashcards'
import AppIcon from '@/components/base/AppIcon.vue'
import AudioButton from '@/components/base/AudioButton.vue'

const flash = useFlashcardsStore()

onMounted(() => {
  if (!flash.started) {
    flash.reset()
  }
})
</script>

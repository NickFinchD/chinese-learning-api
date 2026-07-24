<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Грамматика по HSK
    </h1>

    <div
      v-if="!test.started"
      class="max-w-md"
    >
      <p class="mb-4 text-gray-600 dark:text-gray-400">
        Выберите уровень HSK:
      </p>

      <div class="mb-6 flex flex-wrap gap-2">
        <button
          v-for="level in 6"
          :key="level"
          class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] dark:border-white/10 dark:bg-white/5 dark:text-gray-300"
          :disabled="test.loading"
          @click="test.start(level)"
        >
          HSK {{ level }}
        </button>
      </div>

      <div
        v-if="test.loading"
        class="flex items-center gap-2 text-gray-500 dark:text-gray-400"
      >
        <BaseSpinner />
        Загрузка...
      </div>
    </div>

    <div
      v-else-if="test.quizzes.length === 0"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500 dark:text-gray-400">
        Для уровня HSK {{ test.hsk }} пока нет тестов.
      </p>

      <button
        class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
        @click="test.reset()"
      >
        Выбрать другой уровень
      </button>
    </div>

    <div
      v-else-if="test.isFinished"
      class="max-w-md"
    >
      <div class="animate-pop-in rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-2 flex items-center justify-center gap-2 text-xl font-semibold text-gray-900 dark:text-white">
          <AppIcon
            name="sparkles"
            :size="22"
            class="text-[var(--color-accent)]"
          />
          Тест завершён
        </div>

        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Правильных ответов: {{ test.correctCount }} из {{ test.quizzes.length }}
        </p>

        <button
          class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
          @click="test.reset()"
        >
          Выбрать другой уровень
        </button>
      </div>
    </div>

    <div
      v-else-if="test.currentQuiz"
      class="max-w-md"
    >
      <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">
        Вопрос {{ test.currentIndex + 1 }} из {{ test.quizzes.length }} · HSK {{ test.hsk }}
      </div>

      <div
        :key="test.currentIndex"
        class="animate-fade-in-up rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
      >
        <h2 class="font-hanzi mb-4 text-lg font-semibold text-gray-900 dark:text-white">
          {{ test.currentQuiz.question }}
        </h2>

        <div class="space-y-2">
          <button
            v-for="option in test.currentQuiz.options"
            :key="option.id"
            class="w-full rounded-lg border border-white/50 bg-white/20 p-3 text-left text-gray-800 backdrop-blur-md transition hover:bg-white/40 disabled:cursor-not-allowed dark:border-white/10 dark:bg-white/5 dark:text-gray-200 dark:hover:bg-white/10"
            :disabled="test.answeredResult !== null"
            @click="test.answer(option.id)"
          >
            {{ option.text }}
          </button>
        </div>

        <div
          v-if="test.answeredResult !== null"
          class="animate-pop-in mt-4 text-lg font-semibold"
        >
          <span
            v-if="test.answeredResult"
            class="flex items-center gap-2 text-green-600 dark:text-green-400"
          >
            <AppIcon name="check-circle" />
            Верно!
          </span>

          <span
            v-else
            class="flex items-center gap-2 text-red-600 dark:text-red-400"
          >
            <AppIcon name="x-circle" />
            Неверно
          </span>
        </div>
      </div>

      <button
        v-if="test.answeredResult !== null"
        class="mt-6 w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
        @click="test.next()"
      >
        Далее
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import { useGrammarTestStore } from '@/stores/grammarTest'
import AppIcon from '@/components/base/AppIcon.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const test = useGrammarTestStore()

// The store survives navigation, so returning to this page after leaving
// mid-test would otherwise resume the old level instead of showing the
// picker — start every visit fresh.
onMounted(() => {
  test.reset()
})
</script>

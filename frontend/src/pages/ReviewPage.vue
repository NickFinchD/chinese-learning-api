<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold">
      Повторение
    </h1>

    <div
      v-if="review.loading && !review.sessionStarted"
      class="text-gray-500"
    >
      Загрузка...
    </div>

    <template v-else>
      <div
        v-if="!review.sessionStarted"
        class="max-w-md"
      >
        <div class="mb-6 grid grid-cols-3 gap-4">
          <div class="rounded-xl border border-gray-200 bg-white p-4 text-center shadow-sm">
            <div class="text-2xl font-bold">
              {{ review.statistics?.total_words ?? 0 }}
            </div>
            <div class="text-sm text-gray-500">
              Всего
            </div>
          </div>

          <div class="rounded-xl border border-gray-200 bg-white p-4 text-center shadow-sm">
            <div class="text-2xl font-bold text-blue-600">
              {{ review.statistics?.ready_for_review ?? 0 }}
            </div>
            <div class="text-sm text-gray-500">
              Готово
            </div>
          </div>

          <div class="rounded-xl border border-gray-200 bg-white p-4 text-center shadow-sm">
            <div class="text-2xl font-bold text-green-600">
              {{ review.statistics?.reviewed_words ?? 0 }}
            </div>
            <div class="text-sm text-gray-500">
              Повторено
            </div>
          </div>
        </div>

        <button
          class="w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="!review.statistics?.ready_for_review"
          @click="review.startSession()"
        >
          Начать повторение
        </button>
      </div>

      <div
        v-else-if="review.isFinished"
        class="max-w-md"
      >
        <div class="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm">
          <div class="mb-2 text-xl font-semibold">
            Сессия завершена 🎉
          </div>

          <p class="mb-6 text-gray-600">
            Вы повторили слов: {{ review.words.length }}.
          </p>

          <button
            class="w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
            @click="review.resetSession()"
          >
            Назад к обзору
          </button>
        </div>
      </div>

      <div
        v-else-if="review.currentWord"
        class="max-w-md"
      >
        <div class="mb-2 text-sm text-gray-500">
          Слово {{ review.currentIndex + 1 }} из {{ review.words.length }}
        </div>

        <div class="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm">
          <div class="mb-4 text-4xl font-bold">
            {{ review.currentWord.hanzi }}
          </div>

          <div v-if="review.showAnswer">
            <div class="mb-1 text-gray-500">
              {{ review.currentWord.pinyin }}
            </div>

            <div class="text-lg text-gray-700">
              {{ review.currentWord.translation }}
            </div>
          </div>

          <button
            v-else
            class="rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-200"
            @click="review.reveal()"
          >
            Показать ответ
          </button>
        </div>

        <div
          v-if="review.showAnswer"
          class="mt-6 flex gap-4"
        >
          <button
            class="flex-1 rounded-xl bg-red-100 px-4 py-3 font-semibold text-red-700 transition hover:bg-red-200"
            @click="review.answer(false)"
          >
            Неверно
          </button>

          <button
            class="flex-1 rounded-xl bg-green-100 px-4 py-3 font-semibold text-green-700 transition hover:bg-green-200"
            @click="review.answer(true)"
          >
            Верно
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import { useReviewStore } from '@/stores/review'

const review = useReviewStore()

onMounted(() => {
  review.loadStatistics()
})
</script>

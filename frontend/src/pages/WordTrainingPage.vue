<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold">
      Слово → перевод
    </h1>

    <div
      v-if="!training.started && !training.notEnoughWords && !training.allLearned"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-600">
        Тренировка на основе слов, сохранённых в «Словаре».
      </p>

      <button
        class="w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="training.loading"
        @click="training.start()"
      >
        Начать тренировку
      </button>
    </div>

    <div
      v-else-if="training.notEnoughWords"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500">
        Нужно сохранить хотя бы 2 слова, чтобы начать тренировку.
      </p>

      <div class="flex flex-wrap gap-4">
        <button
          class="rounded-xl bg-gray-100 px-4 py-3 font-semibold text-gray-700 transition hover:bg-gray-200"
          :disabled="training.loading"
          @click="training.start()"
        >
          Проверить ещё раз
        </button>

        <RouterLink
          to="/app/vocabulary"
          class="inline-block rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
        >
          Перейти в словарь
        </RouterLink>
      </div>
    </div>

    <div
      v-else-if="training.allLearned"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500">
        Все сохранённые слова уже изучены 🎉 Сохраните новые слова в «Словаре», чтобы продолжить тренировку.
      </p>

      <div class="flex flex-wrap gap-4">
        <button
          class="rounded-xl bg-gray-100 px-4 py-3 font-semibold text-gray-700 transition hover:bg-gray-200"
          :disabled="training.loading"
          @click="training.start()"
        >
          Проверить ещё раз
        </button>

        <RouterLink
          to="/app/vocabulary"
          class="inline-block rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
        >
          Перейти в словарь
        </RouterLink>
      </div>
    </div>

    <div
      v-else-if="training.isFinished"
      class="max-w-md"
    >
      <div class="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm">
        <div class="mb-2 text-xl font-semibold">
          Тренировка завершена 🎉
        </div>

        <p class="mb-6 text-gray-600">
          Правильных ответов: {{ training.correctCount }} из {{ training.questions.length }}
        </p>

        <button
          class="w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
          @click="training.reset()"
        >
          Пройти ещё раз
        </button>
      </div>
    </div>

    <div
      v-else-if="training.currentQuestion"
      class="max-w-md"
    >
      <div class="mb-2 text-sm text-gray-500">
        Слово {{ training.currentIndex + 1 }} из {{ training.questions.length }}
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6 text-center shadow-sm">
        <div class="mb-1 text-4xl font-bold">
          {{ training.currentQuestion.word.hanzi }}
        </div>

        <div class="mb-6 text-gray-500">
          {{ training.currentQuestion.word.pinyin }}
        </div>

        <div class="space-y-2 text-left">
          <button
            v-for="option in training.currentQuestion.options"
            :key="option.id"
            class="w-full rounded-lg border p-3 transition disabled:cursor-not-allowed"
            :class="optionClass(option.id)"
            :disabled="training.answeredWordId !== null"
            @click="training.answer(option.id)"
          >
            {{ option.translation }}
          </button>
        </div>

        <div
          v-if="training.lastProgress"
          class="mt-4 text-sm"
        >
          <span
            v-if="training.lastProgress.learned"
            class="font-semibold text-green-600"
          >
            🎉 Слово изучено! Смотрите вкладку «Изучено» в словаре.
          </span>

          <span
            v-else
            class="text-gray-500"
          >
            Этап {{ training.lastProgress.stage }} из {{ training.lastProgress.max_stage }}
          </span>
        </div>
      </div>

      <button
        v-if="training.answeredWordId !== null"
        class="mt-6 w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
        @click="training.next()"
      >
        Далее
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'

import { useWordTrainingStore } from '@/stores/wordTraining'

const training = useWordTrainingStore()

onMounted(() => {
  if (!training.started) {
    training.reset()
  }
})

function optionClass(optionId: number) {
  const answered = training.answeredWordId

  if (answered === null) {
    return 'hover:bg-gray-50'
  }

  const correctId = training.currentQuestion?.word.id

  if (optionId === correctId) {
    return 'border-green-500 bg-green-50'
  }

  if (optionId === answered) {
    return 'border-red-500 bg-red-50'
  }

  return 'opacity-60'
}
</script>

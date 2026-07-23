<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Слово → перевод
    </h1>

    <div
      v-if="!training.started && !training.notEnoughWords && !training.allLearned"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-600 dark:text-gray-400">
        Тренировка на основе слов, сохранённых в «Словаре».
      </p>

      <button
        class="w-full rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90 disabled:cursor-not-allowed disabled:opacity-50"
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
      <p class="mb-6 text-gray-500 dark:text-gray-400">
        Нужно сохранить хотя бы 2 слова, чтобы начать тренировку.
      </p>

      <div class="flex flex-wrap gap-4">
        <button
          class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="training.loading"
          @click="training.start()"
        >
          Проверить ещё раз
        </button>

        <RouterLink
          to="/app/vocabulary"
          class="inline-block rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
        >
          Перейти в словарь
        </RouterLink>
      </div>
    </div>

    <div
      v-else-if="training.allLearned"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500 dark:text-gray-400">
        Все сохранённые слова уже изучены 🎉 Сохраните новые слова в «Словаре», чтобы продолжить тренировку.
      </p>

      <div class="flex flex-wrap gap-4">
        <button
          class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="training.loading"
          @click="training.start()"
        >
          Проверить ещё раз
        </button>

        <RouterLink
          to="/app/vocabulary"
          class="inline-block rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
        >
          Перейти в словарь
        </RouterLink>
      </div>
    </div>

    <div
      v-else-if="training.isFinished"
      class="max-w-md"
    >
      <div class="rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          Тренировка завершена 🎉
        </div>

        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Правильных ответов: {{ training.correctCount }} из {{ training.questions.length }}
        </p>

        <button
          class="w-full rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
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
      <div class="mb-2 text-sm text-gray-500 dark:text-gray-400">
        Слово {{ training.currentIndex + 1 }} из {{ training.questions.length }}
      </div>

      <div class="rounded-xl border border-white/50 bg-white/30 p-6 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-1 text-4xl font-bold text-gray-900 dark:text-white">
          {{ training.currentQuestion.word.hanzi }}
        </div>

        <div class="mb-6 text-gray-500 dark:text-gray-400">
          {{ training.currentQuestion.word.pinyin }}
        </div>

        <div class="space-y-2 text-left">
          <button
            v-for="option in training.currentQuestion.options"
            :key="option.id"
            class="w-full rounded-lg border p-3 text-gray-800 backdrop-blur-md transition disabled:cursor-not-allowed dark:text-gray-200"
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
            class="font-semibold text-green-600 dark:text-green-400"
          >
            🎉 Слово изучено! Смотрите вкладку «Изучено» в словаре.
          </span>

          <span
            v-else
            class="text-gray-500 dark:text-gray-400"
          >
            Этап {{ training.lastProgress.stage }} из {{ training.lastProgress.max_stage }}
          </span>
        </div>
      </div>

      <button
        v-if="training.answeredWordId !== null"
        class="mt-6 w-full rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
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
    return 'border-white/50 bg-white/20 hover:bg-white/40 dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10'
  }

  const correctId = training.currentQuestion?.word.id

  if (optionId === correctId) {
    return 'border-green-500 bg-green-50 dark:border-green-500/50 dark:bg-green-500/10'
  }

  if (optionId === answered) {
    return 'border-red-500 bg-red-50 dark:border-red-500/50 dark:bg-red-500/10'
  }

  return 'border-white/50 bg-white/20 opacity-60 dark:border-white/10 dark:bg-white/5'
}
</script>

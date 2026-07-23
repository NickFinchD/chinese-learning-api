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
          class="rounded-full border border-white/50 bg-white/30 px-4 py-3 font-semibold text-gray-700 backdrop-blur-md transition hover:border-[#41b3a3] hover:text-[#41b3a3] dark:border-white/10 dark:bg-white/5 dark:text-gray-300"
          :disabled="test.loading"
          @click="test.start(level)"
        >
          HSK {{ level }}
        </button>
      </div>

      <div
        v-if="test.loading"
        class="text-gray-500 dark:text-gray-400"
      >
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
      <div class="rounded-xl border border-white/50 bg-white/30 p-8 text-center shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div class="mb-2 text-xl font-semibold text-gray-900 dark:text-white">
          Тест завершён 🎉
        </div>

        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Правильных ответов: {{ test.correctCount }} из {{ test.quizzes.length }}
        </p>

        <button
          class="w-full rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
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

      <div class="rounded-xl border border-white/50 bg-white/30 p-6 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
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
          class="mt-4 text-lg font-semibold"
        >
          <span
            v-if="test.answeredResult"
            class="text-green-600 dark:text-green-400"
          >
            ✅ Верно!
          </span>

          <span
            v-else
            class="text-red-600 dark:text-red-400"
          >
            ❌ Неверно
          </span>
        </div>
      </div>

      <button
        v-if="test.answeredResult !== null"
        class="mt-6 w-full rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90"
        @click="test.next()"
      >
        Далее
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGrammarTestStore } from '@/stores/grammarTest'

const test = useGrammarTestStore()
</script>

<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold">
      Грамматика по HSK
    </h1>

    <div
      v-if="!test.started"
      class="max-w-md"
    >
      <p class="mb-4 text-gray-600">
        Выберите уровень HSK:
      </p>

      <div class="mb-6 flex flex-wrap gap-2">
        <button
          v-for="level in 6"
          :key="level"
          class="rounded-xl border border-gray-300 px-4 py-3 font-semibold transition hover:border-blue-500 hover:text-blue-600"
          :disabled="test.loading"
          @click="test.start(level)"
        >
          HSK {{ level }}
        </button>
      </div>

      <div
        v-if="test.loading"
        class="text-gray-500"
      >
        Загрузка...
      </div>
    </div>

    <div
      v-else-if="test.quizzes.length === 0"
      class="max-w-md"
    >
      <p class="mb-6 text-gray-500">
        Для уровня HSK {{ test.hsk }} пока нет тестов.
      </p>

      <button
        class="rounded-xl bg-gray-100 px-4 py-3 font-semibold text-gray-700 transition hover:bg-gray-200"
        @click="test.reset()"
      >
        Выбрать другой уровень
      </button>
    </div>

    <div
      v-else-if="test.isFinished"
      class="max-w-md"
    >
      <div class="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm">
        <div class="mb-2 text-xl font-semibold">
          Тест завершён 🎉
        </div>

        <p class="mb-6 text-gray-600">
          Правильных ответов: {{ test.correctCount }} из {{ test.quizzes.length }}
        </p>

        <button
          class="w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
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
      <div class="mb-2 text-sm text-gray-500">
        Вопрос {{ test.currentIndex + 1 }} из {{ test.quizzes.length }} · HSK {{ test.hsk }}
      </div>

      <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <h2 class="mb-4 text-lg font-semibold">
          {{ test.currentQuiz.question }}
        </h2>

        <div class="space-y-2">
          <button
            v-for="option in test.currentQuiz.options"
            :key="option.id"
            class="w-full rounded-lg border p-3 text-left transition hover:bg-gray-50 disabled:cursor-not-allowed"
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
            class="text-green-600"
          >
            ✅ Верно!
          </span>

          <span
            v-else
            class="text-red-600"
          >
            ❌ Неверно
          </span>
        </div>
      </div>

      <button
        v-if="test.answeredResult !== null"
        class="mt-6 w-full rounded-xl bg-blue-600 px-4 py-3 font-semibold text-white transition hover:bg-blue-700"
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

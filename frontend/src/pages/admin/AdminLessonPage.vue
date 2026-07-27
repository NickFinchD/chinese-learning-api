<template>
  <div>
    <div class="mb-8 flex items-center gap-3">
      <RouterLink
        :to="lesson ? { name: 'admin-course', params: { id: lesson.course_id } } : '/app/admin/courses'"
        class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
      >
        <AppIcon name="arrow-left" />
      </RouterLink>

      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
        {{ lesson ? `${lesson.lesson_number}. ${lesson.title}` : 'Загрузка...' }}
      </h1>
    </div>

    <div
      v-if="loading"
      class="flex justify-center py-12"
    >
      <BaseSpinner :size="28" />
    </div>

    <div
      v-else
      class="max-w-2xl space-y-2"
    >
      <div
        v-for="(step, index) in steps"
        :key="step.id"
        class="flex items-center gap-3 rounded-xl border border-white/50 bg-white/30 p-3 backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
      >
        <span class="shrink-0 rounded-full bg-[var(--color-primary)]/15 px-2.5 py-1 text-xs font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]">
          {{ stepTypeLabel[step.step_type] }}
        </span>

        <span class="min-w-0 flex-1 truncate text-gray-800 dark:text-gray-200">
          {{ step.label || `#${step.entity_id}` }}
        </span>

        <button
          type="button"
          title="Вверх"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="index === 0"
          @click="moveStep(index, -1)"
        >
          <AppIcon
            name="chevron-down"
            :size="16"
            class="rotate-180"
          />
        </button>

        <button
          type="button"
          title="Вниз"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="index === steps.length - 1"
          @click="moveStep(index, 1)"
        >
          <AppIcon
            name="chevron-down"
            :size="16"
          />
        </button>

        <button
          type="button"
          title="Удалить"
          class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-red-500 transition hover:bg-red-500/10"
          @click="removeStep(step)"
        >
          <AppIcon
            name="trash"
            :size="16"
          />
        </button>
      </div>

      <p
        v-if="steps.length === 0"
        class="py-10 text-center text-gray-500 dark:text-gray-400"
      >
        В этом уроке пока нет шагов
      </p>

      <button
        type="button"
        class="mt-2 rounded-full bg-[var(--color-primary)] px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90"
        @click="openAddStep"
      >
        + Добавить шаг
      </button>

      <p
        v-if="error"
        class="text-sm text-red-600 dark:text-red-400"
      >
        {{ error }}
      </p>
    </div>

    <!-- Add-step modal: pick a step type, then search and click an entity. -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition-opacity duration-150"
        enter-from-class="opacity-0"
        leave-active-class="transition-opacity duration-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="addStepOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
          @click.self="addStepOpen = false"
        >
          <div class="max-h-[80vh] w-full max-w-lg overflow-y-auto rounded-2xl border border-white/50 bg-white/90 p-6 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/95">
            <div class="mb-4 flex items-center justify-between">
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                Добавить шаг
              </h2>

              <button
                type="button"
                class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition hover:bg-white/60 dark:text-gray-400 dark:hover:bg-white/10"
                @click="addStepOpen = false"
              >
                <AppIcon name="x" />
              </button>
            </div>

            <div class="mb-4 flex flex-wrap gap-2">
              <button
                v-for="type in stepTypes"
                :key="type"
                type="button"
                class="rounded-full px-4 py-2 text-sm font-medium transition"
                :class="pickerType === type
                  ? 'bg-[var(--color-primary)] text-white'
                  : 'bg-white/40 text-gray-700 hover:bg-white/60 dark:bg-white/5 dark:text-gray-300 dark:hover:bg-white/10'"
                @click="selectPickerType(type)"
              >
                {{ stepTypeLabel[type] }}
              </button>
            </div>

            <BaseInput
              v-model="pickerSearch"
              placeholder="Поиск..."
            />

            <div class="mt-3 max-h-80 space-y-1 overflow-y-auto">
              <button
                v-for="result in pickerResults"
                :key="result.id"
                type="button"
                class="block w-full truncate rounded-lg p-2.5 text-left text-sm text-gray-800 transition hover:bg-white/60 dark:text-gray-200 dark:hover:bg-white/10"
                @click="pickEntity(result.id)"
              >
                {{ result.label }}
              </button>

              <p
                v-if="!pickerLoading && pickerResults.length === 0"
                class="p-2.5 text-sm text-gray-500 dark:text-gray-400"
              >
                Ничего не найдено
              </p>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

import { getLesson } from '@/services/lessons'
import { addLessonStep, getLessonSteps, removeLessonStep, reorderLessonSteps } from '@/services/admin/lessons'
import type { AdminStep, StepType } from '@/services/admin/lessons'
import { adminListWords } from '@/services/admin/words'
import { adminListQuizzes } from '@/services/admin/quizzes'
import { adminListGrammarNotes } from '@/services/admin/grammar'
import { adminListSentences } from '@/services/admin/sentences'
import type { Lesson } from '@/types/lesson'

import AppIcon from '@/components/base/AppIcon.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseSpinner from '@/components/base/BaseSpinner.vue'

const route = useRoute()
const lessonId = Number(route.params.id)

const lesson = ref<Lesson | null>(null)
const steps = ref<AdminStep[]>([])
const loading = ref(true)
const error = ref('')

const stepTypes: StepType[] = ['word', 'quiz', 'grammar', 'sentence_builder']
const stepTypeLabel: Record<StepType, string> = {
  word: 'Слово',
  quiz: 'Тест',
  grammar: 'Грамматика',
  sentence_builder: 'Предложение',
}

async function loadSteps() {
  const response = await getLessonSteps(lessonId)
  steps.value = response.data
}

async function moveStep(index: number, direction: number) {
  const target = index + direction

  if (target < 0 || target >= steps.value.length) return

  const next = [...steps.value]
  const [moved] = next.splice(index, 1)
  next.splice(target, 0, moved)
  steps.value = next

  try {
    await reorderLessonSteps(lessonId, next.map(s => s.id))
  } catch {
    error.value = 'Не удалось сохранить порядок шагов'
    await loadSteps()
  }
}

async function removeStep(step: AdminStep) {
  try {
    await removeLessonStep(lessonId, step.id)
    await loadSteps()
  } catch {
    error.value = 'Не удалось удалить шаг'
  }
}

const addStepOpen = ref(false)
const pickerType = ref<StepType>('word')
const pickerSearch = ref('')
const pickerResults = ref<{ id: number, label: string }[]>([])
const pickerLoading = ref(false)

function openAddStep() {
  addStepOpen.value = true
  pickerSearch.value = ''
  selectPickerType('word')
}

function selectPickerType(type: StepType) {
  pickerType.value = type
  runSearch()
}

async function runSearch() {
  pickerLoading.value = true

  try {
    const search = pickerSearch.value || undefined

    if (pickerType.value === 'word') {
      const response = await adminListWords({ search, limit: 15 })
      pickerResults.value = response.data.items.map(w => ({ id: w.id, label: `${w.hanzi} — ${w.translation}` }))
    } else if (pickerType.value === 'quiz') {
      const response = await adminListQuizzes({ search, limit: 15 })
      pickerResults.value = response.data.items.map(q => ({ id: q.id, label: q.question }))
    } else if (pickerType.value === 'grammar') {
      const response = await adminListGrammarNotes({ search, limit: 15 })
      pickerResults.value = response.data.items.map(n => ({ id: n.id, label: n.title }))
    } else {
      const response = await adminListSentences({ search, limit: 15 })
      pickerResults.value = response.data.items.map(e => ({ id: e.id, label: e.translation }))
    }
  } finally {
    pickerLoading.value = false
  }
}

watch(pickerSearch, () => {
  runSearch()
})

async function pickEntity(entityId: number) {
  try {
    await addLessonStep(lessonId, pickerType.value, entityId)
    addStepOpen.value = false
    await loadSteps()
  } catch {
    error.value = 'Не удалось добавить шаг'
  }
}

onMounted(async () => {
  loading.value = true

  try {
    const [lessonResponse] = await Promise.all([
      getLesson(lessonId),
      loadSteps(),
    ])

    lesson.value = lessonResponse.data
  } finally {
    loading.value = false
  }
})
</script>
